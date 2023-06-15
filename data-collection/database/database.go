package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUsername = "glenncolombie"
	dbPassword = "Nestrix123"
	dbHost     = "db4free.net"
	dbPort     = "3306"
	dbName     = "nestrixdb"
	name       = "jan"
)

type Rekening struct {
	RekeningNummer string
	UserID         string
	Saldo          float64
}

type Transactie struct {
	ID       string
	Receiver string
	Date     string
	Bedrag   float64
}

type Gebruiker struct {
	ID        string
	LastName  string
	FirstName string
	BirthDate string
}

// GetRekeningData returns a greeting for the named person.
func GetRekeningData() string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

func GetAllData() {
	fmt.Println("Go MySQL inside getall")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user, err := getUser(db, name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
		} else {
			fmt.Println("Error retrieving user:", err)
		}
		return
	}

	fmt.Println("User:")
	fmt.Printf("%+v\n", user)

	accounts, err := getRekening(db, user.ID)
	if err != nil {
		fmt.Println("Error retrieving accounts:", err)
		return
	}

	fmt.Println("Accounts:")
	for _, account := range accounts {
		fmt.Printf("%+v\n", account)
	}

	for _, account := range accounts {
		transactions, err := getTransactions(db, user.ID)
		if err != nil {
			fmt.Println("Error retrieving transactions:", err)
			return
		}

		fmt.Printf("Transactions for Account %s:\n", account.RekeningNummer)
		if len(transactions) == 0 {
			fmt.Println("No transactions found.")
		} else {
			for _, transaction := range transactions {
				fmt.Printf("%+v\n", transaction)
			}
		}
	}

}

func getRekening(db *sql.DB, userId string) ([]Rekening, error) {
	rows, err := db.Query("SELECT Rekeningnummer, GebruikerId, Saldo FROM Rekening WHERE GebruikerId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Rekeningen := []Rekening{}
	for rows.Next() {
		Rekening := Rekening{}
		err := rows.Scan(
			&Rekening.RekeningNummer,
			&Rekening.UserID,
			&Rekening.Saldo,
		)
		if err != nil {
			return nil, err
		}
		Rekeningen = append(Rekeningen, Rekening)
	}

	return Rekeningen, nil
}

func getTransactions(db *sql.DB, rekeningNummer string) ([]Transactie, error) {

	rows, err := db.Query("SELECT Id, RekeningnummerBegunstigde, Datum, Bedrag FROM Transactie WHERE Id = ?", rekeningNummer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []Transactie{}
	for rows.Next() {
		transaction := Transactie{}
		err := rows.Scan(
			&transaction.ID,
			&transaction.Receiver,
			&transaction.Date,
			&transaction.Bedrag,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func getUser(db *sql.DB, loginName string) (Gebruiker, error) {
	query := "SELECT Id, Voornaam, Familienaam, Geboortedatum FROM Gebruiker WHERE Voornaam = ?"
	row := db.QueryRow(query, loginName)

	user := Gebruiker{}
	var firstName, lastName sql.NullString
	err := row.Scan(
		&user.ID,
		&firstName,
		&lastName,
		&user.BirthDate,
	)
	if err != nil {
		return Gebruiker{}, err
	}

	// Check if user exists
	if lastName.Valid {
		// User exists, assign the values
		user.LastName = lastName.String
	} else {
		// User does not exist, return an error
		return Gebruiker{}, fmt.Errorf("User not found in database")
	}

	// Check if Voornaam is NULL
	if firstName.Valid {
		// Voornaam is not NULL, use the value
		user.FirstName = firstName.String
	} else {
		// Voornaam is NULL, assign a blank
		user.FirstName = ""
	}

	return user, nil
}
