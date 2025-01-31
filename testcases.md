# Test Cases

## Add expenses

expense-tracker add --description "Groceries" --amount 45.50 --category "Food"
expense-tracker add --description "Netflix" --amount 15.99 --category "Entertainment"
expense-tracker add --description "Gas" --amount 60.00 --category "Transport"

## List all

expense-tracker list

## Update first expense

expense-tracker update --id 1 --description "Weekly Groceries" --amount 50.00

## Delete second expense

expense-tracker delete --id 2

## List after changes

expense-tracker list

## Set budget for current month

expense-tracker budget --month 1 --amount 500

## Set budget for next year

expense-tracker budget --month 3 --year 2026 --amount 300

## Try invalid budget

expense-tracker budget --month 13 --amount 100

## Total summary

expense-tracker summary

## Monthly summary

expense-tracker summary --month 1

## Category filter

expense-tracker summary --category "Transport"

## Budget check

expense-tracker summary --month 1 --budget 100

## CSV export

expense-tracker export --file expenses.csv

## Invalid amount

expense-tracker add --description "Test" --amount -5

## Missing required field

expense-tracker add --description "Missing amount"

## Non-existent ID

expense-tracker update --id 999 --description "Ghost expense"
expense-tracker delete --id 999

## Invalid month filter

expense-tracker summary --month 13
