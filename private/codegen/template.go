package codegen

import "text/template"

// DefaultTemplate is the template used by default for generating code.
var DefaultTemplate = template.Must(template.New("defaultTemplate").Parse(`// Code generated by "{{.CommandLine}}"; DO NOT EDIT

package {{.Package}}

import ({{range .Imports}}
    {{.}}{{end}}
	"github.com/jjeffery/errors"
)
{{range .QueryTypes -}}
{{- if .Method.Get}}
// Get a {{.Singular}} by its primary key. Returns nil if not found.
func (q {{.TypeName}}) Get({{.RowType.IDParams}}) (*{{.RowType.Name}}, error) {
	var row {{.RowType.Name}}
	n, err := q.schema.Select(q.db, &row, {{.QuotedTableName}}, {{.RowType.IDArgs}})
	if err != nil {
		return nil, errors.Wrap(err, "cannot get {{.Singular}}").With(
            {{.RowType.IDKeyvals}}
		)
	}
	if n == 0 {
		return nil, nil
	}
	return &row, nil
}
{{end -}}
{{- if .Method.Select}}
// Select a list of {{.Plural}} from an SQL query.
func (q {{.TypeName}}) Select(query string, args ...interface{}) ([]*{{.RowType.Name}}, error) {
	var rows []*{{.RowType.Name}}
	_, err := q.schema.Select(q.db, &rows, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query {{.Plural}}").With(
			"query", query,
			"args", args,
		)
	}
	return rows, nil
}
{{end -}}
{{- if .Method.SelectOne}}
// SelectOne selects a {{.Singular}} from an SQL query. Returns nil if the query returns no rows.
// If the query returns one or more rows the value for the first is returned and any subsequent
// rows are discarded.
func (q {{.TypeName}}) SelectOne(query string, args ...interface{}) (*{{.RowType.Name}}, error) {
	var row {{.RowType.Name}}
	n, err := q.schema.Select(q.db, &row, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query one {{.Singular}}").With(
			"query", query,
			"args", args,
		)
	}
	if n == 0 {
		return nil, nil
	}
	return &row, nil
}
{{end -}}
{{- if .Method.Insert}}
// Insert a {{.Singular}} row.
func (q {{.TypeName}}) Insert(row *{{.RowType.Name}}) error {
	err := q.schema.Insert(q.db, row, {{.QuotedTableName}})
	if err != nil {
		return errors.Wrap(err, "cannot insert {{.Singular}}").With(
            {{range .RowType.LogProps}}"{{.}}", row.{{.}}, {{end}}
		)
	}
	return nil
}
{{end -}}
{{- if .Method.Update}}
// Update an existing {{.Singular}} row. Returns the number of rows updated,
// which should be zero or one.
func (q {{.TypeName}}) Update(row *{{.RowType.Name}}) (int, error) {
	n, err := q.schema.Update(q.db, row, {{.QuotedTableName}})
	if err != nil {
		return 0, errors.Wrap(err, "cannot update {{.Singular}}").With(
            {{range .RowType.LogProps}}"{{.}}", row.{{.}}, {{end}}
		)
	}
	return n, nil
}
{{end -}}
{{- if .Method.Upsert}}
// Attempt to update a {{.Singular}} row, and if it does not exist then insert it.
func (q {{.TypeName}}) Upsert(row *{{.RowType.Name}}) error {
	n, err := q.schema.Update(q.db, row, {{.QuotedTableName}})
    if err != nil {
		return errors.Wrap(err, "cannot update {{.Singular}} for upsert").With(
            {{range .RowType.LogProps}}"{{.}}", row.{{.}}, {{end}}
        )
    }
    if n > 0 {
        // update successful, row updated
        return nil
    }
	if err := q.schema.Insert(q.db, row, {{.QuotedTableName}}); err != nil {
		return errors.Wrap(err, "cannot insert {{.Singular}} for upsert").With(
            {{range .RowType.LogProps}}"{{.}}", row.{{.}}, {{end}}
		)
	}
	return nil
}
{{end -}}
{{- if .Method.Delete}}
// Delete a {{.Singular}} row. Returns the number of rows deleted, which should
// be zero or one.
func (q {{.TypeName}}) Delete(row *{{.RowType.Name}}) (int, error) {
	n, err := q.schema.Delete(q.db, row, {{.QuotedTableName}})
	if err != nil {
		return 0, errors.Wrap(err, "cannot delete {{.Singular}}").With(
            {{range .RowType.LogProps}}"{{.}}", row.{{.}}, {{end}}
		)
	}
	return n, nil
}
{{end -}}
{{- end}}`))
