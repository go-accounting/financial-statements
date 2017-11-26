package financialstatements

import (
	"sort"
	"strings"
	"time"
)

type Balance struct {
	Account *Account
	Value   int64
}

func (rg *ReportGenerator) BalanceSheet(at time.Time) ([]*Balance, error) {
	e, err := rg.ds.Balances(nil, time.Time{}, at)
	if err != nil {
		return nil, err
	}
	result := []*Balance{}
	for k, v := range e {
		a, err := rg.ds.Account(k)
		if err != nil {
			return nil, err
		}
		if !a.BalanceSheet || a.Summary {
			continue
		}
		if !a.IncreaseOnDebit {
			v = -v
		}
		result = append(result, &Balance{a, v})
	}
	aa, err := rg.ds.Accounts()
	if err != nil {
		return nil, err
	}
	m := map[string]int64{}
	for _, b := range result {
		m[b.Account.Number] = b.Value
	}
	for _, a := range aa {
		if a.Summary {
			result = append(result, &Balance{a, rg.summary(a, aa, m)})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Account.Number, result[j].Account.Number) < 0
	})
	return result, nil
}

func (rg *ReportGenerator) summary(a *Account, aa []*Account, m map[string]int64) int64 {
	if v, ok := m[a.Number]; ok {
		return v
	}
	s := int64(0)
	for _, each := range aa {
		if rg.ds.IsParent(a.Id, each.Id) {
			s += rg.summary(each, aa, m)
		}
	}
	m[a.Number] = s
	return s
}
