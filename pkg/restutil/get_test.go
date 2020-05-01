package restutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pasdam/mockit/matchers/argument"
	"github.com/pasdam/mockit/mockit"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type mocks struct {
		getErr     error
		readAllErr error
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    []byte
		wantErr string
	}{
		{
			name: "Should propagate error if http.Get raises it",
			mocks: mocks{
				getErr:     errors.New("some-get-error"),
				readAllErr: nil,
			},
			args: args{
				url: "some-get-error-url",
			},
			want:    nil,
			wantErr: "unable to perform request, some-get-error",
		},
		{
			name: "Should propagate error if ioutil.Readall raises it",
			mocks: mocks{
				getErr:     nil,
				readAllErr: errors.New("some-read-all-error"),
			},
			args: args{
				url: "some-read-all-error-url",
			},
			want:    nil,
			wantErr: "unable to read the response body, some-read-all-error",
		},
		{
			name: "Should return body content if no errors occur",
			mocks: mocks{
				getErr:     nil,
				readAllErr: nil,
			},
			args: args{
				url: "some-successful-url",
			},
			want:    []byte("some-response-body"),
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/"+tt.args.url, r.RequestURI)
				fmt.Fprintln(w, string(tt.want))
			}))
			defer ts.Close()

			url := ts.URL + "/" + tt.args.url

			var wantErr error = tt.mocks.getErr
			if tt.mocks.getErr != nil {
				mockit.MockFunc(t, http.Get).With(url).Return(nil, tt.mocks.getErr)

			} else if tt.mocks.readAllErr != nil {
				wantErr = tt.mocks.readAllErr
				mockit.MockFunc(t, ioutil.ReadAll).With(argument.Any).Return(nil, tt.mocks.readAllErr)
			}

			got, err := Get(url)

			if wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.Nil(t, err)
			}
			if len(tt.want) > 0 {
				assert.Equal(t, string(tt.want)+"\n", string(got))
			} else {
				assert.Empty(t, got)
			}
		})
	}
}
