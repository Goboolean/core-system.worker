//go:build !develop

package fetcher

var providerRepo = map[Spec]jobProvider{
	{Task: "backtest", ProductType: "stock"}:      initializePastStock,
	{Task: "realtimeTrade", ProductType: "stock"}: initializeRealtimeStock,
}
