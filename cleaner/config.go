package cleaner

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the environment configuration
var Config Definition

// Definition contains all environment variables used during the runtime.
type Definition struct {
	GitlabToken string `envconfig:"GITLAB_TOKEN" required:"true"`
	GitlabURL   string `envconfig:"GITLAB_URL" default:"https://gitlab.com/"`
}

// InitializeConfig unmarshals supplied environment variables
// as a exposed Config variable.
func InitializeConfig() {
	if err := envconfig.Process("", &Config); err != nil {
		log.Fatal(err)
	}
}
