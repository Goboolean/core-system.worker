//go:build develop

package executer

var providerRepo = map[Spec]jobProvider{
	{OutputType: "candlestick"}: initalizeMock,
}
