package verifiable

import (
	"time"
)

// Iden3Credential is struct that represents claim json-ld document
type Iden3Credential struct {
	ID                string                 `json:"id"`
	Context           []string               `json:"@context"`
	Type              []string               `json:"@type"`
	Expiration        time.Time              `json:"expiration,omitempty"`
	Updatable         bool                   `json:"updatable"`
	Version           uint32                 `json:"version"`
	RevNonce          uint64                 `json:"rev_nonce"`
	SubjectPosition   string                 `json:"subject_position,omitempty"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	CredentialStatus  *CredentialStatus      `json:"credentialStatus,omitempty"`
	CredentialSchema  struct {
		ID   string `json:"@id"`
		Type string `json:"type"`
	} `json:"credentialSchema"`
	Proof interface{} `json:"proof,omitempty"`
}

// CredentialStatus contains type and revocation Url
type CredentialStatus struct {
	ID   string               `json:"id"`
	Type CredentialStatusType `json:"type"`
}

//nolint:gosec //reason: no need for security
// SparseMerkleTreeProof is CredentialStatusType
const SparseMerkleTreeProof CredentialStatusType = "SparseMerkleTreeProof"

// CredentialStatusType type for understanding revocation type
type CredentialStatusType string

// JSONSchemaValidator2018 JSON schema
const JSONSchemaValidator2018 = "JsonSchemaValidator2018"

// RevocationStatus status of revocation nonce. Info required to check revocation state of claim in circuits
type RevocationStatus struct {
	Issuer struct {
		State              *string `json:"state"`
		RootOfRoots        *string `json:"root_of_roots,omitempty"`
		ClaimsTreeRoot     *string `json:"claims_tree_root,omitempty"`
		RevocationTreeRoot *string `json:"revocation_tree_root,omitempty"`
	} `json:"issuer"`
	MTP `json:"mtp"`
}
