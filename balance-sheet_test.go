package financialstatements

import (
	"testing"
)

func TestBalanceSheet(t *testing.T) {
	rg, s, d := testData()
	bs, err := rg.BalanceSheet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 2 {
		t.Error("Balance must have two entries")
	}
	if bs[0].Account.Number != "1" {
		t.Error("Balance's entry must have account number")
	}
	if bs[1].Account.Number != "2" {
		t.Error("Balance's entry must have account number")
	}
	if bs[0].Value != 1 {
		t.Error("Balance's value must be 1")
	}
	if bs[1].Value != 1 {
		t.Error("Balance's value must be 1")
	}

	s.transactions = append(s.transactions,
		&Transaction{Id: "2", Date: d, Entries: Entries{"2": 1, "1": -1}})
	bs, err = rg.BalanceSheet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 2 {
		t.Error("Balance must have two entries")
	}
	if bs[0].Account.Number != "1" {
		t.Error("Balance's entry must have account number")
	}
	if bs[1].Account.Number != "2" {
		t.Error("Balance's entry must have account number")
	}
	if bs[0].Value != 0.0 {
		t.Error("Balance's value must be 0")
	}
	if bs[1].Value != 0.0 {
		t.Error("Balance's value must be 0")
	}
	s.transactions = append(s.transactions,
		&Transaction{Id: "3", Date: d, Entries: Entries{"1": 2, "2": -2}})
	bs, err = rg.BalanceSheet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 2 {
		t.Error("Balance must have two entries")
	}
	if bs[0].Account.Number != "1" {
		t.Error("Balance's entry must have account number")
	}
	if bs[1].Account.Number != "2" {
		t.Error("Balance's entry must have account number")
	}
	if bs[0].Value != 2 {
		t.Error("Balance's value must be 2, but was", bs[0].Value)
	}
	if bs[1].Value != 2 {
		t.Error("Balance's value must be 2, but was", bs[1].Value)
	}
	s.transactions = append(s.transactions,
		&Transaction{Id: "4", Date: d, Entries: Entries{"2": 2, "1": -2}, Removes: "3"})
	bs, err = rg.BalanceSheet(d)
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) != 2 {
		t.Error("Balance must have two entries")
	}
	if bs[0].Account.Number != "1" {
		t.Error("Balance's entry must have account number")
	}
	if bs[1].Account.Number != "2" {
		t.Error("Balance's entry must have account number")
	}
	if bs[0].Value != 0 {
		t.Error("Balance's value must be 0 but was", bs[0].Value)
	}
	if bs[1].Value != 0 {
		t.Error("Balance's value must be 0 but was", bs[1].Value)
	}
}
