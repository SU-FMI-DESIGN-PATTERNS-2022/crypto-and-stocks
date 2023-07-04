package env_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
)

var _ = Describe("Config", func() {
	const invalid = "invalid"

	Context("LoadServerConfig", func() {
		const serverPort = "PRICES_SERVER_PORT"

		AfterEach(func() {
			os.Unsetenv(serverPort)
		})

		When("invalid port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(serverPort, invalid)
			})

			It("should return an error", func() {
				_, err := env.LoadServerConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load server config"))
			})
		})

		When("no port env is provided", func() {
			It(`should return the default port "8081" and no error`, func() {
				serverConfig, err := env.LoadServerConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(serverConfig.Port).To(Equal(8081))
			})
		})

		When("port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(serverPort, "8080")
			})

			It("should return the value that is provided and no error", func() {
				serverConfig, err := env.LoadServerConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(serverConfig.Port).To(Equal(8080))
			})
		})
	})

	Context("LoadWebSocetConfig", func() {
		const (
			wsCryptoURL = "WS_CRYPTO_URL"
			wsStockURL  = "WS_STOCK_URL"
			wsKey       = "WS_KEY"
			wsSecret    = "WS_SECRET"
		)

		BeforeEach(func() {
			os.Setenv(wsCryptoURL, "crypto-url")
			os.Setenv(wsStockURL, "stock-url")
			os.Setenv(wsKey, "key")
			os.Setenv(wsSecret, "secret")
		})

		When("no crypto url env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(wsCryptoURL)
			})

			It(`should return the default crypto url "wss://stream.data.alpaca.markets/v1beta1/crypto" and no error`, func() {
				wsConfig, err := env.LoadWebSocetConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(wsConfig.CryptoURL).To(Equal("wss://stream.data.alpaca.markets/v1beta1/crypto"))
			})
		})

		When("no stock url env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(wsStockURL)
			})

			It(`should return the default stock url "wss://stream.data.alpaca.markets/v2/iex" and no error`, func() {
				wsConfig, err := env.LoadWebSocetConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(wsConfig.StockURL).To(Equal("wss://stream.data.alpaca.markets/v2/iex"))
			})
		})

		When("no key env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(wsKey)
			})

			It("should return an error", func() {
				_, err := env.LoadWebSocetConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load web socket config"))
			})
		})

		When("no secret env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(wsSecret)
			})

			It("should return an error", func() {
				_, err := env.LoadWebSocetConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load web socket config"))
			})
		})

		When("every env is provided", func() {
			It("should return the values that are provided and no error", func() {
				wsConfig, err := env.LoadWebSocetConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(wsConfig.CryptoURL).To(Equal("crypto-url"))
				Expect(wsConfig.StockURL).To(Equal("stock-url"))
				Expect(wsConfig.CryptoQuotes).To(Equal([]string{"BTCUSD", "ETHUSD", "ADAUSD", "DOTUSD", "USDTUSD", "SOLUSD", "MATICUSD", "LINKUSD", "ATOMUSD", "BMBUSD", "LTCUSD"}))
				Expect(wsConfig.StockQuotes).To(Equal([]string{"AAPL", "AMZN"}))
				Expect(wsConfig.Key).To(Equal("key"))
				Expect(wsConfig.Secret).To(Equal("secret"))
			})
		})
	})
})
