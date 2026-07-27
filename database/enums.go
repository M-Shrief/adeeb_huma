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
		*rs = []RoleEnum{}
		return nil
	}

	var rawString string

	// 1. Handle different input types from the driver
	switch v := value.(type) {
	case string:
		rawString = v
	case []byte:
		rawString = string(v)
	default:
		return fmt.Errorf("failed to scan Roles, expected string or []byte, got %T", value)
	}

	// 2. Handle empty cases
	if rawString == "" || rawString == "{}" {
		*rs = []RoleEnum{}
		return nil
	}

	// 3. Parse PostgreSQL Array Literal format: "{val1,val2}"
	// Remove surrounding curly braces
	if len(rawString) >= 2 && rawString[0] == '{' && rawString[len(rawString)-1] == '}' {
		rawString = rawString[1 : len(rawString)-1]
	}
	// else {
	// // Fallback if it's just a CSV without braces
	// return fmt.Errorf("unexpected array format: %s", rawString)
	// }

	// 4. Split and Validate
	// Note: This simple split fails if roles contain commas or escaped quotes.
	// For robust parsing of quoted strings, use lib/pq's internal parser or strings.Split with care.
	if rawString == "" {
		*rs = []RoleEnum{}
		return nil
	}

	strs := strings.Split(rawString, ",")
	result := make([]RoleEnum, 0, len(strs))

	for _, s := range strs {
		s = strings.TrimSpace(s)

		// Handle quoted strings if your data contains spaces/special chars (e.g. "{role one,role two}")
		if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
			s = s[1 : len(s)-1]
		}

		if s == "" {
			continue
		}

		role := RoleEnum(s)
		if !IsValidRoleEnum(role) {
			return fmt.Errorf("invalid role scanned: %s", role)
		}
		result = append(result, role)
	}

	*rs = result
	return nil
}

// Helper to validate
func IsValidRoleEnum(r RoleEnum) bool {
	switch r {
	case RoleEnum_Analytics, RoleEnum_Management, RoleEnum_DBA, RoleEnum_Normal, RoleEnum_Banned:
		return true
	}
	return false
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
