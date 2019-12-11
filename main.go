package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"lib/sheetsService"

	"github.com/manifoldco/promptui"
)

const (
	defaultFileName = "data.csv"
	fileNameMessage = "The CSV file to read in."
	spreadsheetID   = "1q1WQ1bABbCHZBWrdEZVIZR33dSKGiIDotOPGyf38X6k"
)

//Transaction Define the data we will extract from the csv
type Transaction struct {
	Date        time.Time
	Description string
	Category    string
	Withdrawals int
	Deposits    int
	Balance     int
}

func ParseMoneyFloatStringToInt(money string) int {
	//remove any money symbols
	money = strings.Replace(money, "$", "", -1)
	//remove any commas
	money = strings.Replace(money, ",", "", -1)
	money_flt, _ := strconv.ParseFloat(money, 64)
	var money_int int = int(money_flt)
	if money_flt > 0 {
		money_int = int(money_flt * 100)
	}
	return money_int
}

func ParseCentsIntToStringMoney(money int) string {
	cents := money % 100
	dollars := money / 100
	return fmt.Sprintf("$%d.%d", dollars, cents)
}

func main() {

	fmt.Println("Welcome to Budget Sheets Manager.")
	var fileName string
	var skipHeader bool

	flag.StringVar(&fileName, "file", defaultFileName, fileNameMessage)
	flag.StringVar(&fileName, "f", defaultFileName, fileNameMessage)
	flag.BoolVar(&skipHeader, "s", true, "Flag to skip first row of the csv.")
	flag.Parse()

	fmt.Printf("You asked me to parse file: %v\n", fileName)

	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	//Consume the first row since its a header.
	if skipHeader {
		reader.Read()
	}
	var transactions []Transaction
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		date, _ := time.Parse("01/02/2006", line[0])
		description := line[1]
		withdrawals := ParseMoneyFloatStringToInt(line[2])
		deposits := ParseMoneyFloatStringToInt(line[3])
		balance := ParseMoneyFloatStringToInt(line[4])

		transactions = append(transactions, Transaction{
			Date:        date,
			Description: description,
			Withdrawals: withdrawals,
			Deposits:    deposits,
			Balance:     balance,
		})
	}
	//transactionJSON, _ := json.Marshal(transactions)
	//fmt.Println(string(transactionJSON))
	/*for _, transaction := range transactions {
		fmt.Printf("%v, %v, %v, %v, %v\n", transaction.Date, transaction.Description, transaction.Withdrawals, transaction.Deposits, transaction.Balance)
	}
	*/

	// Get sheets service
	sheetsService.GetService()

	sheetsTransactions := getTransactionsFromSheets("Transactions", "A:F")
	fmt.Println("Sanity Check returned: ", sanityCheckSheetsData(sheetsTransactions))

	transactions = getNewTransactions(transactions, sheetsTransactions)

	transactions = getCategories(transactions)
	log.Println("Have", len(transactions), "records to add to sheets.")
	for _, transaction := range transactions {
		fmt.Printf("%v, %v, %v, %v, %v, %v\n", transaction.Date, transaction.Description, transaction.Category, transaction.Withdrawals, transaction.Deposits, transaction.Balance)
	}

	//reverse order so that we can just append them and the order remains sane
	// really we should sort by date desc. But we will get there
	// TODO: Sort by date desc
	transactions = reverseOrder(transactions)

	rows := buildRowsFromTransactions(transactions)

	book, rangeStr := "Transactions", "A2:F2"
	sheetsService.AppendRows(spreadsheetID, book, rangeStr, rows)
}

func sanityCheckSheetsData(trans []Transaction) bool {
	balance := trans[0].Balance
	for idx, trans := range trans {
		if idx > 0 {
			balance = balance - trans.Withdrawals + trans.Deposits
			if balance != trans.Balance {
				return false
			}
		}
	}
	return true
}

func createNewMonthBudgetSheet(date time.Time) {

}

func buildRowsFromTransactions(trans []Transaction) [][]interface{} {
	rows := make([][]interface{}, 2)
	for _, transaction := range trans {
		rows = append(rows, []interface{}{
			transaction.Date.Format("1/2/2006"),
			transaction.Description,
			transaction.Category,
			ParseCentsIntToStringMoney(transaction.Withdrawals),
			ParseCentsIntToStringMoney(transaction.Deposits),
			ParseCentsIntToStringMoney(transaction.Balance),
		})
	}
	return rows
}

func getTransactionsFromSheets(book string, rangeStr string) []Transaction {
	//Read what transactions we already saved
	numRows, values := sheetsService.GetRows(spreadsheetID, book, rangeStr)

	var sheetsTransactions []Transaction
	log.Println("Got", numRows, "records from sheets.")
	for _, row := range values {
		if len(row) > 0 {
			date, _ := time.Parse("1/2/2006", row[0].(string))
			sheetsTransactions = append(sheetsTransactions, Transaction{
				Date:        date,
				Description: row[1].(string),
				Withdrawals: ParseMoneyFloatStringToInt(row[2].(string)),
				Deposits:    ParseMoneyFloatStringToInt(row[3].(string)),
				Balance:     ParseMoneyFloatStringToInt(row[4].(string)),
			})
		}
	}
	return sheetsTransactions
}

