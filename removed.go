package financialstatements

type removed map[string]string

func (r removed) found(t *Transaction, len func() int, id func(int) string, remove func(int)) bool {
	if r[t.Id] != "" {
		delete(r, t.Id)
		return true
	}
	if t.Removes != "" {
		found := false
		for i := 0; i < len(); i++ {
			if id(i) == t.Removes {
				remove(i)
				found = true
				break
			}
		}
		if !found {
			r[t.Id] = t.Id
		}
		return true
	}
	return false
}
