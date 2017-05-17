package resource_types

import (
	"database/sql"
	"fmt"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/resource_types"
	"github.com/go-openapi/runtime/middleware"
)

func BuildResourceTypesPostHandler(db *sql.DB) func(resource_types.PostResourceTypesParams) middleware.Responder {

	// Return the handler function.
	return func(params resource_types.PostResourceTypesParams) middleware.Responder {
		resourceTypeIn := params.ResourceTypeIn

		// Start a transaction for this request.
		tx, err := db.Begin()
		if err != nil {
			reason := err.Error()
			return resource_types.NewPostResourceTypesInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Check for a duplicate name.
		duplicate, err := permsdb.GetResourceTypeByName(tx, resourceTypeIn.Name)
		if err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewPostResourceTypesInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}
		if duplicate != nil {
			tx.Rollback()
			reason := fmt.Sprintf("a resource type named %s already exists", *resourceTypeIn.Name)
			return resource_types.NewPostResourceTypesBadRequest().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		// Save the resource type.
		resourceTypeOut, err := permsdb.AddNewResourceType(tx, resourceTypeIn)
		if err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewPostResourceTypesInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			reason := err.Error()
			return resource_types.NewPostResourceTypesInternalServerError().WithPayload(
				&models.ErrorOut{Reason: &reason},
			)
		}

		return resource_types.NewPostResourceTypesCreated().WithPayload(resourceTypeOut)
	}
}
