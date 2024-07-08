//go:build !develop

package fetcher

var providerRepo = map[Spec]jobProvider{
	{Task: "backTest", ProductType: "stock"}: InitializePastStock,
}
