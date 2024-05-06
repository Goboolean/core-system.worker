package mongo

type StockDocument struct {
	Symbol    string
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Average   float32
	Volume    int64
	Timestamp int64
}
