// Code generated by "sqlr-gen"; DO NOT EDIT

package testdata

import (
	"github.com/jjeffery/errors"
)

// get retrieves a Document by its primary key. Returns nil if not found.
func (q *DocumentQuery) get(id string) (*Document, error) {
	var row Document
	n, err := q.schema.Select(q.db, &row, "documents", id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get Document").With(
			"id", id,
		)
	}
	if n == 0 {
		return nil, nil
	}
	return &row, nil
}

// selectRows returns a list of Documents from an SQL query.
func (q *DocumentQuery) selectRows(query string, args ...interface{}) ([]*Document, error) {
	var rows []*Document
	_, err := q.schema.Select(q.db, &rows, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query Documents").With(
			"query", query,
			"args", args,
		)
	}
	return rows, nil
}

// selectRow selects a Document from an SQL query. Returns nil if the query returns no rows.
// If the query returns one or more rows the value for the first is returned and any subsequent
// rows are discarded.
func (q *DocumentQuery) selectRow(query string, args ...interface{}) (*Document, error) {
	var row Document
	n, err := q.schema.Select(q.db, &row, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "cannot query one Document").With(
			"query", query,
			"args", args,
		)
	}
	if n == 0 {
		return nil, nil
	}
	return &row, nil
}

// insert inserts a Document row.
func (q *DocumentQuery) insert(row *Document) error {
	_, err := q.schema.Exec(q.db, row, "insert into documents({}) values({})")
	if err != nil {
		return errors.Wrap(err, "cannot insert Document").With(
			"ID", row.ID,
		)
	}
	return nil
}

// update updates an existing Document row. Returns the number of rows updated,
// which should be zero or one.
func (q *DocumentQuery) update(row *Document) (int, error) {
	n, err := q.schema.Exec(q.db, row, "update documents set {} where {}")
	if err != nil {
		return 0, errors.Wrap(err, "cannot update Document").With(
			"ID", row.ID,
		)
	}
	return n, nil
}

// upsert attempts to update a Document row, and if it does not exist then insert it.
func (q *DocumentQuery) upsert(row *Document) error {
	n, err := q.schema.Exec(q.db, row, "update documents set {} where {}")
	if err != nil {
		return errors.Wrap(err, "cannot update Document for upsert").With(
			"ID", row.ID,
		)
	}
	if n > 0 {
		// update successful, row updated
		return nil
	}
	if _, err := q.schema.Exec(q.db, row, "insert into documents({}) values({})"); err != nil {
		return errors.Wrap(err, "cannot insert Document for upsert").With(
			"ID", row.ID,
		)
	}
	return nil
}

// delete deletes a Document row. Returns the number of rows deleted, which should
// be zero or one.
func (q *DocumentQuery) delete(row *Document) (int, error) {
	n, err := q.schema.Exec(q.db, row, "delete from documents where {}")
	if err != nil {
		return 0, errors.Wrap(err, "cannot delete Document").With(
			"ID", row.ID,
		)
	}
	return n, nil
}
