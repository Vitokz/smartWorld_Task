package server

import (
	"context"
	"encoding/json"
	"github.com/Vitokz/smartWorld_Task/handler"
	myJwt "github.com/Vitokz/smartWorld_Task/internal/services/jwt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Router  *echo.Echo
	Handler *handler.Handler
}

func (s *Server) Routing() {
	gUser := s.Router.Group("/user")
	gUser.POST("/login", s.Handler.Login)
	gUser.GET("/logout", s.Handler.Logout, s.AuthMiddleware)
	gUser.POST("/register", s.Handler.Registration)

	gLibrary := s.Router.Group("/library")
	gLibrary.GET("/take_book/:bookId", s.Handler.TakeBook, s.AuthMiddleware)

	gAdmin := s.Router.Group("/admin")
	gAdmin.GET("/block_user/:userId", s.Handler.BlockUser, s.AuthMiddleware)

	gRating := s.Router.Group("/rating")
	gRating.GET("/all_time", s.Handler.RatingAllTime)
	gRating.GET("/ten_days", s.Handler.RatingReserved)
}

func NewServer(hdlr *handler.Handler) *Server {
	serv := &Server{
		Router:  newRouter(),
		Handler: hdlr,
	}

	serv.Routing()

	return serv
}

func newRouter() *echo.Echo {
	router := echo.New()
	router.Use(middleware.Logger())

	return router
}

func extractToken(req *http.Request) string {
	headers := req.Header
	bearToken := headers.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			req    = c.Request()
			bearer = extractToken(req)
			claims = new(myJwt.Claims)
		)

		tkn, err := jwt.ParseWithClaims(bearer, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Handler.Config.Jwt.SigningKey), nil
		})
		if err != nil && claims.StandardClaims.ExpiresAt > time.Now().Unix() {
			return c.JSON(http.StatusUnauthorized, "not valid jwt token")
		}

		if tkn != nil && tkn.Valid {
			claimsByte, err := json.Marshal(*claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed marshal claims")
			}

			c.Set("claims", claimsByte)

			return next(c)
		}

		id, err := strconv.Atoi(claims.Id)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.ErrUnauthorized)
		}
		tokens, err := s.Handler.Repository.TakeUserJWTByID(context.Background(), id)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.ErrUnauthorized)
		}

		refTkn, err := jwt.ParseWithClaims(tokens.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Handler.Config.Jwt.SigningKey), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "tokens are expired")
		}

		if refTkn != nil && refTkn.Valid {
			claimsByte, err := json.Marshal(*claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "failed marshal claims")
			}

			c.Set("claims", claimsByte)

			newAccToken, newRefreshToken, err := s.Handler.RefreshTokens(*claims)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.ErrUnauthorized)
			}

			c.Set("NewAccessToken", newAccToken)
			c.Set("NewRefreshToken", newRefreshToken)

			return next(c)
		}

		return next(c)
	}
}
