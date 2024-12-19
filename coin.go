package monitor

// List of supported cryptocurrencies.
const (
	OSMO Coin = iota
	USD
)

// Coin represents a cryptocurrency.
type Coin int

// String returns the string representation of the Coin.
func (c Coin) String() string {
	switch c {
	case OSMO:
		return "osmo"
	case USD:
		return "usd"
	}
	return ""
}

// Pair represents a cryptocurrency pair.
type Pair struct {
	Base  Coin // Base Coin
	Quote Coin // Quote Coin
}

// Pairs is a list of cryptocurrency pairs.
type Pairs []Pair
