package routes

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qustavo/dotsql"
	wisdomType "github.com/wisdom-oss/commonTypes/v2"

	"github.com/getkin/kin-openapi/openapi3"

	"microservice/globals"
)

var apiContract *openapi3.T

func TestMain(m *testing.M) {
	godotenv.Load("../.env", ".env")

	var c wisdomType.EnvironmentConfiguration
	err := c.PopulateFromFilePath("./resources/environment.json")
	if os.IsNotExist(err) {
		err = c.PopulateFromFilePath("../resources/environment.json")
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	globals.Environment, err = c.ParseEnvironment()
	if err != nil {
		panic(err)
	}
	pgxAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/wisdom",
		globals.Environment["PG_USER"], globals.Environment["PG_PASS"],
		globals.Environment["PG_HOST"], globals.Environment["PG_PORT"])
	pgxConfig, err := pgxpool.ParseConfig(pgxAddress)
	if err != nil {
		panic(err)
	}
	globals.Db, err = pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}

	globals.SqlQueries, err = dotsql.LoadFromFile("./resources/queries.sql")
	if os.IsNotExist(err) {
		globals.SqlQueries, err = dotsql.LoadFromFile("../resources/queries.sql")
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	apiContract, err = openapi3.NewLoader().LoadFromFile("./openapi.yaml")
	if os.IsNotExist(err) {
		apiContract, err = openapi3.NewLoader().LoadFromFile("../openapi.yaml")
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	// now run the tests
	os.Exit(m.Run())
}
