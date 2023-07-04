package database_test

import (
	"os"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mongo Envs", func() {
	const (
		host     = "MONGO_HOST"
		port     = "MONGO_PORT"
		user     = "MONGO_USER"
		database = "MONGO_DATABASE"
		password = "MONGO_PASSWORD"
		options  = "MONGO_OPTIONS"
	)
	const invalid = "invalid"

	AfterEach(func() {
		os.Unsetenv(host)
		os.Unsetenv(port)
		os.Unsetenv(user)
		os.Unsetenv(database)
		os.Unsetenv(password)
		os.Unsetenv(options)
	})

	When("invalid mongo host env is provided", func() {
		BeforeEach(func() {
			os.Setenv(host, invalid)
		})

		It("should return an error", func() {
			_, err := env.LoadMongoDBConfig()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to load mongo config"))
		})
	})

	When("invalid mongo db name env is provided", func() {
		BeforeEach(func() {
			os.Setenv(database, invalid)
		})

		It("should return an error", func() {
			_, err := env.LoadMongoDBConfig()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to load mongo config"))
		})
	})

	When("no mongo url env is provided", func() {
		It("should return the default mongo url and no error", func() {
			mongoConfig, err := env.LoadMongoDBConfig()

			Expect(err).ToNot(HaveOccurred())
			Expect(mongoConfig.Host).To(Equal("mongodb://localhost:27017"))
		})
	})

	When("no mongo db name env is provided", func() {
		It("should return the default mongo db name and no error", func() {
			mongoConfig, err := env.LoadMongoDBConfig()

			Expect(err).ToNot(HaveOccurred())
			Expect(mongoConfig.Database).To(Equal("crypto-and-stocks"))
		})
	})

	When("mongo url env is provided", func() {
		BeforeEach(func() {
			os.Setenv(host, "mongodb://localhost:27018")
		})

		It("should return the value that is provided and no error", func() {
			mongoConfig, err := env.LoadMongoDBConfig()

			Expect(err).ToNot(HaveOccurred())
			Expect(mongoConfig.Host).To(Equal("mongodb://localhost:27018"))
		})
	})

	When("mongo db name env is provided", func() {
		BeforeEach(func() {
			os.Setenv(database, "crypto-and-stocks-test")
		})

	})

	When("mongo port env is provided", func() {

	})

})
