package orm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB

// Init establishes a connection to the database.
func Init(driverName, dataSourceName string) {
	var err error
	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// We can check the connection here with db.Ping()
}

// Close closes the database connection.
func Close() {
	db.Close()
}

// RecordSet represents a set of records of a model.
type RecordSet struct {
	model *ModelInfo
}

// NewRecordSet creates a new RecordSet for the given model.
func NewRecordSet(modelName string) *RecordSet {
	return &RecordSet{
		model: GetModel(modelName),
	}
}

// Create creates a new record in the database.
func (rs *RecordSet) Create(data map[string]interface{}) (map[string]interface{}, error) {
	var columns []string
	var values []interface{}
	var placeholders []string
	for name, value := range data {
		columns = append(columns, name)
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		rs.model.Name,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	_, err := db.Exec(query, values...)
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
	rows, err := db.Query(query)
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

// Write updates records in the database.
// A proper implementation would take a query condition.
func (rs *RecordSet) Write(data map[string]interface{}) error {
	// To be implemented
	return nil
}

// Unlink deletes records from the database.
// A proper implementation would take a query condition.
func (rs *RecordSet) Unlink() error {
	// To be implemented
	return nil
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

	_, err := db.Exec(query)
	return err
}
