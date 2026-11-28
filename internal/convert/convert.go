package convert

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

type setlistItemKind int

const (
	kindScore setlistItemKind = iota
	kindBookmark
	kindPlaceholder
)

type setlistItem struct {
	Kind  setlistItemKind
	Title string
	Path  string
}

const forScore12PlusSuffix = " (forScore12plus)"

const (
	// MaxSetlistEntries is the supported setlist size. Real rehearsal and
	// concert setlists are far smaller; this bound keeps malformed input cheap.
	MaxSetlistEntries = 100

	// MaxInputBytes bounds the input file before any decompression.
	MaxInputBytes int64 = 2 << 20

	// MaxDecompressedBytes bounds the native plist after gzip decompression.
	MaxDecompressedBytes int64 = 4 << 20

	// MaxDepth bounds nesting in the NSKeyedArchiver object graph.
	MaxDepth = 64

	// MaxStringBytes bounds each title and path.
	MaxStringBytes = 16 << 10
)

type limits struct {
	maxEntries           int
	maxInputBytes        int64
	maxDecompressedBytes int64
	maxDepth             int
	maxStringBytes       int
}

var defaultLimits = limits{
	maxEntries:           MaxSetlistEntries,
	maxInputBytes:        MaxInputBytes,
	maxDecompressedBytes: MaxDecompressedBytes,
	maxDepth:             MaxDepth,
	maxStringBytes:       MaxStringBytes,
}

// NativeOutputPath returns a sibling path that keeps the input extension and adds the forScore 12+ suffix.
func NativeOutputPath(input string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, stem+forScore12PlusSuffix+ext)
}

// TextOutputPath returns a sibling Open Setlist path: stem (forScore12plus).4ss.
func TextOutputPath(input string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	return filepath.Join(dir, stem+forScore12PlusSuffix+".4ss")
}

func SetlistTitle(input string) string {
	base := filepath.Base(input)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// ReadInput reads an input file without exceeding MaxInputBytes.
func ReadInput(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, MaxInputBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > MaxInputBytes {
		return nil, fmt.Errorf("input file exceeds %d bytes", MaxInputBytes)
	}
	return data, nil
}

// FromText builds an Open Setlist XML document with one placeholder per non-blank line.
func FromText(data []byte, title string) ([]byte, error) {
	return fromText(data, title, defaultLimits)
}

func fromText(data []byte, title string, lim limits) ([]byte, error) {
	if int64(len(data)) > lim.maxInputBytes {
		return nil, fmt.Errorf("text input exceeds %d bytes", lim.maxInputBytes)
	}

	var items []setlistItem
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" {
			continue
		}
		if len(line) > lim.maxStringBytes {
			return nil, fmt.Errorf("line %d exceeds %d bytes", i+1, lim.maxStringBytes)
		}
		items = append(items, setlistItem{Kind: kindPlaceholder, Title: line})
		if len(items) > lim.maxEntries {
			return nil, entryLimitError(len(items), lim.maxEntries)
		}
	}
	return marshalOpenSetlist(title, items)
}

func entryLimitError(count, max int) error {
	return fmt.Errorf("setlist has %d entries; the supported maximum is %d. Split the setlist or report a case for a higher limit", count, max)
}

// FromNative builds an Open Setlist XML document from a gzip-compressed native forScore setlist.
func FromNative(data []byte, title string) ([]byte, error) {
	return fromNative(data, title, defaultLimits)
}

func fromNative(data []byte, title string, lim limits) ([]byte, error) {
	if int64(len(data)) > lim.maxInputBytes {
		return nil, fmt.Errorf("input file exceeds %d bytes", lim.maxInputBytes)
	}
	if isOpenSetlistXML(data) {
		return nil, fmt.Errorf("file is already an Open Setlist XML document")
	}
	if !isGzip(data) {
		return nil, fmt.Errorf("not a gzip-compressed forScore setlist (expected native .4ss)")
	}

	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}
	defer zr.Close()
	raw, err := io.ReadAll(io.LimitReader(zr, lim.maxDecompressedBytes+1))
	if err != nil {
		return nil, fmt.Errorf("gzip: %w", err)
	}
	if int64(len(raw)) > lim.maxDecompressedBytes {
		return nil, fmt.Errorf("decompressed setlist exceeds %d bytes", lim.maxDecompressedBytes)
	}

	if !bytes.HasPrefix(raw, []byte("bplist")) {
		return nil, fmt.Errorf("not a binary plist after gzip decompress")
	}

	var archive map[string]any
	if _, err := plist.Unmarshal(raw, &archive); err != nil {
		return nil, fmt.Errorf("plist: %w", err)
	}

	items, err := itemsFromArchive(archive, lim)
	if err != nil {
		return nil, err
	}

	return marshalOpenSetlist(title, items)
}

