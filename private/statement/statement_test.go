package statement_test

import (
	"bytes"
	"database/sql/driver"
	"testing"

	"github.com/jjeffery/sqlrow/private/column"
	"github.com/jjeffery/sqlrow/private/dialect"
	"github.com/jjeffery/sqlrow/private/naming"
	"github.com/jjeffery/sqlrow/private/statement"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestStatementExec(t *testing.T) {
	tests := []struct {
		row          interface{}
		query        string
		sql          string
		dialect      dialect.Dialect
		namer        *column.Namer
		args         []driver.Value
		rowsAffected int64
		lastInsertId int64
	}{
		{
			row: struct {
				ID   int
				Name string
			}{
				ID:   1,
				Name: "xxx",
			},
			dialect:      dialect.For("mysql"),
			namer:        column.NewNamer(naming.Snake),
			query:        "insert into table1({}) values({})",
			sql:          "insert into table1(`id`,`name`) values(?,?)",
			args:         []driver.Value{1, "xxx"},
			rowsAffected: 1,
		},
		{
			row: struct {
				ID       int    `sql:"primary key"`
				Name     string `snake:"the_name"`
				OtherCol int
			}{
				ID:       2,
				Name:     "yy",
				OtherCol: 1,
			},
			dialect:      dialect.For("postgres"),
			namer:        column.NewNamer(naming.Snake),
			query:        "update table1 set {} where {}",
			sql:          `update table1 set "the_name"=$1,"other_col"=$2 where "id"=$3`,
			args:         []driver.Value{"yy", 1, 2},
			rowsAffected: 1,
		},
	}

	for i, tt := range tests {
		// func so that we can defer each loop iteration
		func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()

			mock.ExpectExec(toRE(tt.sql)).
				WithArgs(tt.args...).
				WillReturnResult(sqlmock.NewResult(tt.lastInsertId, tt.rowsAffected))

			stmt, err := statement.Prepare(tt.row, tt.query)
			if err != nil {
				t.Errorf("%d: error=%v", i, err)
				return
			}

			rowCount, err := stmt.Exec(db, tt.dialect, tt.namer, tt.row)
			if err != nil {
				t.Errorf("%d: error=%v", i, err)
				return
			}
			if want, got := int(tt.rowsAffected), rowCount; want != got {
				t.Errorf("%d: want=%d got=%d", i, want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		}()
	}
}

// toRe converts a string to a regular expression.
// The sqlmock uses REs, but we want to check the exact SQL.
func toRE(s string) string {
	var buf bytes.Buffer
	for _, ch := range s {
		switch ch {
		case '?', '(', ')', '\\', '.', '+', '$', '^':
			buf.WriteRune('\\')
		}
		buf.WriteRune(ch)
	}
	return buf.String()
}