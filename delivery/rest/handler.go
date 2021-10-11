package rest

import (
	"net/http"
	"strconv"

	"github.com/egaevan/merchant/model"
	"github.com/egaevan/merchant/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ProductUsecae usecase.ProductUsecae
}

func NewHandler(e *echo.Echo, productUsecae usecase.ProductUsecae) {
	handler := &Handler{
		ProductUsecae: productUsecae,
	}

	e.GET("/product", handler.GetProduct)
	e.GET("/product/:productID", handler.GetOneProduct)
	e.POST("/product", handler.SendProduct)
	e.PATCH("/product/:productID", handler.UpdateProduct)
	e.DELETE("/product/:productID", handler.DeleteProduct)
}

func (h *Handler) GetOneProduct(c echo.Context) error {
	ctx := c.Request().Context()
	productIDParam := c.Param("productID")

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecae.GetOneProduct(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.ProductUsecae.GetProduct(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) SendProduct(c echo.Context) error {
	ctx := c.Request().Context()
	dataReq := model.Product{}

	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecae.SendProduct(ctx, dataReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()
	dataReq := model.ProductUpdate{}
	productIDParam := c.Param("productID")

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	res, err := h.ProductUsecae.UpdateProduct(ctx, dataReq, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()
	productIDParam := c.Param("productID")

	if productIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	productID, err := strconv.Atoi(productIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	err = h.ProductUsecae.DeleteProduct(ctx, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, responseError{
		Message: "Product has been deleted",
	})
}

type responseError struct {
	Message string `json:"message"`
}
