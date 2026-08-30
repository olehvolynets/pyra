package dishes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"pyra/pkg/nutrition"

	"github.com/google/uuid"
)

var ErrNotNumber = errors.New("must be a number")

func dishRef(r *http.Request) (nutrition.DishRef, error) {
	paramUID := r.PathValue("uid")
	if _, err := uuid.Parse(paramUID); err != nil {
		return nutrition.DishRef{}, err
	}
	paramVersion := r.PathValue("version")

	parsedVersion, err := strconv.ParseUint(paramVersion, 10, 64)
	if err != nil || parsedVersion == 0 {
		return nutrition.DishRef{}, fmt.Errorf("invalid version: %s", paramVersion)
	}

	return nutrition.DishRef{
		UID: nutrition.DishUID(paramUID),
		Version: nutrition.DishVersion(parsedVersion),
	}, nil
}
