// utils/frontmatter.go — Unified utils package with adrg/frontmatter integration
// It preserves ExtractFrontMatter(map[string]string) API so you can drop‑in replace
// the old file, **plus** adds typed structs, directory scan, and flexible sorting.

package utils

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

// -----------------------------------------------------------------------------
// Section 1: Generic ExtractFrontMatter (map[string]string) — Backwards compatible
// -----------------------------------------------------------------------------

func ExtractFrontMatter(data []byte) map[string]string {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF}) // strip BOM
	data = bytes.TrimLeft(data, " \t\r\n")

	raw := make(map[string]interface{})
	if _, err := frontmatter.Parse(bytes.NewReader(data), &raw); err != nil && !errors.Is(err, frontmatter.ErrNotFound) {
		fmt.Printf("front‑matter parse error: %v\n", err)
		return nil
	}

	lower := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		lower[strings.ToLower(k)] = v
	}
	return convertMapInterfaceToString(lower)
}

// -----------------------------------------------------------------------------
// Section 2: Typed helpers — Rows, scanning, sorting
// -----------------------------------------------------------------------------

type FrontMatter struct {
	Title       string    `yaml:"title" toml:"title" json:"title"`
	Author      string    `yaml:"author" toml:"author" json:"author"`
	Description string    `yaml:"description" toml:"description" json:"description"`
	Tags        []string  `yaml:"tags" toml:"tags" json:"tags"`
	ParsedDate  time.Time `yaml:"-" toml:"-" json:"-"` // 手动解析
	Slug        string    `yaml:"slug" toml:"slug" json:"slug"`
}

type Row struct {
	FileName string
	Hash     string
	FM       FrontMatter
	Body     string
}

func ParseFile(path string, root string) (Row, error) {
	var row Row
	data, err := os.ReadFile(path)
	if err != nil {
		return row, err
	}

	var raw map[string]interface{}
	rest, err := frontmatter.Parse(bytes.NewReader(data), &raw)
	if err != nil && !errors.Is(err, frontmatter.ErrNotFound) {
		return row, err
	}

	var fm FrontMatter
	if val, ok := raw["title"]; ok {
		fm.Title, _ = val.(string)
	}
	if val, ok := raw["author"]; ok {
		fm.Author, _ = val.(string)
	}
	if val, ok := raw["description"]; ok {
		fm.Description, _ = val.(string)
	}
	if val, ok := raw["tags"]; ok {
		if arr, ok := val.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					fm.Tags = append(fm.Tags, s)
				}
			}
		}
	}
	if val, ok := raw["slug"]; ok {
		fm.Slug, _ = val.(string)
	}
	if val, ok := raw["date"]; ok {
		switch v := val.(type) {
		case string:
			formats := []string{"2006-01-02", "2006-01-02 15:04", time.RFC3339, time.RFC1123Z}
			for _, f := range formats {
				t, err := time.Parse(f, v)
				if err == nil {
					fm.ParsedDate = t
					break
				}
			}
		case time.Time:
			fm.ParsedDate = v
		}
	}

	body, _ := io.ReadAll(bytes.NewReader(rest))
	hash := md5.Sum(data)

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return row, err
	}

	row.FileName = rel
	row.FM = fm
	row.Body = string(body)
	row.Hash = hex.EncodeToString(hash[:])
	return row, nil
}

func ScanDir(root string) ([]Row, error) {
	var rows []Row
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return walkErr
		}
		r, perr := ParseFile(path, root)
		if perr != nil {
			return perr
		}
		rows = append(rows, r)
		return nil
	})
	return rows, err
}

// SortKey enumerates columns available for sorting.
// Extend as you add more fields.

type SortKey string

const (
	ByTitle  SortKey = "title"
	ByAuthor SortKey = "author"
	ByDate   SortKey = "date"
)

