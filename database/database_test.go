package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMySQLDatabaseConnection(t *testing.T) {
	db := New(MySQL, &Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
	assert.NotNil(t, db)
}

func TestMySQLDatabaseMigration(t *testing.T) {
	db := New(MySQL, &Config{
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseHost:     "localhost",
		DatabasePort:     "3306",
		DatabaseName:     "grpc_microservices",
	})
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
	err := Migrate(db, &TestTable1{}, &TestTable2{}, &TestTable3{})
	assert.NoError(t, err)

	// drop tables after migration test
	err = db.Migrator().DropTable("test_table1", "test_table2", "test_table3")
	assert.NoError(t, err)
}

func TestPostgresDatabaseConnection(t *testing.T) {
	db := New(Postgres, &Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
	assert.NotNil(t, db)
}

func TestPostgresDatabaseMigration(t *testing.T) {
	db := New(Postgres, &Config{
		DatabaseUser:     "postgres",
		DatabasePassword: "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     "5432",
		DatabaseName:     "grpc_auth_service",
	})
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
	err := Migrate(db, &TestTable1{}, &TestTable2{}, &TestTable3{})
	assert.NoError(t, err)

	// drop tables after migration test
	err = db.Migrator().DropTable("test_table1", "test_table2", "test_table3")
	assert.NoError(t, err)
}
