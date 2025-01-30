package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// Expense represents a single expense entry.
type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
}

var expenses []Expense

const expensesFile = "expenses.json"

var nextID = 1

var description string
var amount float64
var id int
var month int

var rootCmd = &cobra.Command{
	Use:   "expense-tracker",
	Short: "A simple expense manager CLI application",
}

// Loadexpenses loads expenses from a JSON file
func loadExpenses() error {
	if _, err := os.Stat(expensesFile); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		return os.WriteFile(expensesFile, []byte("[]"), 0644)
	}

	data, err := os.ReadFile(expensesFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &expenses); err != nil {
		return err
	}

	// Update nextID
	maxID := 0
	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}
	nextID = maxID + 1
	return nil
}

// saves expenses to a JSON file
func saveExpenses() error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(expensesFile, data, 0644)
}

func withExpensePersistence(fn func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(cmd, args)
		if err := saveExpenses(); err != nil {
			log.Printf("Error saving expenses: %v\n", err)
		}
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new expense",
	Run: withExpensePersistence(func(cmd *cobra.Command, args []string) {
		if amount <= 0 {
			fmt.Println("Amount must be a positive value.")
			return
		}
		expense := Expense{
			ID:          nextID,
			Date:        time.Now(),
			Description: description,
			Amount:      amount,
		}
		expenses = append(expenses, expense)
		nextID++
		fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)
	}),
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing expense",
	Run: withExpensePersistence(func(cmd *cobra.Command, args []string) {
		if amount <= 0 {
			fmt.Println("Amount must be a positive value.")
			return
		}
		var updated bool
		for i, expense := range expenses {
			if expense.ID == id {
				if description != "" {
					expenses[i].Description = description
				}
				if amount != 0 {
					expenses[i].Amount = amount
				}
				updated = true
				fmt.Printf("Expense updated successfully (ID: %d)\n", id)
				break
			}
		}
		if !updated {
			fmt.Printf("Expense with ID %d not found\n", id)
		}
	}),
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an expense",
	Run: withExpensePersistence(func(cmd *cobra.Command, args []string) {
		var foundIndex = -1
		for i, expense := range expenses {
			if expense.ID == id {
				foundIndex = i
				break
			}
		}
		if foundIndex != -1 {
			expenses = append(expenses[:foundIndex], expenses[foundIndex+1:]...)
			fmt.Printf("Expense deleted successfully (ID: %d)\n", id)
		} else {
			fmt.Printf("Expense with ID %d not found\n", id)
		}
	}),
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all expenses",
	Run: func(cmd *cobra.Command, args []string) {
		if len(expenses) == 0 {
			fmt.Println("No expenses recorded yet.")
			return
		}
		fmt.Println("ID\tDate\t\tDescription\tAmount")
		for _, expense := range expenses {
			fmt.Printf("%d\t%s\t%s\t%.2f\n", expense.ID, expense.Date.Format("2006-01-02"), expense.Description, expense.Amount)
		}
	},
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarize expenses",
	Run: func(cmd *cobra.Command, args []string) {
		totalExpenses := 0.0
		if month != 0 {
			totalMonthlyExpenses := 0.0
			for _, expense := range expenses {
				if expense.Date.Month() == time.Month(month) && expense.Date.Year() == time.Now().Year() {
					totalMonthlyExpenses += expense.Amount
				}
			}
			fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(month), totalMonthlyExpenses)

		} else {
			for _, expense := range expenses {
				totalExpenses += expense.Amount
			}
			fmt.Printf("Total expenses: $%.2f\n", totalExpenses)
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(summaryCmd)

	// Add flags for the "add" command
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the expense")
	addCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Amount spent")

	// Mark flags as required
	err := addCmd.MarkFlagRequired("description")
	if err != nil {
		log.Fatalf("Error marking description flag as required: %v", err)
	}
	err = addCmd.MarkFlagRequired("amount")
	if err != nil {
		log.Fatalf("Error marking amount flag as required: %v", err)
	}

	// Flags for the "update" command
	updateCmd.Flags().IntVarP(&id, "id", "i", 0, "Expense ID to update")
	updateCmd.Flags().StringVarP(&description, "description", "d", "", "New description for the expense")
	updateCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "New amount for the expense")
	err = updateCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatalf("Error marking id flag as required for update: %v", err)
	}

	// Flags for the "delete" command
	deleteCmd.Flags().IntVarP(&id, "id", "i", 0, "Expense ID to delete")
	err = deleteCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatalf("Error marking id flag as required for delete: %v", err)
	}

	// Flags for the "summary" command
	summaryCmd.Flags().IntVarP(&month, "month", "m", 0, "Month to summarize expenses for (1-12)")
}

func main() {
	if err := loadExpenses(); err != nil {
		log.Fatalf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
