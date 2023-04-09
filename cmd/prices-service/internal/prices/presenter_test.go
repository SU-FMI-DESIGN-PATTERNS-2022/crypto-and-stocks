package prices_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices"
	mock_prices "github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/prices/mocks"
)

var _ = Describe("Presenter", func() {
	const msg = "hi"
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

	FContext("StockHandler", func() {
		When("", func() {
			BeforeEach(func() {
				gomock.InOrder(
					mockUpgarder.EXPECT().Upgrade(nil, nil, nil).Return(mockConn, nil),
					mockBus.EXPECT().Subscribe("stocks", gomock.Any()).DoAndReturn(
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
			It("", func() {
				presenter.StockHandler(nil, nil)
				Expect(response).To(Equal(msg))
			})
		})
	})
})