func isGzip(data []byte) bool {
	return len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b
}

func isOpenSetlistXML(data []byte) bool {
	s := strings.TrimSpace(string(data))
	return strings.HasPrefix(s, "<?xml") || strings.HasPrefix(s, "<forScore")
}

func itemsFromArchive(archive map[string]any, lim limits) ([]setlistItem, error) {
	objects, ok := archive["$objects"].([]any)
	if !ok {
		return nil, fmt.Errorf("NSKeyedArchiver: missing $objects")
	}
	top, ok := archive["$top"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("NSKeyedArchiver: missing $top")
	}
	resolved, err := resolveArchive(objects, top["setlist"], lim.maxDepth)
	if err != nil {
		return nil, fmt.Errorf("NSKeyedArchiver: %w", err)
	}
	if resolved == nil {
		return nil, fmt.Errorf("NSKeyedArchiver: missing $top.setlist")
	}
	rawItems, ok := resolved.([]any)
	if !ok {
		return nil, fmt.Errorf("NSKeyedArchiver: $top.setlist is not an array")
	}
	if len(rawItems) > lim.maxEntries {
		return nil, entryLimitError(len(rawItems), lim.maxEntries)
	}

	out := make([]setlistItem, 0, len(rawItems))
	for i, raw := range rawItems {
		dict, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("setlist item %d is not a dictionary", i+1)
		}
		item, err := classifyItem(dict, lim)
		if err != nil {
			return nil, fmt.Errorf("setlist item %d: %w", i+1, err)
		}
		out = append(out, item)
	}
	return out, nil
}

func classifyItem(dict map[string]any, lim limits) (setlistItem, error) {
	title, _ := dict["Title"].(string)
	path, _ := dict["FilePath"].(string)
	if err := checkString("Title", title, lim); err != nil {
		return setlistItem{}, err
	}
	if err := checkString("FilePath", path, lim); err != nil {
		return setlistItem{}, err
	}

	if isTruthy(dict["Placeholder"]) {
		if title == "" {
			return setlistItem{}, fmt.Errorf("placeholder is missing Title")
		}
		return setlistItem{Kind: kindPlaceholder, Title: title}, nil
	}

	if isBookmark(dict) {
		if title == "" || path == "" {
			return setlistItem{}, fmt.Errorf("bookmark requires Title and FilePath")
		}
		return setlistItem{Kind: kindBookmark, Title: title, Path: path}, nil
	}

	if strings.HasSuffix(strings.ToLower(path), ".pdf") {
		return setlistItem{Kind: kindScore, Title: title, Path: path}, nil
	}

	return setlistItem{}, fmt.Errorf("unrecognized item (Title=%q FilePath=%q)", title, path)
}

func checkString(label, value string, lim limits) error {
	if len(value) > lim.maxStringBytes {
		return fmt.Errorf("%s exceeds %d bytes", label, lim.maxStringBytes)
	}
	return nil
}

func isBookmark(dict map[string]any) bool {
	if _, ok := dict["Bookmark"]; ok {
		if !isTruthy(dict["Bookmark"]) {
			if s, ok := dict["Bookmark"].(string); ok {
				return s != ""
			}
			return false
		}
		return true
	}
	for _, key := range []string{"Type", "Kind"} {
		if s, ok := dict[key].(string); ok && strings.EqualFold(s, "bookmark") {
			return true
		}
	}
	return false
}

func isTruthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case int64:
		return t != 0
	case uint64:
		return t != 0
	case int:
		return t != 0
	case float64:
		return t != 0
	case string:
		return strings.EqualFold(t, "true") || t == "1"
	default:
		return false
	}
}

