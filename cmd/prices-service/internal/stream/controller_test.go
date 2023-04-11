package stream_test

import (
	"errors"

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

			It("should return an error", func() {
				err := controller.StartStreamsToWrite()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(responseErrMsg))
				Expect(err.Error()).To(ContainSubstring("crypto"))
			})
		})

		When("starting stocks stream fails", func() {
			BeforeEach(func() {
				mockCryptoStream.EXPECT().Start(gomock.Any()).Return(nil)
				mockStockStream.EXPECT().Start(gomock.Any()).Return(errors.New(responseErrMsg))
			})

			It("should return an error", func() {
				err := controller.StartStreamsToWrite()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(responseErrMsg))
				Expect(err.Error()).To(ContainSubstring("stocks"))
			})
		})
	})
})
