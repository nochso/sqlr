Introduction
============

Package `sqlr` is designed to reduce the effort required to implement
common operations performed with SQL databases. It is intended for programmers
who are comfortable with writing SQL, but would like assistance with the
sometimes tedious process of preparing SQL queries for tables that have a
large number of columns, or have a variable number of input parameters.

This package is designed to work seamlessly with the standard library
"database/sql" package. It does not provide any layer on top of `*sql.DB`
or `*sql.Tx`. If the calling program has a need to execute queries independently
of this package, it can use "database/sql" directly, or make use of any other
third party package that uses "database/sql".

The following features are provided to simplify writing SQL database queries:

Prepare SQL from row structures 
    It is error-prone and tedious to write SQL queries with long lists of column
    names and placeholder variables. Package sqlr provides a way to map column lists
    in the SQL with a Go language structure. The result is a way to write concise 
    SQL queries even for tables with large number of columns::

        schema.Exec(db, row, "insert into table_name({}) values({})")

Autoincrement column values
    When inserting rows, if a column is defined as an autoincrement column, then the 
    generated value will be retrieved from the database server, and the corresponding 
    field in the row structure will be updated::

        type Row struct {
            Id         int     `sql:"primary key autoincrement"`
            GivenName  string
            FamilyName string
        }
    
Null columns
    Package sqlr provides a convenient mechanism to map NULL values in the database to
    the equivalent empty value in the Go struct field, for example mapping NULL to zero
    for integer values, or NULL to the empty string for string values. This is not always
    applicable as NULL and the empty string are not necessarily the same thing, but in many
    cases there is no ambiguity, and the result is simpler, smaller code::

        type Row struct {
            // ... lots of other fields and then ...

            PhoneNumber string  `sql:"null"` // stores NULL for empty string
            FaxNumber   string  `sql:"null"`
            AddressID   int     `sql:"null"` // stores NULL for zero
        }

JSON columns
    A convenient mechanism for marshaling complex objects as JSON text for storage in 
    an SQL database is supported::

        type Row struct {
            // ... lots of other fields and then ...

            Cmplx  *MyComplexDataStructure `sql:"json"` // will be stored as JSON text
        }

WHERE IN Clauses with multiple values
    When an SQL query contains a `WHERE IN` clause, it usually requires additional string 
    manipulation to match the number of placeholders in the query with args. 
    This package simplifies queries with a variable number of arguments: when processing
    an SQL query, it detects if any of the arguments are slices and adjusts the query
    accordingly::

        id := []int {1, 3, 5, 7, 9}
        _, err := schema.Select(db, &rows, `select {} from widgets where id in (?)`, ids)
    
    
Code generation
    The `sqlr-gen` command is provided to help with generating repetitive code.


Installing
----------

To obtain a copy of the `sqlr` package, use ``go get``:

.. code-block:: sh

    go get github.com/jjeffery/sqlr

Note that additional setup is required if you wish to run the tests
against database servers. The setup required is discussed in :ref:`tests`.

Source Code
-----------

The source code is available at https://github.com/jjeffery/sqlr.
