# Currency Converter TUI

## Overview

The Currency Converter TUI is a terminal-based application built in Go that fetches real-time exchange rates, saves them locally, and allows users to convert an amount of money from one currency to another. The app features a user-friendly text-based interface, utilizing the tool from `charmbracelet`, such as `bubbletea`, `huh` and `lipgloss` libraries for building modern, interactive TUI (Text User Interface) applications. This tool is designed for learning purposes, but it also provides a quick and efficient currency conversion directly from the terminal.

## Features

- **Fetch Real-Time Exchange Rates:** The app fetches the latest exchange rates from a public API and saves them locally in a JSON file for offline use.

    <a href="https://www.exchangerate-api.com">Rates By Exchange Rate API</a>

- **Currency Conversion:** Converts an amount from one currency to multiple other currencies.

- **Progress Bar:** A progress bar is displayed during the fetching and saving of exchange rates, providing visual feedback to users.

- **Automatic Updates:** The app checks if the exchange rates are outdated and automatically updates them if necessary.

## How-to-use the App

### 1. Choose Origin Currency

First, select the origin currency. The current supported options are:

+ **USD** - United States Dollar
+ **JPY** - Japanese Yen
+ **BGN** - Bulgarian Lev
+ **EUR** - Euro

### 2. Select Target Currencies

You can select one or more target currencies for conversion. You’ll be provided with a list of available currencies based on the fetched exchange rates.

### 3. Input the Amount to Convert

After selecting the currencies, enter the amount of money you want to convert. The app will display the converted amounts for each selected target currency.

### 4. Confirmation

Once everything is set up, you will be asked to confirm whether you'd like to proceed with the conversion or exit the app.

## Future Improvements

- [ ] Progress bar implementation.
- [ ] Organise the target currencies in a better way.
- [ ] Make it pretty!!!
- [ ] Option to save conversion history for future reference.
- [ ] More advanced CLI features such as customizable themes and configurations.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
