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

func TestPrepare(t *testing.T) {
	tests := []struct {
		row     interface{}
		sql     string
		queries map[string]string
	}{
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "insert into tbl({}) values({})",
			queries: map[string]string{
				"mysql":    "insert into tbl(`name`) values(?)",
				"postgres": `insert into tbl("name") values($1)`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "insert into tbl({all}) values({})",
			queries: map[string]string{
				"mysql":    "insert into tbl(`id`,`name`) values(?,?)",
				"postgres": `insert into tbl("id","name") values($1,$2)`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key"`
				Name string
			}{},
			sql: "insert into tbl({}) values({})",
			queries: map[string]string{
				"mysql":    "insert into tbl(`id`,`name`) values(?,?)",
				"postgres": `insert into tbl("id","name") values($1,$2)`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "update tbl set {} where {}",
			queries: map[string]string{
				"mysql":    "update tbl set `name`=? where `id`=?",
				"postgres": `update tbl set "name"=$1 where "id"=$2`,
			},
		},
		{
			row: struct {
				ID    string `sql:"primary key auto increment"`
				Hash  string `sql:"pk"`
				Name  string
				Count int
			}{},
			sql: "update [xxx]\nset\n{}\nwhere {}",
			queries: map[string]string{
				"mysql":    "update `xxx` set `name`=?,`count`=? where `id`=? and `hash`=?",
				"postgres": `update "xxx" set "name"=$1,"count"=$2 where "id"=$3 and "hash"=$4`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "delete from tbl where {}",
			queries: map[string]string{
				"mysql":    "delete from tbl where `id`=?",
				"postgres": `delete from tbl where "id"=$1`,
			},
		},
		{
			row: struct {
				ID    string `sql:"primary key auto increment"`
				Hash  string `sql:"pk"`
				Name  string
				Count int
			}{},
			sql: "delete from `xxx`\n-- this is a comment\nwhere {}",
			queries: map[string]string{
				"mysql":    "delete from `xxx` where `id`=? and `hash`=?",
				"postgres": `delete from "xxx" where "id"=$1 and "hash"=$2`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "select {} from tbl where {}",
			queries: map[string]string{
				"mysql":    "select `id`,`name` from tbl where `id`=?",
				"postgres": `select "id","name" from tbl where "id"=$1`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Name string
			}{},
			sql: "select {alias t} from tbl t where {pk,alias t}",
			queries: map[string]string{
				"mysql":    "select t.`id`,t.`name` from tbl t where t.`id`=?",
				"postgres": `select t."id",t."name" from tbl t where t."id"=$1`,
			},
		},
		{
			row: struct {
				ID   string `sql:"primary key auto increment"`
				Home struct {
					Postcode string
				}
			}{},
			sql: "select {alias t} from tbl t where {pk,alias t}",
			queries: map[string]string{
				"mysql":    "select t.`id`,t.`home_postcode` from tbl t where t.`id`=?",
				"postgres": `select t."id",t."home_postcode" from tbl t where t."id"=$1`,
			},
		},
		{
			row: struct {
				ID    string `sql:"primary key auto increment"`
				Hash  string `sql:"pk"`
				Name  string
				Count int
			}{},
			sql: "select {} from `xxx`\nwhere {}",
			queries: map[string]string{
				"mysql":    "select `id`,`hash`,`name`,`count` from `xxx` where `id`=? and `hash`=?",
				"postgres": `select "id","hash","name","count" from "xxx" where "id"=$1 and "hash"=$2`,
			},
		},
	}

	for i, tt := range tests {
		for dialectName, query := range tt.queries {
			dia := dialect.For(dialectName)
			namer := column.NewNamer(naming.Snake)
			stmt, err := statement.Prepare(tt.row, tt.sql)
			if err != nil {
				t.Errorf("%d: expected no error: got %v", i, err)
				continue
			}
			if got, want := stmt.SQLString(dia, namer), query; got != want {
				t.Errorf("%d: %s: expected=%q, actual=%q", i, dia.Name(), want, got)
			}
		}
	}
}

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

			rowCount, err := stmt.Exec(db, tt.dialect, tt.namer, tt.row, nil)
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
