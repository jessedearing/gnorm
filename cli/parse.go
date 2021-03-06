package cli // import "gnorm.org/gnorm/cli"

import (
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"gnorm.org/gnorm/environ"
	"gnorm.org/gnorm/run"
)

func parseFile(env environ.Values, file string) (*run.Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithMessage(err, "can't open config file")
	}
	defer f.Close()
	return parse(env, f)
}

// parse reads the configuration file and returns a gnorm config value.
func parse(env environ.Values, r io.Reader) (*run.Config, error) {
	c := Config{}
	m, err := toml.DecodeReader(r, &c)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing config file")
	}
	undec := m.Undecoded()
	if len(undec) > 0 {
		log.Println("Warning: unknown values present in config file:", undec)
	}

	if len(c.Schemas) == 0 {
		return nil, errors.New("no schemas specified in config")
	}

	if c.NameConversion == "" {
		return nil, errors.New("no NameConversion specified in config")
	}
	if len(c.ExcludeTables) > 0 && len(c.IncludeTables) > 0 {
		return nil, errors.New("both include tables and exclude tables")
	}

	include, err := parseTables(c.IncludeTables, c.Schemas)
	if err != nil {
		return nil, err
	}

	exclude, err := parseTables(c.ExcludeTables, c.Schemas)
	if err != nil {
		return nil, err
	}

	cfg := &run.Config{
		ConnStr:         c.ConnStr,
		Schemas:         c.Schemas,
		NullableTypeMap: c.NullableTypeMap,
		TypeMap:         c.TypeMap,
		TemplateDir:     c.TemplateDir,
		PostRun:         c.PostRun,
		ExcludeTables:   exclude,
		IncludeTables:   include,
	}

	switch strings.ToLower(c.DBType) {
	case "":
		return nil, errors.New("no DBType specificed")
	case "postgres":
		cfg.DBType = run.Postgres
	case "mysql":
		cfg.DBType = run.Mysql
	default:
		return nil, errors.Errorf("unsupported dbtype %q", c.DBType)
	}

	t, err := template.New("NameConversion").Funcs(environ.FuncMap).Parse(c.NameConversion)
	if err != nil {
		return nil, errors.WithMessage(err, "error parsing NameConversion template")
	}
	cfg.NameConversion = t

	if c.SchemaPath != "" {
		t, err := template.New("SchemaPath").Funcs(environ.FuncMap).Parse(c.SchemaPath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.SchemaPath = t
	}

	if c.TablePath != "" {
		t, err := template.New("TablePath").Funcs(environ.FuncMap).Parse(c.TablePath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.TablePath = t
	}

	if c.EnumPath != "" {
		t, err := template.New("EnumPath").Funcs(environ.FuncMap).Parse(c.EnumPath)
		if err != nil {
			return nil, errors.WithMessage(err, "error parsing SchemaPath template")
		}
		cfg.EnumPath = t
	}

	if cfg.EnumPath == nil && cfg.TablePath == nil && cfg.SchemaPath == nil {
		return nil, errors.New("no output paths defined, so no output will be generated")
	}

	cfg.ConnStr = os.Expand(c.ConnStr, func(s string) string {
		return env.Env[s]
	})
	return cfg, nil
}

// parseTables takes a list of tablenames in "<schema.>table" format and spits
// out a map of schema to list of tables.  Tables with no schema apply to all
// schemas.  Tables with a schema apply to only that schema.  Tables that
// specify a schema not in the list of schemas given are an error.
func parseTables(tables, schemas []string) (map[string][]string, error) {
	out := make(map[string][]string, len(schemas))
	for _, s := range schemas {
		out[s] = nil
	}
	for _, t := range tables {
		vals := strings.Split(t, ".")
		switch len(vals) {
		case 1:
			// just the table name, so it goes for all schemas
			for schema := range out {
				out[schema] = append(out[schema], t)
			}
		case 2:
			// schema and table
			list, ok := out[vals[0]]
			if !ok {
				return nil, errors.Errorf("%q specified for tables but schema %q not in schema list", t, vals[0])
			}
			out[vals[0]] = append(list, vals[1])
		default:
			// too many periods... bad format
			return nil, errors.Errorf(`badly formatted table: %q, should be just "table" or "table.schema"`, t)
		}
	}

	return out, nil
}
