package slogpfx

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"testing"
	"testing/slogtest"

	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(slog.NewJSONHandler(&buf, nil), &HandlerOptions{
		PrefixKeys: nil,
		PrefixFormatter: func(prefixes []slog.Value) string {
			return "ABC"
		},
	})
	results := func() []map[string]any {
		var results []map[string]any
		dec := json.NewDecoder(&buf)
		for {
			var m map[string]any
			if err := dec.Decode(&m); err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}
			results = append(results, m)
		}
		return results
	}
	err := slogtest.TestHandler(h, results)
	require.NoError(t, err)
}
