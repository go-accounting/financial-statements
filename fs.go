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
	Id      string
	Date    time.Time
	Memo    string
	Entries Entries
	Removes string
	Created time.Time
}

type Entries map[string]int64

type Account struct {
	Id                   string
	Number               string
	Name                 string
	IncreaseOnDebit      bool
	Summary              bool
	BalanceSheet         bool
	IncomeStatement      bool
	IncomeStatementGroup string
}

func NewReportGenerator(ds DataSource) *ReportGenerator {
	return &ReportGenerator{ds}
}
