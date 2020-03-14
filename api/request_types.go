package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/cyverse-de/requests/db"
	"github.com/cyverse-de/requests/model"
	"github.com/labstack/echo"
)

// validateRequestTypeName returns an error if a request type name is invalid.
func validateRequestTypeName(name string) error {
	re := regexp.MustCompile("^[\\w-]+$")
	if !re.MatchString(name) {
		return fmt.Errorf("request type names may only contain alphanumerics, underscores, and hyphens")
	}
	return nil
}

// GetRequestTypesHandler handles GET requests to the /request-types endpoint.
func (a *API) GetRequestTypesHandler(ctx echo.Context) error {

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Obtain the list of request types.
	requestTypes, err := db.ListRequestTypes(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, model.RequestTypeListing{
		RequestTypes: requestTypes,
	})
}

// RegisterRequestTypeHandler handles POST requests to the /request-types/{name} endpoint.
func (a *API) RegisterRequestTypeHandler(ctx echo.Context) error {
	name := ctx.Param("name")

	// Validate the request type name.
	err := validateRequestTypeName(name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// If a request type with the same name already exists, return it.
	requestType, err := db.GetRequestType(tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}
	if requestType != nil {
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return err
		}
		return ctx.JSON(http.StatusOK, requestType)
	}

	// The request type doesn't exist yet. Add it.
	requestType, err = db.AddRequestType(tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Return the response.
	return ctx.JSON(http.StatusOK, requestType)
}

// GetRequestTypeHandler handles GET requests to the /request-types/{name} endpoint.
func (a *API) GetRequestTypeHandler(ctx echo.Context) error {
	name := ctx.Param("name")

	// Start a transaction.
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}

	// Get the request type.
	requestType, err := db.GetRequestType(tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Return the response.
	if requestType == nil {
		return ctx.JSON(http.StatusNotFound, ErrorResponse{
			Message: "not found",
		})
	}
	return ctx.JSON(http.StatusOK, requestType)
}
