package schemas

import (
	openapi_v2 "github.com/googleapis/gnostic/openapiv2"
	"github.com/pkg/errors"
	"github.com/stackrox/helmtest/internal/rox-imported/gziputil"
	"io/fs"
	"path"
	"strings"
	"sync"
)

// Registry allows locating OpenAPI schemas by their name.
type Registry struct {
	mutex sync.Mutex
	schemas map[string]*schemaEntry
}

func NewRegistry(loadBuiltin bool) *Registry {
	r := &Registry{
		schemas: make(map[string]*schemaEntry),
	}
	return r
}

func (r *Registry) GetSchema(schemaName string) (*schema, error) {

}

type schemaEntry struct {
	name     string
	schema   *schema
	loadErr  error
	loadOnce sync.Once
}

func (e *schemaEntry) get() (*schema, error) {
	e.loadOnce.Do(func() {
		e.schema, e.loadErr = e.load()
	})
	return e.schema, e.loadErr
}

func (e *schemaEntry) load() (*schema, error) {
	schemaBytes, err := fs.ReadFile(openAPISchemaFS, path.Join("openapi-schemas", e.name+".json.gz"))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid name %q", e.name)
	}

	openapiDoc, err := gziputil.Decompress(schemaBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading openapi docs %s", e.name)
	}

	doc, err := openapi_v2.ParseDocument(openapiDoc)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing OpenAPI doc for %s", e.name)
	}
	schema, err := newSchema(doc)
	if err != nil {
		return nil, errors.Wrapf(err, "creating OpenAPI schema from document for %s", e.name)
	}

	return schema, nil
}

func getSchemaEntry(name string) *schemaEntry {
	name = strings.ToLower(name)

	allSchemasMutex.Lock()
	defer allSchemasMutex.Unlock()

	entry := allSchemas[name]
	if entry != nil {
		return entry
	}
	entry = &schemaEntry{
		name: name,
	}
	allSchemas[name] = entry
	return entry
}