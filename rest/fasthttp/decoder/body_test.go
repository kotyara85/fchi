package decoder

import (
	"bytes"
	"errors"
	"github.com/swaggest/rest"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validatorFunc func(in rest.ParamIn, namedData map[string]interface{}) error

func (v validatorFunc) ValidateRequestData(in rest.ParamIn, namedData map[string]interface{}) error {
	return v(in, namedData)
}

func Test_decodeJSONBody(t *testing.T) {
	createBody := bytes.NewReader(
		[]byte(`{"amount": 123,"customerId": "248df4b7-aa70-47b8-a036-33ac447e668d","type": "withdraw"}`))
	createReq, err := http.NewRequest(http.MethodPost, "/US/order/348df4b7-aa70-47b8-a036-33ac447e668d", createBody)
	assert.NoError(t, err)

	type Input struct {
		Amount     int    `json:"amount"`
		CustomerID string `json:"customerId"`
		Type       string `json:"type"`
	}

	i := Input{}
	assert.NoError(t, decodeJSONBody(createReq, &i, nil))
	assert.Equal(t, 123, i.Amount)
	assert.Equal(t, "248df4b7-aa70-47b8-a036-33ac447e668d", i.CustomerID)
	assert.Equal(t, "withdraw", i.Type)

	vl := validatorFunc(func(in rest.ParamIn, namedData map[string]interface{}) error {
		return nil
	})

	i = Input{}
	_, err = createBody.Seek(0, io.SeekStart)
	assert.NoError(t, err)
	assert.NoError(t, decodeJSONBody(createReq, &i, vl))
	assert.Equal(t, 123, i.Amount)
	assert.Equal(t, "248df4b7-aa70-47b8-a036-33ac447e668d", i.CustomerID)
	assert.Equal(t, "withdraw", i.Type)
}

func Test_decodeJSONBody_emptyBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "any", nil)
	require.NoError(t, err)

	var i []int

	err = decodeJSONBody(req, &i, nil)
	assert.EqualError(t, err, "missing request body to decode json")
}

func Test_decodeJSONBody_badContentType(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "any", bytes.NewBufferString("123"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "text/plain")

	var i []int

	err = decodeJSONBody(req, &i, nil)
	assert.EqualError(t, err, "request with \"application/json\" content type expected \"text/plain\" received")
}

func Test_decodeJSONBody_decodeFailed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "any", bytes.NewBufferString("abc"))
	require.NoError(t, err)

	var i []int

	err = decodeJSONBody(req, &i, nil)
	assert.EqualError(t, err, "failed to decode json: invalid character 'a' looking for beginning of value")
}

func Test_decodeJSONBody_unmarshalFailed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "any", bytes.NewBufferString("123"))
	require.NoError(t, err)

	var i []int

	err = decodeJSONBody(req, &i, nil)
	assert.EqualError(t, err, "json: cannot unmarshal number into Go value of type []int")
}

func Test_decodeJSONBody_validateFailed(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "any", bytes.NewBufferString("[123]"))
	require.NoError(t, err)

	var i []int

	vl := validatorFunc(func(in rest.ParamIn, namedData map[string]interface{}) error {
		return errors.New("failed")
	})

	err = decodeJSONBody(req, &i, vl)
	assert.EqualError(t, err, "failed")
}
