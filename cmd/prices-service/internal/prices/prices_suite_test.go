package prices_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prices Suite")
}
