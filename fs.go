package financialstatements

import "time"

type ReportGenerator struct {
	ds DataSource
}

type DataSource interface {
	Transactions(accounts []string, from, to time.Time) (chan *Transaction, chan error)
	Balances(accounts []string, from, to time.Time) (Entries, error)
	Accounts() ([]*Account, error)
	Account(string) (*Account, error)
	IsParent(parent, child string) bool
}

type Transaction struct {
	Id      string    `json:"_id"`
	Date    time.Time `json:"date"`
	Memo    string    `json:"memo"`
	Entries Entries   `json:"entries"`
	Removes string    `json:"-"`
	Created time.Time `json:"created"`
}

type Entries map[string]int64

type Account struct {
	Id                   string `json:"_id"`
	Number               string `json:"number"`
	Name                 string `json:"name"`
	IncreaseOnDebit      bool   `json:"increaseOnDebit"`
	Summary              bool   `json:"-"`
	BalanceSheet         bool   `json:"-"`
	IncomeStatement      bool   `json:"-"`
	IncomeStatementGroup string `json:"-"`
}

func NewReportGenerator(ds DataSource) *ReportGenerator {
	return &ReportGenerator{ds}
}
