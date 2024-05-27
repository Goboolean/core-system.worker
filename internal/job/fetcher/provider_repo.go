//go:build !develop

package fetcher

var providerRepo = map[Spec]jobProvider{
	{Task: "backtest", ProductType: "stock"}:      initalizePastStock,
	{Task: "realtimeTrade", ProductType: "stock"}: initalizeRealtimeStock,
}
