package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Person struct {
	Name         string
	Phone        string
	Balance      float64
	Transactions Transactions
}

func NewPerson(name, phone string) Person {
	return Person{
		Name:         name,
		Phone:        phone,
		Balance:      0,
		Transactions: Transactions{},
	}
}

func (p *Person) NewTransaction(amount float64, description string, isDebt bool) {
	if isDebt {
		p.Balance -= amount
	} else {
		p.Balance += amount
	}
	p.Transactions.NewTransaction(amount, description, isDebt)
}

type Transaction struct {
	Amount      float64
	Description string
	IsDebt      bool
	DateTime    time.Time
}

func (t Transaction) DisplayTransaction() {
	if t.IsDebt {
		FPrint("\t%v | Debt\t|\t₺%.2f | %v\n", t.DateTime.Format(time.DateTime), t.Amount, t.Description)
	} else {
		FPrint("\t%v | Payment\t|\t₺%.2f | %v\n", t.DateTime.Format(time.DateTime), t.Amount, t.Description)
	}
}

type People []Person

func (ppl *People) NewPerson(name, phone string) {
	*ppl = append(*ppl, NewPerson(name, phone))
}

type Transactions []Transaction

func (tnxs *Transactions) NewTransaction(amount float64, description string, isDebt bool) {
	*tnxs = append(*tnxs, Transaction{
		Amount:      amount,
		Description: description,
		IsDebt:      isDebt,
		DateTime:    time.Now(),
	})
}

func (tnxs Transactions) ListTransactions() {
	if len(tnxs) == 0 {
		LPrint("** The user does not have any transaction.")
		return
	}
	for _, t := range tnxs {
		t.DisplayTransaction()
	}
}

var listOfPeople People = make(People, 1, 20)

var reader = bufio.NewReader(os.Stdin)

const SaveLocation = "data.json"

var run bool = true

func main() {

	fmt.Println("Welcome to Simple Ledger v1.0 App...")
	fmt.Println("Developed by Evrim Altay KOLUAÇIK")
	fmt.Println("\t8 February 2025 - SBux Boğaçayı, Konyaaltı, Antalya / TÜRKİYE")

	listOfPeople = Load()

	menu()
}

func menu() {
	var choice int

	for run {
		fmt.Println(`
1-) Add Debt
2-) Add Payment
3-) List Users
4-) Add a User
5-) Edit a User
6-) Remove a User
7-) Reports
8-) Exit`)
		choice = GetChoice("Choose")

		switch choice {
		case 1:
			AddDebt()
		case 2:
			AddPayment()
		case 3:
			ListPeople()
		case 4:
			AddPerson()
		case 5:
			EditPerson()
		case 6:
			RemovePerson()
		case 7:
			SeeReport()
		default:
			run = false
			fmt.Println("Bye...")
		}

	}
}

func listPeopleWith(key string) {
	if len(listOfPeople) == 0 {
		EPrint("There is nobody registered to the system!")
		return
	}
	for index, person := range listOfPeople {
		switch key {
		case "balance":
			LPrint("#", index, " - ", person.Name, " (₺", person.Balance, ")")
		case "transaction":
			LPrint("#", index, " - ", person.Name, " (", len(person.Transactions), " transactions)")
		default:
			LPrint("#", index, " - ", person.Name, " (", key, ")")
		}
	}
}

func personSelector() (*Person, bool) {
	if len(listOfPeople) == 0 {
		return &Person{}, false
	}

	if len(listOfPeople) == 1 {
		return &listOfPeople[0], true
	}

	listPeopleWith("balance")
	personIndex := GetChoice("Choose a User")
	somethingHasFound := false
	var foundPerson *Person
	for index := range listOfPeople {
		if index == personIndex {
			somethingHasFound = true
			foundPerson = &listOfPeople[index]
			break
		}
	}

	return foundPerson, somethingHasFound
}

func personSelectorToRemove() (int, bool) {
	if len(listOfPeople) == 0 {
		return -1, false
	}

	if len(listOfPeople) == 1 {
		return 0, true
	}

	listPeopleWith("balance")
	personIndex := GetChoice("Choose a User to Remove")
	somethingHasFound := false
	var foundPersonIndex int
	for index := range listOfPeople {
		if index == personIndex {
			somethingHasFound = true
			foundPersonIndex = index
			break
		}
	}

	return foundPersonIndex, somethingHasFound
}

func AddDebt() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("User not found!")
		return
	}

	FLPrint("\n## Add Debt to User '%v'", foundPerson.Name)

	amount, err := GetFloat("Debt Amount")

	if err != nil {
		EPrint(err.Error() + " Cancelling...")
		return
	}

	desciption := GetString("Debt Description")

	foundPerson.NewTransaction(amount, desciption, true)
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Saving Error!")
		return
	}

	LPrint("++ Debt record has been added.")
}

func AddPayment() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("User not found!")
		return
	}

	FLPrint("\n## Add Payment to User '%v'", foundPerson.Name)

	amount, err := GetFloat("Payment Amount")

	if err != nil {
		EPrint(err.Error() + " Cancelling...")
		return
	}

	desciption := GetString("Payment Description")

	foundPerson.NewTransaction(amount, desciption, false)
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Saving Error!")
		return
	}

	LPrint("++ Payment record has been added.")
}

