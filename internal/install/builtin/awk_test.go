package builtin

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHost struct{}

func (testHost) Log(tag string, message ...string)                     {}
func (testHost) File(ctx context.Context, path string) []byte          { return nil }
func (testHost) Write(ctx context.Context, path string) io.WriteCloser { return nil }
func (testHost) MkDir(ctx context.Context, path string)                {}
func (testHost) Remove(ctx context.Context, path string)               {}
func (testHost) Move(ctx context.Context, from, to string) error       { return nil }

func TestAwk(t *testing.T) {
	// github release parsing from garden.io install script.
	json := `{"id":46049006,"tag_name":"0.12.24","update_url":"/garden-io/garden/releases/tag/0.12.24","update_authenticity_token":"SLV3kvPgqL6L1eX1Xd1w6NQq+WN5ullhmregGuptKSNkkjqL1I3MeAtDVewrOkT/I5y3u2/ZhJeaLA/uGtUbYQ==","delete_url":"/garden-io/garden/releases/tag/0.12.24","delete_authenticity_token":"ztsz0EsObINs7+G4n6d+Bl7YcTI+DcZowoNULGRTgqtdnyEhB9ruBJhY2OixNjC1+Hnni29HFlP3IG5TvLNZwQ==","edit_url":"/garden-io/garden/releases/edit/0.12.24"}`

	buf := &bytes.Buffer{}

	host := testHost{}
	ios := IOs{
		In:  strings.NewReader(json),
		Out: buf,
	}

	err := Awk(context.Background(), host, ios, []string{
		`-F[,:}]`,
		// \042 is an octal escaped quote ". Its up to the program to interpret it.
		"{for(i=1;i<=NF;i++){if($i~/tag_name\\042/){print $(i+1)}}}",
	})

	assert.NoError(t, err, "awk run should not error")
	assert.Equal(t, "\"0.12.24\"\n", buf.String(), "json parse failed")
}
