package operations

func (o *Operation) newID(operatonType string) int {
	o.lastID[operatonType] += 1
	return o.lastID[operatonType]
}
