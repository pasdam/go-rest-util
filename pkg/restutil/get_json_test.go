package restutil_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/pasdam/go-rest-util/pkg/restutil"
	"github.com/pasdam/mockit/mockit"
	"github.com/stretchr/testify/assert"
)

func TestGetJSON(t *testing.T) {
	type body struct {
		Field string `json:"field"`
	}
	type mocks struct {
		body         body
		getErr       error
		unmarshalErr error
	}
	type args struct {
		url string
	}
	tests := []struct {
		name  string
		mocks mocks
		args  args
	}{
		{
			name: "Should propagate error if Get raises it",
			mocks: mocks{
				body:         body{},
				getErr:       errors.New("some-get-error"),
				unmarshalErr: nil,
			},
			args: args{
				url: "some-get-error-url",
			},
		},
		{
			name: "Should propagate error if json.Unmarshal raises it",
			mocks: mocks{
				body:         body{},
				getErr:       nil,
				unmarshalErr: errors.New("some-unmarshal-error"),
			},
			args: args{
				url: "some-unmarshal-error-url",
			},
		},
		{
			name: "Should parse JSON correctly",
			mocks: mocks{
				body: body{
					Field: "some-success-value",
				},
				getErr:       nil,
				unmarshalErr: nil,
			},
			args: args{
				url: "some-success-url",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := body{}
			headers := http.Header{}

			bodyJSON, _ := json.Marshal(tt.mocks.body)
			mockit.MockFunc(t, restutil.Get).With(tt.args.url, headers).Return(bodyJSON, tt.mocks.getErr)
			if tt.mocks.unmarshalErr != nil {
				mockit.MockFunc(t, json.Unmarshal).With(bodyJSON, &b).Return(tt.mocks.unmarshalErr)
			}

			var wantErr error
			if tt.mocks.getErr != nil {
				wantErr = tt.mocks.getErr
			} else {
				wantErr = tt.mocks.unmarshalErr
			}

			err := restutil.GetJSON(tt.args.url, headers, &b)

			assert.Equal(t, wantErr, err)
			assert.Equal(t, tt.mocks.body, b)
		})
	}
}
