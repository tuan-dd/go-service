package persistence

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/tuan-dd/go-pkg/common/response"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/generated/ent"
	"github.com/tuan-dd/go-service/url/internal/infrastructure/global"
)

func ConvertEntErr(err error, messagePrefix string) *response.AppError {
	if err == nil {
		return nil
	}

	// 1. Check for ent.NotFoundError
	if ent.IsNotFound(err) {
		return response.NotFound(messagePrefix + " not found")
	}

	// 2. Check for ent.ValidationError
	if ent.IsValidationError(err) {
		return response.QueryInvalid(messagePrefix + ": Validation failed.")
	}

	// 3. Check for ent.ConstraintError (most complex path)
	if ent.IsConstraintError(err) {
		underlyingErr := errors.Unwrap(err)
		if underlyingErr == nil { // Should not happen if IsConstraintError is true, but good practice
			underlyingErr = err
		}

		// Attempt to unwrap to *mysql.MySQLError (MySQL)
		var mysqlErr *mysql.MySQLError
		if errors.As(underlyingErr, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062: // ER_DUP_ENTRY
				return response.Duplicate(messagePrefix + " already exists")
			case 1451: // ER_ROW_IS_REFERENCED_2 (cannot delete/update parent)
				return response.Conflict(messagePrefix + " cannot be modified because it is being used by other records")
			case 1452: // ER_NO_REFERENCED_ROW_2 (cannot add/update child)
				return response.QueryInvalid(messagePrefix + " contains invalid reference data")
			case 1048: // ER_BAD_NULL_ERROR
				return response.QueryInvalid(messagePrefix + " is missing required information")
			case 1216, 1217: // ER_NO_REFERENCED_ROW_1, ER_ROW_IS_REFERENCED
				return response.QueryInvalid(messagePrefix + " contains invalid reference data")
			case 1205, 1213: // ER_LOCK_DEADLOCK
				global.Log.Error("Database deadlock occurred", mysqlErr)
				return response.ServerError("Operation temporarily unavailable, please try again")
			default:
				global.Log.Error("Database constraint error", mysqlErr)
				return response.ServerError("Unable to complete the operation")
			}
		}
	}

	// 4. Check for ent.NotLoadedError
	if ent.IsNotLoaded(err) {
		// This is typically a server-side logic error.
		// Log it for developer attention.
		var nle *ent.NotLoadedError
		edgeName := "unknown edge"
		if errors.As(err, &nle) {
			edgeName = nle.Error()
		}
		developerMsg := fmt.Sprintf("Server logic error: edge '%s' was not loaded before access.", edgeName)
		global.Log.Error(developerMsg, err)
		return response.ServerError("")
	}

	// 5. Check for ent.NotSingularError
	if ent.IsNotSingular(err) {
		global.Log.Warn("NotSingularError encountered", err)
		return response.NotFound(messagePrefix + ": The requested item is not uniquely identifiable or not found.")
	}

	// Fallback for other/unknown errors (as in user's original function)
	return response.ConvertDatabaseError(err)
}
