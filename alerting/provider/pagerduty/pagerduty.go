package pagerduty

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TwiN/gatus/v5/alerting/alert"
	"github.com/TwiN/gatus/v5/client"
	"github.com/TwiN/gatus/v5/config/endpoint"
	"github.com/TwiN/logr"
	"gopkg.in/yaml.v3"
)

const restAPIURL = "https://events.pagerduty.com/v2/enqueue"

var (
	ErrIntegrationKeyNotSet   = errors.New("integration-key must have exactly 32 characters")
	ErrDuplicateGroupOverride = errors.New("duplicate group override")
)

type Config struct {
	IntegrationKey string `yaml:"integration-key"`
}

func (cfg *Config) Validate() error {
	if len(cfg.IntegrationKey) != 32 {
		return ErrIntegrationKeyNotSet
	}
	return nil
}

func (cfg *Config) Merge(override *Config) {
	if len(override.IntegrationKey) > 0 {
		cfg.IntegrationKey = override.IntegrationKey
	}
}

type AlertProvider struct {
	DefaultConfig Config       `yaml:",inline"`
	DefaultAlert  *alert.Alert `yaml:"default-alert,omitempty"`
	Overrides     []Override   `yaml:"overrides,omitempty"`
}

type Override struct {
	Group  string `yaml:"group"`
	Config `yaml:",inline"`
}

func (provider *AlertProvider) Validate() error {
	registeredGroups := make(map[string]bool)
	for _, override := range provider.Overrides {
		if isAlreadyRegistered := registeredGroups[override.Group]; isAlreadyRegistered || override.Group == "" {
			return ErrDuplicateGroupOverride
		}
		registeredGroups[override.Group] = true
	}
	return provider.DefaultConfig.Validate()
}

func (provider *AlertProvider) Send(ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result, resolved bool) error {
	cfg, err := provider.GetConfig(ep.Group, alert)
	if err != nil {
		return err
	}

	requestBody := provider.buildRequestBody(cfg, ep, alert, result, resolved)
	buffer := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, restAPIURL, buffer)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := client.GetHTTPClient(nil)
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read response body
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read response body: %w", readErr)
	}

	if response.StatusCode > 399 {
		return fmt.Errorf("call to provider alert returned status code %d: %s", response.StatusCode, string(body))
	}

	if alert.IsSendingOnResolved() {
		if resolved {
			alert.ResolveKey = ""
		} else {
			var payload pagerDutyResponsePayload
			if err = json.Unmarshal(body, &payload); err != nil {
				logr.Errorf("[pagerduty.Send] Failed to decode PagerDuty response: %s", err.Error())
			} else {
				alert.ResolveKey = payload.DedupKey
			}
		}
	}

	return nil
}

type Body struct {
	RoutingKey  string  `json:"routing_key"`
	DedupKey    string  `json:"dedup_key"`
	EventAction string  `json:"event_action"`
	Payload     Payload `json:"payload"`
}

type Payload struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
}

func (provider *AlertProvider) buildRequestBody(cfg *Config, ep *endpoint.Endpoint, alert *alert.Alert, result *endpoint.Result, resolved bool) []byte {
	eventAction := "trigger"
	resolveKey := ""
	msg := ""
	if resolved {
		eventAction = "resolve"
		resolveKey = alert.ResolveKey
		msg = "RESOLVED"
	} else {
		msg = "TRIGGERED"
	}

	message := fmt.Sprintf("%s: %s - %s", msg, ep.DisplayName(), alert.GetDescription())

	body, _ := json.Marshal(Body{
		RoutingKey:  cfg.IntegrationKey,
		DedupKey:    resolveKey,
		EventAction: eventAction,
		Payload: Payload{
			Summary:  message,
			Source:   "Gatus",
			Severity: provider.pagerDutySeverity(result),
		},
	})
	return body
}

func (provider *AlertProvider) pagerDutySeverity(result *endpoint.Result) string {
	if result.Success {
		return "info"
	}
	return "critical"
}

func (provider *AlertProvider) GetConfig(group string, alert *alert.Alert) (*Config, error) {
	cfg := provider.DefaultConfig
	for _, override := range provider.Overrides {
		if group == override.Group {
			cfg.Merge(&override.Config)
			break
		}
	}

	if len(alert.ProviderOverride) != 0 {
		overrideConfig := Config{}
		if err := yaml.Unmarshal(alert.ProviderOverrideAsBytes(), &overrideConfig); err != nil {
			return nil, err
		}
		cfg.Merge(&overrideConfig)
	}

	err := cfg.Validate()
	return &cfg, err
}

func (provider *AlertProvider) GetDefaultAlert() *alert.Alert {
	return provider.DefaultAlert
}

func (provider *AlertProvider) ValidateOverrides(group string, alert *alert.Alert) error {
	_, err := provider.GetConfig(group, alert)
	return err
}

type pagerDutyResponsePayload struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	DedupKey string `json:"dedup_key"`
}
