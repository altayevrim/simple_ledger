package main

import "fmt"

type Person struct {
	name         string
	phone        string
	balance      float64
	transactions Transactions
}

func NewPerson(name, phone string) Person {
	return Person{
		name:         name,
		phone:        phone,
		balance:      0,
		transactions: Transactions{},
	}
}

func (p *Person) NewTransaction(amount float64, description string, isDebt bool) {
	if isDebt {
		p.balance -= amount
	} else {
		p.balance += amount
	}
	p.transactions.NewTransaction(amount, description, isDebt)
}

type Transaction struct {
	amount      float64
	description string
	isDebt      bool
}

func (t Transaction) DisplayTransaction() {
	if t.isDebt {
		FPrint("\tBORÇ - ₺%.2f - %v\n", t.amount, t.description)
	} else {
		FPrint("\tALACAK - ₺%.2f - %v\n", t.amount, t.description)
	}
}

type People []Person

func (ppl *People) NewPerson(name, phone string) *People {

}

type Transactions []Transaction

func (tnxs *Transactions) NewTransaction(amount float64, description string, isDebt bool) {
	*tnxs = append(*tnxs, Transaction{
		amount:      amount,
		description: description,
		isDebt:      isDebt,
	})
}

func (tnxs Transactions) ListTransactions() {
	for _, t := range tnxs {
		t.DisplayTransaction()
	}
}

var listOfPeople People = make(People, 2, 20)

func main() {

	fmt.Println("Borç Takip v1.0 Uygulamama Hoş Geldiniz...")
	fmt.Println("Geliştiren: Evrim Altay KOLUAÇIK")

	listOfPeople[0] = NewPerson("Ferad Altılar", "+905340364488")
	listOfPeople[0].NewTransaction(500, "Deneme", false)

	menu()
}

func menu() {
	run := true
	var choice int

	for run {
		fmt.Println(`
1-) Borç Ekle
2-) Ödeme Ekle
3-) Kişi Listele
4-) Kişi Ekle
5-) Kişi Düzenle
6-) Kişi Sil
7-) Borçluları Gör
8-) Çıkış`)
		choice = GetChoice("Seçin")

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
			fmt.Println("Görüşürüz...")
		}

	}
}

func listPeopleWith(key string) {
	if len(listOfPeople) == 0 {
		EPrint("Sisteme kayıtlı hiçkimse yok!")
		return
	}
	for index, person := range listOfPeople {
		switch key {
		case "balance":
			LPrint("#", index, " - ", person.name, " (₺", person.balance, ")")
		case "transaction":
			LPrint("#", index, " - ", person.name, " (", len(person.transactions), " işlem)")
		default:
			LPrint("#", index, " - ", person.name, " (", key, ")")
		}
	}
}

func AddDebt() {

}

func AddPayment() {

}

func ListPeople() {
	listPeopleWith("balance")
	personIndex := GetChoice("Kişi Seçin")
	somethingHasFound := false
	var foundPerson Person
	for index, person := range listOfPeople {
		if index == personIndex {
			somethingHasFound = true
			foundPerson = person
			break
		}
	}

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
	} else {
		ViewPersonDetails(foundPerson)
	}
}

func ViewPersonDetails(p Person) {
	LPrint()
	FLPrint("Kişi: %v", p.name)
	FLPrint("Bakiye: %.2f", p.balance)
	LPrint("İşlemler:")
	p.transactions.ListTransactions()
	Enter2Continue()
}

func AddPerson() {

}

func EditPerson() {

}

func RemovePerson() {

}

func SeeReport() {

}

// Line Print
func LPrint(data ...any) {
	fmt.Println(data...)
}

// Error Print
func EPrint(data any) {
	fmt.Printf("** %v\n\n", data)
}

func FPrint(format string, params ...any) {
	fmt.Printf(format, params...)
}

func FLPrint(format string, params ...any) {
	fmt.Printf(format+"\n", params...)
}

func GetChoice(prompt string) (choice int) {
	if prompt != "" {
		FPrint("%v: ", prompt)
	}
	fmt.Scanln(&choice)
	return choice
}

func Enter2Continue() {
	LPrint()
	LPrint()
	LPrint("Devam etmek için enter tuşuna basın...")
	fmt.Scanln()
}
