package address

type Addr struct {
	Address string
	Source  string
}

func New(address string) *Addr {
	return &Addr{
		Address: address,
		Source:  address,
	}
}

// This allows to attach a custom source to an address
func NewSourced(address, source string) *Addr {
	return &Addr{
		Address: address,
		Source:  source,
	}
}
