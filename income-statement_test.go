package financialstatements

import (
	"fmt"
	"testing"
)

func TestIncomeStatement(t *testing.T) {
	rg, s, d := testData()
	s.accounts = append(s.accounts,
		&Account{Id: "3", Number: "3", Name: "Revenue", IncomeStatement: true, IncreaseOnDebit: false})
	s.accounts = append(s.accounts,
		&Account{Id: "4", Number: "4", Name: "Expenses", IncomeStatement: true, IncreaseOnDebit: true})
	s.transactions = append(s.transactions,
		&Transaction{Id: "2", Date: d, Entries: Entries{"4": 1, "1": -1}})
	s.transactions = append(s.transactions,
		&Transaction{Id: "3", Date: d, Entries: Entries{"1": 2, "3": -2}})
	is, err := rg.IncomeStatement(d, d)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range []struct{ expected, actual interface{} }{
		{false, is.NonOperatingRevenue == nil},
		{false, is.NonOperatingExpense == nil},
		{2, is.NonOperatingRevenue.Balance},
		{1, is.NonOperatingExpense.Balance},
	} {
		if fmt.Sprint(s.expected) != fmt.Sprint(s.actual) {
			t.Errorf("case %v: expected %v but was %v", i, s.expected, s.actual)
		}
	}
}
