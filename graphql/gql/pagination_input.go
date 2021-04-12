package gql

func (p *PaginationInput) Bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}

	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}

	return skipValue, takeValue
}
