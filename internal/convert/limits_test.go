package convert

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"howett.net/plist"
)

func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func archiveWithItems(t *testing.T, items []any) []byte {
	t.Helper()
	archive := map[string]any{
		"$objects": []any{items},
		"$top":     map[string]any{"setlist": plist.UID(0)},
	}
	raw, err := plist.Marshal(archive, plist.BinaryFormat)
	if err != nil {
		t.Fatal(err)
	}
	return gzipBytes(t, raw)
}

func itemDict(title, path string) map[string]any {
	return map[string]any{"Title": title, "FilePath": path}
}

func TestFromTextEntryLimit(t *testing.T) {
	var lines []string
	for i := 0; i < MaxSetlistEntries; i++ {
		lines = append(lines, "Score "+strings.Repeat("x", i%3))
	}
	if _, err := FromText([]byte(strings.Join(lines, "\n")), "set"); err != nil {
		t.Fatalf("exactly %d entries should succeed: %v", MaxSetlistEntries, err)
	}

	lines = append(lines, "One too many")
	if _, err := FromText([]byte(strings.Join(lines, "\n")), "set"); err == nil {
		t.Fatalf("%d entries should fail", MaxSetlistEntries+1)
	} else if !strings.Contains(err.Error(), "supported maximum is 100") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromTextIgnoresBlankLinesForLimit(t *testing.T) {
	var lines []string
	for i := 0; i < MaxSetlistEntries; i++ {
		lines = append(lines, "Score", "", "  ")
	}
	if _, err := FromText([]byte(strings.Join(lines, "\n")), "set"); err != nil {
		t.Fatalf("blank lines should not count towards the limit: %v", err)
	}
}

func TestFromNativeEntryLimit(t *testing.T) {
	items := make([]any, MaxSetlistEntries)
	for i := range items {
		items[i] = itemDict("Score", "Score.pdf")
	}
	if _, err := FromNative(archiveWithItems(t, items), "set"); err != nil {
		t.Fatalf("exactly %d native items should succeed: %v", MaxSetlistEntries, err)
	}

	items = append(items, itemDict("Extra", "Extra.pdf"))
	if _, err := FromNative(archiveWithItems(t, items), "set"); err == nil {
		t.Fatalf("%d native items should fail", MaxSetlistEntries+1)
	} else if !strings.Contains(err.Error(), "supported maximum is 100") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFromTextInputLimit(t *testing.T) {
	lim := limits{maxInputBytes: 4, maxEntries: 100, maxDecompressedBytes: 4, maxDepth: 8, maxStringBytes: 8}
	if _, err := fromText([]byte("12345"), "set", lim); err == nil {
		t.Fatal("expected input size error")
	}
}

func TestFromNativeInputAndDecompressionLimits(t *testing.T) {
	lim := limits{maxInputBytes: 4, maxEntries: 100, maxDecompressedBytes: 16, maxDepth: 8, maxStringBytes: 8}
	if _, err := fromNative([]byte("12345"), "set", lim); err == nil {
		t.Fatal("expected input size error")
	}

	payload := gzipBytes(t, []byte(strings.Repeat("x", 32)))
	if _, err := fromNative(payload, "set", lim); err == nil {
		t.Fatal("expected decompression size error")
	}
}

func TestReadInputLimit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.4ss")
	data := bytes.Repeat([]byte("x"), int(MaxInputBytes)+1)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadInput(path); err == nil {
		t.Fatal("expected input size error")
	}
}

func TestCyclicObjectReference(t *testing.T) {
	archive := map[string]any{
		"$objects": []any{
			map[string]any{
				"NS.keys":    []any{"loop"},
				"NS.objects": []any{plist.UID(0)},
			},
		},
		"$top": map[string]any{"setlist": plist.UID(0)},
	}
	_, err := itemsFromArchive(archive, defaultLimits)
	if err == nil || !strings.Contains(err.Error(), "cyclic object reference") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestObjectGraphDepthLimit(t *testing.T) {
	archive := map[string]any{
		"$objects": []any{
			[]any{[]any{[]any{}}},
		},
		"$top": map[string]any{"setlist": plist.UID(0)},
	}
	lim := defaultLimits
	lim.maxDepth = 2
	if _, err := itemsFromArchive(archive, lim); err == nil || !strings.Contains(err.Error(), "maximum depth") {
		t.Fatalf("expected depth error, got %v", err)
	}
}

func TestFromNativeMalformedInputs(t *testing.T) {
	for _, tc := range []struct {
		name string
		data []byte
		want string
	}{
		{"not gzip", []byte("not a setlist"), "gzip-compressed"},
		{"gzip not bplist", gzipBytes(t, []byte("not a plist")), "binary plist"},
		{"gzip invalid plist", gzipBytes(t, []byte("bplist00bad")), "plist"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FromNative(tc.data, "set")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want error containing %q", err, tc.want)
			}
		})
	}
}

func TestClassifyItemAndXMLOutput(t *testing.T) {
	archive := map[string]any{
		"$objects": []any{
			[]any{
				map[string]any{"Title": "Bookmark", "FilePath": "Score.pdf", "Bookmark": true},
				itemDict("Upper", "Score.PDF"),
				map[string]any{"Title": "Placeholder", "Placeholder": true},
			},
		},
		"$top": map[string]any{"setlist": plist.UID(0)},
	}
	items, err := itemsFromArchive(archive, defaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 || items[0].Kind != kindBookmark || items[1].Kind != kindScore || items[2].Kind != kindPlaceholder {
		t.Fatalf("unexpected items: %+v", items)
	}

	out, err := FromText([]byte("A & B <C> \"D\"\n"), `Title & "quote"`)
	if err != nil {
		t.Fatal(err)
	}
	var parsed struct {
		Title string `xml:"title,attr"`
		Items []struct {
			Title string `xml:"title,attr"`
		} `xml:"placeholder"`
	}
	if err := xml.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("XML output is invalid: %v\n%s", err, out)
	}
	if parsed.Title != `Title & "quote"` || len(parsed.Items) != 1 || parsed.Items[0].Title != "A & B <C> \"D\"" {
		t.Fatalf("escaping failed: %+v", parsed)
	}
}

func TestStringLimit(t *testing.T) {
	lim := defaultLimits
	lim.maxStringBytes = 4
	if _, err := fromText([]byte("12345\n"), "set", lim); err == nil {
		t.Fatal("expected text string limit error")
	}
	archive := map[string]any{
		"$objects": []any{[]any{itemDict("long title", "Score.pdf")}},
		"$top":     map[string]any{"setlist": plist.UID(0)},
	}
	if _, err := itemsFromArchive(archive, lim); err == nil {
		t.Fatal("expected native string limit error")
	}
}
