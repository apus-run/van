package gormutil

import "gorm.io/gorm"

// RecordExists checks if a record exists in the database.
func RecordExists(db *gorm.DB, where interface{}, tableName ...string) (exists bool, err error) {
	if len(tableName) > 0 {
		db = db.Table(tableName[0])
	} else {
		db = db.Model(where)
	}
	return Exists(db.Where(where))
}

// Exists checks if any record exists in the database.
func Exists(db *gorm.DB) (exists bool, err error) {
	n, err := Count(db.Limit(1))
	return n > 0, err
}

// Count counts the number of records in the database.
func Count(db *gorm.DB) (n int64, err error) {
	return n, db.Count(&n).Error
}

// LimitAndOffset applies limit and offset to the database query.
func LimitAndOffset(limit, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			db = db.Offset(offset)
		}
		if limit > 0 {
			db = db.Limit(limit)
		}
		return db
	}
}
