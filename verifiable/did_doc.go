package verifiable

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// DIDDocument defines current supported did doc model.
type DIDDocument struct {
	Context            interface{}                `json:"@context"`
	ID                 string                     `json:"id"`
	Service            []interface{}              `json:"service,omitempty"`
	VerificationMethod []CommonVerificationMethod `json:"verificationMethod,omitempty"`
	Authentication     []Authentication           `json:"authentication,omitempty"`
	KeyAgreement       []interface{}              `json:"keyAgreement,omitempty"`
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
	ID                   string                 `json:"id"`
	Type                 string                 `json:"type"`
	Controller           string                 `json:"controller"`
	PublicKeyJwk         map[string]interface{} `json:"publicKeyJwk,omitempty"`
	PublicKeyMultibase   string                 `json:"publicKeyMultibase,omitempty"`
	PublicKeyHex         string                 `json:"publicKeyHex,omitempty"`
	PublicKeyBase58      string                 `json:"publicKeyBase58,omitempty"`
	EthereumAddress      string                 `json:"ethereumAddress,omitempty"`
	BlockchainAccountID  string                 `json:"blockchainAccountId,omitempty"`
	StateContractAddress string                 `json:"stateContractAddress,omitempty"`
	IdentityState
}

type Authentication struct {
	CommonVerificationMethod
	did string
}

func (a *Authentication) IsDID() bool {
	return a.did != ""
}

func (a *Authentication) DID() string {
	return a.did
}

func (a *Authentication) UnmarshalJSON(b []byte) error {
	if b == nil {
		return nil
	}
	type Alias Authentication
	switch b[0] {
	case '{':
		tmp := Alias{}
		err := json.Unmarshal(b, &tmp)
		if err != nil {
			return errors.Errorf("invalid json payload for authentication: %v", err)
		}
		*a = (Authentication)(tmp)
	case '"':
		err := json.Unmarshal(b, &a.did)
		if err != nil {
			return fmt.Errorf("faild parse did: %v", err)
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

// StateInfo is information about identity state
type StateInfo struct {
	ID                  string `json:"id"`
	State               string `json:"state"`
	ReplacedByState     string `json:"replacedByState"`
	CreatedAtTimestamp  string `json:"createdAtTimestamp"`
	ReplacedAtTimestamp string `json:"replacedAtTimestamp"`
	CreatedAtBlock      string `json:"createdAtBlock"`
	ReplacedAtBlock     string `json:"replacedAtBlock"`
}

// GistInfo representation state of gist root.
type GistInfo struct {
	Root                string `json:"root"`
	ReplacedByRoot      string `json:"replacedByRoot"`
	CreatedAtTimestamp  string `json:"createdAtTimestamp"`
	ReplacedAtTimestamp string `json:"replacedAtTimestamp"`
	CreatedAtBlock      string `json:"createdAtBlock"`
	ReplacedAtBlock     string `json:"replacedAtBlock"`
}

// IdentityState representation all info about identity.
type IdentityState struct {
	Published *bool      `json:"published,omitempty"`
	Info      *StateInfo `json:"info,omitempty"`
	Global    *GistInfo  `json:"global,omitempty"`
}
