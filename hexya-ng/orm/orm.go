package orm

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

// Init establishes a connection to the database.
func Init(driverName, dataSourceName string) {
	var err error
	DB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// We can check the connection here with DB.Ping()
}

// Close closes the database connection.
func Close() {
	DB.Close()
}

// RecordSetI is the interface for a RecordSet.
type RecordSetI interface {
	Create(data map[string]interface{}) (map[string]interface{}, error)
	Read() ([]map[string]interface{}, error)
	ReadOne(id int) (map[string]interface{}, error)
	Write(id int, data map[string]interface{}) error
	Unlink(id int) error
	CreateTable() error
}

// RecordSet represents a set of records of a model.
type RecordSet struct {
	model *ModelInfo
}

// NewRecordSet creates a new RecordSet for the given model.
func NewRecordSet(modelName string) RecordSetI {
	return &RecordSet{
		model: GetModel(modelName),
	}
}

// Create creates a new record in the database.
func (rs *RecordSet) Create(data map[string]interface{}) (map[string]interface{}, error) {
	var columns []string
	for name := range data {
		columns = append(columns, name)
	}
	sort.Strings(columns)

	var values []interface{}
	var placeholders []string
	for _, col := range columns {
		values = append(values, data[col])
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		rs.model.Name,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := DB.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	// This is a simplification. A real implementation would return the created record with its ID.
	return data, nil
}

// Read reads records from the database.
// A proper implementation would take a query condition.
func (rs *RecordSet) Read() ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", rs.model.Name)
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		results = append(results, m)
	}

	return results, nil
}

// ReadOne reads a single record from the database.
func (rs *RecordSet) ReadOne(id int) (map[string]interface{}, error) {
	var columns []string
	for name := range rs.model.Fields {
		columns = append(columns, name)
	}
	sort.Strings(columns)
	query := fmt.Sprintf("SELECT id, %s FROM %s WHERE id = ?", strings.Join(columns, ", "), rs.model.Name)
	row := DB.QueryRow(query, id)

	values := make([]interface{}, len(columns)+1)
	valuePtrs := make([]interface{}, len(columns)+1)
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := row.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["id"] = values[0]
	for i, colName := range columns {
		m[colName] = values[i+1]
	}

	return m, nil
}

// Write updates a record in the database.
func (rs *RecordSet) Write(id int, data map[string]interface{}) error {
	var columns []string
	for name := range data {
		columns = append(columns, name)
	}
	sort.Strings(columns)

	var setClauses []string
	var values []interface{}
	for _, col := range columns {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
		values = append(values, data[col])
	}
	values = append(values, id)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		rs.model.Name,
		strings.Join(setClauses, ", "),
	)

	_, err := DB.Exec(query, values...)
	return err
}

// Unlink deletes a record from the database.
func (rs *RecordSet) Unlink(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", rs.model.Name)
	_, err := DB.Exec(query, id)
	return err
}

// CreateTable creates the table for the model if it doesn't exist.
// This is a simplification for testing purposes.
func (rs *RecordSet) CreateTable() error {
	var columns []string
	for name, field := range rs.model.Fields {
		// This is a very basic type mapping.
		var sqlType string
		switch field.Type {
		case "Char":
			sqlType = "TEXT"
		case "Integer":
			sqlType = "INTEGER"
		default:
			sqlType = "TEXT"
		}
		columns = append(columns, fmt.Sprintf("%s %s", name, sqlType))
	}

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, %s)",
		rs.model.Name,
		strings.Join(columns, ", "),
	)

	_, err := DB.Exec(query)
	return err
}
