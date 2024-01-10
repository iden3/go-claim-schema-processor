package verifiable

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	mt "github.com/iden3/go-merkletree-sql/v2"
	tst "github.com/iden3/go-schema-processor/v2/testing"
	"github.com/iden3/iden3comm/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type credStatusResolverMock struct {
}

type mockNoStateError struct {
}

func (m mockNoStateError) Error() string {
	return "execution reverted: Identity does not exist"
}

func (m mockNoStateError) ErrorCode() int {
	return 3
}

func (m credStatusResolverMock) GetStateInfoByID(id *big.Int) (StateInfo, error) {
	if id.String() == "29305636064099160210536948077705157048478988844998217946273455478812643842" {
		return StateInfo{
			State: "4191494968776819400863455954888115392137551122958477943242938172592557294132",
		}, nil
	}

	if id.String() == "25116094451735045024912155729979573740232593171393457835171656777831420418" {
		return StateInfo{}, mockNoStateError{}
	}
	return StateInfo{}, nil
}

func (m credStatusResolverMock) GetRevocationStatus(id *big.Int, nonce uint64) (RevocationStatus, error) {
	return RevocationStatus{}, nil
}

func (m credStatusResolverMock) GetRevocationStatusByIDAndState(id *big.Int, state *big.Int, nonce uint64) (RevocationStatus, error) {
	return RevocationStatus{}, nil
}

