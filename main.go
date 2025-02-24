package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/shopspring/decimal"
)

type ExchangeRates struct {
	Base  string             `json:"base_code"`
	Rates map[string]float64 `json:"rates"`
}

func fetchAndSaveExchangeRates(filename string) error {
	url := "https://open.er-api.com/v6/latest/USD"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read GET response: %v", err)
	}

	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		log.Fatalf("Failed to save JSON file: %v", err)
	}

	fmt.Println("Exchange rates saved to", filename)
	return nil
}

var (
	originCurrency    string
	convertedCurrency string
)

func main() {
	filePath := "exchangeRates.json"

	f := fetchAndSaveExchangeRates(filePath)
	if f != nil {
		fmt.Println(f)
		return
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
			huh.NewSelect[string]().
				Title("Choose currency to convert to:").
				Options(
					huh.NewOption("United States Dollar", "USD"),
					huh.NewOption("Japanese Yen", "JPY"),
					huh.NewOption("Bulgarian Lev", "BGN"),
					huh.NewOption("Euro", "EUR"),
				).
				Value(&convertedCurrency),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter amount to convert: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}

	text = strings.TrimSpace(text)

	amountToConvert, err := decimal.NewFromString(text)
	if err != nil {
		log.Fatalf("Error parsing decimal: %v", err)
	}

	fmt.Println("Amount in decimal:", amountToConvert)
}
