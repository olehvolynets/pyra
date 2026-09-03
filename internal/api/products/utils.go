package products

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"uuid"

	"pyra/pkg/nutrition"
)

var ErrNotNumber = errors.New("must be a number")

func productRef(r *http.Request) (nutrition.ProductRef, error) {
	paramUID := r.PathValue("uid")
	if _, err := uuid.Parse(paramUID); err != nil {
		return nutrition.ProductRef{}, err
	}
	paramVersion := r.PathValue("version")

	parsedVersion, err := strconv.ParseUint(paramVersion, 10, 64)
	if err != nil || parsedVersion == 0 {
		return nutrition.ProductRef{}, fmt.Errorf("invalid version: %s", paramVersion)
	}

	return nutrition.ProductRef{
		UID: nutrition.ProductUID(paramUID),
		Version: nutrition.ProductVersion(parsedVersion),
	}, nil
}
