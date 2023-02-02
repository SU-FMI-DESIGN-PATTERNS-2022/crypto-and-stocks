package crypto_prices_repository

import (
	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/internal/database"
)

type CryptoPrices struct {
	database.Prices
	Exchange string `bson:"exchange,omitempty"`
}
