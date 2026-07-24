package mdcoach

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/OneOfOne/xxhash"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
)

// Environment holds normalized execution configuration and utility methods.
type Environment struct {
	Debug       bool
	WorkingPath string
	CachePath   string
	Output      string // Output directory.
	Overwrite   bool
	db          *sql.DB
	dbPutToken  *sql.Stmt
	dbDelToken  *sql.Stmt
	dbPutMeta   *sql.Stmt
}

// SetupDatabase creates
func (e *Environment) setupDatabase() (err error) {
	f := filepath.Join(e.CachePath, `cache-1-1-4.sqlite3`)
	_, exists := os.Stat(f)
	e.db, err = sql.Open("sqlite3", f)
	if err != nil {
		return err
	}
	if exists != nil {
		e.db.Exec(`CREATE TABLE tokens(token varchar(16) PRIMARY KEY)`)
		e.db.Exec(`CREATE TABLE indexed(uid varchar(16) PRIMARY KEY, type varchar(16), view varchar(128), print varchar(128), title varchar(128), author varchar(128), description text)`)
	}
	e.dbPutToken, _ = e.db.Prepare(`INSERT INTO tokens VALUES(?)`)
	e.dbDelToken, _ = e.db.Prepare(`DELETE FROM tokens WHERE token=?`)
	e.dbPutMeta, _ = e.db.Prepare(`INSERT INTO indexed VALUES(?,?,?,?,?,?,?)`)
	return err
}

// NewEnvironment configures a passable environment object.
func NewEnvironment(output string) *Environment {
	e := new(Environment)
	e.Output = output
	e.CachePath = filepath.Join(filepath.Dir(output), `.cache`)
	info, err := os.Stat(e.CachePath)
	if err != nil || !info.IsDir() {
		log.Printf(`Setting up cache at "%s".`, e.CachePath)
		err = os.MkdirAll(e.CachePath, 0700)
		if err != nil {
			log.Fatalf(`Directory "%s" is inaccessible.`, e.CachePath)
		}
		err = IOextractAssets(e.CachePath, ``)
	}
	e.WorkingPath, _ = os.Getwd()
	if err := e.setupDatabase(); err != nil {
		log.Fatal(err)
	}
	return e
}

// Open exposes environment as a vialbe http.FileSystem interface.
func (e *Environment) Open(name string) (http.File, error) {
	return os.Open(filepath.Clean(filepath.Join(e.WorkingPath, name)))
}

// Token provides fast hash from input.
func (e *Environment) Token(diff string) string {
	h := xxhash.New64()
	h.WriteString(diff)
	return fmt.Sprintf("%x", h.Sum64())
}

// EnsureToken returns deterministic cache token and true, if it was just generated.
func (e *Environment) EnsureToken(diff string) (token string, created bool) {
	token = e.Token(diff)
	_, err := e.dbPutToken.Exec(token)
	return token, err == nil
}

// DeleteToken removes a token.
func (e *Environment) DeleteToken(token string) {
	e.dbDelToken.Exec(token)
}

// Create creates a file with overwrite warning.
func (e *Environment) Create(file string, cache bool) (path string, w io.WriteCloser) {
	if cache {
		token, _ := e.EnsureToken(file)
		file = filepath.Clean(filepath.Join(e.CachePath, token+filepath.Ext(file)))
	} else {
		file = filepath.Clean(file)
		if !e.Overwrite {
			if _, err := os.Stat(file); err == nil {
				response := `N`
				fmt.Printf(`Overwrite "%s"? [y/N]      `, file)
				if fmt.Scanln(&response); !strings.HasPrefix(strings.ToLower(response), `y`) {
					log.Fatalf(`Cannot overwrite file "%s"!`, file)
				}
			}
		}
	}
	handle, err := os.Create(file)
	if err != nil {
		log.Fatalf(`Could not create file "%s"!`, file)
	}
	return file, handle
}

// SaveMeta caches data for MakeIndex.
func (e *Environment) SaveMeta(view, print string, meta map[string]interface{}) {
	g := func(key string) string {
		if v, ok := meta[key].(string); ok {
			return v
		}
		return ``
	}
	if token, created := e.EnsureToken(filepath.Base(meta[`source`].(string))); created {
		e.dbPutMeta.Exec(token, g(`type`), view, print,
			g(`title`), g(`author`), g(`description`))
	}
}

// MakeIndex generates a directory index using cached meta data.
func (e *Environment) MakeIndex(path string) error {
	meta := make(map[string]interface{})
	meta[`sections`] = make([]map[string]interface{}, 0)
	var tp, view, print, title, author, description string
	rows, err := e.db.Query(`SELECT type, view, print, title, author, description FROM indexed`)
	for rows.Next() {
		rows.Scan(&tp, &view, &print, &title, &author, &description)
		m := make(map[string]interface{})
		m[`view`] = view
		m[`print`] = print
		m[`title`] = title
		m[`description`] = description
		if tp == `syllabus` {
			meta[`syllabus`] = m
		} else {
			meta[`sections`] = append(meta[`sections`].([]map[string]interface{}), m)
		}
	}
	_, handle := e.Create(filepath.Join(path, `index.html`), false)
	defer handle.Close()
	meta[`stylesheets`] = []string{`.cache/index.css`, `.cache/pygments.css`, `.cache/emote.css`}
	err = template.Must(template.New("head").Parse(tmplHead)).Execute(handle, meta)
	if err != nil {
		return err
	}
	err = template.Must(template.New("index").Parse(tmplIndex)).Execute(handle, meta)
	if err != nil {
		return err
	}
	io.WriteString(handle, tmplFoot)
	return nil
}

// Loader connects document parser to a file system using a closure.
func (e *Environment) Loader(metaCallback func(map[string]interface{})) func(string) ([]byte, error) {
	recursionGuard := 0
	return func(file string) ([]byte, error) {
		recursionGuard++
		if recursionGuard > 256 {
			return []byte{}, fmt.Errorf(`maximum recursion level of 256 reached at "%s"`, file)
		}
		result := make(map[string]interface{})
		index, buf := 0, bytes.NewBuffer(nil)
		handle, err := os.Open(filepath.Clean(filepath.Join(e.WorkingPath, file)))
		if err != nil {
			return []byte{}, err
		}
		_, err = io.Copy(buf, handle)
		if err != nil {
			return []byte{}, err
		}
		handle.Close()
		b := buf.Bytes()
		if len(b) > 7 && bytes.Equal(b[:4], []byte("---\n")) {
			if i := bytes.Index(b, slideBoundary); i > -1 {
				index = i + len(slideBoundary)
			}
		}
		if index > 0 {
			result = ParseFrontMatter(b[:index])
		}
		result[`source`] = file
		metaCallback(result)
		return b[index:], nil
	}
}

// Close saves token map to cache.
func (e *Environment) Close() {
	e.db.Close()
}
