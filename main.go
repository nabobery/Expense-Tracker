package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"time"
)

// Expense represents a single expense entry.
type Expense struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category,omitempty"`
}

// Budget represents a monthly budget
type Budget struct {
	Month  int     `json:"month"`
	Year   int     `json:"year"`
	Amount float64 `json:"amount"`
}

var expenses []Expense

var budgets []Budget

const expensesFile = "expenses.json"

const budgetsFile = "budgets.json"

var nextID = 1

var description string
var amount float64
var id int
var month int
var year int
var category string
var budgetAmount float64
var exportFile string

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

// loadBudgets loads budgets from a JSON file
func loadBudgets() error {
	if _, err := os.Stat(budgetsFile); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		return os.WriteFile(budgetsFile, []byte("[]"), 0644)
	}

	data, err := os.ReadFile(budgetsFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &budgets); err != nil {
		return err
	}
	return nil
}

// saves budgets to a JSON file
func saveBudgets() error {
	data, err := json.MarshalIndent(budgets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(budgetsFile, data, 0644)
}

func withExpensePersistence(fn func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(cmd, args)
		if err := saveExpenses(); err != nil {
			log.Printf("Error saving expenses: %v\n", err)
		}
	}
}

func withBudgetPersistence(fn func(cmd *cobra.Command, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(cmd, args)
		if err := saveBudgets(); err != nil {
			log.Printf("Error saving budgets: %v\n", err)
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
			Category:    category,
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
				if category != "" {
					expenses[i].Category = category
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
		var filteredExpenses []Expense
		if category != "" {
			for _, expense := range expenses {
				if expense.Category == category {
					filteredExpenses = append(filteredExpenses, expense)
				}
			}
		} else {
			filteredExpenses = expenses
		}

		if month != 0 {
			if month < 1 || month > 12 {
				fmt.Println("Invalid month. Please enter a value between 1 and 12.")
				return
			}
			totalMonthlyExpenses := 0.0
			currentYear := time.Now().Year()
			for _, expense := range filteredExpenses {
				if expense.Date.Month() == time.Month(month) && expense.Date.Year() == time.Now().Year() {
					totalMonthlyExpenses += expense.Amount
				}
			}
			fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(month), totalMonthlyExpenses)
			// Check against both stored budget and provided budget flag
			if budgetAmount > 0 {
				if totalMonthlyExpenses > budgetAmount {
					fmt.Printf("Warning: Expenses exceed provided budget of $%.2f for %s\n", budgetAmount, time.Month(month))
				}
			} else if currentBudget := getBudget(month, currentYear); currentBudget != nil {
				if totalMonthlyExpenses > currentBudget.Amount {
					fmt.Printf("Warning: You have exceeded your stored budget of $%.2f for %s\n", currentBudget.Amount, time.Month(month))
				}
			}

		} else {
			for _, expense := range filteredExpenses {
				totalExpenses += expense.Amount
			}
			fmt.Printf("Total expenses: $%.2f\n", totalExpenses)
		}

	},
}

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Set or view monthly budget",
	Run: withBudgetPersistence(func(cmd *cobra.Command, args []string) {
		if budgetAmount <= 0 {
			fmt.Println("Budget amount must be a positive value.")
			return
		}
		if month == 0 {
			fmt.Println("Month is required to set a budget.")
			return
		}
		if month < 1 || month > 12 {
			fmt.Println("Invalid month. Please enter a value between 1 and 12.")
			return
		}
		if year == 0 {
			year = time.Now().Year()
		}
		budget := Budget{
			Month:  month,
			Year:   year,
			Amount: budgetAmount,
		}
		setBudget(budget)
		fmt.Printf("Budget set successfully for %s %d: $%.2f\n", time.Month(month), year, budgetAmount)
	}),
}

func getBudget(month int, year int) *Budget {
	for _, budget := range budgets {
		if budget.Month == month && budget.Year == year {
			return &budget
		}
	}
	return nil
}

func setBudget(budget Budget) {
	for i, b := range budgets {
		if b.Month == budget.Month && b.Year == budget.Year {
			budgets[i] = budget
			return
		}
	}
	budgets = append(budgets, budget)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export expenses to a CSV file",
	Run: func(cmd *cobra.Command, args []string) {
		if exportFile == "" {
			fmt.Println("Please specify an export file using --file")
			return
		}
		file, err := os.Create(exportFile)
		if err != nil {
			log.Fatalf("Could not create export file: %v", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"ID", "Date", "Description", "Amount", "Category"}
		if err := writer.Write(header); err != nil {
			log.Fatalf("Could not write header to csv: %v", err)
		}

		for _, expense := range expenses {
			row := []string{
				strconv.Itoa(expense.ID),
				expense.Date.Format("2006-01-02"),
				expense.Description,
				strconv.FormatFloat(expense.Amount, 'f', 2, 64),
				expense.Category,
			}
			if err := writer.Write(row); err != nil {
				log.Fatalf("Could not write row to csv: %v", err)
			}
		}
		fmt.Printf("Expenses exported to %s\n", exportFile)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(budgetCmd)
	rootCmd.AddCommand(exportCmd)

	// Add flags for the "add" command
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the expense")
	addCmd.Flags().Float64VarP(&amount, "amount", "a", 0, "Amount spent")
	addCmd.Flags().StringVarP(&category, "category", "c", "", "Category of the expense")

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
	updateCmd.Flags().StringVarP(&category, "category", "c", "", "New category for the expense")
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

	summaryCmd.Flags().IntVarP(&month, "month", "m", 0, "Month to to filter expenses by (1-12)")
	summaryCmd.Flags().StringVarP(&category, "category", "c", "", "Category to filter expenses by")
	summaryCmd.Flags().Float64VarP(&budgetAmount, "budget", "b", 0, "Budget amount to check against")

	// Flags for the "budget" command
	budgetCmd.Flags().IntVarP(&month, "month", "m", 0, "Month to set budget for (1-12)")
	budgetCmd.Flags().IntVarP(&year, "year", "y", 0, "Year to set budget for")
	budgetCmd.Flags().Float64VarP(&budgetAmount, "amount", "a", 0, "Budget amount")
	err = budgetCmd.MarkFlagRequired("amount")
	if err != nil {
		log.Fatalf("Error marking amount flag as required for budget: %v", err)
	}

	// Flags for the "export" command
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "File to export expenses to")
	err = exportCmd.MarkFlagRequired("file")
	if err != nil {
		log.Fatalf("Error marking file flag as required for export: %v", err)
	}
}

func main() {
	if err := loadExpenses(); err != nil {
		log.Fatalf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}
	if err := loadBudgets(); err != nil {
		log.Fatalf("Error loading budgets: %v\n", err)
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
