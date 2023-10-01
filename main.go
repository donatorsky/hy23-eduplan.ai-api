package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type APIErrorResponse struct {
	Details string `json:"details"`
}

func (r *APIErrorResponse) Error() string {
	return r.Details
}

func NewAPIErrorResponse(details string) *APIErrorResponse {
	return &APIErrorResponse{
		Details: details,
	}
}

func NewInternalServerErrorResponse() *APIErrorResponse {
	return &APIErrorResponse{
		Details: "Internal server error",
	}
}

type University struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Voivodeship string `json:"voivodeship"`
	City        string `json:"city"`
}

func main() {
	e := setup()

	e.Logger.Fatal(e.Start(":20292"))
}

func setup() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORS(),
	)

	db, err := sql.Open("sqlite3", "./database.sqlite")
	if err != nil {
		e.Logger.Fatalf("Failed to connect to the database: %s", err.Error())
	}
	// defer db.Close()

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	e.GET("/universities-by-specialization", func(c echo.Context) error {
		specializationName := c.QueryParam("name")

		if specializationName == "" {
			return c.JSON(http.StatusUnprocessableEntity, NewAPIErrorResponse("Query parameter 'name' is required"))
		}

		specializationsQuery, specializationsArgs, err := squirrel.StatementBuilder.
			Select("university_id").
			From("specializations").
			Where("name = ? COLLATE NOCASE", specializationName).
			ToSql()
		if err != nil {
			c.Logger().Error(err.Error())

			return c.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse())
		}

		rows, err := sq.
			Select(
				"id",
				"name",
				"link",
				"voivodeship",
				"city",
			).
			From("universities").
			Where(fmt.Sprintf("id IN (%s)", specializationsQuery), specializationsArgs...).
			RunWith(db).
			QueryContext(c.Request().Context())
		if err != nil {
			c.Logger().Error(err.Error())

			return c.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse())
		}

		defer rows.Close()

		universities := make([]University, 0)

		for rows.Next() {
			var university University

			if err := rows.Scan(
				&university.Id,
				&university.Name,
				&university.Link,
				&university.Voivodeship,
				&university.City,
			); err != nil {
				c.Logger().Error(err.Error())

				return c.JSON(http.StatusInternalServerError, NewInternalServerErrorResponse())
			}

			universities = append(universities, university)
		}

		if err := rows.Err(); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, universities)
	})

	return e
}
