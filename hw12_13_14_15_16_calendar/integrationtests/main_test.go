package integrationtests //nolint

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const delay = 5 * time.Second

var (
	configFile string
	config     *Config
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/integrationtests/config.yml", "Path to configuration file")
}

func TestMain(m *testing.M) {
	log.Printf("wait %s for service availability...", delay)
	time.After(delay)

	flag.Parse()
	var err error
	config, err = LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s", err.Error())
		os.Exit(1)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}
//nolint