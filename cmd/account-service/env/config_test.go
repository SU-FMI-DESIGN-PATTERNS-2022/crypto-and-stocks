package env_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/env"
)

var _ = Describe("Config", func() {
	const invalid = "invalid"

	Context("LoadServerConfig", func() {
		const serverPort = "ACCOUNT_SERVER_PORT"

		AfterEach(func() {
			os.Unsetenv(serverPort)
		})

		When("invalid port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(serverPort, invalid)
			})

			It("should return an error", func() {
				_, err := env.LoadServerConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load server config"))
			})
		})

		When("no port env is provided", func() {
			It(`should return the default port "8080" and no error`, func() {
				serverConfig, err := env.LoadServerConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(serverConfig.Port).To(Equal(8080))
			})
		})

		When("port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(serverPort, "8081")
			})

			It("should return the value that is provided and no error", func() {
				serverConfig, err := env.LoadServerConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(serverConfig.Port).To(Equal(8081))
			})
		})
	})

	Context("LoadPostgreDBConfig", func() {
		const (
			postgreHost     = "POSTGRE_HOST"
			postgrePort     = "POSTGRE_PORT"
			postgreName     = "POSTGRE_NAME"
			postgreUser     = "POSTGRE_USER"
			postgrePassword = "POSTGRE_PASSWORD"
		)

		BeforeEach(func() {
			os.Setenv(postgreHost, "host")
			os.Setenv(postgrePort, "8080")
			os.Setenv(postgreName, "name")
			os.Setenv(postgreUser, "user")
			os.Setenv(postgrePassword, "pass")
		})

		When("no host env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(postgreHost)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("no port env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(postgrePort)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("invalid port env is provided", func() {
			BeforeEach(func() {
				os.Setenv(postgrePort, invalid)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("no name env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(postgreName)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("no user env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(postgreUser)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("no password env is provided", func() {
			BeforeEach(func() {
				os.Unsetenv(postgrePassword)
			})

			It("should return an error", func() {
				_, err := env.LoadPostgreDBConfig()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to load postgre database config"))
			})
		})

		When("every env is provided", func() {
			It("should return the values that are provided and no error", func() {
				postgreConfig, err := env.LoadPostgreDBConfig()

				Expect(err).ToNot(HaveOccurred())
				Expect(postgreConfig.Host).To(Equal("host"))
				Expect(postgreConfig.Port).To(Equal(8080))
				Expect(postgreConfig.Name).To(Equal("name"))
				Expect(postgreConfig.User).To(Equal("user"))
				Expect(postgreConfig.Password).To(Equal("pass"))
			})
		})
	})
})