func ListPeople() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("User not found!")
		return
	}

	ViewPersonDetails(*foundPerson)

}

func ViewPersonDetails(p Person) {
	LPrint()
	FLPrint("User: %v", p.Name)
	FLPrint("Phone: %v", p.Phone)
	LPrint("------------------------")
	LPrint("Transactions:")
	p.Transactions.ListTransactions()
	LPrint("------------------------")
	FLPrint("Balance: %.2f", p.Balance)
	Enter2Continue()
}

func AddPerson() {
	LPrint("Please enter user details to add a new user")
	name := GetString("Name")

	if name == "" {
		EPrint("Cannot leave the name blank, cancelling...")
		return
	}

	phone := GetString("Phone")

	listOfPeople.NewPerson(name, phone)
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Saving Error!")
		return
	}

	LPrint("++ New user has been added.")
	Enter2Continue()
}

func EditPerson() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("User not found!")
		return
	}

	FLPrint("\n## Edit the User '%v'", foundPerson.Name)

	name := GetString("New Name")

	if name == "" {
		EPrint("Cannot leave the name blank, cancelling...")
		return
	}

	phone := GetString("New Phone")

	foundPerson.Name = name
	if phone != "" {
		foundPerson.Phone = phone
	}
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Saving Error!")
		return
	}

	LPrint("++ The user has been updated.")
	Enter2Continue()
}

func RemovePerson() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	foundPersonIndex, somethingHasFound := personSelectorToRemove()

	if !somethingHasFound {
		EPrint("User not found!")
		return
	}

	FLPrint("\n++ The User Will Be Removed '%v'!!!", listOfPeople[foundPersonIndex].Name)

	decision := strings.ToLower(GetString("Are you sure? (Yes / No)"))

	if decision == "e" || decision == "evet" || decision == "yes" || decision == "y" {
		listOfPeople = removePersonFromList(listOfPeople, foundPersonIndex)
		isSaved := Save(listOfPeople)

		if !isSaved {
			EPrint("Saving Error!")
			return
		}

		FLPrint("++ The user '%v' has been removed...", listOfPeople[foundPersonIndex].Name)
	} else {
		LPrint("** Cancelled...")
	}
	Enter2Continue()
}

func removePersonFromList(_listOfPeople People, pos int) People {
	if pos == 0 {
		return _listOfPeople[1:]
	} else if pos == len(_listOfPeople)-1 {
		return _listOfPeople[:len(_listOfPeople)-1]
	}

	return append(_listOfPeople[:pos], _listOfPeople[pos+1:]...)
}

func SeeReport() {
	if len(listOfPeople) == 0 {
		EPrint("There is no user in the system, please add one to continue...")
		return
	}
	LPrint("1- Balance Report")
	LPrint("2- Transaction Report")
	choice := GetChoice("Please Choose the Report Type")
	switch choice {
	case 1:
		LPrint("## Report by Balance")
		listPeopleWith("balance")
		Enter2Continue()
	case 2:
		LPrint("## Report by Transaction Count")
		listPeopleWith("transaction")
		Enter2Continue()
	default:
		EPrint("There is no such report...")
	}
}

// Line Print
func LPrint(data ...any) {
	fmt.Println(data...)
}

// Error Print
func EPrint(data any) {
	fmt.Printf("** %v\n\n", data)
	Enter2Continue()
}

// Format Print
func FPrint(format string, params ...any) {
	fmt.Printf(format, params...)
}

// Format and Line Print
func FLPrint(format string, params ...any) {
	fmt.Printf(format+"\n", params...)
}

// Plain Print
func PPrint(data ...any) {
	fmt.Print(data...)
}

func GetChoice(prompt string) (choice int) {
	if prompt != "" {
		FPrint("%v: ", prompt)
	}
	fmt.Scanln(&choice)
	return choice
}

func GetFloat(prompt string) (value float64, err error) {
	if prompt != "" {
		FPrint("%v: ₺", prompt)
	}
	fmt.Scanln(&value)

	if value <= 0 {
		return 0, errors.New(prompt + " amount cannot be 0 or negative.")
	}

	return value, nil
}

func GetString(prompt string) string {
	FPrint("%v: ", prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		EPrint(prompt + " couldn't get! Error: " + err.Error())
		return ""
	}

	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")

	return value
}

func Enter2Continue() {
	LPrint()
	LPrint()
	LPrint("Please enter to continue...")
	fmt.Scanln()
}

func Save(p People) bool {
	if len(p) == 0 {
		EPrint("Nothing to save!")
		return false
	}

	jsonContent, err := json.Marshal(p)

	if err != nil {
		EPrint("Couldn't save to the disk! JSON Error: " + err.Error())
		return false
	}

	err = os.WriteFile(SaveLocation, jsonContent, 0644)

	if err != nil {
		EPrint("Couldn't save to the disk! Write Error: " + err.Error())
		return false
	}

	return true
}

func Load() People {
	jsonContent, err := os.ReadFile(SaveLocation)

	if err != nil {
		EPrint("No save file found. Making a fresh start...")
		return People{}
	}
	var _listOfPeople People
	err = json.Unmarshal(jsonContent, &_listOfPeople)

	if err != nil {
		EPrint("Corrupt save file. The program will stop. Please backup your file '" + SaveLocation + "' and remove it in order to continue.")
		panic("Terminating...")
	}

	return _listOfPeople
}
