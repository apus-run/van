package gormutil

import (
	"github.com/go-sql-driver/mysql"
)

const (
	errMySQLDuplicatedRecord          uint16 = 1062
	errMySQLForeignKeyConstraintFails uint16 = 1452
)

// IsMySQLDuplicatedRecordErr checks if the error is a MySQL duplicate record error.
func IsMySQLDuplicatedRecordErr(err error) bool {
	mErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return mErr.Number == errMySQLDuplicatedRecord
}

// IsMySQLForeignKeyConstraintFailsError checks if the error is a MySQL foreign key constraint fails error.
func IsMySQLForeignKeyConstraintFailsError(err error) bool {
	mErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return mErr.Number == errMySQLForeignKeyConstraintFails
}
