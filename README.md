# Expense-Tracker

https://roadmap.sh/projects/expense-tracker

Expense Tracker is a command-line application that helps you manage your expenses efficiently. It allows you to add, update, delete, and view expenses, as well as generate summaries.

## Features

- **Add Expense:** Add an expense with a description and amount.
- **Update Expense:** Modify an existing expense.
- **Delete Expense:** Remove an expense.
- **View Expenses:** List all expenses.
- **Expense Summary:** View a summary of all expenses.
- **Monthly Summary:** View a summary of expenses for a specific month (of the current year).

## Additional Features (Optional)

- **Expense Categories:** Add categories to expenses and filter by them.
- **Budgeting:** Set monthly budgets and receive warnings when exceeding them.
- **CSV Export:** Export expenses to a CSV file.

## Requirements

- Go (latest version recommended)
- Cobra library for command-line argument parsing
- A text editor or IDE for Go development

## Installation

1. Clone the repository:

   ```bash
      git clone https://github.com/nabobery/Expense-Tracker.git
   ```

2. Navigate to the project directory:

   ```bash
   cd Expense-Tracker
   ```

3. Install the Cobra library:

   ```bash
   go get -u github.com/spf13/cobra
   ```

4. Build the application:

   ```bash
   # For Linux/Mac
   go build -o expense-tracker

   # For Windows
   go build -o expense-tracker.exe main.go
   ```

## Usage

The application uses the following commands:

### `add`

Adds a new expense.

```bash
expense-tracker add --description "Lunch" --amount 20
```

Output:

```bash
Expense added successfully (ID: 1)
```

### `update`

Updates an existing expense.

```bash
expense-tracker update --id 1 --description "Lunch Meeting" --amount 25
```

Output:

```bash
Expense updated successfully (ID: 1)
```

### `delete`

Deletes an expense.

```bash
expense-tracker delete --id 1
```

Output:

```bash
Expense deleted successfully (ID: 1)
```

### `list`

Lists all expenses.

```bash
expense-tracker list
```

Output:

```bash
ID Date Description Amount
1 2025-01-29 Lunch Meeting $25
```

### `summary`

Provides a summary of all expenses.

```bash
expense-tracker summary
```

Output:

```bash
Total expenses: $25
```

### `summary --month`

Provides a summary of expenses for a specific month.

```bash
expense-tracker summary --month 1
```

Output:

```bash
Total expenses for January: $25
```

## Data Storage

The application stores expense data in a JSON file. Each expense has the following structure:

```json
{
  "id": 1,
  "date": "2024-08-06",
  "description": "Lunch",
  "amount": 20
}
```

## Error Handling

The application handles invalid inputs and edge cases, such as:

- Negative amounts
- Non-existent expense IDs
- Invalid month values

## Extending the Project

- Add expense categories and allow users to filter expenses by category.
- Allow users to set a budget for each month and show a warning when the user exceeds the budget
- Allow users to export expenses to a CSV file.

## Contributing

Contributions to the Expense Tracker project are welcome! If you find a bug or want to suggest an improvement, please open an issue or submit a pull request.

## License

This project is open-source and free to use under the [MIT License](LICENSE). Contributions are welcome!
