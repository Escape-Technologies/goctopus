package address

type Sourced struct {
	Address string
	Source  string
}

func NewSourced(address, source string) *Sourced {
	return &Sourced{
		Address: address,
		Source:  source,
	}
}
