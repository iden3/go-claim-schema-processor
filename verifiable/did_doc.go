package verifiable

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// DIDDocument defines current supported did doc model.
type DIDDocument struct {
	Context            []string                 `json:"@context"`
	ID                 string                   `json:"id"`
	Service            []interface{}            `json:"service"`
	VerificationMethod CommonVerificationMethod `json:"verificationMethod"`
	Authentication     []Authentication         `json:"authentication"`
	KeyAgreement       []interface{}            `json:"keyAgreement"`
}

// Service describes standard DID document service field.
type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// PushService describes the services of push notifications
type PushService struct {
	Service
	Metadata PushMetadata `json:"metadata"`
}

// PushMetadata describes the structure of the data for push notifications
type PushMetadata struct {
	Devices []EncryptedDeviceMetadata `json:"devices"`
}

// EncryptedDeviceMetadata describes the structure of encrypted device metadata
type EncryptedDeviceMetadata struct {
	Ciphertext string `json:"ciphertext"` // base64 encoded
	Alg        string `json:"alg"`
}

// DeviceMetadata describes the structure of device metadata
type DeviceMetadata struct {
	AppID     string `json:"app_id"`
	PushToken string `json:"push_token"`
}

// CommonVerificationMethod DID doc verification method.
type CommonVerificationMethod struct {
	ID                  string                 `json:"id"`
	Type                string                 `json:"type"`
	Controller          string                 `json:"controller"`
	PublicKeyJwk        map[string]interface{} `json:"publicKeyJwk"`
	PublicKeyMultibase  string                 `json:"publicKeyMultibase,omitempty"`
	PublicKeyHex        string                 `json:"publicKeyHex,omitempty"`
	EthereumAddress     string                 `json:"ethereumAddress,omitempty"`
	BlockchainAccountID string                 `json:"blockchainAccountId,omitempty"`
}

type Authentication struct {
	CommonVerificationMethod
	did string
}

func (a *Authentication) UnmarshalJSON(b []byte) error {
	switch b[0] {
	case '{':
		tmp := Authentication{}
		err := json.Unmarshal(b, &tmp)
		if err != nil {
			return errors.Errorf("invalid json payload for authentication: %w", err)
		}
		*a = (Authentication)(tmp)
	case '"':
		err := json.Unmarshal(b, &a.did)
		if err != nil {
			return fmt.Errorf("faild parse did: %w", err)
		}
	default:
		return errors.New("authentication is invalid")
	}
	return nil
}

func (a *Authentication) MarshalJSON() ([]byte, error) {
	if a.did == "" {
		return json.Marshal(a.CommonVerificationMethod)
	} else {
		return json.Marshal(a)
	}
}
