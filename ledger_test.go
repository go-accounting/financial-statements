package financialstatements

import (
	"strconv"
	"testing"
)

func TestLedger(t *testing.T) {
	rg, s, d := testData()
	ledger, err := rg.Ledger("1", d, d)
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Account.Name != "Assets" {
		t.Error("Account name must be 'Assets'")
	}
	if len(ledger.Entries) != 1 {
		t.Errorf("Ledger must have one entry, but was %v", len(ledger.Entries))
	} else {
		entry := ledger.Entries[0]
		if entry.Counterpart.Number != "2" {
			t.Error("Counterpart must be account #2")
		}
		if entry.Balance != 1 {
			t.Error("Entry's balance must be 1")
		}
	}
	if ledger.Balance != 0 {
		t.Error("Ledger's balance must be 0")
	}

	ledger, err = rg.Ledger("1", d.AddDate(0, 0, 1), d.AddDate(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Account.Name != "Assets" {
		t.Error("Account name must be 'Assets'")
	}
	if len(ledger.Entries) != 0 {
		t.Error("Ledger must have zero entries")
	}
	if ledger.Balance != 1 {
		t.Errorf("Ledger's balance must be 1, but was %v", ledger.Balance)
	}
	for i := 0; i < 4; i++ {
		s.transactions = append(s.transactions,
			&Transaction{Id: strconv.Itoa(i + 2), Date: d, Entries: Entries{"1": 1, "2": -1}})
	}
	ledger, err = rg.Ledger("1", d, d.AddDate(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	if l := len(ledger.Entries); l != 5 {
		t.Errorf("Ledger must have five entries and have %v", l)
	}
}