func getMonthNames(date time.Time) (string, string) {
	thisMonth := date.Month()
	nextMonth := thisMonth + 1
	if nextMonth > 12 {
		nextMonth = nextMonth - 12
	}
	return thisMonth.String(), nextMonth.String()
}
func getCategoryValues() []string {
	numRows, values := sheetsService.GetRows(spreadsheetID, "Categories", "A:A")
	log.Println("Got", numRows, "categories.")
	var strValues []string
	for _, value := range values {
		if len(value) > 0 {
			strValues = append(strValues, value[0].(string))
		}
	}

	return strValues
}
func dateEqual(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day()
}
func dateGT(date1, date2 time.Time) bool {
	return date1.Year() > date2.Year() && date1.Month() > date2.Month() && date1.Day() > date2.Day()
}
func dateGTE(date1, date2 time.Time) bool {
	return date1.Year() >= date2.Year() && date1.Month() >= date2.Month() && date1.Day() >= date2.Day()
}
func getMinMaxDate(data []Transaction) (time.Time, time.Time) {
	minDate, maxDate := data[0].Date, data[0].Date
	for _, value := range data {
		//If value.Date is more than maxDate, its our new max
		if dateGT(value.Date, maxDate) {
			maxDate = value.Date
		}
		//if minDate is more than value.Date, value.Date is the new min
		if dateGT(minDate, value.Date) {
			minDate = value.Date
		}
	}
	return minDate, maxDate

}
func reverseOrder(trans []Transaction) []Transaction {
	for i := len(trans)/2 - 1; i >= 0; i-- {
		opp := len(trans) - 1 - i
		trans[i], trans[opp] = trans[opp], trans[i]
	}
	return trans
}
func transactionEqual(t1, t2 Transaction) bool {
	return dateEqual(t1.Date, t2.Date) && t1.Description == t2.Description && t1.Withdrawals == t2.Withdrawals && t1.Deposits == t2.Deposits && t1.Balance == t2.Balance
}

func getCategories(transactions []Transaction) []Transaction {
	values := getCategoryValues()
	for i := range transactions {
		money := 0.0
		if transactions[i].Deposits > 0 {
			money = float64(transactions[i].Deposits) / 100.0
		}
		if transactions[i].Withdrawals > 0 {
			money = 0.0 - float64(transactions[i].Withdrawals)/100.0
		}
		thisMonth, nextMonth := getMonthNames(transactions[i].Date)
		strValues := append(values, "Income for "+thisMonth)
		strValues = append(strValues, "Income for "+nextMonth)
		promptStr := fmt.Sprintf("%v %v %v", transactions[i].Date.Format("1/2/2006"), transactions[i].Description, money)
		/*templates := promptui.SelectTemplates{
			Active:   `💰{{ .Title | green | bold }}`,
			Inactive: `{{ .Title | green }}`,
			Selected: `{{ "✔" | green | bold }} {{ "Category" | bold }}: {{ .Title | green }}`,
		}*/
		prompt := promptui.Select{
			Label: promptStr,
			Items: strValues,
			//Templates: &templates,
			Searcher: func(input string, idx int) bool {
				category := strings.ToLower(strValues[idx])
				if strings.Contains(input, " ") {
					for _, word := range strings.Split(strings.Trim(input, " "), " ") {
						if !strings.Contains(category, word) {
							return false
						}
					}
					// we made it through the loop, so we must have all the "words"
					return true
				}
				if strings.Contains(category, input) {
					return true
				}
				return false
			},
			// TODO: I would really like to figure out how to change this to allow for adding.
		}
		_, categoryStr, err := prompt.Run()
		if err != nil {
			log.Panic("Failed to set category.", err)
			break
		}
		//log.Println("Selected",categoryStr)
		transactions[i].Category = categoryStr

	}

	return transactions
}

func getTransactionsForDate(data []Transaction, date time.Time) []Transaction {
	var relevantData []Transaction
	for _, value := range data {
		if dateGTE(value.Date, date) {
			relevantData = append(relevantData, value)
		}
	}
	return relevantData
}

func getNewTransactions(loaded []Transaction, fromSheets []Transaction) []Transaction {
	//We don't have any saved data, so just give back the loaded data.
	if len(fromSheets) == 0 {
		return loaded
	}
	//Get the last date from the already saved data,
	//we only really care about stuff posted on that date and after.
	_, lastSavedDate := getMinMaxDate(fromSheets)
	relevantSavedData := getTransactionsForDate(fromSheets, lastSavedDate)
	//Empty slice to put new transactions in.
	var newData []Transaction
	for _, loadedvalue := range loaded {
		if dateEqual(loadedvalue.Date, lastSavedDate) {
			//is the data we need to compare
			alreadysaved := false
			for _, savedvalue := range relevantSavedData {
				if transactionEqual(loadedvalue, savedvalue) {
					alreadysaved = true
					break
				}
			}
			//if we haven't set alreadysaved, it wasn't in the relevantSavedData, so add it to the new
			if !alreadysaved {
				newData = append(newData, loadedvalue)
			}
		}
		if dateGT(loadedvalue.Date, lastSavedDate) {
			//save this because it is new.
			newData = append(newData, loadedvalue)
		}
	}

	return newData
}
