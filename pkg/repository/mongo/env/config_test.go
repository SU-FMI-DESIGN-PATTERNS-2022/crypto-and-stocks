package env_test

import (
	"os"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/repository/mongo/env"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("LoadMongoDBConfig", func() {
		const (
			mongoHost         = "MONGO_HOST"
			mongoPort         = "MONGO_PORT"
			mongoLocalDriver  = "MONGO_LOCAL_DRIVER"
			mongoRemoteDriver = "MONGO_REMOTE_DRIVER"
			mongoUser         = "MONGO_USER"
			mongoDatabase     = "MONGO_DATABASE"
			mongoPassword     = "MONGO_PASSWORD"
			mongoOptions      = "MONGO_OPTIONS"
		)

		BeforeEach(func() {
			os.Setenv(mongoHost, "host")
			os.Setenv(mongoPort, "8080")
			os.Setenv(mongoLocalDriver, "local-driver")
			os.Setenv(mongoRemoteDriver, "remote-driver")
			os.Setenv(mongoUser, "user")
			os.Setenv(mongoDatabase, "database")
			os.Setenv(mongoPassword, "pass")
			os.Setenv(mongoOptions, "options")
		})

		When("no host env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoHost)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("no port env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoPort)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("invalid port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(mongoPort, "invalid")
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("no local driver env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoLocalDriver)
			})

			It("should return the values that are provided and no error", func() {
				mongoConfig, err := env.LoadMongoDBConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(mongoConfig.Host).To(Equal("host"))
				Expect(mongoConfig.Port).To(Equal(8080))
				Expect(mongoConfig.LocalDriver).To(Equal(""))
				Expect(mongoConfig.RemoteDriver).To(Equal("remote-driver"))
				Expect(mongoConfig.User).To(Equal("user"))
				Expect(mongoConfig.Database).To(Equal("database"))
				Expect(mongoConfig.Password).To(Equal("pass"))
				Expect(mongoConfig.Options).To(Equal("options"))
			})
		})

		When("no remote driver env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoRemoteDriver)
			})

			It("should return the values that are provided and no error", func() {
				mongoConfig, err := env.LoadMongoDBConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(mongoConfig.Host).To(Equal("host"))
				Expect(mongoConfig.Port).To(Equal(8080))
				Expect(mongoConfig.LocalDriver).To(Equal("local-driver"))
				Expect(mongoConfig.RemoteDriver).To(Equal(""))
				Expect(mongoConfig.User).To(Equal("user"))
				Expect(mongoConfig.Database).To(Equal("database"))
				Expect(mongoConfig.Password).To(Equal("pass"))
				Expect(mongoConfig.Options).To(Equal("options"))
			})
		})

		When("both local and remote driver envs are not provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoRemoteDriver)
				os.Unsetenv(mongoLocalDriver)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("at least one of the drevers(local or remote) should be set"))
			})
		})

		When("no user env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoUser)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("no database env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoDatabase)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("no password env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoPassword)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("no options env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(mongoOptions)
			})

			It("should return an error", func() {
				_, err := env.LoadMongoDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load mongo database config"))
			})
		})

		When("every env is provided", func() {
			It("should return the values that are provided and no error", func() {
				mongoConfig, err := env.LoadMongoDBConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(mongoConfig.Host).To(Equal("host"))
				Expect(mongoConfig.Port).To(Equal(8080))
				Expect(mongoConfig.LocalDriver).To(Equal("local-driver"))
				Expect(mongoConfig.RemoteDriver).To(Equal("remote-driver"))
				Expect(mongoConfig.User).To(Equal("user"))
				Expect(mongoConfig.Database).To(Equal("database"))
				Expect(mongoConfig.Password).To(Equal("pass"))
				Expect(mongoConfig.Options).To(Equal("options"))
			})
		})
	})
})
