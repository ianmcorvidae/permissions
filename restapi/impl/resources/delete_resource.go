package resources

import (
	"database/sql"
	"fmt"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/resources"

	"github.com/cyverse-de/logcabin"
	"github.com/go-openapi/runtime/middleware"
)

func BuildDeleteResourceHandler(db *sql.DB) func(resources.DeleteResourceParams) middleware.Responder {

	// Return the handler function.
	return func(params resources.DeleteResourceParams) middleware.Responder {

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewDeleteResourceInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Verify that the resource exists.
		exists, err := permsdb.ResourceExists(tx, &params.ID)
		if err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewDeleteResourceInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if !exists {
			tx.Rollback()
			reason := fmt.Sprintf("resource, %s, not found", params.ID)
			return resources.NewDeleteResourceNotFound().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Delete the resource.
		err = permsdb.DeleteResource(tx, &params.ID)
		if err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewDeleteResourceInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Commit the transaction.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logcabin.Error.Print(err)
			reason := err.Error()
			return resources.NewDeleteResourceInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		return resources.NewDeleteResourceOK()
	}
}