// SortRows orders rows in place by key; desc==true flips descending.
func SortRows(rows []Row, key SortKey, desc bool) {
	var less func(i, j int) bool
	switch key {
	case ByAuthor:
		less = func(i, j int) bool { return strings.ToLower(rows[i].FM.Author) < strings.ToLower(rows[j].FM.Author) }
	case ByDate:
		less = func(i, j int) bool { return rows[i].FM.ParsedDate.Before(rows[j].FM.ParsedDate) }
	default: // ByTitle
		less = func(i, j int) bool { return strings.ToLower(rows[i].FM.Title) < strings.ToLower(rows[j].FM.Title) }
	}
	sort.Slice(rows, func(i, j int) bool {
		if desc {
			return !less(i, j)
		}
		return less(i, j)
	})
}

// -----------------------------------------------------------------------------
// CSV helpers
// -----------------------------------------------------------------------------

func WriteCSV(rows []Row, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	if err := writer.Write([]string{"FileName", "Title", "Author", "Date", "Tags", "MD5"}); err != nil {
		return err
	}
	for _, r := range rows {
		if err := writer.Write([]string{
			r.FileName,
			r.FM.Title,
			r.FM.Author,
			r.FM.ParsedDate.Format("2006-01-02 15:04"),
			strings.Join(r.FM.Tags, ", "),
			r.Hash,
		}); err != nil {
			return err
		}
	}
	return nil
}

func ExportToCSV(rows []Row, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteCSV(rows, f)
}

// -----------------------------------------------------------------------------
// Sync rows to SQLite DB (new)
// -----------------------------------------------------------------------------

func SyncToDB(rows []Row, db *sql.DB) error {
	_, _ = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_uploads_filename ON uploads(filename);`)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(`
        INSERT OR IGNORE INTO uploads (filename, filepath, user_id, author, created_at)
        VALUES (?, ?, NULL, ?, ?);
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range rows {
		created := r.FM.ParsedDate.Format("2006-01-02 15:04")
		if r.FM.ParsedDate.IsZero() {
			created = time.Now().Format("2006-01-02 15:04")
		}
		if _, err := stmt.Exec(r.FileName, "content/"+r.FileName, r.FM.Author, created); err != nil {
			fmt.Println("跳过:", r.FileName, err)
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Section 3: Value‑conversion helpers (mostly unchanged from your legacy code)
// -----------------------------------------------------------------------------

func convertInterfaceArrayToString(arr []interface{}) string {
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		switch v := item.(type) {
		case string:
			out = append(out, v)
		case bool:
			out = append(out, strconv.FormatBool(v))
		case float32:
			out = append(out, strconv.FormatFloat(float64(v), 'f', -1, 32))
		case float64:
			out = append(out, strconv.FormatFloat(v, 'f', -1, 64))
		default: // all int kinds via reflect
			rv := reflect.ValueOf(v)
			if rv.Kind() >= reflect.Int && rv.Kind() <= reflect.Int64 {
				out = append(out, strconv.FormatInt(rv.Int(), 10))
			} else {
				fmt.Printf("unhandled array elem type %T\n", v)
			}
		}
	}
	return strings.Join(out, ", ")
}

func convertMapInterfaceToString(m map[string]interface{}) map[string]string {
	res := make(map[string]string, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			res[k] = val
		case bool:
			res[k] = strconv.FormatBool(val)
		case float32:
			res[k] = strconv.FormatFloat(float64(val), 'f', -1, 32)
		case float64:
			res[k] = strconv.FormatFloat(val, 'f', -1, 64)
		case time.Time:
			res[k] = val.Format(time.RFC3339)
		case []interface{}:
			res[k] = convertInterfaceArrayToString(val)
		default:
			rv := reflect.ValueOf(val)
			if rv.Kind() >= reflect.Int && rv.Kind() <= reflect.Int64 {
				res[k] = strconv.FormatInt(rv.Int(), 10)
			} else {
				fmt.Printf("skip key %s (type %T)\n", k, val)
			}
		}
	}
	return res
}

// -----------------------------------------------------------------------------
// End of file — drop this into your existing utils package and remove the old
// implementation.  All previous calls to utils.ExtractFrontMatter still work,
// and you now also have ParseFile / ScanDir / SortRows for list pages.
