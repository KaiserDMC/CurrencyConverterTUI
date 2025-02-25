package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/shopspring/decimal"
)

type ExchangeRates struct {
	Base               string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
}

func fetchAndSaveExchangeRates(filename string) error {
	url := "https://open.er-api.com/v6/latest/USD"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to read GET response: %v", err)
	}

	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		return fmt.Errorf("Failed to save JSON file: %v", err)
	}

	fmt.Println("Exchange rates saved to", filename)
	return nil
}

func readExchangeRates(filename string) (*ExchangeRates, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to read exchangeRates file: %v", err)
	}

	var exchangeData ExchangeRates
	err = json.Unmarshal(data, &exchangeData)
	if err != nil {
		return nil, fmt.Errorf("Failed to write json data: %v", err)
	}

	return &exchangeData, nil
}

var (
	originCurrency    string
	convertedCurrency []string
	inputText         string
)

func main() {
	filePath := "exchangeRates.json"

	wd, _ := os.Getwd()
	if _, err := os.Stat(wd + "/" + filePath); errors.Is(err, os.ErrNotExist) {
		if err := fetchAndSaveExchangeRates(filePath); err != nil {
			fmt.Println(err)
			return
		}
	}

	exchangeData, err := readExchangeRates(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check if update is needed
	currentTimeUnix := time.Now().Unix()
	if currentTimeUnix >= exchangeData.TimeNextUpdateUnix {
		fmt.Println("Exchange rates are outdated. Updating...")
		if err := fetchAndSaveExchangeRates(filePath); err != nil {
			fmt.Println(err)
			return
		}

		exchangeData, err = readExchangeRates(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Exchange rates are up to date. No update needed.")
	}

	var currencyOptions []huh.Option[string]
	for currency := range exchangeData.Rates {
		currencyOptions = append(currencyOptions, huh.NewOption(currency, currency))
	}

	form := huh.NewForm(

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose origin currency:").
				Options(
					huh.NewOption("United States Dollar", "USD"),
					huh.NewOption("Japanese Yen", "JPY"),
					huh.NewOption("Bulgarian Lev", "BGN"),
					huh.NewOption("Euro", "EUR"),
				).
				Value(&originCurrency),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose currency to convert to:").
				Options(currencyOptions...).
				Height(25).
				Filterable(true).
				Value(&convertedCurrency),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	huh.NewInput().
		Title("What amount do you wish to convert?").
		Prompt("?").
		Value(&inputText).
		Run()

	amountToConvert, err := decimal.NewFromString(inputText)
	if err != nil {
		log.Fatalf("Error parsing decimal: %v", err)
	}

	fmt.Println("Amount in decimal:", amountToConvert)
}
