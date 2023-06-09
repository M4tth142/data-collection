package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"github.com/PuerkitoBio/goquery"
)

type BankAccount struct {
	AccountNumber string
	HolderName    string
	Balance       float64
}

type StockMarketInfo struct {
	Symbol string
	Price  float64
}

func fetchBankAccountDetails() ([]BankAccount, error) {
	// Connect to the database
	// Replace <database_connection_string> with the actual connection string
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/bank")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database to fetch bank account details
	rows, err := db.Query("SELECT account_number, holder_name, balance FROM bank_accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the resultset and populate bank account details
	var accounts []BankAccount
	for rows.Next() {
		var account BankAccount
		if err := rows.Scan(&account.AccountNumber, &account.HolderName, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func scrapeStockMarketInfo() ([]StockMarketInfo, error) {
	// Make an HTTP GET request to the website
	doc, err := goquery.NewDocument("https://example.com/stock-market")
	if err != nil {
		return nil, err
	}

	// Process the webpage and extract stock market information
	var stockInfo []StockMarketInfo
	doc.Find("table.stock-info tbody tr").Each(func(i int, s *goquery.Selection) {
		symbol := s.Find("td.symbol").Text()
		priceStr := s.Find("td.price").Text()
		// Convert priceStr to float64
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			// Handle parsing error
			log.Printf("Error parsing price for symbol %s: %v\n", symbol, err)
			return
		}

		stock := StockMarketInfo{
			Symbol: symbol,
			Price:  price,
		}
		stockInfo = append(stockInfo, stock)
	})

	return stockInfo, nil
}
func main() {
	// Fetch bank account details
	accounts, err := fetchBankAccountDetails()
	if err != nil {
		log.Fatal("Error fetching bank account details:", err)
	}

	fmt.Println("Bank Account Details:")
	for _, account := range accounts {
		fmt.Printf("Account Number: %s, Holder Name: %s, Balance: %.2f\n", account.AccountNumber, account.HolderName, account.Balance)
	}
	fmt.Println()

	// Scrape stock market information
	stocks, err := scrapeStockMarketInfo()
	if err != nil {
		log.Fatal("Error scraping stock market information:", err)
	}

	fmt.Println("Stock Market Information:")
	for _, stock := range stocks {
		fmt.Printf("Symbol: %s, Price: %.2f\n", stock.Symbol, stock.Price)
	}
}
