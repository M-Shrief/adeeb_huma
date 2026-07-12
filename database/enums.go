package database

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CreateEnumsIfNotExists() {

	_ = Conn.Transaction(func(tx *gorm.DB) error {
		// if tx.Statement.ConnPool == nil {
		// 	logger.Error().Msg("transaction connection pool is nil")
		// 	return errors.New("transaction connection pool is nil")
		// }
		if err := tx.Exec(RoleEnum_CreateSQLIfNotExists).Error; err != nil {
			return err
		}
		return nil
	})

	// if err != nil {
	// 	logger.Error().Err(err).Msg("Transaction failed")
	// }
}

type RoleEnum string

const (
	RoleEnum_Normal     RoleEnum = "Normal"
	RoleEnum_Banned     RoleEnum = "Banned"
	RoleEnum_Management RoleEnum = "Management"
	RoleEnum_DBA        RoleEnum = "DBA"
	RoleEnum_Analytics  RoleEnum = "Analytics"
)

var RoleEnum_CreateSQLIfNotExists = fmt.Sprintf(
	`DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type 
			WHERE typname = '%s' 
			AND typnamespace = 'public'::regnamespace
		) THEN
			CREATE TYPE %s AS ENUM ('%s', '%s', '%s', '%s', '%s');
		END IF;
	END $$;
	`,
	`role_enum`, `role_enum`, RoleEnum_Normal, RoleEnum_Banned, RoleEnum_Management, RoleEnum_DBA, RoleEnum_Analytics,
)

// 2. Implement the driver.Valuer interface to save to DB
func (r RoleEnum) Value() (driver.Value, error) {
	return string(r), nil
}

// 3. Implement the sql.Scanner interface to read from DB
func (r *RoleEnum) Scan(value interface{}) error {
	if value == nil {
		*r = ""
		return nil
	}
	// Postgres often returns []byte, MySQL might return string
	switch v := value.(type) {
	case []byte:
		*r = RoleEnum(v)
	case string:
		*r = RoleEnum(v)
	default:
		return errors.New("invalid type for RoleEnum")
	}
	return nil
}
