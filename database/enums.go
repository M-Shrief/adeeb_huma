package database

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

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

		if err := tx.Exec(OutfitTypeEnum_CreateSQLIfNotExists).Error; err != nil {
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

type RolesType []RoleEnum

func (rs *RolesType) Scan(value interface{}) error {
	if value == nil {
		*rs = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Roles: %v", value)
	}
	str := string(bytes)
	if str == "{}" {
		*rs = RolesType{}
		return nil
	}
	// Remove surrounding braces and split
	trimmed := strings.Trim(str, "{}")
	parts := strings.Split(trimmed, ",")
	result := make(RolesType, 0, len(parts))
	for _, p := range parts {
		result = append(result, RoleEnum(p))
	}
	*rs = result
	return nil
}

func (rs RolesType) Value() (driver.Value, error) {
	if len(rs) == 0 {
		return "{}", nil
	}
	var sb strings.Builder
	sb.WriteString("{")
	for i, r := range rs {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(string(r))
	}
	sb.WriteString("}")
	return sb.String(), nil
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

type OutfitTypeEnum string

const (
	OutfitTypeEnum_Tshirt7    OutfitTypeEnum = "تيشيرت - لياقة 7"
	OutfitTypeEnum_TshirtHalf OutfitTypeEnum = "تيشيرت - نص لياقة "
	OutfitTypeEnum_TshirtPolo OutfitTypeEnum = "تشيرت - لياقة بولو"
	OutfitTypeEnum_Jacket     OutfitTypeEnum = "جاكيت"
	OutfitTypeEnum_Sweetshirt OutfitTypeEnum = "سويت شيرت"
	OutfitTypeEnum_Pullover   OutfitTypeEnum = "بلوفر"
)

var OutfitTypeEnum_CreateSQLIfNotExists = fmt.Sprintf(
	`DO $$ BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_type 
			WHERE typname = '%s' 
			AND typnamespace = 'public'::regnamespace
		) THEN
			CREATE TYPE %s AS ENUM ('%s', '%s', '%s', '%s', '%s', '%s');
		END IF;
	END $$;
	`,
	`outfit_type_enum`, `outfit_type_enum`,
	OutfitTypeEnum_Tshirt7, OutfitTypeEnum_TshirtHalf,
	OutfitTypeEnum_TshirtPolo, OutfitTypeEnum_Sweetshirt, OutfitTypeEnum_Jacket, OutfitTypeEnum_Pullover,
)

// Implement the driver.Valuer interface to save to DB
func (ot OutfitTypeEnum) Value() (driver.Value, error) {
	return string(ot), nil
}

// Implement the sql.Scanner interface to read from DB
func (ot *OutfitTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*ot = ""
		return nil
	}
	// Postgres often returns []byte, MySQL might return string
	switch v := value.(type) {
	case []byte:
		*ot = OutfitTypeEnum(v)
	case string:
		*ot = OutfitTypeEnum(v)
	default:
		return errors.New("invalid type for OutfitTypeEnum")
	}
	return nil
}
