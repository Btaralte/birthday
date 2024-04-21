package api

import (
	"birthdayreminder/config"
	"birthdayreminder/db"
	"birthdayreminder/models"
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	e      *echo.Echo
	db     *db.DBService
	config *config.Config
}

func NewServer(config *config.Config) (*Server, error) {
	e := echo.New()
	db, err := db.NewDBService(config)
	if err != nil {
		return nil, err
	}
	s := &Server{
		e:      e,
		db:     db,
		config: config,
	}
	s.routes()
	return s, nil
}

func (s *Server) Start() {
	address := fmt.Sprintf(":%d", s.config.Port)
	s.e.Logger.Fatal(s.e.Start(address))
}

func (s *Server) routes() {
	s.e.GET("/api/birthday/listAll", s.getAllBirthdays)
	s.e.POST("/api/birthday/create", s.insertBirhtDayHandler)
}
func (s *Server) insertBirhtDayHandler(c echo.Context) error {
	var birthday models.BirthDay
	if err := c.Bind(&birthday); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Input")
	}
	if err := s.db.InsertBirthday(context.Background(), &birthday); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusCreated, birthday)
}

func (s *Server) getAllBirthdays(c echo.Context) error {
	birthdays, err := s.db.GetAllBirthDays(context.Background())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	var json = make(map[string]interface{})
	if birthdays != nil {
		json["birthDays"] = birthdays
	} else {
		json["birthDays"] = make([]interface{}, 0)
	}
	return c.JSON(http.StatusOK, json)
}
