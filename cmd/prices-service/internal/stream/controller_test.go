package stream_test

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream"
	mock_stream "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/stream/mocks"
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/database"
)

var _ = Describe("Controller", func() {
	const (
		cryptoTopic    = "crypto"
		stocksTopic    = "stocks"
		responseErrMsg = "response fails"
		handleErrMsg   = "handling message fails"
	)

	var (
		ctrl             *gomock.Controller
		mockEventBus     *mock_stream.MockEventBus
		mockCryptoStream *mock_stream.MockPriceStream
		mockStockStream  *mock_stream.MockPriceStream
		controller       *stream.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockEventBus = mock_stream.NewMockEventBus(ctrl)
		mockCryptoStream = mock_stream.NewMockPriceStream(ctrl)
		mockStockStream = mock_stream.NewMockPriceStream(ctrl)
		controller = stream.NewController(mockCryptoStream, mockStockStream, mockEventBus)
	})

	Context("StartStreamsToWrite", func() {
		When("starting crypto stream fails", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
				mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
			})

			It("should receive error", func() {
				errCh := controller.StartStreamsToWrite()
				err := errors.New(responseErrMsg)
				Eventually(errCh).Should(Receive(&err))
			})
		})

		When("starting stocks stream fails", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
				mockStockStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
			})

			It("should receive error", func() {
				errCh := controller.StartStreamsToWrite()
				err := errors.New(responseErrMsg)
				Eventually(errCh).Should(Receive(&err))
			})
		})

		When("starting both crypto and stocks streams, both fail", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
				mockStockStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
			})

			It("should receive two errors", func() {
				errCh := controller.StartStreamsToWrite()
				err := errors.New(responseErrMsg)
				Eventually(errCh).Should(Receive(&err))
				Eventually(errCh).Should(Receive(&err))
			})
		})

		When("starting both crypto and stocks streams succeeds", func() {
			var (
				cryptoJSON []byte
				stocksJSON []byte
				failJSON   []byte
			)

			When("publishing crypto price fails", func() {
				BeforeEach(func() {
					failJSON = nil
					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
				})

				It("should receive error", func() {
					errCh := controller.StartStreamsToWrite()
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing crypto price succeeds", func() {
				BeforeEach(func() {
					cryptoJSON = getCryptoPriceAsJSON()

					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(cryptoJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockEventBus.EXPECT().Publish(cryptoTopic, gomock.Any())
				})

				It("should not receive error", func() {
					errCh := controller.StartStreamsToWrite()
					Consistently(errCh).ShouldNot(Receive())
				})
			})

			When("publishing stocks price fails", func() {
				BeforeEach(func() {
					failJSON = nil
					mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
				})

				It("should receive error", func() {
					errCh := controller.StartStreamsToWrite()
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing stocks price succeeds", func() {
				BeforeEach(func() {
					stocksJSON = getStocksPriceAsJSON()

					mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(stocksJSON)
					})
					mockEventBus.EXPECT().Publish(stocksTopic, gomock.Any())
				})

				It("should not receive error", func() {
					errCh := controller.StartStreamsToWrite()
					Consistently(errCh).ShouldNot(Receive())
				})
			})

			When("publishing both crypto and stocks prices fails", func() {
				BeforeEach(func() {
					failJSON = nil

					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
				})

				It("should receive two errors", func() {
					errCh := controller.StartStreamsToWrite()
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing both crypto and stocks prices succeeds", func() {
				BeforeEach(func() {
					cryptoJSON = getCryptoPriceAsJSON()
					stocksJSON = getStocksPriceAsJSON()

					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(cryptoJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(stocksJSON)
					})
					mockEventBus.EXPECT().Publish(cryptoTopic, gomock.Any())
					mockEventBus.EXPECT().Publish(stocksTopic, gomock.Any())
				})

				It("should not receive error", func() {
					errCh := controller.StartStreamsToWrite()
					Consistently(errCh).ShouldNot(Receive())
				})
			})
		})
	})

	// TODO:
	// Context("StopStreams", func() {
	// 	When("stopping crypto stream", func() {
	// 		BeforeEach(func() {
	// 			mockCryptoStream.EXPECT().Stop()
	// 			mockStockStream.EXPECT().Stop()
	// 		})

	// 		It("should not panic", func() {
	// 			Expect(controller.StopStreams).ToNot(Panic())
	// 		})
	// 	})
	// })
})

// TODO: reduce the bellow functions
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

func getCryptoPriceAsJSON() []byte {
	price := []database.CryptoPrices{convertToCryptoPrice(getCryptoResponse())}
	res, err := json.Marshal(price)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return res
}

func getStocksPriceAsJSON() []byte {
	price := []database.StockPrices{convertToStockPrice(getStockResponse())}
	res, err := json.Marshal(price)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return res
}

// TODO: create getFailJSON() function
