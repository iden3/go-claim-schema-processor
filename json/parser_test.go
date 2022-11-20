package json

import (
	"encoding/json"
	"testing"

	"github.com/iden3/go-schema-processor/processor"
	"github.com/iden3/go-schema-processor/utils"
	"github.com/stretchr/testify/require"
)

func TestGetSerializedData(t *testing.T) {
	testCases := []struct {
		title       string
		indexFields []string
		indexValues []uint32
		valueFields []string
		valueValues []uint32
		expected    processor.ParsedSlots
		expectedErr string
	}{
		{
			title:       "index fills into one slot",
			indexFields: []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8"},
			indexValues: []uint32{1, 2, 3, 4, 5, 6, 7, 8},
			expected: processor.ParsedSlots{
				IndexA: []byte{
					1,
					2,
					3,
					4,
					5,
					6,
					7,
					8,
				},
				IndexB: []byte{},
				ValueA: []byte{},
				ValueB: []byte{},
			},
		},
		{
			title:       "index does not fills into one slot",
			indexFields: []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8"},
			indexValues: []uint32{
				4026531841, 1138881939, 2042196113, 674490440,
				2172737629, 3092268470, 3778125865, 811880050},
			expected: processor.ParsedSlots{
				IndexA: []byte{
					1, 0, 0, 240,
					147, 245, 225, 67,
					145, 112, 185, 121,
					72, 232, 51, 40,
					93, 88, 129, 129,
					182, 69, 80, 184,
					41, 160, 49, 225,
				},
				IndexB: []byte{114, 78, 100, 48},
				ValueA: []byte{},
				ValueB: []byte{},
			},
		},
		{
			title: "index overflows",
			indexFields: []string{
				"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8",
				"f9", "f10", "f11", "f12", "f13", "f14", "f15", "f16",
			},
			indexValues: []uint32{
				4026531841, 1138881939, 2042196113, 674490440,
				2172737629, 3092268470, 3778125865, 811880050,
				4026531841, 1138881939, 2042196113, 674490440,
				2172737629, 3092268470, 3778125865, 811880050,
			},
			expectedErr: "slots overflow",
		},
		{
			title:       "value fills into one slot",
			valueFields: []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8"},
			valueValues: []uint32{1, 2, 3, 4, 5, 6, 7, 8},
			expected: processor.ParsedSlots{
				IndexA: []byte{},
				IndexB: []byte{},
				ValueA: []byte{
					1,
					2,
					3,
					4,
					5,
					6,
					7,
					8,
				},
				ValueB: []byte{},
			},
		},
		{
			title:       "value does not fills into one slot",
			valueFields: []string{"f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8"},
			valueValues: []uint32{
				4026531841, 1138881939, 2042196113, 674490440,
				2172737629, 3092268470, 3778125865, 811880050},
			expected: processor.ParsedSlots{
				IndexA: []byte{},
				IndexB: []byte{},
				ValueA: []byte{
					1, 0, 0, 240,
					147, 245, 225, 67,
					145, 112, 185, 121,
					72, 232, 51, 40,
					93, 88, 129, 129,
					182, 69, 80, 184,
					41, 160, 49, 225,
				},
				ValueB: []byte{114, 78, 100, 48},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			var schema = &CommonJSONSerializationSchema{
				Index: struct {
					Type    string   `json:"type"`
					Default []string `json:"default"`
				}{Default: tc.indexFields},
				Value: struct {
					Type    string   `json:"type"`
					Default []string `json:"default"`
				}{Default: tc.valueFields},
			}
			var inputData = make(map[string]interface{})
			for i, k := range tc.indexFields {
				inputData[k] = tc.indexValues[i]
			}
			for i, k := range tc.valueFields {
				inputData[k] = tc.valueValues[i]
			}
			input, err := json.Marshal(inputData)
			require.NoError(t, err)
			slots, err := utils.FillClaimSlots(input, schema.Index.Default, schema.Value.Default)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, slots)
			}
		})
	}
}
