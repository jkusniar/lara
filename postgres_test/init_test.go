/*
   Copyright (C) 2016-2017 Contributors as noted in the AUTHORS file

   This file is part of lara, veterinary practice support software.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jkusniar/lara"
	"github.com/jkusniar/lara/postgres"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	connServ = `postgres://%s:%s@%s:%d?sslmode=disable`
)

func createDB(dbUser, dbPass, dbHost, dbName string, dbPort uint64) {
	db, err := sql.Open("postgres", fmt.Sprintf(connServ, dbUser, dbPass, dbHost, dbPort))
	if err != nil {
		log.Panicf("Error opening test database: %+v", err)
	}
	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
		log.Panicf("Error droping database %s: %+v", dbName, err)
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH OWNER = %s TEMPLATE = template0 ENCODING = 'UTF8'",
		dbName, dbUser)); err != nil {
		log.Panicf("Error creating database %s: %+v", dbName, err)
	}
}

func runSql(db *sql.DB, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(b))

	return err
}

var (
	ownerService   lara.OwnerService
	patientService lara.PatientService
	recordService  lara.RecordService
	searchService  lara.SearchService
	userService    lara.UserService
	productService lara.ProductService
	reportService  lara.ReportService
	titleService   lara.TitleService
	unitService    lara.UnitService
	genderService  lara.GenderService
	speciesService lara.SpeciesService
	breedService   lara.BreedService
	addressService lara.AddressService
	tagService     lara.TagService
	testCtx        context.Context
)

func TestMain(m *testing.M) {
	dbHost := fromEnv("POSTGRES_HOST", "localhost")
	dbUser := fromEnv("POSTGRES_USER", "postgres")
	dbPass := fromEnv("POSTGRES_PASSWORD", "")
	dbName := fromEnv("POSTGRES_DB", "lara_test")
	dbPort, err := strconv.ParseUint(os.Getenv("POSTGRES_PORT"), 10, 64)
	if err != nil {
		dbPort = 5432
	}

	createDB(dbUser, dbPass, dbHost, dbName, dbPort)

	db, err := postgres.Open(dbUser, dbPass, dbHost,
		dbName, uint(dbPort), "disable")
	if err != nil {
		log.Fatalf("FATAL: %+v", err)
	}
	defer db.Close()

	if err := runSql(db, "schema.sql"); err != nil {
		log.Panicf("error creating DB schema: %+v", err)
	}

	if err := runSql(db, "test_data.sql"); err != nil {
		log.Panicf("error creating test data: %+v", err)
	}

	searchService = &postgres.SearchService{DB: db}
	ownerService = &postgres.OwnerService{DB: db}
	patientService = &postgres.PatientService{DB: db}
	recordService = &postgres.RecordService{DB: db}
	userService = &postgres.UserService{DB: db, Pass: &passwordMock{}}
	productService = &postgres.ProductService{DB: db}
	addressService = &postgres.AddressService{DB: db}
	sls := postgres.SimpleLovService{DB: db}
	titleService = &sls
	unitService = &sls
	genderService = &sls
	speciesService = &sls
	breedService = &sls
	loc, _ := time.LoadLocation("Europe/Bratislava") // time.Location for unit tests
	reportService = &postgres.ReportService{DB: db, Loc: loc}
	tagService = &postgres.TagService{DB: db}

	// test user in context
	u, _ := lara.MakeUser("testuser",
		[]string{
			lara.ViewRecord.String(),
			lara.EditRecord.String(),
			lara.ViewReports.String()})
	testCtx = lara.ContextWithUser(context.Background(), u)

	log.SetOutput(ioutil.Discard) //calm down logger in tests
	result := m.Run()

	db.Close()
	os.Exit(result)
}

func fromEnv(envName, defVal string) (result string) {
	result = os.Getenv(envName)
	if result == "" {
		result = defVal
	}
	return
}

// checkErrCode checks error's code by calling it's Code() method.
// Actual error's code is returned as well
func checkErrCode(err error, expCode int) (codeEqual bool, actualCode int) {
	type codedError interface {
		Code() int
	}

	if ce, ok := errors.Cause(err).(codedError); ok {
		codeEqual = ce.Code() == expCode
		actualCode = ce.Code()
	}

	return
}

func timeEmpty(t time.Time) bool {
	return t.String() == "0001-01-01 00:00:00 +0000 UTC"
}
