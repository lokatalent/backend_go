package util

import (
	"github.com/labstack/echo/v4"
)

// ErrInternalServer logs server error before returning the actual error
func ErrInternalServer(ctx echo.Context, err error) error {
	ctx.Logger().Error(err)

	return err
}