func TestW3CCredential_ValidateBJJSignatureProof(t *testing.T) {
	in := `{
    "id": "urn:uuid:3a8d1822-a00e-11ee-8f57-a27b3ddbdc29",
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://schema.iden3.io/core/jsonld/iden3proofs.jsonld",
        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    ],
    "type": [
        "VerifiableCredential",
        "KYCAgeCredential"
    ],
    "expirationDate": "2361-03-21T21:14:48+02:00",
    "issuanceDate": "2023-12-21T16:35:46.737547+02:00",
    "credentialSubject": {
        "birthday": 19960424,
        "documentType": 2,
        "id": "did:polygonid:polygon:mumbai:2qH2mPVRN7ZDCnEofjeh8Qd2Uo3YsEhTVhKhjB8xs4",
        "type": "KYCAgeCredential"
    },
    "credentialStatus": {
        "id": "https://rhs-staging.polygonid.me/node?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
        "revocationNonce": 74881362,
        "statusIssuer": {
            "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf/claims/revocation/status/74881362",
            "revocationNonce": 74881362,
            "type": "SparseMerkleTreeProof"
        },
        "type": "Iden3ReverseSparseMerkleTreeProof"
    },
    "issuer": "did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf",
    "credentialSchema": {
        "id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
        "type": "JsonSchema2023"
    },
    "proof": [
        {
            "type": "BJJSignature2021",
            "issuerData": {
                "id": "did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf",
                "state": {
                    "claimsTreeRoot": "d946e9cb604bceb0721e4548c291b013647eb56a2cd755b965e6c3b840026517",
                    "value": "f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e"
                },
                "authCoreClaim": "cca3371a6cb1b715004407e325bd993c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d7d1691a4202c0a1e580da2a87118c26a399849c42e52c4d97506a5bf5985923e6ec8ef6caeb482daa0d7516a864ace8fba2854275781583934349b51ba70c190000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "mtp": {
                    "existence": true,
                    "siblings": []
                },
                "credentialStatus": {
                    "id": "https://rhs-staging.polygonid.me/node?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
                    "revocationNonce": 0,
                    "statusIssuer": {
                        "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf/claims/revocation/status/0",
                        "revocationNonce": 0,
                        "type": "SparseMerkleTreeProof"
                    },
                    "type": "Iden3ReverseSparseMerkleTreeProof"
                }
            },
            "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a000000000000000000000000000000021264874acc807e8862077487500a0e9b550a84d667348fc936a4dd0e730b00d4bfb0b3fc0b67c4437ee22848e5de1a7a71748c428358625a5fbac1cebf982000000000000000000000000000000000000000000000000000000000000000005299760400000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "signature": "1783ff1c8207d3047a2ba6baa341dc8a6cb095e5683c6fb619ba4099d3332d2b209dca0a0676e41d4675154ea07662c7d9e14a7ee57259f85f3596493ac71a01"
        }
    ]
}`
	var vc W3CCredential
	err := json.Unmarshal([]byte(in), &vc)
	require.NoError(t, err)

	resolverURL := "http://my-universal-resolver/1.0/identifiers"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://my-universal-resolver/1.0/identifiers/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
		httpmock.NewStringResponder(200, `{"@context":"https://w3id.org/did-resolution/v1","didDocument":{"@context":["https://www.w3.org/ns/did/v1","https://schema.iden3.io/core/jsonld/auth.jsonld"],"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","verificationMethod":[{"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf#stateInfo","type":"Iden3StateInfo2023","controller":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","stateContractAddress":"80001:0x134B1BE34911E39A8397ec6289782989729807a4","published":true,"info":{"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","state":"34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409","replacedByState":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1703174663","replacedAtTimestamp":"0","createdAtBlock":"43840767","replacedAtBlock":"0"},"global":{"root":"92c4610a24247a4013ce6de4903452d164134a232a94fd1fe37178bce4937006","replacedByRoot":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1704439557","replacedAtTimestamp":"0","createdAtBlock":"44415346","replacedAtBlock":"0"}}]},"didResolutionMetadata":{"contentType":"application/did+ld+json","retrieved":"2024-01-05T08:05:13.413770024Z","pattern":"^(did:polygonid:.+)$","driverUrl":"http://driver-did-polygonid:8080/1.0/identifiers/","duration":429,"did":{"didString":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","methodSpecificId":"polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","method":"polygonid"}},"didDocumentMetadata":{}}`))

	// httpmock.RegisterMatcherResponder("POST", "http://my-rpc/v2/1111",
	// 	httpmock.BodyContainsString(`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"data":"0xb4bdea550010961e749448c0c935c85ae263d271b383a2f1fa92ebb74ac9b652efab1202","from":"0x0000000000000000000000000000000000000000","to":"0x134b1be34911e39a8397ec6289782989729807a4"},"latest"]}`),
	// 	httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x0010961e749448c0c935c85ae263d271b383a2f1fa92ebb74ac9b652efab120209444d55300819a07594d731c0bf7d3713f5e9324e0435f926c3ef1d8e4a823400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000065846207000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029cf4ff0000000000000000000000000000000000000000000000000000000000000000"}`))

	httpmock.RegisterResponder("GET", "https://rhs-staging.polygonid.me/node/34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409",
		httpmock.NewStringResponder(200, `{"node":{"hash":"34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409","children":["4436ea12d352ddb84d2ac7a27bbf7c9f1bfc7d3ff69f3e6cf4348f424317fd0b","0000000000000000000000000000000000000000000000000000000000000000","37eabc712cdaa64793561b16b8143f56f149ad1b0c35297a1b125c765d1c071e"]},"status":"OK"}`))

	config := []W3CProofVerificationOpt{WithResolver(credStatusResolverMock{}), WithPackageManager(*iden3comm.NewPackageManager())}
	isValid, err := vc.VerifyProof(context.Background(), BJJSignatureProofType, resolverURL, config...)
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestW3CCredential_ValidateBJJSignatureProofGenesis(t *testing.T) {
	in := `{
    "id": "urn:uuid:b7a1e232-a0d3-11ee-bc8a-a27b3ddbdc29",
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://schema.iden3.io/core/jsonld/iden3proofs.jsonld",
        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    ],
    "type": [
        "VerifiableCredential",
        "KYCAgeCredential"
    ],
    "expirationDate": "2361-03-21T21:14:48+02:00",
    "issuanceDate": "2023-12-22T16:09:27.444712+02:00",
    "credentialSubject": {
        "birthday": 19960424,
        "documentType": 2,
        "id": "did:polygonid:polygon:mumbai:2qJm6vBXtHWMqm9A9f5zihRNVGptHAHcK8oVxGUTg8",
        "type": "KYCAgeCredential"
    },
    "credentialStatus": {
        "id": "https://rhs-staging.polygonid.me/node?state=da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817",
        "revocationNonce": 1102174849,
        "statusIssuer": {
            "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks/claims/revocation/status/1102174849",
            "revocationNonce": 1102174849,
            "type": "SparseMerkleTreeProof"
        },
        "type": "Iden3ReverseSparseMerkleTreeProof"
    },
    "issuer": "did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks",
    "credentialSchema": {
        "id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
        "type": "JsonSchema2023"
    },
    "proof": [
        {
            "type": "BJJSignature2021",
            "issuerData": {
                "id": "did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks",
                "state": {
                    "claimsTreeRoot": "aec50251fdc67959254c74ab4f2e746a7cd1c6f494c8ac028d655dfbccea430e",
                    "value": "da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817"
                },
                "authCoreClaim": "cca3371a6cb1b715004407e325bd993c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c08ac5cc7c5aa3e8190e188cf8d1737c92d16188541b582ef676c55b3a842c06c4985e9d4771ee6d033c2021a3d177f7dfa51859d99a9a476c2a910e887dc8240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "mtp": {
                    "existence": true,
                    "siblings": []
                },
                "credentialStatus": {
                    "id": "https://rhs-staging.polygonid.me/node?state=da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817",
                    "revocationNonce": 0,
                    "statusIssuer": {
                        "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks/claims/revocation/status/0",
                        "revocationNonce": 0,
                        "type": "SparseMerkleTreeProof"
                    },
                    "type": "Iden3ReverseSparseMerkleTreeProof"
                }
            },
            "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a00000000000000000000000000000002128aa2ae20d4f8f7b9d673e06498fa410f3c5a790194f3b9284a2018f30d0037d1e542f1b72c9d5ca4b46d93710fbfa23a7c9c36eb3ca0eb0f9548ad9c140c000000000000000000000000000000000000000000000000000000000000000081dab14100000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "signature": "2a2e4d79f3aa440154643252d1b9074f9651fffcd653fb2fcadc07f55cd1f9a20a812dd7df8ba8775653984cfb7120f999751f9c25473fd634c7f2d88419c102"
        }
    ]
}`
	var vc W3CCredential
	err := json.Unmarshal([]byte(in), &vc)
	require.NoError(t, err)

	resolverURL := "http://my-universal-resolver/1.0/identifiers"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://my-universal-resolver/1.0/identifiers/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks?state=da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817",
		httpmock.NewStringResponder(200, `{"@context":"https://w3id.org/did-resolution/v1","didDocument":{"@context":["https://www.w3.org/ns/did/v1","https://schema.iden3.io/core/jsonld/auth.jsonld"],"id":"did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks","verificationMethod":[{"id":"did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks#stateInfo","type":"Iden3StateInfo2023","controller":"did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks","stateContractAddress":"80001:0x134B1BE34911E39A8397ec6289782989729807a4","published":false,"global":{"root":"92c4610a24247a4013ce6de4903452d164134a232a94fd1fe37178bce4937006","replacedByRoot":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1704439557","replacedAtTimestamp":"0","createdAtBlock":"44415346","replacedAtBlock":"0"}}]},"didResolutionMetadata":{"contentType":"application/did+ld+json","retrieved":"2024-01-05T08:02:25.986085836Z","pattern":"^(did:polygonid:.+)$","driverUrl":"http://driver-did-polygonid:8080/1.0/identifiers/","duration":434,"did":{"didString":"did:polygonid:polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks","methodSpecificId":"polygon:mumbai:2qLx3hTJBV8REpNDK2RiG7eNBVzXMoZdPfi2uhF7Ks","method":"polygonid"}},"didDocumentMetadata":{}}`))

	httpmock.RegisterMatcherResponder("POST", "http://my-rpc/v2/1111",
		httpmock.BodyContainsString(`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"data":"0xb4bdea55000e3717b8601710678ac6bc754dc7876d513fffe8e2bf4dbb52cc0cd9ba1202","from":"0x0000000000000000000000000000000000000000","to":"0x134b1be34911e39a8397ec6289782989729807a4"},"latest"]}`),
		httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"error":{"code":3,"message":"execution reverted: Identity does not exist","data":"0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000174964656e7469747920646f6573206e6f74206578697374000000000000000000"}}`))

	httpmock.RegisterResponder("GET", "https://rhs-staging.polygonid.me/node/da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817",
		httpmock.NewStringResponder(200, `{"node":{"hash":"da6184809dbad90ccc52bb4dbfe2e8ff3f516d87c74d75bcc68a67101760b817","children":["aec50251fdc67959254c74ab4f2e746a7cd1c6f494c8ac028d655dfbccea430e","0000000000000000000000000000000000000000000000000000000000000000","0000000000000000000000000000000000000000000000000000000000000000"]},"status":"OK"}`))

	config := []W3CProofVerificationOpt{WithResolver(credStatusResolverMock{}), WithPackageManager(*iden3comm.NewPackageManager())}
	isValid, err := vc.VerifyProof(context.Background(), BJJSignatureProofType, resolverURL, config...)
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestW3CCredential_ValidateIden3SparseMerkleTreeProof(t *testing.T) {
	in := `{
    "id": "urn:uuid:3a8d1822-a00e-11ee-8f57-a27b3ddbdc29",
    "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://schema.iden3.io/core/jsonld/iden3proofs.jsonld",
        "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    ],
    "type": [
        "VerifiableCredential",
        "KYCAgeCredential"
    ],
    "expirationDate": "2361-03-21T21:14:48+02:00",
    "issuanceDate": "2023-12-21T16:35:46.737547+02:00",
    "credentialSubject": {
        "birthday": 19960424,
        "documentType": 2,
        "id": "did:polygonid:polygon:mumbai:2qH2mPVRN7ZDCnEofjeh8Qd2Uo3YsEhTVhKhjB8xs4",
        "type": "KYCAgeCredential"
    },
    "credentialStatus": {
        "id": "https://rhs-staging.polygonid.me/node?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
        "revocationNonce": 74881362,
        "statusIssuer": {
            "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf/claims/revocation/status/74881362",
            "revocationNonce": 74881362,
            "type": "SparseMerkleTreeProof"
        },
        "type": "Iden3ReverseSparseMerkleTreeProof"
    },
    "issuer": "did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf",
    "credentialSchema": {
        "id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
        "type": "JsonSchema2023"
    },
    "proof": [
        {
            "type": "BJJSignature2021",
            "issuerData": {
                "id": "did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf",
                "state": {
                    "claimsTreeRoot": "d946e9cb604bceb0721e4548c291b013647eb56a2cd755b965e6c3b840026517",
                    "value": "f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e"
                },
                "authCoreClaim": "cca3371a6cb1b715004407e325bd993c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d7d1691a4202c0a1e580da2a87118c26a399849c42e52c4d97506a5bf5985923e6ec8ef6caeb482daa0d7516a864ace8fba2854275781583934349b51ba70c190000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "mtp": {
                    "existence": true,
                    "siblings": []
                },
                "credentialStatus": {
                    "id": "https://rhs-staging.polygonid.me/node?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
                    "revocationNonce": 0,
                    "statusIssuer": {
                        "id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf/claims/revocation/status/0",
                        "revocationNonce": 0,
                        "type": "SparseMerkleTreeProof"
                    },
                    "type": "Iden3ReverseSparseMerkleTreeProof"
                }
            },
            "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a000000000000000000000000000000021264874acc807e8862077487500a0e9b550a84d667348fc936a4dd0e730b00d4bfb0b3fc0b67c4437ee22848e5de1a7a71748c428358625a5fbac1cebf982000000000000000000000000000000000000000000000000000000000000000005299760400000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "signature": "1783ff1c8207d3047a2ba6baa341dc8a6cb095e5683c6fb619ba4099d3332d2b209dca0a0676e41d4675154ea07662c7d9e14a7ee57259f85f3596493ac71a01"
        },
        {
            "type": "Iden3SparseMerkleTreeProof",
            "issuerData": {
                "id": "did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf",
                "state": {
                    "txId": "0x7ab71a8c5e91064e21beb586012f8b89932c255e243c496dec895a501a42e243",
                    "blockTimestamp": 1703174663,
                    "blockNumber": 43840767,
                    "rootOfRoots": "37eabc712cdaa64793561b16b8143f56f149ad1b0c35297a1b125c765d1c071e",
                    "claimsTreeRoot": "4436ea12d352ddb84d2ac7a27bbf7c9f1bfc7d3ff69f3e6cf4348f424317fd0b",
                    "revocationTreeRoot": "0000000000000000000000000000000000000000000000000000000000000000",
                    "value": "34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409"
                }
            },
            "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a000000000000000000000000000000021264874acc807e8862077487500a0e9b550a84d667348fc936a4dd0e730b00d4bfb0b3fc0b67c4437ee22848e5de1a7a71748c428358625a5fbac1cebf982000000000000000000000000000000000000000000000000000000000000000005299760400000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "mtp": {
                "existence": true,
                "siblings": [
                    "0",
                    "10581662619345074277108685138429405012286849178024033034405862946888154171097"
                ]
            }
        }
    ]
}`
	var vc W3CCredential
	err := json.Unmarshal([]byte(in), &vc)
	require.NoError(t, err)

	resolverURL := "http://my-universal-resolver/1.0/identifiers"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://my-universal-resolver/1.0/identifiers/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf?state=34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409",
		httpmock.NewStringResponder(200, `{"@context":"https://w3id.org/did-resolution/v1","didDocument":{"@context":["https://www.w3.org/ns/did/v1","https://schema.iden3.io/core/jsonld/auth.jsonld"],"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","verificationMethod":[{"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf#stateInfo","type":"Iden3StateInfo2023","controller":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","stateContractAddress":"80001:0x134B1BE34911E39A8397ec6289782989729807a4","published":true,"info":{"id":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","state":"34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409","replacedByState":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1703174663","replacedAtTimestamp":"0","createdAtBlock":"43840767","replacedAtBlock":"0"},"global":{"root":"92c4610a24247a4013ce6de4903452d164134a232a94fd1fe37178bce4937006","replacedByRoot":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1704439557","replacedAtTimestamp":"0","createdAtBlock":"44415346","replacedAtBlock":"0"}}]},"didResolutionMetadata":{"contentType":"application/did+ld+json","retrieved":"2024-01-05T07:53:42.67771172Z","pattern":"^(did:polygonid:.+)$","driverUrl":"http://driver-did-polygonid:8080/1.0/identifiers/","duration":442,"did":{"didString":"did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","methodSpecificId":"polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf","method":"polygonid"}},"didDocumentMetadata":{}}`))

	config := []W3CProofVerificationOpt{WithResolver(credStatusResolverMock{}), WithPackageManager(*iden3comm.NewPackageManager())}
	isValid, err := vc.VerifyProof(context.Background(), Iden3SparseMerkleTreeProofType, resolverURL, config...)
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestW3CCredential_ValidateBJJSignatureProofAgentStatus(t *testing.T) {
	in := `{
        "id": "urn:uuid:79d93584-ae2c-11ee-8050-a27b3ddbdc28",
        "@context": [
            "https://www.w3.org/2018/credentials/v1",
            "https://schema.iden3.io/core/jsonld/iden3proofs.jsonld",
            "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
        ],
        "type": [
            "VerifiableCredential",
            "KYCAgeCredential"
        ],
        "expirationDate": "2361-03-21T21:14:48+02:00",
        "issuanceDate": "2024-01-08T15:47:34.113565+02:00",
        "credentialSubject": {
            "birthday": 19960424,
            "documentType": 2,
            "id": "did:polygonid:polygon:mumbai:2qFDziX3k3h7To2jDJbQiXFtcozbgSNNvQpb6TgtPE",
            "type": "KYCAgeCredential"
        },
        "credentialStatus": {
            "id": "http://localhost:8001/api/v1/agent",
            "revocationNonce": 3262660310,
            "type": "Iden3commRevocationStatusV1.0"
        },
        "issuer": "did:polygonid:polygon:mumbai:2qJp131YoXVu8iLNGfL3TkQAWEr3pqimh2iaPgH3BJ",
        "credentialSchema": {
            "id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
            "type": "JsonSchema2023"
        },
        "proof": [
            {
                "type": "BJJSignature2021",
                "issuerData": {
                    "id": "did:polygonid:polygon:mumbai:2qJp131YoXVu8iLNGfL3TkQAWEr3pqimh2iaPgH3BJ",
                    "state": {
                        "claimsTreeRoot": "b35562873d9870f20e3d44dd94502f4156785a4b09d7906914758a7e0ed26829",
                        "value": "2de39210318bbc7fc79e24150c2790089c8385d7acffc0f0ebf1641b95087e0f"
                    },
                    "authCoreClaim": "cca3371a6cb1b715004407e325bd993c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000167c1d2857ca6579d6e995198876cdfd4edb4fe2eeedeadbabaaed3008225205e7b8ab88a60b9ef0999be82625e0831872d8aca16b2932852c3731e9df69970a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                    "mtp": {
                        "existence": true,
                        "siblings": []
                    },
                    "credentialStatus": {
                        "id": "http://localhost:8001/api/v1/agent",
                        "revocationNonce": 0,
                        "type": "Iden3commRevocationStatusV1.0"
                    }
                },
                "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a00000000000000000000000000000002123cbcd9d0f3a493561510c72b47afcb02e2f09b3855291c6b77d224260d0014f503c3ab03eebe757d5b50b570186a69d90c49904155f5fc71e0e7f5b8aa120000000000000000000000000000000000000000000000000000000000000000d63e78c200000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "signature": "56ab45ad828c4860d02e111b2732c969005046ee26dbc7d1e5bd6a6c6604ed81c3f55ffb9349f4d407f59e2e210f6d256a328d30edae2c7c95dd057240ee8902"
            },
            {
                "type": "Iden3SparseMerkleTreeProof",
                "issuerData": {
                    "id": "did:polygonid:polygon:mumbai:2qJp131YoXVu8iLNGfL3TkQAWEr3pqimh2iaPgH3BJ",
                    "state": {
                        "txId": "0x02f1af6a616715ccb7511176ca53d39a28c55201effca0b43a343ee6e9dc8c97",
                        "blockTimestamp": 1704721690,
                        "blockNumber": 44542683,
                        "rootOfRoots": "eaa48e4a7d3fe2fabbd939c7df1048c3f647a9a7c9dfadaae836ec78ba673229",
                        "claimsTreeRoot": "d9597e2fef206c9821f2425e513a68c8c793bc93c9216fb883fedaaf72abf51c",
                        "revocationTreeRoot": "0000000000000000000000000000000000000000000000000000000000000000",
                        "value": "96161f3fbbdd68c72bc430dae474e27b157586b33b9fbf4a3f07d75ce275570f"
                    }
                },
                "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a00000000000000000000000000000002123cbcd9d0f3a493561510c72b47afcb02e2f09b3855291c6b77d224260d0014f503c3ab03eebe757d5b50b570186a69d90c49904155f5fc71e0e7f5b8aa120000000000000000000000000000000000000000000000000000000000000000d63e78c200000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "mtp": {
                    "existence": true,
                    "siblings": [
                        "18730028644149260049434737497088408840959357817865392043806470281178241979827"
                    ]
                }
            }
        ]
    }`
	var vc W3CCredential
	err := json.Unmarshal([]byte(in), &vc)
	require.NoError(t, err)

	resolverURL := "http://my-universal-resolver/1.0/identifiers"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://my-universal-resolver/1.0/identifiers/did%3Apolygonid%3Apolygon%3Amumbai%3A2qJp131YoXVu8iLNGfL3TkQAWEr3pqimh2iaPgH3BJ?state=2de39210318bbc7fc79e24150c2790089c8385d7acffc0f0ebf1641b95087e0f",
		httpmock.NewStringResponder(200, `{"didDocument":{"@context":["https://www.w3.org/ns/did/v1","https://schema.iden3.io/core/jsonld/auth.jsonld"],"id":"did:polygonid:polygon:mumbai:2qEChbFATnamWnToMgNycnVi4W9Xw5772qX61qwki6","verificationMethod":[{"id":"did:polygonid:polygon:mumbai:2qEChbFATnamWnToMgNycnVi4W9Xw5772qX61qwki6#stateInfo","type":"Iden3StateInfo2023","controller":"did:polygonid:polygon:mumbai:2qEChbFATnamWnToMgNycnVi4W9Xw5772qX61qwki6","stateContractAddress":"80001:0x134B1BE34911E39A8397ec6289782989729807a4","published":false,"global":{"root":"ff3e987dc4c279af0e77ac2b1983ed8cf627bfeebbc6d5d56be2526cc7286621","replacedByRoot":"0000000000000000000000000000000000000000000000000000000000000000","createdAtTimestamp":"1704719148","replacedAtTimestamp":"0","createdAtBlock":"44541667","replacedAtBlock":"0"}}]}}`))

	httpmock.RegisterResponder("POST", "http://localhost:8001/api/v1/agent",
		httpmock.NewStringResponder(200, `{"body":{"issuer":{"claimsTreeRoot":"d9597e2fef206c9821f2425e513a68c8c793bc93c9216fb883fedaaf72abf51c","revocationTreeRoot":"0000000000000000000000000000000000000000000000000000000000000000","rootOfRoots":"eaa48e4a7d3fe2fabbd939c7df1048c3f647a9a7c9dfadaae836ec78ba673229","state":"96161f3fbbdd68c72bc430dae474e27b157586b33b9fbf4a3f07d75ce275570f"},"mtp":{"existence":false,"siblings":[]}},"from":"did:polygonid:polygon:mumbai:2qJp131YoXVu8iLNGfL3TkQAWEr3pqimh2iaPgH3BJ","id":"9ece0dad-9267-4a52-b611-f0615b0143fb","thid":"8bdc87dc-1755-41d5-b483-26562836068e","to":"did:polygonid:polygon:mumbai:2qFDziX3k3h7To2jDJbQiXFtcozbgSNNvQpb6TgtPE","typ":"application/iden3comm-plain-json","type":"https://iden3-communication.io/revocation/1.0/status"}`))

	pckManager := iden3comm.NewPackageManager()
	err = pckManager.RegisterPackers(&PlainMessagePacker{})
	require.NoError(t, err)
	config := []W3CProofVerificationOpt{WithResolver(credStatusResolverMock{}), WithPackageManager(*pckManager)}
	isValid, err := vc.VerifyProof(context.Background(), BJJSignatureProofType, resolverURL, config...)
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestW3CCredential_JSONUnmarshal(t *testing.T) {
	in := `{
    "id": "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/claim/52cec4e3-7d1d-11ed-ade2-0242ac180007",
    "@context": [
      "https://www.w3.org/2018/credentials/v1",
      "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/iden3credential-v2.json-ld",
      "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld"
    ],
    "type": [
      "VerifiableCredential",
      "KYCAgeCredential"
    ],
    "expirationDate": "2361-03-21T19:14:48Z",
    "issuanceDate": "2022-12-16T08:40:41.515927692Z",
    "credentialSubject": {
      "birthday": 19960424,
      "documentType": 2,
      "id": "did:iden3:polygon:mumbai:x3YTKSK1fWBbQAmMhArxvFBcG8tL7m2ZMFh5LSyjH",
      "type": "KYCAgeCredential"
    },
    "credentialStatus": {
      "id": "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/identities/did%3Aiden3%3Apolygon%3Amumbai%3AwvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG/claims/revocation/status/1529060834",
      "revocationNonce": 1529060834,
      "type": "SparseMerkleTreeProof"
    },
    "issuer": "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
    "credentialSchema": {
      "id": "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
      "type": "JsonSchemaValidator2018"
    },
    "proof": [
      {
        "type": "BJJSignature2021",
        "issuerData": {
          "id": "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
          "state": {
            "claimsTreeRoot": "93121670a2a82d42adb3eae22d609c2495ee675d36feaaef75bd030b3e98f621",
            "value": "fab7bdf8551406b0bc2df0dabf811449d74628f02e98b2e4ea02f01b996a4e05"
          },
          "authCoreClaim": "013fd3f623559d850fb5b02ff012d0e20000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001409ffecd5566451e39ee1cf7ff2e5b369ef6a708e51f80d7ba282e5c1f6d80eb88eb6df418a768c1f9dc4cc1c6109564f6d5a36d74a7085d9f90c66ae03641c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
          "mtp": {
            "existence": true,
            "siblings": []
          },
          "credentialStatus": {
            "id": "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/identities/did%3Aiden3%3Apolygon%3Amumbai%3AwvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG/claims/revocation/status/0",
            "revocationNonce": 0,
            "type": "SparseMerkleTreeProof"
          }
        },
        "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a0000000000000000000000000000000112b4f1183b6a0708a8addd31c093004ac2e40ab1b291ad6d208244032b0c006947c37450a6a4c50a586e8a253dc8385d8d1ee77b37f464fe5052dc2f0dd8020000000000000000000000000000000000000000000000000000000000000000e29d235b00000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "signature": "b36ed82e13d2868d6b5c5dff0f461e309e1af4cf3fdc9822fd0f86b76c820f19cd728d06ff22c259d4aeef3406c3d44577014fbd0e8fb14330022de77bda8302"
      },
      {
        "type": "Iden3SparseMerkleProof",
        "issuerData": {
          "id": "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
          "state": {
            "txId": "0x705881f799496f399321f7b3b0f9aab80e358e5fdacb877ef18f10afc8be156e",
            "blockTimestamp": 1671180108,
            "blockNumber": 29756768,
            "rootOfRoots": "db07217f60526821e8c079802ebfbfb9cd07e42d4220ff72f264d9bddbe87d2f",
            "claimsTreeRoot": "447b1dfd065752d099c4c8eeb181dfe1363c64491eb413f01d6e60daf6bc792e",
            "revocationTreeRoot": "0000000000000000000000000000000000000000000000000000000000000000",
            "value": "0bc71a0bdbf1a3e8513069b170c6b62601288fcf231f874b52e4e546dddcbb2d"
          }
        },
        "coreClaim": "c9b2370371b7fa8b3dab2a5ba81b68382a0000000000000000000000000000000112b4f1183b6a0708a8addd31c093004ac2e40ab1b291ad6d208244032b0c006947c37450a6a4c50a586e8a253dc8385d8d1ee77b37f464fe5052dc2f0dd8020000000000000000000000000000000000000000000000000000000000000000e29d235b00000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "mtp": {
          "existence": true,
          "siblings": [
            "0",
            "13291429422163653257975736723599735973011351095941906941706092370486076739639",
            "13426716414767621234869633661856285788095461522423569801792562280466318278688"
          ]
        }
      }
    ]
  }`
	var vc W3CCredential
	err := json.Unmarshal([]byte(in), &vc)
	require.NoError(t, err)

	want := W3CCredential{
		ID: "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/claim/52cec4e3-7d1d-11ed-ade2-0242ac180007",
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/iden3credential-v2.json-ld",
			"https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v3.json-ld",
		},
		Type: []string{
			"VerifiableCredential",
			"KYCAgeCredential",
		},
		Expiration: &[]time.Time{
			time.Date(2361, 3, 21, 19, 14, 48, 0, time.UTC)}[0],
		IssuanceDate: &[]time.Time{
			time.Date(2022, 12, 16, 8, 40, 41, 515927692, time.UTC)}[0],
		CredentialSubject: map[string]any{
			"birthday":     float64(19960424),
			"documentType": float64(2),
			"id":           "did:iden3:polygon:mumbai:x3YTKSK1fWBbQAmMhArxvFBcG8tL7m2ZMFh5LSyjH",
			"type":         "KYCAgeCredential",
		},
		CredentialStatus: map[string]any{
			"id":              "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/identities/did%3Aiden3%3Apolygon%3Amumbai%3AwvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG/claims/revocation/status/1529060834",
			"revocationNonce": float64(1529060834),
			"type":            "SparseMerkleTreeProof",
		},
		Issuer: "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
		CredentialSchema: CredentialSchema{
			ID:   "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json/KYCAgeCredential-v3.json",
			Type: "JsonSchemaValidator2018",
		},
		Proof: CredentialProofs{
			&BJJSignatureProof2021{
				Type: BJJSignatureProofType,
				IssuerData: IssuerData{
					ID: "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
					State: State{
						ClaimsTreeRoot: &[]string{"93121670a2a82d42adb3eae22d609c2495ee675d36feaaef75bd030b3e98f621"}[0],
						Value:          &[]string{"fab7bdf8551406b0bc2df0dabf811449d74628f02e98b2e4ea02f01b996a4e05"}[0],
					},
					AuthCoreClaim: "013fd3f623559d850fb5b02ff012d0e20000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001409ffecd5566451e39ee1cf7ff2e5b369ef6a708e51f80d7ba282e5c1f6d80eb88eb6df418a768c1f9dc4cc1c6109564f6d5a36d74a7085d9f90c66ae03641c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
					MTP:           mustProof(t, true, []*mt.Hash{}),
					CredentialStatus: map[string]interface{}{
						"id":              "http://ec2-34-247-165-109.eu-west-1.compute.amazonaws.com:8888/api/v1/identities/did%3Aiden3%3Apolygon%3Amumbai%3AwvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG/claims/revocation/status/0",
						"revocationNonce": float64(0),
						"type":            "SparseMerkleTreeProof",
					},
				},
				CoreClaim: "c9b2370371b7fa8b3dab2a5ba81b68382a0000000000000000000000000000000112b4f1183b6a0708a8addd31c093004ac2e40ab1b291ad6d208244032b0c006947c37450a6a4c50a586e8a253dc8385d8d1ee77b37f464fe5052dc2f0dd8020000000000000000000000000000000000000000000000000000000000000000e29d235b00000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				Signature: "b36ed82e13d2868d6b5c5dff0f461e309e1af4cf3fdc9822fd0f86b76c820f19cd728d06ff22c259d4aeef3406c3d44577014fbd0e8fb14330022de77bda8302",
			},
			&Iden3SparseMerkleProof{
				Type: Iden3SparseMerkleProofType,
				IssuerData: IssuerData{
					ID: "did:iden3:polygon:mumbai:wvEkzpApgwGHrSTxEFG6V6HrTCa5R2rwQ3XWAkrnG",
					State: State{
						TxID:               &[]string{"0x705881f799496f399321f7b3b0f9aab80e358e5fdacb877ef18f10afc8be156e"}[0],
						BlockTimestamp:     &[]int{1671180108}[0],
						BlockNumber:        &[]int{29756768}[0],
						RootOfRoots:        &[]string{"db07217f60526821e8c079802ebfbfb9cd07e42d4220ff72f264d9bddbe87d2f"}[0],
						ClaimsTreeRoot:     &[]string{"447b1dfd065752d099c4c8eeb181dfe1363c64491eb413f01d6e60daf6bc792e"}[0],
						RevocationTreeRoot: &[]string{"0000000000000000000000000000000000000000000000000000000000000000"}[0],
						Value:              &[]string{"0bc71a0bdbf1a3e8513069b170c6b62601288fcf231f874b52e4e546dddcbb2d"}[0],
						Status:             "",
					},
					AuthCoreClaim:    "",
					MTP:              nil,
					CredentialStatus: nil,
				},
				CoreClaim: "c9b2370371b7fa8b3dab2a5ba81b68382a0000000000000000000000000000000112b4f1183b6a0708a8addd31c093004ac2e40ab1b291ad6d208244032b0c006947c37450a6a4c50a586e8a253dc8385d8d1ee77b37f464fe5052dc2f0dd8020000000000000000000000000000000000000000000000000000000000000000e29d235b00000000281cdcdf0200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				MTP: mustProof(t, true, []*mt.Hash{
					mustHash(t, "0"),
					mustHash(t,
						"13291429422163653257975736723599735973011351095941906941706092370486076739639"),
					mustHash(t,
						"13426716414767621234869633661856285788095461522423569801792562280466318278688"),
				}),
			},
		},
	}
	require.Equal(t, want, vc)
}

func TestW3CCredential_MerklizationWithEmptyID(t *testing.T) {
	defer tst.MockHTTPClient(t, map[string]string{
		"https://www.w3.org/2018/credentials/v1":              "../merklize/testdata/httpresp/credentials-v1.jsonld",
		"https://example.com/schema-delivery-address.json-ld": "../json/testdata/schema-delivery-address.json-ld",
	})()

	vcData, err := os.ReadFile("../json/testdata/non-merklized-1.json-ld")
	require.NoError(t, err)
	var vc W3CCredential
	err = json.Unmarshal(vcData, &vc)
	require.NoError(t, err)

	want := W3CCredential{
		ID: "",
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://example.com/schema-delivery-address.json-ld",
		},
		Type: []string{
			"VerifiableCredential",
			"DeliverAddressMultiTestForked",
		},
		CredentialSubject: map[string]any{
			"type":             "DeliverAddressMultiTestForked",
			"price":            "123.52",
			"isPostalProvider": false,
			"postalProviderInformation": map[string]any{
				"insured": true,
				"weight":  "1.3",
			},
		},
		CredentialStatus: nil,
		Issuer:           "",
		CredentialSchema: CredentialSchema{
			ID:   "",
			Type: "",
		},
	}
	require.Equal(t, want, vc)

	ctx := context.Background()
	mz, err := vc.Merklize(ctx)
	require.NoError(t, err)
	path, err := mz.ResolveDocPath("credentialSubject.price")
	require.NoError(t, err)
	_, err = mz.Entry(path)
	require.NoError(t, err)
}
