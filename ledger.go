package financialstatements

import (
	"sort"
	"time"
)

type Ledger struct {
	Account *Account       `json:"account"`
	Balance int64          `json:"balance"`
	Entries []*LedgerEntry `json:"entries"`
}

type LedgerEntry struct {
	Id          string      `json:"_id"`
	Date        time.Time   `json:"date"`
	Memo        string      `json:"memo"`
	Counterpart Counterpart `json:"counterpart"`
	Balance     int64       `json:"balance"`
	created     time.Time
}

type Counterpart struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

func (rg *ReportGenerator) Ledger(accountId string, from, to time.Time) (*Ledger, error) {
	a, err := rg.ds.Account(accountId)
	if err != nil {
		return nil, err
	}
	e, err := rg.ds.Balances([]string{accountId}, time.Time{}, from.AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}
	balance := e[accountId]
	if !a.IncreaseOnDebit {
		balance = -balance
	}
	result := &Ledger{a, balance, nil}
	tch, errch := rg.ds.Transactions([]string{accountId}, from, to)
	removed := removed{}
	for t := range tch {
		if removed.found(t,
			func() int { return len(result.Entries) },
			func(i int) string { return result.Entries[i].Id },
			func(i int) { result.Entries = append(result.Entries[:i], result.Entries[i+1:]...) }) {
			continue
		}
		counterpart := Counterpart{}
		for k, v := range t.Entries {
			if sign(v) == sign(t.Entries[accountId]) {
				continue
			}
			if counterpart.Name == "" {
				c, err := rg.ds.Account(k)
				if err != nil {
					return nil, err
				}
				counterpart.Number = c.Number
				counterpart.Name = c.Name
			} else {
				counterpart.Number = ""
				counterpart.Name = "many"
			}
		}
		result.Entries = append(result.Entries,
			&LedgerEntry{t.Id, t.Date, t.Memo, counterpart, t.Entries[accountId], t.Created})
	}
	if err := <-errch; err != nil {
		return nil, err
	}
	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].created.Before(result.Entries[j].created)
	})
	runningBalance := balance
	for _, e := range result.Entries {
		if a.IncreaseOnDebit {
			runningBalance += e.Balance
		} else {
			runningBalance -= e.Balance
		}
		e.Balance = runningBalance
	}
	return result, nil
}

func sign(i int64) int64 {
	if i < 0 {
		return -1
	} else if i > 0 {
		return 1
	}
	return 0
}
