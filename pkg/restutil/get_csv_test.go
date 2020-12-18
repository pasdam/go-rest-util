package restutil_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/pasdam/go-rest-util/pkg/restutil"
	"github.com/pasdam/mockit/mockit"
	"github.com/stretchr/testify/assert"
)

func Test_GetCSV(t *testing.T) {
	type mocks struct {
		responseBody string
		responseErr  error
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    args
		want    [][]string
		wantErr error
	}{
		{
			name: "Should return csv rows if no error occurs",
			mocks: mocks{
				responseBody: "row0-col0,row0-col1\nrow1-col0,row1-col1",
				responseErr:  nil,
			},
			args: args{
				url: "some-url",
			},
			want: [][]string{
				{"row0-col0", "row0-col1"},
				{"row1-col0", "row1-col1"},
			},
			wantErr: nil,
		},
		{
			name: "Should propagate error if http.Get raises it",
			mocks: mocks{
				responseBody: "",
				responseErr:  errors.New("some-get-error"),
			},
			args: args{
				url: "some-get-error-url",
			},
			want:    nil,
			wantErr: errors.New("unable to perform request, some-get-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := mockit.MockFunc(t, http.Get).With(tt.args.url)
			var response *http.Response
			if tt.mocks.responseErr == nil {
				response = &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(tt.mocks.responseBody)),
				}
			}
			stub.Return(response, tt.mocks.responseErr)

			got, err := restutil.GetCSV(tt.args.url)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
