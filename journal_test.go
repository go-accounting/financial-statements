package financialstatements

import (
	"testing"
)

func TestJournal(t *testing.T) {
	rg, s, d := testData()
	journal, err := rg.Journal(d, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(journal) != 1 {
		t.Error("Journal must have one entry")
	}
	if journal[0].Id != "1" {
		t.Error("Journal's entry must use transaction's key")
	}
	s.transactions = append(s.transactions,
		&Transaction{Id: "2", Date: d, Entries: Entries{"2": 1, "1": -1}})
	journal, err = rg.Journal(d, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(journal) != 2 {
		t.Error("Journal must have two entries")
	}
	if journal[0].Id != "1" {
		t.Error("Journal's entry must use transaction's key")
	}
	if journal[1].Id != "2" {
		t.Error("Journal's entry must use transaction's key")
	}
	journal, err = rg.Journal(d.AddDate(0, 0, 1), d.AddDate(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	if len(journal) != 0 {
		t.Error("Journal must have no entries")
	}
}
