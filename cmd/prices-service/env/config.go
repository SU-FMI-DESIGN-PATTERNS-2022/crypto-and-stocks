package env

type WebSocetConfig struct {
	CryptoURL string
	// StockURL     string
	CryptoQuotes []string
	// StockQuotes  []string
	Key    string
	Secret string
}

func LoadWebSocetConfig() WebSocetConfig {
	return WebSocetConfig{
		CryptoURL: "wss://stream.data.alpaca.markets/v1beta1/crypto",
		// StockURL:     "wss://stream.data.alpaca.markets/v2/iex",
		CryptoQuotes: []string{"BTCUSD", "ETHUSD", "ADAUSD", "DOTUSD", "USDTUSD", "SOLUSD", "MATICUSD", "LINKUSD", "ATOMUSD", "BMBUSD", "LTCUSD"},
		// StockQuotes:  []string{},
		Key:    "AK0NAAOFH7FREK1FKCXU",
		Secret: "IYUve7og11abpgu4Qv8MIGCoFTd8HdALaLg6aaZ5",
	}
}
