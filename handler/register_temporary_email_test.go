package handler

import (
	"net/http"
	"testing"
)

func RestRegisterTemporaryEmail(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"正しいリクエストの時は201となる": {
			reqFile: "testdata/register_temporary_email/201_req.json.golden",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/register_temporary_email/201_rsp.json.golden",
			},
		},
		"リクエストデータが正しくない場合は401エラーを返す": {
			reqFile: "testdata/register_temporary_email/400_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/register_temporary_email/400_rsp.json.golden",
			},
		},
		"登録済みのメールアドレスは409エラーを返す": {
			reqFile: "testdata/register_temporary_email/409_req.json.golden",
			want: want{
				status:  http.StatusConflict,
				rspFile: "testdata/register_temporary_email/409_rsp.json.golden",
			},
		},
	}

	println(tests)
}
