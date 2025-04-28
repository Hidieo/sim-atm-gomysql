package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var currentAccountID int

func main() {
	var err error

	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/belajar_golang_db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n ")
		fmt.Println("=== ATM CLI ===")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Exit")
		fmt.Print("Choose option: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			register(reader)
		case "2":
			login(reader)
		case "3":
			fmt.Println("Terima kasih! Program selesai.")
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func register(reader *bufio.Reader) {
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("Enter PIN: ")
	pin, _ := reader.ReadString('\n')

	name = strings.TrimSpace(name)
	pin = strings.TrimSpace(pin)

	result, err := db.Exec("INSERT INTO accounts (name, pin, balance) VALUES (?, ?, 0.0)", name, pin)
	if err != nil {
		log.Println("Registration failed:", err)
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to get new account ID:", err)
		return
	}

	fmt.Println("\n=== Registration Successful ===")
	fmt.Println("Account ID :", newID)
	fmt.Println("Account Name :", name)
}

func login(reader *bufio.Reader) {
	fmt.Print("Enter account ID: ")
	idStr, _ := reader.ReadString('\n')
	fmt.Print("Enter PIN: ")
	pin, _ := reader.ReadString('\n')

	idStr = strings.TrimSpace(idStr)
	pin = strings.TrimSpace(pin)
	id, _ := strconv.Atoi(idStr)

	row := db.QueryRow("SELECT id FROM accounts WHERE id = ? AND pin = ?", id, pin)
	err := row.Scan(&currentAccountID)
	if err != nil {
		fmt.Println("Login failed")
		return
	}

	fmt.Println("Login successful!")

	for {
		fmt.Println("\n ")
		fmt.Println("\n ")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Deposit")
		fmt.Println("3. Withdraw")
		fmt.Println("4. Transfer")
		fmt.Println("5. Transaction History")
		fmt.Println("6. Logout")
		fmt.Print("Choose option: ")
		opt, _ := reader.ReadString('\n')
		opt = strings.TrimSpace(opt)

		switch opt {
		case "1":
			checkBalance()
		case "2":
			deposit(reader)
		case "3":
			withdraw(reader)
		case "4":
			transfer(reader)
		case "5":
			transactionHistory()
		case "6":
			return
		default:
			fmt.Println("Invalid option")
		}
	}
}

func checkBalance() {
	var balance float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = ?", currentAccountID).Scan(&balance)
	fmt.Printf("Current Balance: Rp %.2f\n", balance)
}

func deposit(reader *bufio.Reader) {
	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)

	_, err := db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, currentAccountID)
	if err == nil {
		db.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, 'deposit', ?)", currentAccountID, amount)
		fmt.Println("Deposit successful!")
	} else {
		fmt.Println("Deposit failed:", err)
	}
}

func withdraw(reader *bufio.Reader) {
	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)

	var balance float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = ?", currentAccountID).Scan(&balance)
	if balance < amount {
		fmt.Println("Insufficient balance")
		return
	}

	_, err := db.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, currentAccountID)
	if err == nil {
		db.Exec("INSERT INTO transactions (account_id, type, amount) VALUES (?, 'withdraw', ?)", currentAccountID, amount)
		fmt.Println("Withdraw successful!")
	} else {
		fmt.Println("Withdraw failed:", err)
	}
}

func transfer(reader *bufio.Reader) {
	fmt.Print("Enter target account ID: ")
	targetStr, _ := reader.ReadString('\n')
	fmt.Print("Enter amount: ")
	amountStr, _ := reader.ReadString('\n')
	targetID, _ := strconv.Atoi(strings.TrimSpace(targetStr))
	amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
	var balance float64
	db.QueryRow("SELECT balance FROM accounts WHERE id = ?", currentAccountID).Scan(&balance)
	if balance < amount {
		fmt.Println("Insufficient balance")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, currentAccountID)
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, targetID)
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, target_id) VALUES (?, 'transfer_out', ?, ?)", currentAccountID, amount, targetID)
	_, err = tx.Exec("INSERT INTO transactions (account_id, type, amount, target_id) VALUES (?, 'transfer_in', ?, ?)", targetID, amount, currentAccountID)

	if err != nil {
		tx.Rollback()
		fmt.Println("Transfer failed")
		return
	}
	tx.Commit()
	fmt.Println("Transfer successful!")
}

func transactionHistory() {
	fmt.Println("=== Transaction History ===")
	rows, err := db.Query(`
        SELECT type, amount, target_id, created_at
        FROM transactions
        WHERE account_id = ? AND (type = 'transfer_in' OR type = 'transfer_out')
        ORDER BY created_at DESC
    `, currentAccountID)

	if err != nil {
		fmt.Println("Error fetching transactions:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var txType string
		var amount float64
		var targetID sql.NullInt64
		var createdAt string

		err := rows.Scan(&txType, &amount, &targetID, &createdAt)
		if err != nil {
			fmt.Println("Error scanning transaction:", err)
			return
		}

		if txType == "transfer_out" {
			fmt.Printf("[OUT] Rp %.2f ➔ to Account ID %d at %s\n", amount, targetID.Int64, createdAt)
		} else if txType == "transfer_in" {
			fmt.Printf("[IN] Rp %.2f ⇦ from Account ID %d at %s\n", amount, targetID.Int64, createdAt)
		}
	}
}
