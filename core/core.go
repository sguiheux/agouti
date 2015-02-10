// Package core is a WebDriver API for Go.
package core

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti/api"
)

// NewWebDriver returns an instance of a WebDriver specified by
// a templated URL and command. The URL should be the location of the
// WebDriver Wire Protocol web service brought up by the command. The
// command should be provided as a list of arguments (which are each
// templated). The optional timeout specifies how long to wait for the
// web service to become available. Default timeout is 5 seconds.
//
// Valid template parameters are:
//   {{.Host}} - local address to bind to (usually 127.0.0.1)
//   {{.Port}} - arbitrary free port on the local address
//   {{.Address}} - {{.Host}}:{{.Port}}
//
// Selenium JAR example:
//   command := []string{"java", "-jar", "selenium-server.jar", "-port", "{{.Port}}"}
//   core.NewWebDriver("http://{{.Address}}/wd/hub", command)
func NewWebDriver(url string, command []string, timeout ...time.Duration) WebDriver {
	return &webDriver{api.NewWebDriver(url, command, timeout...)}
}

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
func ChromeDriver() WebDriver {
	return NewWebDriver("http://{{.Address}}", []string{"chromedriver", "--silent", "--port={{.Port}}"})
}

// PhantomJS returns an instance of a PhantomJS WebDriver.
// The return error is deprecated and will always be nil.
func PhantomJS() (WebDriver, error) {
	return NewWebDriver("http://{{.Address}}", []string{"phantomjs", "--webdriver={{.Address}}"}), nil
}

// Selenium returns an instance of a Selenium WebDriver.
// The return error is deprecated and will always be nil.
func Selenium() (WebDriver, error) {
	return NewWebDriver("http://{{.Address}}/wd/hub", []string{"selenium-server", "-port", "{{.Port}}"}), nil
}

// NewPage opens a Page using the provided WebDriver URL.
func NewPage(url string, desired Capabilities) (Page, error) {
	session, err := api.Open(url, api.Capabilities(desired.(capabilities)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %s", err)
	}
	return newPage(session), nil
}

// SauceLabs opens a Sauce Labs session and returns a Page. Does not support Sauce Connect.
func SauceLabs(name, platform, browser, version, username, key string) (Page, error) {
	url := "http://ondemand.saucelabs.com/wd/hub"
	capabilities := api.Capabilities{
		"name":        name,
		"platform":    platform,
		"browserName": browser,
		"version":     version,
		"username":    username,
		"accessKey":   key,
	}

	session, err := api.Open(url, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to Sauce Labs: %s", err)
	}
	return newPage(session), nil
}
