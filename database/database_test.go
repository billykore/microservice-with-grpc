package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDatabaseConnection(t *testing.T) {
	db, err := New(&Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestDatabaseMigration(t *testing.T) {
	db, err := New(&Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	type TestTable1 struct {
		gorm.Model
		Foo string
		Bar string
	}
	type TestTable2 struct {
		gorm.Model
		Foo string
		Bar string
	}
	type TestTable3 struct {
		gorm.Model
		Foo string
		Bar string
	}
	err = Migrate(db, &TestTable1{}, &TestTable2{}, &TestTable3{})
	assert.NoError(t, err)

	// drop tables after migration test
	err = db.Migrator().DropTable("test_table1", "test_table2", "test_table3")
	assert.NoError(t, err)
}

func TestDropTestTable(t *testing.T) {
	db, err := New(&Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = db.Migrator().DropTable("customers")
	if err != nil {
		panic(err)
	}
}
