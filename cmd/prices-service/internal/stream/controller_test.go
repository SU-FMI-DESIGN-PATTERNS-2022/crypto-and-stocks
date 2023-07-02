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
		var errCh chan error

		BeforeEach(func() {
			errCh = make(chan error, 1)
		})

		JustBeforeEach(func() {
			go func() {
				if err := controller.StartStreamsToWrite(); err != nil {
					errCh <- err
				}
			}()
		})

		When("starting crypto stream fails", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
				mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
			})

			It("should return error", func() {
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
				err := errors.New(responseErrMsg)
				Eventually(errCh).Should(Receive(&err))
			})
		})

		When("starting both crypto and stocks streams, both fail", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
				mockStockStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
			})

			It("should receive error", func() {
				err := errors.New(responseErrMsg)
				Eventually(errCh).Should(Receive(&err))
			})
		})

		When("starting both crypto and stocks streams succeeds", func() {
			var (
				cryptoJSON []byte
				stocksJSON []byte
				failJSON   []byte
			)

			BeforeEach(func() {
				cryptoJSON = convertToCryptoJSON(getCryptoResponse())
				stocksJSON = convertToStockJSON(getStockResponse())
				failJSON = getFailJSON()
			})

			When("publishing crypto price fails", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
				})

				It("should receive error", func() {
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing crypto price succeeds", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(cryptoJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockEventBus.EXPECT().Publish(cryptoTopic, getCryptoResponse())
				})

				It("should not receive error", func() {
					Consistently(errCh).ShouldNot(Receive())
				})
			})

			When("publishing stocks price fails", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
				})

				It("should receive error", func() {
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing stocks price succeeds", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(stocksJSON)
					})
					mockEventBus.EXPECT().Publish(stocksTopic, getStockResponse())
				})

				It("should not receive error", func() {
					Consistently(errCh).ShouldNot(Receive())
				})
			})

			When("publishing both crypto and stocks prices fails", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(failJSON)
					})
				})

				It("should receive error", func() {
					err := errors.New(handleErrMsg)
					Eventually(errCh).Should(Receive(&err))
				})
			})

			When("publishing both crypto and stocks prices succeeds", func() {
				BeforeEach(func() {
					mockCryptoStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(cryptoJSON)
					})
					mockStockStream.EXPECT().Start(gomock.Any()).DoAndReturn(func(msgHandler func([]byte) error) error {
						return msgHandler(stocksJSON)
					})
					mockEventBus.EXPECT().Publish(cryptoTopic, getCryptoResponse())
					mockEventBus.EXPECT().Publish(stocksTopic, getStockResponse())
				})

				It("should not receive error", func() {
					Consistently(errCh).ShouldNot(Receive())
				})
			})
		})
	})

	Context("StopStreams", func() {
		BeforeEach(func() {
			mockCryptoStream.EXPECT().Stop().MinTimes(1)
			mockStockStream.EXPECT().Stop().MinTimes(1)
		})

		When("stopping streams after start", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
				mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
			})

			It("should not publish", func() {
				errCh := make(chan error)
				go func() {
					if err := controller.StartStreamsToWrite(); err != nil {
						errCh <- err
					}
				}()
				time.Sleep(time.Millisecond)
				controller.StopStreams()

				mockEventBus.EXPECT().Publish(cryptoTopic, gomock.Any()).MaxTimes(0)
				mockEventBus.EXPECT().Publish(stocksTopic, gomock.Any()).MaxTimes(0)
				Consistently(errCh).ShouldNot(Receive())
			})
		})

		When("stopping streams before start", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil).MaxTimes(0)
				mockStockStream.EXPECT().Start(gomock.Any()).Return(nil).MaxTimes(0)
			})

			It("should not panic", func() {
				Expect(controller.StopStreams).ToNot(Panic())
			})
		})

		When("stopping streams after already being closed", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
				mockStockStream.EXPECT().Start(gomock.Any()).Return(nil)
			})

			It("should not panic", func() {
				go controller.StartStreamsToWrite()
				time.Sleep(time.Millisecond)

				controller.StopStreams()
				Expect(controller.StopStreams).ToNot(Panic())
			})
		})
	})
})

func getCryptoResponse() stream.CryptoResponse {
	date, err := time.Parse("2006-01-02", "2023-02-16")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return stream.CryptoResponse{
		Response: stream.Response{
			Symbol: "symbol",
			Date:   date,
		},
		Exchange: "excahnge",
	}
}

func getStockResponse() stream.StockResponse {
	date, err := time.Parse("2006-01-02", "2023-02-16")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return stream.StockResponse{
		Response: stream.Response{
			Symbol: "symbol",
			Date:   date,
		},
		AskExchange: "ask exchange",
		BidExchange: "bid exchange",
		Conditions:  []string{"condition1", "condition2"},
		Tape:        "tape",
	}
}

func convertToCryptoJSON(cryptoResp stream.CryptoResponse) []byte {
	responses := []stream.CryptoResponse{cryptoResp}
	res, err := json.Marshal(responses)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return res
}

func convertToStockJSON(stockResp stream.StockResponse) []byte {
	responses := []stream.StockResponse{stockResp}
	res, err := json.Marshal(responses)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return res
}

func getFailJSON() []byte {
	const failObject = `{
		"type": "fail"
	}`
	return []byte(failObject)
}
