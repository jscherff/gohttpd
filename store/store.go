// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	`fmt`
	`strings`
)

// DataStore is an interface that represents a data store.
type DataStore interface {
	Register(schemaName string)
	Prepare(queryFile string) (err error)
	Select(queryName string, dest, arg interface{}) (err error)
	Insert(queryName string, arg interface{}) (id int64, err error)
	Exec(queryName string, arg interface{}) (rows int64, err error)
	Get(queryName string, dest, arg interface{}) (err error)
	String() (info string)
	Close()
}

// New creates a new DataStore instance using only the the provided driver
// name and configuration file/string using the registered factory method
// associated with the driver name from the factories registry.
func New(driver, config string) (DataStore, error) {
	if factory, ok := factories[driver]; !ok {
		return nil, fmt.Errorf(`driver %q not found`, driver)
	} else {
		return factory(config)
	}
}

// query contains SQL Xquery components needed for building prepared statements.
type query struct {
	Table string
	Command string
	Filters []string
	Columns []string
	sqlStmt string
}

// table returns the lowercase name of the table.
func (this *query) table() (string) {
	return strings.ToLower(this.Table)
}

// command returns the uppercase SQL command.
func (this *query) command() (string) {
	return strings.ToUpper(this.Command)
}

// columns is a list of column names for INSERT, SELECT, and UPDATE
// statements.
func (this *query) columns() (string) {
	return strings.Join(this.Columns, `, `)
}

// filters is a list of columns used in the conditions clause of a SQL
// statement. The interface currently only supportes ANDed conditions.
func (this *query) filters() (string) {
	var filters []string
	for _, column := range this.Filters {
		filters = append(filters, fmt.Sprintf(`%[1]v = :%[1]v`, column))
	}
	return strings.Join(filters, ` AND `)
}

// params is a list of named parameters for INSERT statements.
func (this *query) params() (string) {
	var params []string
	for _, column := range this.Columns {
		if column == `*` {
			continue
		}
		params = append(params, fmt.Sprintf(`:%v`, column))
	}
	return strings.Join(params, `, `)
}

// setters is a list of column assignments for UPDATE statements.
func (this *query) setters() (string) {
	var setters []string
	for _, column := range this.Columns {
		if column == `*` {
			continue
		}
		setters = append(setters, fmt.Sprintf(`%[1]v = :%[1]v`, column))
	}
	return strings.Join(setters, `, `)
}

// String implements the Stringer interface for the Query object and returns
// the complete SQL statement string assembled from the statement components.
func (this *query) String() (string) {

	if this.sqlStmt != `` {
		return this.sqlStmt
	}

	if this.table() == `` || this.command() == `` {
		return ``
	}

	switch this.command() {

	case `INSERT`, `REPLACE`:
		this.sqlStmt = fmt.Sprintf(`%s INTO %s (%s) VALUES (%s)`,
			this.command(),
			this.table(),
			this.columns(),
			this.params(),
		)

	case `SELECT`:
		this.sqlStmt = fmt.Sprintf(`%s %s FROM %s`,
			this.command(),
			this.columns(),
			this.table(),
		)

	case `UPDATE`:
		this.sqlStmt = fmt.Sprintf(`%s %s SET %s`,
			this.command(),
			this.table(),
			this.setters(),
		)

	case `DELETE`:
		this.sqlStmt = fmt.Sprintf(`DELETE FROM %s`,
			this.table(),
		)

	default:
		return ``
	}

	if len(this.Filters) > 0 {
		this.sqlStmt += fmt.Sprintf(` WHERE %s`,
			this.filters(),
		)
	}

	return this.sqlStmt
}

type queries struct {
	Driver string
	Schema string
	Query map[string]*query
}

type Queries interface {
	Build() 