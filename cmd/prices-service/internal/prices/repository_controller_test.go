package prices_test

import (
	"errors"
	"time"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	mock_prices "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices/mocks"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RepositoryController", func() {
	const (
		cryptoTopic     = "crypto"
		stocksTopic     = "stocks"
		subscribeErrMsg = "subscribe failed"
		errMsg          = "some error"
	)

	var (
		ctrl           *gomock.Controller
		mockCryptoRepo *mock_prices.MockPricesRepository[database.CryptoPrices]
		mockStocksRepo *mock_prices.MockPricesRepository[database.StockPrices]
		mockEventBus   *mock_prices.MockEventBus
		repoController *prices.RepositoryController
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockCryptoRepo = mock_prices.NewMockPricesRepository[database.CryptoPrices](ctrl)
		mockStocksRepo = mock_prices.NewMockPricesRepository[database.StockPrices](ctrl)
		mockEventBus = mock_prices.NewMockEventBus(ctrl)

		repoController = prices.NewRepositoryController(mockCryptoRepo, mockStocksRepo)
	})

	Context("ListenForStoring", func() {
		When("subscribing to crypto messages fails", func() {
			BeforeEach(func() {
				mockEventBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).Return(errors.New(subscribeErrMsg))
			})

			It("should return an error", func() {
				err := repoController.ListenForStoring(mockEventBus)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(subscribeErrMsg))
				Expect(err.Error()).To(ContainSubstring("crypto"))
			})
		})

		When("subscribing to stocks messages fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockEventBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).Return(nil),
					mockEventBus.EXPECT().Subscribe(stocksTopic, gomock.Any()).Return(errors.New(subscribeErrMsg)),
				)
			})

			It("should return an error", func() {
				err := repoController.ListenForStoring(mockEventBus)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(subscribeErrMsg))
				Expect(err.Error()).To(ContainSubstring("stocks"))
			})
		})

		When("subscribing to crypto and stocks messages succeed and handiling from event bus messages of type", func() {
			var (
				cryptoResp stream.CryptoResponse
				stockResp  stream.StockResponse
			)

			BeforeEach(func() {
				cryptoResp = getCryptoResponse()
				stockResp = getStockResponse()

				gomock.InOrder(
					mockEventBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(stream.CryptoResponse))
							Expect(ok).To(BeTrue())
							f(cryptoResp)
							return nil
						}),

					mockEventBus.EXPECT().Subscribe(stocksTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(stream.StockResponse))
							Expect(ok).To(BeTrue())
							f(stockResp)
							return nil
						}),
				)
			})

			Context("crypto", func() {
				When("storing the converted price fails", func() {
					BeforeEach(func() {
						mockCryptoRepo.EXPECT().StoreEntry(convertToCryptoPrice(cryptoResp)).Return(errors.New(errMsg))
						mockStocksRepo.EXPECT().StoreEntry(gomock.Any()).Return(nil)
					})

					It("should not return an error", func() {
						err := repoController.ListenForStoring(mockEventBus)

						Expect(err).ToNot(HaveOccurred())
					})
				})

				When("storing the converted price succeed", func() {
					BeforeEach(func() {
						mockCryptoRepo.EXPECT().StoreEntry(convertToCryptoPrice(cryptoResp)).Return(nil)
						mockStocksRepo.EXPECT().StoreEntry(gomock.Any()).Return(nil)
					})

					It("should not return an error", func() {
						err := repoController.ListenForStoring(mockEventBus)

						Expect(err).ToNot(HaveOccurred())
					})
				})
			})

			Context("stocks", func() {
				When("storing the converted price fails", func() {
					BeforeEach(func() {
						mockStocksRepo.EXPECT().StoreEntry(convertToStockPrice(stockResp)).Return(errors.New(errMsg))
						mockCryptoRepo.EXPECT().StoreEntry(gomock.Any()).Return(nil)
					})

					It("should not return an error", func() {
						err := repoController.ListenForStoring(mockEventBus)

						Expect(err).ToNot(HaveOccurred())
					})
				})

				When("storing the converted price succeed", func() {
					BeforeEach(func() {
						mockStocksRepo.EXPECT().StoreEntry(convertToStockPrice(stockResp)).Return(nil)
						mockCryptoRepo.EXPECT().StoreEntry(gomock.Any()).Return(nil)
					})

					It("should not return an error", func() {
						err := repoController.ListenForStoring(mockEventBus)

						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})
	})
})

func getCryptoResponse() stream.CryptoResponse {
	return stream.CryptoResponse{
		Response: stream.Response{
			Symbol: "symbol",
			Date:   time.Now(),
		},
		Exchange: "excahnge",
	}
}

func convertToCryptoPrice(cryptoResp stream.CryptoResponse) database.CryptoPrices {
	return database.CryptoPrices{
		Prices: database.Prices{
			Symbol:   cryptoResp.Symbol,
			BidPrice: cryptoResp.BidPrice,
			BidSize:  cryptoResp.BidSize,
			AskPrice: cryptoResp.AskPrice,
			AskSize:  cryptoResp.AskSize,
			Date:     cryptoResp.Date,
		},
		Exchange: cryptoResp.Exchange,
	}
}

func getStockResponse() stream.StockResponse {
	return stream.StockResponse{
		Response: stream.Response{
			Symbol: "symbol",
			Date:   time.Now(),
		},
		AskExchange: "ask exchange",
		BidExchange: "bid exchange",
		Conditions:  []string{"condition1", "condition2"},
		Tape:        "tape",
	}
}

func convertToStockPrice(stockResp stream.StockResponse) database.StockPrices {
	return database.StockPrices{
		Prices: database.Prices{
			Symbol:   stockResp.Symbol,
			BidPrice: stockResp.BidPrice,
			BidSize:  stockResp.BidSize,
			AskPrice: stockResp.AskPrice,
			AskSize:  stockResp.AskSize,
			Date:     stockResp.Date,
		},
		AskExchange: stockResp.AskExchange,
		BidExchange: stockResp.BidExchange,
		TradeSize:   stockResp.TradeSize,
		Conditions:  stockResp.Conditions,
		Tape:        stockResp.Tape,
	}
}