type resolver struct {
	objects  []any
	memo     map[int]any
	visiting map[int]bool
	maxDepth int
}

func resolveArchive(objects []any, v any, maxDepth int) (any, error) {
	r := resolver{
		objects:  objects,
		memo:     make(map[int]any),
		visiting: make(map[int]bool),
		maxDepth: maxDepth,
	}
	return r.resolve(v, 0)
}

func (r *resolver) resolve(v any, depth int) (any, error) {
	if depth > r.maxDepth {
		return nil, fmt.Errorf("object graph exceeds maximum depth %d", r.maxDepth)
	}

	switch t := v.(type) {
	case plist.UID:
		idx := int(t)
		if idx < 0 || idx >= len(r.objects) {
			return nil, fmt.Errorf("UID %d out of range", idx)
		}
		if r.visiting[idx] {
			return nil, fmt.Errorf("cyclic object reference at UID %d", idx)
		}
		if cached, ok := r.memo[idx]; ok {
			return cached, nil
		}
		r.visiting[idx] = true
		resolved, err := r.resolve(r.objects[idx], depth+1)
		delete(r.visiting, idx)
		if err != nil {
			return nil, err
		}
		r.memo[idx] = resolved
		return resolved, nil
	case []any:
		out := make([]any, len(t))
		for i, element := range t {
			resolved, err := r.resolve(element, depth+1)
			if err != nil {
				return nil, fmt.Errorf("array element %d: %w", i, err)
			}
			out[i] = resolved
		}
		return out, nil
	case map[string]any:
		if _, ok := t["$classname"]; ok {
			return t, nil
		}
		keys, hasKeys := t["NS.keys"]
		objects, hasObjects := t["NS.objects"]
		if hasKeys && hasObjects {
			resolvedKeys, err := r.resolve(keys, depth+1)
			if err != nil {
				return nil, err
			}
			resolvedValues, err := r.resolve(objects, depth+1)
			if err != nil {
				return nil, err
			}
			keyList, ok := resolvedKeys.([]any)
			if !ok {
				return nil, fmt.Errorf("NSDictionary keys are not an array")
			}
			valueList, ok := resolvedValues.([]any)
			if !ok {
				return nil, fmt.Errorf("NSDictionary objects are not an array")
			}
			if len(keyList) != len(valueList) {
				return nil, fmt.Errorf("NSDictionary key/value count mismatch")
			}
			out := make(map[string]any, len(keyList))
			for i, key := range keyList {
				keyString, ok := key.(string)
				if !ok {
					return nil, fmt.Errorf("NSDictionary key is not a string")
				}
				out[keyString] = valueList[i]
			}
			return out, nil
		}
		if hasObjects {
			return r.resolve(objects, depth+1)
		}
		out := make(map[string]any, len(t))
		for key, value := range t {
			if key == "$class" {
				continue
			}
			resolved, err := r.resolve(value, depth+1)
			if err != nil {
				return nil, err
			}
			out[key] = resolved
		}
		return out, nil
	default:
		return v, nil
	}
}

func marshalOpenSetlist(title string, items []setlistItem) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" ?>`)
	buf.WriteByte('\n')
	buf.WriteString(`<forScore kind="setlist" version="1.0" title="`)
	if err := xml.EscapeText(&buf, []byte(title)); err != nil {
		return nil, err
	}
	buf.WriteString(`">`)
	buf.WriteByte('\n')
	for _, it := range items {
		name := "score"
		switch it.Kind {
		case kindBookmark:
			name = "bookmark"
		case kindPlaceholder:
			name = "placeholder"
		}
		buf.WriteByte('\t')
		buf.WriteByte('<')
		buf.WriteString(name)
		buf.WriteString(` title="`)
		if err := xml.EscapeText(&buf, []byte(it.Title)); err != nil {
			return nil, err
		}
		buf.WriteByte('"')
		if it.Kind != kindPlaceholder {
			buf.WriteString(` path="`)
			if err := xml.EscapeText(&buf, []byte(it.Path)); err != nil {
				return nil, err
			}
			buf.WriteByte('"')
		}
		buf.WriteString(" />")
		buf.WriteByte('\n')
	}
	buf.WriteString("</forScore>\n")
	return buf.Bytes(), nil
}
