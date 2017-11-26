package financialstatements

import "time"

type JournalEntry struct {
	Id      string
	Date    time.Time
	Memo    string
	Debits  []JournalAmount
	Credits []JournalAmount
}

type JournalAmount struct {
	Account *Account
	Value   int64
}

func (rg *ReportGenerator) Journal(from, to time.Time) ([]*JournalEntry, error) {
	result := []*JournalEntry{}
	ch, errch := rg.ds.Transactions(nil, from, to)
	removed := removed{}
	for t := range ch {
		if removed.found(t,
			func() int { return len(result) },
			func(i int) string { return result[i].Id },
			func(i int) { result = append(result[:i], result[i+1:]...) }) {
			continue
		}
		je := &JournalEntry{t.Id, t.Date, t.Memo, nil, nil}
		for k, v := range t.Entries {
			a, err := rg.ds.Account(k)
			if err != nil {
				return nil, err
			}
			if v > 0 {
				je.Debits = append(je.Debits, JournalAmount{a, v})
			} else {
				je.Credits = append(je.Credits, JournalAmount{a, -v})
			}
		}
		result = append(result, je)
	}
	if err := <-errch; err != nil {
		return nil, err
	}
	return result, nil
}
