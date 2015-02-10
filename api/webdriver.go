package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/sclevine/agouti/api/internal/bus"
	"github.com/sclevine/agouti/api/internal/service"
)

type Capabilities map[string]interface{}

type WebDriver struct {
	Service interface {
		URL() (string, error)
		Start() error
		Stop() error
	}
	sessions []*Session
}

func NewWebDriver(url string, command []string, timeout ...time.Duration) *WebDriver {
	if len(timeout) == 0 {
		timeout = []time.Duration{5 * time.Second}
	}

	driverService := &service.Service{
		URLTemplate: url,
		CmdTemplate: command,
		Timeout:     timeout[0],
	}

	return &WebDriver{Service: driverService}
}

func (w *WebDriver) Open(desired ...Capabilities) (*Session, error) {
	if len(desired) == 0 {
		desired = append(desired, Capabilities{})
	} else if len(desired) > 1 {
		return nil, errors.New("too many arguments")
	}

	url, err := w.Service.URL()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve URL: %s", err)
	}

	busClient, err := bus.Connect(url, desired[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}

	session := &Session{busClient}
	w.sessions = append(w.sessions, session)
	return session, nil
}

func (w *WebDriver) Start() error {
	if err := w.Service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	return nil
}

func (d *WebDriver) Stop() error {
	for _, session := range d.sessions {
		session.Delete()
	}

	if err := d.Service.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %s", err)
	}

	return nil
}
