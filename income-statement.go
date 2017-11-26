package financialstatements

import (
	"sort"
	"strings"
	"time"
)

type IncomeStatementEntry struct {
	Balance int64      `json:"balance"`
	Details []*Balance `json:"details"`
}

type IncomeStatement struct {
	GrossRevenue          *IncomeStatementEntry `json:"grossRevenue"`
	Deduction             *IncomeStatementEntry `json:"deduction"`
	SalesTax              *IncomeStatementEntry `json:"salesTax"`
	NetRevenue            *IncomeStatementEntry `json:"netRevenue"`
	Cost                  *IncomeStatementEntry `json:"cost"`
	GrossProfit           *IncomeStatementEntry `json:"grossProfit"`
	OperatingExpense      *IncomeStatementEntry `json:"operatingExpense"`
	NetOperatingIncome    *IncomeStatementEntry `json:"netOperatingIncome"`
	NonOperatingRevenue   *IncomeStatementEntry `json:"nonOperatingRevenue"`
	NonOperatingExpense   *IncomeStatementEntry `json:"nonOperatingExpense"`
	NonOperatingTax       *IncomeStatementEntry `json:"nonOperatingTax"`
	IncomeBeforeIncomeTax *IncomeStatementEntry `json:"incomeBeforeIncomeTax"`
	IncomeTax             *IncomeStatementEntry `json:"incomeTax"`
	Dividends             *IncomeStatementEntry `json:"dividends"`
	NetIncome             *IncomeStatementEntry `json:"netIncome"`
}

func (rg *ReportGenerator) IncomeStatement(from, to time.Time) (*IncomeStatement, error) {
	e, err := rg.ds.Balances(nil, from, to)
	if err != nil {
		return nil, err
	}
	balances := []*Balance{}
	for k, v := range e {
		a, err := rg.ds.Account(k)
		if err != nil {
			return nil, err
		}
		if !a.IncomeStatement {
			continue
		}
		if !a.IncreaseOnDebit {
			v = -v
		}
		balances = append(balances, &Balance{a, v})
	}
	sort.Slice(balances, func(i, j int) bool {
		return strings.Compare(balances[i].Account.Number, balances[j].Account.Number) < 0
	})
	result := &IncomeStatement{}
	var revenueRoots, expenseRoots []*Account
	accounts, err := rg.ds.Accounts()
	if err != nil {
		return nil, err
	}
	addRoot := func(r *Account) {
		if r.IncreaseOnDebit {
			expenseRoots = append(expenseRoots, r)
		} else {
			revenueRoots = append(revenueRoots, r)
		}
	}
	for _, r := range accounts {
		if !r.IncomeStatement {
			continue
		}
		found := false
		for _, p := range accounts {
			if rg.ds.IsParent(p.Id, r.Id) {
				found = true
				break
			}
		}
		if !found {
			addRoot(r)
		}
	}
	if (len(revenueRoots) + len(expenseRoots)) == 1 {
		parentKey := append(revenueRoots, expenseRoots...)[0].Id
		revenueRoots = revenueRoots[0:0]
		expenseRoots = expenseRoots[0:0]
		for _, r := range accounts {
			if rg.ds.IsParent(parentKey, r.Id) {
				addRoot(r)
			}
		}
	}
	addBalance := func(entry *IncomeStatementEntry, balance *Balance) *IncomeStatementEntry {
		if !balance.Account.Summary && balance.Value != 0 {
			if entry == nil {
				entry = &IncomeStatementEntry{}
			}
			entry.Balance += balance.Value
			entry.Details = append(entry.Details, balance)
		}
		return entry
	}
	isDescendent := func(account *Account, parents []*Account) bool {
		for _, p := range parents {
			if p.Id == account.Id || rg.ds.IsParent(p.Id, account.Id) {
				return true
			}
		}
		return false
	}
	for _, b := range balances {
		account := b.Account
		switch account.IncomeStatementGroup {
		case "operating":
			if isDescendent(account, revenueRoots) {
				result.GrossRevenue = addBalance(result.GrossRevenue, b)
			} else if isDescendent(account, expenseRoots) {
				result.OperatingExpense = addBalance(result.OperatingExpense, b)
			}
		case "deduction":
			result.Deduction = addBalance(result.Deduction, b)
		case "salesTax":
			result.SalesTax = addBalance(result.SalesTax, b)
		case "cost":
			result.Cost = addBalance(result.Cost, b)
		case "nonOperatingTax":
			result.NonOperatingTax = addBalance(result.NonOperatingTax, b)
		case "incomeTax":
			result.IncomeTax = addBalance(result.IncomeTax, b)
		case "dividends":
			result.Dividends = addBalance(result.Dividends, b)
		default:
			if isDescendent(account, revenueRoots) {
				result.NonOperatingRevenue = addBalance(result.NonOperatingRevenue, b)
			} else {
				result.NonOperatingExpense = addBalance(result.NonOperatingExpense, b)
			}
		}
	}

	ze := &IncomeStatementEntry{}
	z := func(e *IncomeStatementEntry) *IncomeStatementEntry {
		if e == nil {
			return ze
		} else {
			return e
		}
	}

	result.NetRevenue = &IncomeStatementEntry{
		Balance: z(result.GrossRevenue).Balance - z(result.Deduction).Balance -
			z(result.SalesTax).Balance}
	result.GrossProfit = &IncomeStatementEntry{
		Balance: z(result.NetRevenue).Balance -
			z(result.Cost).Balance}
	result.NetOperatingIncome = &IncomeStatementEntry{
		Balance: z(result.GrossProfit).Balance -
			z(result.OperatingExpense).Balance}
	result.IncomeBeforeIncomeTax = &IncomeStatementEntry{
		Balance: z(result.NetOperatingIncome).Balance +
			z(result.NonOperatingRevenue).Balance -
			z(result.NonOperatingExpense).Balance - z(result.NonOperatingTax).Balance}
	result.NetIncome = &IncomeStatementEntry{
		Balance: z(result.IncomeBeforeIncomeTax).Balance -
			z(result.IncomeTax).Balance - z(result.Dividends).Balance}

	if result.NetRevenue.Balance == 0 || (z(result.Deduction).Balance == 0 &&
		z(result.SalesTax).Balance == 0) {
		result.NetRevenue = nil
	}
	if result.GrossProfit.Balance == 0 || z(result.Cost).Balance == 0 {
		result.GrossProfit = nil
	}
	if result.NetOperatingIncome.Balance == 0 {
		result.NetOperatingIncome = nil
	}
	if result.IncomeBeforeIncomeTax.Balance == 0 ||
		z(result.NonOperatingTax).Balance == 0 {
		result.IncomeBeforeIncomeTax = nil
	}
	return result, nil
}
