package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/shopspring/decimal"
)

type ExchangeRates struct {
	Base               string                     `json:"base_code"`
	Rates              map[string]decimal.Decimal `json:"rates"`
	TimeNextUpdateUnix int64                      `json:"time_next_update_unix"`
}

var (
	originCurrency    string
	convertedCurrency []string
	inputText         string
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

type model struct {
	table    table.Model
	progress progress.Model
}

type tickMsg time.Time

func fetchAndSaveExchangeRates(filename string) error {
	url := "https://open.er-api.com/v6/latest/USD"

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read GET response: %v", err)
	}

	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		return fmt.Errorf("failed to save JSON file: %v", err)
	}

	fmt.Println("Exchange rates saved to", filename)

	return nil
}

func readExchangeRates(filename string) (*ExchangeRates, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read exchangeRates file: %v", err)
	}

	var exchangeData ExchangeRates
	err = json.Unmarshal(data, &exchangeData)
	if err != nil {
		return nil, fmt.Errorf("failed to write json data: %v", err)
	}

	return &exchangeData, nil
}

func getStoragePath(filename string) string {
	// Get the absolute path of the executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		os.Exit(1)
	}

	// Get the directory where the executable is located
	exeDir := filepath.Dir(exePath)

	// Go back to project root (assuming binary is in /bin)
	projectRoot := filepath.Join(exeDir, "..")

	// Construct the full path
	return filepath.Join(projectRoot, "storage", filename)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.NewStyle().Margin(1, 2).Render(m.table.View()) + "\n(Use ↑/↓ to navigate, q to quit))"
}

func main() {
	filePath := getStoragePath("exchangeRates.json")

	//wd, _ := os.Getwd()
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
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

	var confirmExit bool

	form := huh.NewForm(

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose origin currency:").
				Options(
					huh.NewOption("USD", "USD"),
					huh.NewOption("JPY", "JPY"),
					huh.NewOption("BGN", "BGN"),
					huh.NewOption("EUR", "EUR"),
				).
				Value(&originCurrency),

			huh.NewMultiSelect[string]().
				Title("Choose currency to convert to:").
				Options(currencyOptions...).
				Height(15).
				Filterable(true).
				Value(&convertedCurrency),

			huh.NewInput().
				Title("What amount do you wish to convert?").
				Prompt("?").
				Value(&inputText),

			huh.NewConfirm().
				Title("How do you wish to proceed?").
				Affirmative("Convert").
				Negative("Exit").
				Value(&confirmExit),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	if !confirmExit {
		fmt.Println("Exiting program...")
		os.Exit(0)
	}

	amountToConvert, err := decimal.NewFromString(inputText)
	if err != nil {
		log.Fatalf("Error parsing decimal: %v", err)
	}

	// Actual conversion
	var usdRate decimal.Decimal

	switch originCurrency {
	case "USD":
		usdRate = decimal.NewFromFloat(1)
	case "JPY", "BGN", "EUR":
		usdRate = decimal.NewFromFloat(1).Div(exchangeData.Rates[originCurrency])
	default:
		log.Fatalf("Unsupported origin currency: %s", originCurrency)
	}

	var rows []table.Row
	for _, currency := range convertedCurrency {
		if rate, exists := exchangeData.Rates[currency]; exists {
			convertedValue := amountToConvert.Mul(usdRate).Mul(rate)
			rows = append(rows, table.Row{currency, convertedValue.StringFixed(2)})
		} else {
			rows = append(rows, table.Row{currency, "N/A"})
		}
	}

	columns := []table.Column{
		{Title: "Currency", Width: 10},
		{Title: "Converted Amount", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	p := tea.NewProgram(model{table: t})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
