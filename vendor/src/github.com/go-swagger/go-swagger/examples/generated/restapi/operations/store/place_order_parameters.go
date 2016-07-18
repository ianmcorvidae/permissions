package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-swagger/go-swagger/errors"
	"github.com/go-swagger/go-swagger/httpkit"
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/go-swagger/go-swagger/examples/generated/models"
)

// NewPlaceOrderParams creates a new PlaceOrderParams object
// with the default values initialized.
func NewPlaceOrderParams() PlaceOrderParams {
	var ()
	return PlaceOrderParams{}
}

// PlaceOrderParams contains all the bound params for the place order operation
// typically these are obtained from a http.Request
//
// swagger:parameters placeOrder
type PlaceOrderParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*order placed for purchasing the pet
	  In: body
	*/
	Body *models.Order
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *PlaceOrderParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if httpkit.HasBody(r) {
		defer r.Body.Close()
		var body models.Order
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("body", "body", "", err))
		} else {
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
