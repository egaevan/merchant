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
	UserUsecae    usecase.UserUsecae
}

type responseError struct {
	Message string `json:"message"`
}

const (
	isAdmin int = 1
)

func NewHandler(e *echo.Echo, productUsecae usecase.ProductUsecae, userUsecae usecase.UserUsecae) {
	handler := &Handler{
		ProductUsecae: productUsecae,
		UserUsecae:    userUsecae,
	}

	// Routing Product
	e.GET("/product", handler.GetProduct)
	e.GET("/product/:productID", handler.GetOneProduct)
	e.POST("/product", handler.SendProduct, JwtVerify)
	e.PATCH("/product/:productID", handler.UpdateProduct, JwtVerify)
	e.DELETE("/product/:productID", handler.DeleteProduct, JwtVerify)

	// Routing User
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)
	e.DELETE("/user/:userID", handler.DeleteUser, JwtVerify)
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
	dataReq := model.ProductDetail{}

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
	dataReq := model.ProductDetail{}
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

func (h *Handler) Login(c echo.Context) error {
	dataReq := model.User{}
	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})
		return echo.ErrBadRequest
	}

	user, err := h.UserUsecae.Login(c.Request().Context(), dataReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: err.Error(),
		})
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged in",
		"token":   user.Token,
	})
}

func (h *Handler) Register(c echo.Context) error {
	dataReq := model.User{}
	if err := c.Bind(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid data request",
		})

		return echo.ErrBadRequest
	}

	err := h.UserUsecae.CreateUser(c.Request().Context(), dataReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: err.Error(),
		})
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, "success")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	userIDParam := c.Param("userID")

	userInfo := c.Get("user").(*model.Token)

	if userInfo.Role != isAdmin {
		// unauthorized
		return echo.ErrUnauthorized
	}

	// admin
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "invalid parameter",
		})

		return echo.ErrBadRequest
	}

	err = h.UserUsecae.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseError{
			Message: "internal error",
		})

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, responseError{
		Message: "User has been deleted",
	})
}
