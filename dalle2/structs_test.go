package dalle2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const sampleJson = "{\"n\":4,\"size\":\"1024x1024\",\"response_format\":\"b64_json\",\"user\":\"test_user\"}"

func TestJsonMarshalApiOption(t *testing.T) {
	testOption := &apiOption{}
	if err := json.Unmarshal([]byte(sampleJson), testOption); err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, 4, testOption.N)
	assert.EqualValues(t, LARGE, testOption.Size)
	assert.EqualValues(t, BASE64, testOption.Format)
	assert.EqualValues(t, "test_user", testOption.User)

	b, err := json.Marshal(testOption)
	if err != nil {
		t.Fatalf("error unmarshaling json: %s", err)
	}
	assert.EqualValues(t, sampleJson, string(b))
}

func TestOption(t *testing.T) {
	testOption := &apiOption{}
	for _, o := range []Option{
		WithFormat(BASE64),
		WithNumImages(2),
		WithSize(MEDIUM),
		WithUser("test_user"),
	} {
		o(testOption)
	}

	assert.EqualValues(t, 2, testOption.N)
	assert.EqualValues(t, MEDIUM, testOption.Size)
	assert.EqualValues(t, BASE64, testOption.Format)
	assert.EqualValues(t, "test_user", testOption.User)
}
