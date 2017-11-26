package financialstatements

import (
	"fmt"
	"strings"
	"time"
)

type store struct {
	accounts     []*Account
	transactions []*Transaction
}

func (s *store) Transactions(accounts []string, from, to time.Time) (chan *Transaction, chan error) {
	ch := make(chan *Transaction)
	errch := make(chan error)
	go func() {
		for _, t := range s.transactions {
			if !t.Date.Before(from) && !t.Date.After(to) {
				if len(accounts) == 0 {
					ch <- t
				} else {
				out:
					for _, a := range accounts {
						for k, _ := range t.Entries {
							if a == k {
								ch <- t
								break out
							}
						}
					}
				}
			}
		}
		close(ch)
		errch <- nil
	}()
	return ch, errch
}

func (s *store) Balances(accounts []string, from, to time.Time) (Entries, error) {
	ch, errch := s.Transactions(accounts, from, to)
	result := Entries{}
	for t := range ch {
		for k, v := range t.Entries {
			result[k] += v
		}
	}
	if err := <-errch; err != nil {
		return nil, err
	}
	return result, nil
}

func (s *store) Accounts() ([]*Account, error) {
	return s.accounts, nil
}

func (s *store) Account(id string) (*Account, error) {
	for _, a := range s.accounts {
		if a.Id == id {
			return a, nil
		}
	}
	return nil, fmt.Errorf("Account not found")
}

func (s *store) IsParent(parent, child string) bool {
	a, _ := s.Account(child)
	p, _ := s.Account(parent)
	return a.Number != p.Number && strings.HasPrefix(a.Number, p.Number)
}

func testData() (*ReportGenerator, *store, time.Time) {
	d := time.Date(2014, 5, 1, 0, 0, 0, 0, time.UTC)
	s := &store{
		[]*Account{
			&Account{Id: "1", Number: "1", Name: "Assets", BalanceSheet: true, IncreaseOnDebit: true},
			&Account{Id: "2", Number: "2", Name: "Liabilities", BalanceSheet: true, IncreaseOnDebit: false},
		},
		[]*Transaction{
			&Transaction{Id: "1", Date: d, Entries: Entries{"1": 1, "2": -1}},
		},
	}
	return NewReportGenerator(s), s, d
}
