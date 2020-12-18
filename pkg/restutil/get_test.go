package restutil_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pasdam/go-rest-util/pkg/restutil"
	"github.com/pasdam/mockit/matchers/argument"
	"github.com/pasdam/mockit/mockit"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type mocks struct {
		getErr        error
		newRequestErr error
		readAllErr    error
	}
	type args struct {
		url    string
		header http.Header
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    []byte
		wantErr string
	}{
		{
			name: "Should propagate error if http.Client.Do raises it",
			mocks: mocks{
				getErr: errors.New("some-get-error"),
			},
			args: args{
				url: "some-get-error-url",
				header: map[string][]string{
					"header-1": {"header-1-value"},
					"header-2": {"header-2-value"},
				},
			},
			want:    nil,
			wantErr: "unable to perform request, some-get-error",
		},
		{
			name: "Should propagate error if ioutil.Readall raises it",
			mocks: mocks{
				readAllErr: errors.New("some-read-all-error"),
			},
			args: args{
				url: "some-read-all-error-url",
				header: map[string][]string{
					"header-1": {"header-1-value"},
					"header-2": {"header-2-value"},
				},
			},
			want:    nil,
			wantErr: "unable to read the response body, some-read-all-error",
		},
		{
			name: "Should propagate error if http.NewRequest raises it",
			mocks: mocks{
				newRequestErr: errors.New("some-new-request-error"),
			},
			args: args{
				url: "some-new-request-error-url",
				header: map[string][]string{
					"header-3": {"header-3-value"},
					"header-4": {"header-4-value"},
				},
			},
			want:    nil,
			wantErr: "unable to create request, some-new-request-error",
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
				// assert.Equal(t, tt.args.header, r.Header)
				for key, value := range tt.args.header {
					assert.Equal(t, value[0], r.Header.Get(key))
				}
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/"+tt.args.url, r.RequestURI)
				fmt.Fprintln(w, string(tt.want))
			}))
			defer ts.Close()

			url := ts.URL + "/" + tt.args.url

			var wantErr error = tt.mocks.getErr
			if tt.mocks.getErr != nil {
				c := &http.Client{}
				mockit.MockMethodForAll(t, c, c.Do).With(argument.Any).Return(nil, wantErr)

			} else if tt.mocks.newRequestErr != nil {
				wantErr = tt.mocks.newRequestErr
				mockit.MockFunc(t, http.NewRequest).With("GET", url, nil).Return(nil, wantErr)

			} else if tt.mocks.readAllErr != nil {
				wantErr = tt.mocks.readAllErr
				mockit.MockFunc(t, ioutil.ReadAll).With(argument.Any).Return(nil, wantErr)
			}

			got, err := restutil.Get(url, tt.args.header)

			if wantErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
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
