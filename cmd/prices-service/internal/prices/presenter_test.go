package prices_test

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	mock_prices "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices/mocks"
)

var _ = Describe("Presenter", func() {
	const (
		cryptoTopic     = "crypto"
		stocksTopic     = "stocks"
		msg             = "hi"
		upgradeErrMsg   = "upgrade failed"
		subscribeErrMsg = "subscribe failed"
		writeJSONErrMsg = "write JSON failed"
	)

	var (
		ctrl         *gomock.Controller
		mockUpgarder *mock_prices.MockUpgrader
		mockBus      *mock_prices.MockEventBus
		mockConn     *mock_prices.MockConnection
		presenter    *prices.Presenter
		response     string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockUpgarder = mock_prices.NewMockUpgrader(ctrl)
		mockBus = mock_prices.NewMockEventBus(ctrl)
		mockConn = mock_prices.NewMockConnection(ctrl)
		presenter = prices.NewPresenter(mockUpgarder, mockBus)
		response = ""
	})

	Context("StockHandler", func() {
		When("upgrading the HTTP server connection to the WebSocket protocol fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(nil, errors.New(upgradeErrMsg)),
				)
			})
			It("should return an error", func() {
				presenter.StockHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("Subscribing for responding fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(stocksTopic, gomock.Any()).Return(errors.New(subscribeErrMsg)))
			})
			It("should return an error", func() {
				presenter.StockHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("Writing JSON fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(stocksTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(resp interface{}))
							Expect(ok).To(BeTrue())
							f(msg)
							return nil
						},
					),
					mockConn.EXPECT().WriteJSON(gomock.Any()).Return(errors.New(writeJSONErrMsg)))
			})
			It("should return an error", func() {
				presenter.StockHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("upgrading the HTTP server connection to the WebSocket protocol succeed and subscribing for responding", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(stocksTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(resp interface{}))
							Expect(ok).To(BeTrue())
							f(msg)
							return nil
						},
					),
					mockConn.EXPECT().WriteJSON(gomock.Any()).DoAndReturn(
						func(json string) error {
							response = json
							return nil
						},
					),
				)
			})
			It("should not return an error", func() {
				presenter.StockHandler(nil, nil)
				Expect(response).To(Equal(msg))
			})
		})
	})

	Context("CryptoHandler", func() {
		When("upgrading the HTTP server connection to the WebSocket protocol fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(nil, errors.New(upgradeErrMsg)),
				)
			})
			It("should return an error", func() {
				presenter.CryptoHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("Subscribing for responding fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).Return(errors.New(subscribeErrMsg)))
			})
			It("should return an error", func() {
				presenter.CryptoHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("Writing JSON fails", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(resp interface{}))
							Expect(ok).To(BeTrue())
							f(msg)
							return nil
						},
					),
					mockConn.EXPECT().WriteJSON(gomock.Any()).Return(errors.New(writeJSONErrMsg)))
			})

			It("should return an error", func() {
				presenter.CryptoHandler(nil, nil)
				Expect(response).To(Equal(""))
			})
		})

		When("upgrading the HTTP server connection to the WebSocket protocol succeed and subscribing for responding", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe(cryptoTopic, gomock.Any()).DoAndReturn(
						func(topic string, fn interface{}) error {
							f, ok := fn.(func(resp interface{}))
							Expect(ok).To(BeTrue())
							f(msg)
							return nil
						},
					),
					mockConn.EXPECT().WriteJSON(gomock.Any()).DoAndReturn(
						func(json string) error {
							response = json
							return nil
						},
					),
				)
			})
			It("should not return an error", func() {
				presenter.CryptoHandler(nil, nil)
				Expect(response).To(Equal(msg))
			})
		})
	})
})
