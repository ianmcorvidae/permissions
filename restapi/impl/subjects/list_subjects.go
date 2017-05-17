package subjects

import (
	"database/sql"
	"github.com/cyverse-de/permissions/logger"
	"github.com/cyverse-de/permissions/models"
	permsdb "github.com/cyverse-de/permissions/restapi/impl/db"
	"github.com/cyverse-de/permissions/restapi/operations/subjects"

	"github.com/go-openapi/runtime/middleware"
)

func listSubjectsInternalServerError(reason string) middleware.Responder {
	return subjects.NewListSubjectsInternalServerError().WithPayload(
		&models.ErrorOut{Reason: &reason},
	)
}

func BuildListSubjectsHandler(db *sql.DB) func(subjects.ListSubjectsParams) middleware.Responder {

	// Return the handler function.
	return func(params subjects.ListSubjectsParams) middleware.Responder {

		// Start a transaction for the request.
		tx, err := db.Begin()
		if err != nil {
			logger.Log.Error(err)
			return listSubjectsInternalServerError(err.Error())
		}

		// Obtain the list of subjects.
		result, err := permsdb.ListSubjects(tx, params.SubjectType, params.SubjectID)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return listSubjectsInternalServerError(err.Error())
		}

		// Commit the transaction for the request.
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Log.Error(err)
			return listSubjectsInternalServerError(err.Error())
		}

		// Return the result.
		return subjects.NewListSubjectsOK().WithPayload(&models.SubjectsOut{Subjects: result})
	}
}
