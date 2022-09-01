package merklize

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
	"github.com/piprate/json-gold/ld"
	"github.com/stretchr/testify/require"
)

func getDataset(t testing.TB) *ld.RDFDataset {
	in := `{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://w3id.org/citizenship/v1",
    "https://w3id.org/security/bbs/v1"
  ],
  "id": "https://issuer.oidp.uscis.gov/credentials/83627465",
  "type": ["VerifiableCredential", "PermanentResidentCard"],
  "issuer": "did:example:489398593",
  "identifier": 83627465,
  "name": "Permanent Resident Card",
  "description": "Government of Example Permanent Resident Card.",
  "issuanceDate": "2019-12-03T12:19:52Z",
  "expirationDate": "2029-12-03T12:19:52Z",
  "credentialSubject": [
    {
      "id": "did:example:b34ca6cd37bbf23",
      "type": ["PermanentResident", "Person"],
      "givenName": "JOHN",
      "familyName": "SMITH",
      "gender": "Male",
      "image": "data:image/png;base64,iVBORw0KGgokJggg==",
      "residentSince": "2015-01-01",
      "lprCategory": "C09",
      "lprNumber": "999-999-999",
      "commuterClassification": "C1",
      "birthCountry": "Bahamas",
      "birthDate": "1958-07-17"
    },
    {
      "id": "did:example:b34ca6cd37bbf24",
      "type": ["PermanentResident", "Person"],
      "givenName": "JOHN",
      "familyName": "SMITH",
      "gender": "Male",
      "image": "data:image/png;base64,iVBORw0KGgokJggg==",
      "residentSince": "2015-01-01",
      "lprCategory": "C09",
      "lprNumber": "999-999-999",
      "commuterClassification": "C1",
      "birthCountry": "Bahamas",
      "birthDate": "1958-07-18"
    }
  ]
}`
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(in), &obj)
	if err != nil {
		panic(err)
	}

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Algorithm = "URDNA2015"

	out4, err := proc.Normalize(obj, options)
	require.NoError(t, err)

	out5, ok := out4.(*ld.RDFDataset)
	require.True(t, ok, "%[1]T\n%[1]v", out4)

	return out5
}

func TestEntriesFromRDF(t *testing.T) {
	dataset := getDataset(t)

	entries, err := EntriesFromRDF(dataset)
	require.NoError(t, err)

	if false {
		for i, e := range entries {
			t.Logf("%2v: %v => %v", i, e.key, e.value)
		}
	}

	wantEntries := []RDFEntry{
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://schema.org/birthDate"},
			value: "1958-07-17",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://schema.org/familyName"},
			value: "SMITH",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://schema.org/gender"},
			value: "Male",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://schema.org/givenName"},
			value: "JOHN",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://schema.org/image"},
			value: "data:image/png;base64,iVBORw0KGgokJggg==",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 0},
			value: "http://schema.org/Person",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 1},
			value: "https://w3id.org/citizenship#PermanentResident",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"https://w3id.org/citizenship#birthCountry"},
			value: "Bahamas",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"https://w3id.org/citizenship#commuterClassification"},
			value: "C1",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"https://w3id.org/citizenship#lprCategory"},
			value: "C09",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"https://w3id.org/citizenship#lprNumber"},
			value: "999-999-999",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0,
				"https://w3id.org/citizenship#residentSince"},
			value: "2015-01-01",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://schema.org/birthDate"},
			value: "1958-07-18",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://schema.org/familyName"},
			value: "SMITH",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://schema.org/gender"},
			value: "Male",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://schema.org/givenName"},
			value: "JOHN",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://schema.org/image"},
			value: "data:image/png;base64,iVBORw0KGgokJggg==",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 0},
			value: "http://schema.org/Person",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 1},
			value: "https://w3id.org/citizenship#PermanentResident",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"https://w3id.org/citizenship#birthCountry"},
			value: "Bahamas",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"https://w3id.org/citizenship#commuterClassification"},
			value: "C1",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"https://w3id.org/citizenship#lprCategory"},
			value: "C09",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"https://w3id.org/citizenship#lprNumber"},
			value: "999-999-999",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1,
				"https://w3id.org/citizenship#residentSince"},
			value: "2015-01-01",
		},
		{
			key:   []interface{}{"http://schema.org/description"},
			value: "Government of Example Permanent Resident Card.",
		},
		{
			key:   []interface{}{"http://schema.org/identifier"},
			value: int64(83627465),
		},
		{
			key:   []interface{}{"http://schema.org/name"},
			value: "Permanent Resident Card",
		},
		{
			key: []interface{}{
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 0},
			value: "https://w3id.org/citizenship#PermanentResidentCard",
		},
		{
			key: []interface{}{
				"http://www.w3.org/1999/02/22-rdf-syntax-ns#type", 1},
			value: "https://www.w3.org/2018/credentials#VerifiableCredential",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 0},
			value: "did:example:b34ca6cd37bbf23",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#credentialSubject", 1},
			value: "did:example:b34ca6cd37bbf24",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#expirationDate"},
			value: "2029-12-03T12:19:52Z",
		},
		{
			key: []interface{}{
				"https://www.w3.org/2018/credentials#issuanceDate"},
			value: "2019-12-03T12:19:52Z",
		},
		{
			key:   []interface{}{"https://www.w3.org/2018/credentials#issuer"},
			value: "did:example:489398593",
		},
	}
	require.Equal(t, wantEntries, entries)
}

