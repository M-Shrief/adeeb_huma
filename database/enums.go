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

		if err := tx.Exec(TimePeriodEnum_CreateSQLIfNotExists).Error; err != nil {
			return err
		}

		if err := tx.Exec(StatusEnum_CreateSQLIfNotExists).Error; err != nil {
			return err
		}

		return nil
	})

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

// Implement the driver.Valuer interface to save to DB
func (r RoleEnum) Value() (driver.Value, error) {
	return string(r), nil
}

// Implement the sql.Scanner interface to read from DB
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

type TimePeriodEnum string

const (
	TimePeriodEnum_Undefined  TimePeriodEnum = "غير محدد"
	TimePeriodEnum_Jahli      TimePeriodEnum = "جاهلي"
	TimePeriodEnum_Amoei      TimePeriodEnum = "أموي"
	TimePeriodEnum_Abasi      TimePeriodEnum = "عباسي"
	TimePeriodEnum_Andalusi   TimePeriodEnum = "أندلسي"
	TimePeriodEnum_TurkishEra TimePeriodEnum = "عثماني ومملوكي"
	TimePeriodEnum_Modern     TimePeriodEnum = "حديث"
)

var TimePeriodEnum_CreateSQLIfNotExists = fmt.Sprintf(
	`DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type 
			WHERE typname = '%s' 
			AND typnamespace = 'public'::regnamespace
		) THEN
			CREATE TYPE %s AS ENUM ('%s', '%s', '%s', '%s', '%s', '%s', '%s');
		END IF;
	END $$;
	`,
	`time_period_enum`, `time_period_enum`, TimePeriodEnum_Undefined,
	TimePeriodEnum_Jahli, TimePeriodEnum_Amoei, TimePeriodEnum_Abasi,
	TimePeriodEnum_Andalusi, TimePeriodEnum_TurkishEra, TimePeriodEnum_Modern,
)

// Implement the driver.Valuer interface to save to DB
func (tp TimePeriodEnum) Value() (driver.Value, error) {
	return string(tp), nil
}

// Implement the sql.Scanner interface to read from DB
func (tp *TimePeriodEnum) Scan(value interface{}) error {
	if value == nil {
		*tp = ""
		return nil
	}
	// Postgres often returns []byte, MySQL might return string
	switch v := value.(type) {
	case []byte:
		*tp = TimePeriodEnum(v)
	case string:
		*tp = TimePeriodEnum(v)
	default:
		return errors.New("invalid type for TimePeriodEnum")
	}
	return nil
}

type StatusEnum string

const (
	StatusEnum_InProgress StatusEnum = "in progress"
	StatusEnum_Aborted    StatusEnum = "aborted"
	StatusEnum_Completed  StatusEnum = "completed"
)

var StatusEnum_CreateSQLIfNotExists = fmt.Sprintf(
	`DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type 
			WHERE typname = '%s' 
			AND typnamespace = 'public'::regnamespace
		) THEN
			CREATE TYPE %s AS ENUM ('%s', '%s', '%s');
		END IF;
	END $$;
	`,
	`status_enum`, `status_enum`, StatusEnum_InProgress, StatusEnum_Aborted, StatusEnum_Completed,
)

// Implement the driver.Valuer interface to save to DB
func (s StatusEnum) Value() (driver.Value, error) {
	return string(s), nil
}

// Implement the sql.Scanner interface to read from DB
func (s *StatusEnum) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	// Postgres often returns []byte, MySQL might return string
	switch v := value.(type) {
	case []byte:
		*s = StatusEnum(v)
	case string:
		*s = StatusEnum(v)
	default:
		return errors.New("invalid type for StatusEnum")
	}
	return nil
}
