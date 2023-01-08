package tv

// SocketMessage ...
type SocketMessage struct {
	Message string      `json:"m"`
	Payload interface{} `json:"p"`
}

// QuoteMessage ...
type QuoteMessage struct {
	Symbol string     `mapstructure:"n"`
	Status string     `mapstructure:"s"`
	Data   *QuoteData `mapstructure:"v"`
}

// QuoteData ...
type QuoteData struct {
	Price  *float64 `mapstructure:"lp"`
	Volume *float64 `mapstructure:"volume"`
	Bid    *float64 `mapstructure:"bid"`
	Ask    *float64 `mapstructure:"ask"`
}

// Flags ...
type Flags struct {
	Flags []string `json:"flags"`
}

type Line struct {
	ID          int64
	Pair        string
	Val         float64
	Description string
	Typ         string
	Timeframe   string
}
