package kubenurse

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestServerHandler(t *testing.T) {
	r := require.New(t)

	fakeClient := fake.NewFakeClient()
	kubenurse, err := New(fakeClient)

	r.NoError(err)
	r.NotNil(kubenurse)

	ts := httptest.NewServer(kubenurse.http.Handler)
	defer ts.Close()

	var tests = map[string]struct {
		wantCode int
	}{
		"/": {
			wantCode: http.StatusMovedPermanently,
		},
		"/ready": {
			wantCode: http.StatusOK,
		},
		"/alive": {
			// 500 since servicechecks won't work
			wantCode: http.StatusInternalServerError,
		},
		"/alwayshappy": {
			wantCode: http.StatusOK,
		},
		// TODO: also test that metrics are present
		"/metrics": {
			wantCode: http.StatusOK,
		},
	}

	testClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for path, tc := range tests {
		t.Run(path, func(t *testing.T) {
			r := require.New(t)

			res, err := testClient.Get(ts.URL + path)
			r.NoError(err)
			r.Equal(tc.wantCode, res.StatusCode)
		})
	}
}
