package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config")
}

var _ = Describe("Config", func() {
	Describe("#validate", func() {
		It("returns an error", func() {
			cfg := &Config{}
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
		})
		It("does not return an error", func() {
			test := "test"
			cfg := &Config{Password: test, User: test, Host: test, Name: test}
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return an error for NewFromEnv", func() {
			test := "test"
			os.Setenv("REDSHIFT_PASSWORD", test)
			defer os.Unsetenv("REDSHIFT_PASSWORD")
			os.Setenv("REDSHIFT_USER", test)
			defer os.Unsetenv("REDSHIFT_USER")
			os.Setenv("REDSHIFT_HOST", test)
			defer os.Unsetenv("REDSHIFT_HOST")
			os.Setenv("REDSHIFT_NAME", test)
			defer os.Unsetenv("REDSHIFT_NAME")
			_, err := NewFromEnv()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
