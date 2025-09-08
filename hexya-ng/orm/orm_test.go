package orm

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	DB = db

	NewModel("Course")
	GetModel("Course").AddFields(map[string]*FieldInfo{
		"Name":        {Name: "Name", Type: "Char"},
		"Description": {Name: "Description", Type: "Char"},
	})

	rs := NewRecordSet("Course")
	courseData := map[string]interface{}{
		"Name":        "Test Course",
		"Description": "This is a test course",
	}

	mock.ExpectExec("INSERT INTO Course \\(Description, Name\\) VALUES \\(\\?, \\?\\)").
		WithArgs("This is a test course", "Test Course").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = rs.Create(courseData)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUnlink(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	DB = db

	NewModel("Order")
	GetModel("Order").AddFields(map[string]*FieldInfo{
		"Number": {Name: "Number", Type: "Char"},
	})

	rs := NewRecordSet("Order")

	mock.ExpectExec("DELETE FROM Order WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = rs.Unlink(1)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	DB = db

	NewModel("User")
	GetModel("User").AddFields(map[string]*FieldInfo{
		"Name": {Name: "Name", Type: "Char"},
		"Age":  {Name: "Age", Type: "Integer"},
	})

	rs := NewRecordSet("User")

	rows := sqlmock.NewRows([]string{"id", "Name", "Age"}).
		AddRow(1, "John", 30).
		AddRow(2, "Jane", 25)

	mock.ExpectQuery("SELECT \\* FROM User").WillReturnRows(rows)

	data, err := rs.Read()
	require.NoError(t, err)

	expectedData := []map[string]interface{}{
		{"id": int64(1), "Name": "John", "Age": int64(30)},
		{"id": int64(2), "Name": "Jane", "Age": int64(25)},
	}

	assert.Equal(t, expectedData, data)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestReadOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	DB = db

	NewModel("Product")
	GetModel("Product").AddFields(map[string]*FieldInfo{
		"Name":  {Name: "Name", Type: "Char"},
		"Price": {Name: "Price", Type: "Integer"},
	})

	rs := NewRecordSet("Product")

	rows := sqlmock.NewRows([]string{"id", "Name", "Price"}).AddRow(1, "Laptop", 1200)

	mock.ExpectQuery("SELECT id, Name, Price FROM Product WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	data, err := rs.ReadOne(1)
	require.NoError(t, err)

	expectedData := map[string]interface{}{"id": int64(1), "Name": "Laptop", "Price": int64(1200)}
	assert.Equal(t, expectedData, data)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestWrite(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	DB = db

	NewModel("Customer")
	GetModel("Customer").AddFields(map[string]*FieldInfo{
		"Name":    {Name: "Name", Type: "Char"},
		"Address": {Name: "Address", Type: "Char"},
	})

	rs := NewRecordSet("Customer")
	updateData := map[string]interface{}{
		"Address": "123 Main St",
	}

	mock.ExpectExec("UPDATE Customer SET Address = \\? WHERE id = \\?").
		WithArgs("123 Main St", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = rs.Write(1, updateData)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
