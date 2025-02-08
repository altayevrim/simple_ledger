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
		FPrint("\t%v | Borç\t|\t₺%.2f | %v\n", t.DateTime.Format(time.DateTime), t.Amount, t.Description)
	} else {
		FPrint("\t%v | Alacak\t|\t₺%.2f | %v\n", t.DateTime.Format(time.DateTime), t.Amount, t.Description)
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
		LPrint("** Kullanıcının hiçbir işlemi yok.")
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

	fmt.Println("Borç Takip v1.0 Uygulamama Hoş Geldiniz...")
	fmt.Println("Geliştiren: Evrim Altay KOLUAÇIK")
	fmt.Println("\t8 Şubat 2025 - SBux Boğaçayı, Konyaaltı, Antalya / TÜRKİYE")

	listOfPeople = Load()

	menu()
}

func menu() {
	var choice int

	for run {
		fmt.Println(`
1-) Borç Ekle
2-) Ödeme Ekle
3-) Kişi Listele
4-) Kişi Ekle
5-) Kişi Düzenle
6-) Kişi Sil
7-) Rapor
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
		EPrint("Sisteme kayıtlı hiç kimse yok!")
		return
	}
	for index, person := range listOfPeople {
		switch key {
		case "balance":
			LPrint("#", index, " - ", person.Name, " (₺", person.Balance, ")")
		case "transaction":
			LPrint("#", index, " - ", person.Name, " (", len(person.Transactions), " işlem)")
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
	if len(listOfPeople) == 0 {
		return -1, false
	}

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
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisine Borç Kaydı Ekle", foundPerson.Name)

	amount, err := GetFloat("Borç Tutarı")

	if err != nil {
		EPrint(err.Error() + " İşlem iptal ediliyor.")
		return
	}

	desciption := GetString("Borç Açıklaması")

	foundPerson.NewTransaction(amount, desciption, true)
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Kayıt hatası!")
		return
	}

	LPrint("++ Borç kaydı başarıyla eklendi.")
}

func AddPayment() {
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisine Ödeme Kaydı Ekle", foundPerson.Name)

	amount, err := GetFloat("Ödeme Tutarı")

	if err != nil {
		EPrint(err.Error() + " İşlem iptal ediliyor.")
		return
	}

	desciption := GetString("Ödeme Açıklaması")

	foundPerson.NewTransaction(amount, desciption, false)
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Kayıt hatası!")
		return
	}

	LPrint("++ Ödeme kaydı başarıyla eklendi.")
}

func ListPeople() {
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	ViewPersonDetails(*foundPerson)

}

func ViewPersonDetails(p Person) {
	LPrint()
	FLPrint("Kişi: %v", p.Name)
	FLPrint("Telefon: %v", p.Phone)
	LPrint("------------------------")
	LPrint("İşlemler:")
	p.Transactions.ListTransactions()
	LPrint("------------------------")
	FLPrint("Bakiye: %.2f", p.Balance)
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
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Kayıt hatası!")
		return
	}

	LPrint("++Kişi başarıyla eklendi!")
	Enter2Continue()
}

func EditPerson() {
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	foundPerson, somethingHasFound := personSelector()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n##'%v' Kişisini Düzenle", foundPerson.Name)

	name := GetString("Yeni İsim")

	if name == "" {
		EPrint("İsim alanı boş geçilemez, işlem iptal edildi...")
		return
	}

	phone := GetString("Yeni Telefon")

	foundPerson.Name = name
	if phone != "" {
		foundPerson.Phone = phone
	}
	isSaved := Save(listOfPeople)

	if !isSaved {
		EPrint("Kayıt hatası!")
		return
	}

	LPrint("++ Kişi başarıyla güncellendi!")
	Enter2Continue()
}

func RemovePerson() {
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	foundPersonIndex, somethingHasFound := personSelectorToRemove()

	if !somethingHasFound {
		EPrint("Seçtiğiniz kullanıcı bulunamadı!")
		return
	}

	FLPrint("\n++'%v' Kişisi Silinecek!!!", listOfPeople[foundPersonIndex].Name)

	decision := strings.ToLower(GetString("Silme İşleminden Emin Misiniz? (Evet / Hayır)"))

	if decision == "e" || decision == "evet" {
		listOfPeople = removePersonFromList(listOfPeople, foundPersonIndex)
		isSaved := Save(listOfPeople)

		if !isSaved {
			EPrint("Kayıt hatası!")
			return
		}

		FLPrint("++ '%v' başarıyla silindi...", listOfPeople[foundPersonIndex].Name)
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
	if len(listOfPeople) == 0 {
		EPrint("Sistemde hiç kişi yok! Lütfen yeni bir kişi ekleyin.")
		return
	}
	LPrint("1- Borç Raporu")
	LPrint("2- İşlem Raporu")
	choice := GetChoice("Rapor Tipi Seçin")
	switch choice {
	case 1:
		LPrint("## Borçlara Göre Rapor")
		listPeopleWith("balance")
		Enter2Continue()
	case 2:
		LPrint("## İşlem Sayısına Göre Rapor")
		listPeopleWith("transaction")
		Enter2Continue()
	default:
		EPrint("Seçtiğiniz rapor bulunamadı...")
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

func Save(p People) bool {
	if len(p) == 0 {
		EPrint("Kayıt edilecek hiçbir şey yok!")
		return false
	}

	jsonContent, err := json.Marshal(p)

	if err != nil {
		EPrint("Değişiklikler diske yazılamadı! JSON Hatası: " + err.Error())
		return false
	}

	err = os.WriteFile(SaveLocation, jsonContent, 0644)

	if err != nil {
		EPrint("Değişiklikler diske yazılamadı! Kayıt Hatası: " + err.Error())
		return false
	}

	return true
}

func Load() People {
	jsonContent, err := os.ReadFile(SaveLocation)

	if err != nil {
		EPrint("Kayıt dosyası bulunamadı. Sıfırdan başlanacak")
		return People{}
	}
	var _listOfPeople People
	err = json.Unmarshal(jsonContent, &_listOfPeople)

	if err != nil {
		EPrint("Kayıt dosyasında bir hata var. Sistem çalışmayı sonlandıracak.")
		panic("Lütfen '" + SaveLocation + "' dosyasını silip programı yeniden çalıştırın.")
	}

	return _listOfPeople
}
