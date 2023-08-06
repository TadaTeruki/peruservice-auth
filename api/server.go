package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type Server struct {
	Router *echo.Echo
	db     *sqlx.DB
	config *ServerConfig
}

func NewServer(e *echo.Echo) *Server {
	conf, err := QueryServerConfig()
	if err != nil {
		log.Fatalf("failed to query server config: %v", err)
	}
	db, err := sqlx.Open(
		"postgres",
		"host="+conf.envConf.DBHost+
			" port="+conf.envConf.DBPort+
			" user="+conf.envConf.DBUser+
			" password="+conf.envConf.DBPassWord+
			" dbname="+conf.envConf.DBName+
			" sslmode=disable",
	)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	return &Server{
		Router: e,
		db:     db,
		config: conf,
	}
}

func (s *Server) Start() error {

	// middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// allow origins
	var allow_origins []string
	if s.config.envConf.Mode == "PRODUCTION" {
		allow_origins = s.config.envConf.AuthAllowOrigins
	} else {
		allow_origins = []string{"*"}
	}

	// set CORS
	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allow_origins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
		AllowCredentials: true,
	}))

	// routes
	s.Router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server: ok")
	})

	s.Router.POST("/login", s.Login)
	s.Router.POST("/refresh", s.Refresh)

	return s.Router.Start(":" + s.config.envConf.AuthPort)
}

func (s *Server) Login(c echo.Context) error {
	// bind request to LoginRequest struct
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "incorrect request body",
		})
	}

	// query row from db which has the same admin_id as req.AdminID
	var admin Admin
	err := s.db.Get(&admin, "SELECT * FROM admin WHERE admin_id=$1", req.AdminID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to find accout: " + err.Error(),
		})
	}

	// compare req.Password and admin.Password
	if req.Password != admin.Password {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "invalid password",
		})
	}

	// generate refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"admin_id": admin.AdminID,
		"type":     "refresh",
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.jsonConf.RefreshTokenExpDurationHour)).Unix(),
	})

	// get private key from file
	privateKeySrc, err := os.ReadFile(s.config.envConf.PrivateKeyFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to read private key: " + err.Error(),
		})
	}

	// parse private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeySrc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse private key: " + err.Error(),
		})
	}

	// create refresh token string
	token_str, err := token.SignedString(privateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create refresh token: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		RefreshToken: token_str,
	})
}

func (s *Server) Refresh(c echo.Context) error {
	// get token from header
	refresh_token_str := c.Request().Header.Get("Authorization")
	if refresh_token_str == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "no token provided",
		})
	}

	// remove "Bearer " from token
	refresh_token_str = refresh_token_str[7:]

	// get public key from file
	publicKeySrc, err := os.ReadFile(s.config.envConf.PublicKeyFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to read public key: " + err.Error(),
		})
	}

	// parse token
	refresh_token, err := jwt.Parse(refresh_token_str, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKeySrc)
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid token: " + err.Error(),
		})
	}

	// get admin_id from token
	claims, ok := refresh_token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "no required claims in token",
		})
	}

	// check admin_id
	admin_id, ok := claims["admin_id"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "no required claims in token",
		})
	}

	// check token type (is refresh token?)
	token_type, _ := claims["type"].(string)
	if token_type != "refresh" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid token type",
		})
	}

	// get private key from file
	privateKeySrc, err := os.ReadFile(s.config.envConf.PrivateKeyFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to read private key: " + err.Error(),
		})
	}

	// parse private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeySrc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse private key: " + err.Error(),
		})
	}

	// create access token
	access_token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"admin_id": admin_id,
		"type":     "access",
		"exp":      time.Now().Add(time.Minute * time.Duration(s.config.jsonConf.AccessTokenExpDurationMin)).Unix(),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to parse private key: " + err.Error(),
		})
	}

	// create access token string
	access_token_str, err := access_token.SignedString(privateKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to create access token: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: access_token_str,
	})
}