func TestProof(t *testing.T) {
	dataset := getDataset(t)

	entries, err := EntriesFromRDF(dataset)
	require.NoError(t, err)

	ctx := context.Background()

	mt, err := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 40)
	require.NoError(t, err)

	err = AddEntriesToMerkleTree(ctx, mt, entries)
	require.NoError(t, err)

	// [https://www.w3.org/2018/credentials#credentialSubject 1 http://schema.org/birthDate] => 1958-07-18
	path := Path{
		"https://www.w3.org/2018/credentials#credentialSubject", 1,
		"http://schema.org/birthDate"}

	entry, err := NewRDFEntry(path, "1958-07-18")
	require.NoError(t, err)

	key, val, err := entry.KeyValueHashes()
	require.NoError(t, err)

	p, _, err := mt.GenerateProof(ctx, key, nil)
	require.NoError(t, err)

	ok := merkletree.VerifyProof(mt.Root(), p, key, val)
	require.True(t, ok)
}

func TestProofInteger(t *testing.T) {
	dataset := getDataset(t)

	entries, err := EntriesFromRDF(dataset)
	require.NoError(t, err)

	ctx := context.Background()

	mt, err := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 40)
	require.NoError(t, err)

	err = AddEntriesToMerkleTree(ctx, mt, entries)
	require.NoError(t, err)

	entry, err := NewRDFEntry(Path{"http://schema.org/identifier"}, 83627465)
	require.NoError(t, err)

	key, val, err := entry.KeyValueHashes()
	require.NoError(t, err)

	p, _, err := mt.GenerateProof(ctx, key, nil)
	require.NoError(t, err)

	ok := merkletree.VerifyProof(mt.Root(), p, key, val)
	require.True(t, ok)
}

func TestNewRelationship(t *testing.T) {
	iri := func(in string) ld.IRI {
		i := ld.NewIRI(in)
		return *i
	}
	dataset := getDataset(t)
	if false {
		logDataset(dataset)
	}

	rs, err := newRelationship(dataset.Graphs["@default"])
	require.NoError(t, err)
	wantRS := &relationship{
		parents: map[ld.IRI]quadKey{
			iri("did:example:b34ca6cd37bbf23"): {
				subject:   iri("https://issuer.oidp.uscis.gov/credentials/83627465"),
				predicate: iri("https://www.w3.org/2018/credentials#credentialSubject"),
			},
			iri("did:example:b34ca6cd37bbf24"): {
				subject:   iri("https://issuer.oidp.uscis.gov/credentials/83627465"),
				predicate: iri("https://www.w3.org/2018/credentials#credentialSubject"),
			},
		},
		children: map[ld.IRI][]ld.IRI{
			iri("https://issuer.oidp.uscis.gov/credentials/83627465"): {
				iri("did:example:b34ca6cd37bbf23"),
				iri("did:example:b34ca6cd37bbf24"),
			},
		},
	}
	require.Equal(t, wantRS, rs)
}

func logDataset(in *ld.RDFDataset) {
	for s, gs := range in.Graphs {
		fmt.Printf("[6] %v: %v\n", s, len(gs))
		for i, g := range gs {
			subject := "nil"
			if g.Subject != nil {
				subject = g.Subject.GetValue()
			}
			predicate := "nil"
			if g.Predicate != nil {
				predicate = g.Predicate.GetValue()
			}
			object := "nil"
			var ol2 string
			ol, olOK := g.Object.(*ld.Literal)
			if olOK {
				ol2 = ol.Datatype
			}

			if g.Object != nil {
				object = g.Object.GetValue()
			}
			graph := "nil"
			if g.Graph != nil {
				graph = g.Graph.GetValue()
			}
			fmt.Printf(`[7] %v:
	Subject [%T]: %v
	Predicate [%T]: %v
	Object [%T]: %v %v
	Graph [%T]: %v
`, i,
				g.Subject, subject,
				g.Predicate, predicate,
				g.Object, object, ol2,
				g.Graph, graph)
		}
	}
}
