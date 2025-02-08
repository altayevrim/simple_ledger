package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

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

func (ppl *People) NewPerson(name, phone string) {
	*ppl = append(*ppl, NewPerson(name, phone))
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
	if len(tnxs) == 0 {
		LPrint("** Kullanıcının hiçbir işlemi yok.")
		return
	}
	for _, t := range tnxs {
		t.DisplayTransaction()
	}
}

var listOfPeople People = make(People, 1, 20)

var reader = bufio.NewReader(os.Stdin)

func main() {

	fmt.Println("Borç Takip v1.0 Uygulamama Hoş Geldiniz...")
	fmt.Println("Geliştiren: Evrim Altay KOLUAÇIK")
	fmt.Println("\t8 Şubat 2025 - SBux Boğaçayı, Konyaaltı, Antalya / TÜRKİYE")

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

func personSelector() (*Person, bool) {
	if len(listOfPeople) == 1 {
		return &listOfPeople[0], true
	}

	listPeopleWith("balance")
	personIndex := GetChoice("Kişi Seçin")
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
	if len(listOfPeople) == 1 {
		return 0, true
	}

	listPeopleWith("balance")
	personIndex := GetChoice("Silmek için Kişi Seçin")
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
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisine Borç Kaydı Ekle", foundPerson.name)

	amount, err := GetFloat("Borç Tutarı")

	if err != nil {
		EPrint(err.Error() + " İşlem iptal ediliyor.")
		return
	}

	desciption := GetString("Borç Açıklaması")

	foundPerson.NewTransaction(amount, desciption, true)
}

func AddPayment() {
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisine Ödeme Kaydı Ekle", foundPerson.name)

	amount, err := GetFloat("Ödeme Tutarı")

	if err != nil {
		EPrint(err.Error() + " İşlem iptal ediliyor.")
		return
	}

	desciption := GetString("Ödeme Açıklaması")

	foundPerson.NewTransaction(amount, desciption, false)
}

func ListPeople() {
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	ViewPersonDetails(*foundPerson)

}

func ViewPersonDetails(p Person) {
	LPrint()
	FLPrint("Kişi: %v", p.name)
	FLPrint("Telefon: %v", p.phone)
	FLPrint("Bakiye: %.2f", p.balance)
	LPrint("İşlemler:")
	p.transactions.ListTransactions()
	Enter2Continue()
}

func AddPerson() {
	LPrint("Kişi eklemek için sırasıyla isim ve telefon numarası girin")
	name := GetString("İsim")

	if name == "" {
		EPrint("İsim boş geçilemez, işlem iptal edildi...")
		return
	}

	phone := GetString("Telefon")

	listOfPeople.NewPerson(name, phone)
	LPrint("Kişi başarıyla eklendi!")
	Enter2Continue()
}

func EditPerson() {
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisini Düzenle", foundPerson.name)

	name := GetString("Yeni İsim")

	if name == "" {
		EPrint("İsim alanı boş geçilemez, işlem iptal edildi...")
		return
	}

	phone := GetString("Yeni Telefon")

	foundPerson.name = name
	if phone != "" {
		foundPerson.phone = phone
	}

	LPrint("Kişi başarıyla güncellendi!")
	Enter2Continue()
}

func RemovePerson() {
	foundPersonIndex, somethingHasFound := personSelectorToRemove()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n++'%v' Kişisi Silinecek!!!", listOfPeople[foundPersonIndex].name)

	decision := strings.ToLower(GetString("Silme İşleminden Emin Misiniz? (Evet / Hayır)"))

	if decision == "e" || decision == "evet" {
		FLPrint("'%v' başarıyla silindi...", listOfPeople[foundPersonIndex].name)
		listOfPeople = removePersonFromList(listOfPeople, foundPersonIndex)
	} else {
		LPrint("** İşlem iptal edildi...")
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
		return 0, errors.New(prompt + " değeri 0 veya negatif olamaz.")
	}

	return value, nil
}

func GetString(prompt string) string {
	FPrint("%v: ", prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		EPrint(prompt + " Alınamadı! Hata: " + err.Error())
		return ""
	}

	value = strings.ReplaceAll(value, "\n", "")
	value = strings.ReplaceAll(value, "\r", "")

	return value
}

func Enter2Continue() {
	LPrint()
	LPrint()
	LPrint("Devam etmek için enter tuşuna basın...")
	fmt.Scanln()
}
