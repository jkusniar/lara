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

package http_test

import (
	"io"
	"io/ioutil"
	"log"
	syshttp "net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jkusniar/lara"
	"github.com/jkusniar/lara/http"
	"github.com/jkusniar/lara/mock"
	"github.com/pkg/errors"
)

type testAuthToken struct {
}

func (t *testAuthToken) Create(u *lara.User) (string, error) {
	return u.Login, nil
}

func (t *testAuthToken) Parse(token string) (*lara.User, error) {
	return lara.MakeUser("testuser",
		[]string{
			lara.ViewRecord.String(),
			lara.EditRecord.String(),
			lara.ViewReports.String()}) // authenticate, full permissions
}

func newHttpHandler() syshttp.Handler {
	ownMock := mock.OwnerService{}

	// GetOwner with ID "42" throws error. ID 2 doesn't exist
	ownMock.GetFn = func(id uint64) (*lara.GetOwner, error) {
		switch id {
		case 42:
			return nil, errors.New("get owner by id failed")
		case 2:
			return nil, lara.NewCodedError(404, errors.New("owner with ID 2 not found"))
		default:
			return &lara.GetOwner{
				Versioned: lara.Versioned{ID: 1},
				Owner: lara.Owner{FirstName: "first",
					LastName: "last"},
				Patients: []lara.OwnersPatient{}}, nil
		}
	}
	// ID 42 is SQL error, ID 2 doesn't exist in DB
	ownMock.UpdateFn = func(id uint64, o *lara.UpdateOwner) error {
		switch id {
		case 42:
			return errors.New("update owner failed")
		case 2:
			return lara.NewCodedError(404, errors.New("owner with ID 2 not found"))
		default:
			return nil
		}
	}

	// create user with LastName "error" causes error.
	ownMock.CreateFn = func(o *lara.CreateOwner) (uint64, error) {
		if o.LastName == "error" {
			return 0, errors.New("create owner failed")
		}
		return 2, nil
	}

	searchMock := mock.SearchService{}

	// search for q == "error" causes error
	searchMock.SearchFn = func(q string) (*lara.SearchResult, error) {
		switch q {
		case "error":
			return nil, errors.New("search failed")
		default:
			return &lara.SearchResult{Total: 2,
				Records: []lara.SearchRecord{
					{ID: 1, Name: "Johnny GetOwner"},
					{ID: 2, Name: "Dunco"}}}, nil
		}
	}

	productMock := mock.ProductService{}
	// search for p.Query == "error" causes error
	productMock.SearchFn = func(p *lara.ProductSearchRequest) (*lara.ProductSearchResult,
		error) {
		switch p.Query {
		case "error":
			return nil, errors.New("search failed")
		default:
			return &lara.ProductSearchResult{Total: 2,
				Products: []lara.Product{
					{ID: 1, Name: "Prod1", Unit: "Unit1",
						Price: "1.00"},
					{ID: 2, Name: "Prod2", Unit: "Unit2",
						Price: "2.00"}}}, nil
		}
	}

	reportMock := mock.ReportService{}
	reportMock.GetIncomeStatisticsFn = func(r *lara.ReportRequest) (*lara.IncomeStatistics, error) {
		t, _ := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Apr 22, 2003 at 1:00pm (UTC)")
		loc, _ := time.LoadLocation("Europe/Bratislava")
		if r.ValidFrom.In(loc) != t.In(loc) {
			return nil, errors.New("report failed")
		}

		return &lara.IncomeStatistics{
			Records:         1,
			Income:          "9.42",
			IncomeBilled:    "6.28",
			IncomeNotBilled: "3.14"}, nil
	}

	patientMock := mock.PatientService{}
	patientMock.GetFn = func(id uint64) (*lara.GetPatient, error) {
		return &lara.GetPatient{
			Versioned: lara.Versioned{ID: 1},
			Patient:   lara.Patient{Name: "pet"},
			Records:   []lara.PatientsRecord{},
			Tags:      []lara.PatientsTag{},
		}, nil
	}
	patientMock.CreateFn = func(r *lara.CreatePatient) (uint64, error) {
		return 42, nil
	}
	patientMock.UpdateFn = func(id uint64, p *lara.UpdatePatient) error {
		return nil
	}

	recordMock := mock.RecordService{}
	recordMock.GetFn = func(id uint64) (*lara.GetRecord, error) {
		return &lara.GetRecord{
			Versioned: lara.Versioned{ID: 1},
			Items: []lara.GetRecordItem{{
				ID: 2,
				RecordItem: lara.RecordItem{
					ProductID:    100,
					ProductPrice: "1.00",
					Amount:       "2.00",
					ItemPrice:    "3.00",
					ItemType:     lara.Labor,
				},
				PLU: "10",
			}},
		}, nil
	}
	recordMock.CreateFn = func(*lara.CreateRecord) (uint64, error) {
		return 42, nil
	}
	recordMock.UpdateFn = func(id uint64, r *lara.UpdateRecord) error {
		return nil
	}

	makeSimpleLOVGetAllFn := func(lovType string) func() (*lara.LOVItemList, error) {
		return func() (*lara.LOVItemList, error) {
			return &lara.LOVItemList{Items: []lara.LOVItem{{ID: 42, Name: lovType}}}, nil
		}
	}

	sls := mock.SimpleLovService{}
	sls.GetAllTitlesFn = makeSimpleLOVGetAllFn("title")
	sls.GetAllUnitsFn = makeSimpleLOVGetAllFn("unit")
	sls.GetAllGendersFn = makeSimpleLOVGetAllFn("gender")
	sls.GetAllSpeciesFn = makeSimpleLOVGetAllFn("species")
	sls.GetAllBreedsFn = func(speciesId uint64) (*lara.LOVItemList, error) {
		if speciesId != 42 {
			return nil, errors.New("get by id failed")
		}
		return &lara.LOVItemList{Items: []lara.LOVItem{{ID: speciesId, Name: "breed"}}}, nil
	}

	addressMock := mock.AddressService{}
	addressMock.SearchCityFn = func(query string) (*lara.CityStreetList, error) {
		if query == "fail" {
			return nil, errors.New("search city failed")
		}
		return &lara.CityStreetList{Total: 1, Items: []lara.CityStreet{{ID: 1, Name: "city", ZIP: "000"}}}, nil
	}

	addressMock.SearchStreetForCityFn = func(cityID uint64, query string) (*lara.CityStreetList, error) {
		if query == "fail" {
			return nil, errors.New("search street failed")
		}
		return &lara.CityStreetList{Total: 1, Items: []lara.CityStreet{{ID: 1, Name: "street", ZIP: "000"}}}, nil
	}

	tagMock := mock.TagService{}
	tagMock.GetFn = func(id uint64) (*lara.GetTag, error) {
		return &lara.GetTag{
			Versioned: lara.Versioned{ID: 1},
			Type:      "RFID",
			Value:     "007",
			Data:      []byte{},
		}, nil
	}
	tagMock.CreateFn = func(r *lara.CreateTag) (uint64, error) {
		return 42, nil
	}
	tagMock.UpdateFn = func(id uint64, p *lara.UpdateTag) error {
		return nil
	}
	tagMock.GetPatientByTagFn = func(tagValue string) (*lara.PatientByTag, error) {
		return &lara.PatientByTag{Name: "p", Species: "s", Breed: "b", Gender: "g",
			OwnerID: 1, OwnerName: "n", OwnerAddress: "a"}, nil
	}

	srv := http.Server{
		Token:          &testAuthToken{},
		SearchService:  &searchMock,
		OwnerService:   &ownMock,
		PatientSevice:  &patientMock,
		RecordService:  &recordMock,
		ProductService: &productMock,
		ReportService:  &reportMock,
		TitleService:   &sls,
		UnitService:    &sls,
		GenderService:  &sls,
		SpeciesService: &sls,
		BreedService:   &sls,
		AddressService: &addressMock,
		TagService:     &tagMock,
	}

	return srv.Router()
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) //calm down logger in tests
	os.Exit(m.Run())
}

func TestAllHandlers(t *testing.T) {
	var tests = []struct {
		name            string
		reqMethod       string
		reqURL          string
		reqBody         io.Reader
		expCode         int
		expBody         string
		useBodyContains bool // when true, use strings.Contains to compare acutal/expected resp body
	}{
		//Search
		{"SearchHandler_OK",
			"GET", "/api/v1/search?q=something", nil, 200,
			`{"total":2,"records":[{"id":1,"name":"Johnny GetOwner","address":""},{"id":2,"name":"Dunco","address":""}]}` + "\n",
			false},
		{"SearchHandler_NoQuery",
			"GET", "/api/v1/search", nil, 200,
			`{"total":2,"records":[{"id":1,"name":"Johnny GetOwner","address":""},{"id":2,"name":"Dunco","address":""}]}` + "\n",
			false},
		{"SearchHandler_Error",
			"GET", "/api/v1/search?q=error", nil, 500,
			`search failed`, true},

		// SearchPatientByTagHandler tests
		{"SearchPatientByTagHandler_OK",
			"GET", "/api/v1/search/patient-by-tag/xyz", nil, 200,
			`{"tagType":"","name":"p","species":"s","breed":"b","gender":"g","ownerId":1,"ownerName":"n","ownerAddress":"a"}` + "\n", false},

		// GetOwnerHandler tests
		{"GetOwnerHandler_OK",
			"GET", "/api/v1/owner/1", nil, 200,
			`{"id":1,"version":0,"firstName":"first","lastName":"last","titleId":0,"cityId":0,"streetId":0,"houseNo":"","phone1":"","phone2":"","email":"","note":"","IC":"","DIC":"","ICDPH":"","title":"","city":"","street":"","creator":"","created":"0001-01-01T00:00:00Z","modifier":"","modified":"0001-01-01T00:00:00Z","patients":[]}` + "\n", false},
		{"GetOwnerHandler_BadParam",
			"GET", "/api/v1/owner/Nan", nil, 404,
			"invalid owner ID", true},
		{"GetOwnerHandler_NotFound",
			"GET", "/api/v1/owner/2", nil, 404,
			"owner with ID 2 not found", true},
		{"GetOwnerHandler_SqlError",
			"GET", "/api/v1/owner/42", nil, 500,
			"get owner by id failed", true},

		// CreateOwnerHandler tests
		{"CreateOwnerHandler_OK",
			"POST", "/api/v1/owner",
			strings.NewReader(`{"firstName":"A","lastName":"B"}`),
			200, "2", false},
		{"CreateOwnerHandler_BadJSON",
			"POST", "/api/v1/owner",
			strings.NewReader(`:-){"firstName":"A","lastName":"B"}`),
			400, "json decode error", true},
		{"CreateOwnerHandler_NoBody",
			"POST", "/api/v1/owner", strings.NewReader(""), 400,
			`json decode error: EOF`, true},
		{"CreateOwnerHandler_SqlError",
			"POST", "/api/v1/owner",
			strings.NewReader(`{"firstName":"error","lastName":"error"}`),
			500, "create owner failed", true},

		// UpdateOwnerHandler tests
		{"UpdateOwnerHandler_OK",
			"PUT", "/api/v1/owner/1",
			strings.NewReader(`{"firstName":"A","lastName":"B"}`),
			200, "", false},
		{"UpdateOwnerHandler_BadJSON",
			"PUT", "/api/v1/owner/1",
			strings.NewReader(`:-){"firstName":"A","lastName":"B"}`),
			400, "json decode error", true},
		{"UpdateOwnerHandler_NoBody",
			"PUT", "/api/v1/owner/1", strings.NewReader(""), 400,
			`json decode error: EOF`, true},
		{"UpdateOwnerHandler_BadIdFormat",
			"PUT", "/api/v1/owner/NaN",
			strings.NewReader(`{"firstName":"A","lastName":"B"}`),
			404, "invalid owner ID", true},
		{"UpdateOwnerHandler_NotFound",
			"PUT", "/api/v1/owner/2",
			strings.NewReader(`{"firstName":"A","lastName":"B"}`),
			404, "owner with ID 2 not found", true},
		{"UpdateOwnerHandler_SqlError",
			"PUT", "/api/v1/owner/42",
			strings.NewReader(`{"firstName":"A","lastName":"B"}`),
			500, "update owner failed", true},

		//SearchProduct
		{"SearchProductHandler_OK",
			"POST", "/api/v1/productsearch",
			strings.NewReader(`{"Query":"test"}`), 200,
			`{"total":2,"products":[{"id":1,"name":"Prod1","unit":"Unit1","price":"1.00"},{"id":2,"name":"Prod2","unit":"Unit2","price":"2.00"}]}` + "\n",
			false},
		{"SearchProductHandler_NoBody",
			"POST", "/api/v1/productsearch", strings.NewReader(""), 400,
			`json decode error: EOF`, true},
		{"SearchProductHandler_BadJSON",
			"POST", "/api/v1/productsearch",
			strings.NewReader(`:-){}{}`),
			400, "json decode error", true},
		{"SearchProductHandler_Error",
			"POST", "/api/v1/productsearch",
			strings.NewReader(`{"Query":"error"}`), 500,
			`search failed`, true},

		// Income reports
		{"GetIncomeStatisticsHandler_OK",
			"POST", "/api/v1/report/income",
			strings.NewReader(`{"ValidFrom":"2003-04-22T13:00:00Z"}`), 200,
			`{"records":1,"income":"9.42","incomeBilled":"6.28","incomeNotBilled":"3.14"}` + "\n",
			false},
		{"GetIncomeStatisticsHandler_NoBody",
			"POST", "/api/v1/report/income", strings.NewReader(""), 400,
			`json decode error: EOF`, true},
		{"GetIncomeStatisticsHandler_BadJSON",
			"POST", "/api/v1/report/income",
			strings.NewReader(`:-){}{}`),
			400, "json decode error", true},
		{"GetIncomeStatisticsHandler_Error",
			"POST", "/api/v1/report/income",
			strings.NewReader(`{"ValidFrom":"2001-05-30T09:30:10+02:00"}`), 500,
			`report failed`, true},

		// GetPatientHandler tests
		{"GetPatientHandler_OK",
			"GET", "/api/v1/patient/1", nil, 200,
			`{"id":1,"version":0,"creator":"","created":"0001-01-01T00:00:00Z","modifier":"","modified":"0001-01-01T00:00:00Z","name":"pet","birthDate":"0001-01-01T00:00:00Z","speciesId":0,"breedId":0,"genderId":0,"note":"","dead":false,"species":"","breed":"","gender":"","records":[],"tags":[]}` + "\n", false},
		// failed requests tested by GetOwnerHandler tests

		// GetRecordHandler tests
		{"GetRecordHandler_OK",
			"GET", "/api/v1/record/1", nil, 200,
			`{"id":1,"version":0,"creator":"","created":"0001-01-01T00:00:00Z","modifier":"","modified":"0001-01-01T00:00:00Z","date":"0001-01-01T00:00:00Z","text":"","billed":false,"items":[{"id":2,"productId":100,"productPrice":"1.00","amount":"2.00","itemPrice":"3.00","itemType":"Labor","product":"","unit":"","plu":"10"}],"total":""}` + "\n", false},
		// failed requests tested by GetOwnerHandler tests

		// CreateRecordHandler tests
		{"CreateRecordHandler_OK",
			"POST", "/api/v1/record",
			strings.NewReader(`{"patientId":1,"text":"test","billed":true,"items":[{"productId":2,"productPrice":"1.00","amount":"2.00","itemPrice":"3.00","itemType":"Material"}]}`),
			200, "42", false},
		{"CreateRecordHandler_BadJSON",
			"POST", "/api/v1/record",
			strings.NewReader(`:-)`),
			400, "json decode error", true},
		{"CreateRecordHandler_NoBody",
			"POST", "/api/v1/record", strings.NewReader(""), 400,
			`json decode error: EOF`, true},

		// CreatePatientHandler tests
		{"CreatePatientHandler_OK",
			"POST", "/api/v1/patient",
			strings.NewReader(`{"ownerId":2,"name":"test-pet","record":{"text":"ttt","items":[{"productId":2,"productPrice":"1.00","amount":"2.00","itemPrice":"3.00","itemType":"Material"}]}}`),
			200, "42", false},
		{"CreatePatientHandler_BadJSON",
			"POST", "/api/v1/patient",
			strings.NewReader(`:-)`),
			400, "json decode error", true},
		{"CreatePatientHandler_NoBody",
			"POST", "/api/v1/patient", strings.NewReader(""), 400,
			`json decode error: EOF`, true},

		// UpdatePatientHandler tests
		{"UpdatePatientHandler_OK",
			"PUT", "/api/v1/patient/1",
			strings.NewReader(`{"Name":"A"}`),
			200, "", false},
		// failed requests tested by UpdateOwnerHandler tests

		// UpdateRecordHandler tests
		{"UpdateRecordHandler_OK",
			"PUT", "/api/v1/record/1",
			strings.NewReader(`{"Text":"ttt"}`),
			200, "", false},
		// failed requests tested by UpdateOwnerHandler tests

		// GetAllTitlesHandler test
		{"GetAllTitlesHandler_OK",
			"GET", "/api/v1/title", nil, 200,
			`{"items":[{"id":42,"name":"title"}]}` + "\n", false},

		// GetAllUnitsHandler test
		{"GetAllUnitsHandler_OK",
			"GET", "/api/v1/unit", nil, 200,
			`{"items":[{"id":42,"name":"unit"}]}` + "\n", false},

		// GetAllGendersHandler test
		{"GetAllGendersHandler_OK",
			"GET", "/api/v1/gender", nil, 200,
			`{"items":[{"id":42,"name":"gender"}]}` + "\n", false},

		// GetAllSpeciesHandler test
		{"GetAllSpeciesHandler_OK",
			"GET", "/api/v1/species", nil, 200,
			`{"items":[{"id":42,"name":"species"}]}` + "\n", false},

		// GetAllBreedsBySpeciesHandler tests
		{"GetAllBreedsBySpeciesHandler_OK",
			"GET", "/api/v1/breed/by-species/42", nil, 200,
			`{"items":[{"id":42,"name":"breed"}]}` + "\n", false},
		{"GetSpeciesHandler_BadParam",
			"GET", "/api/v1/breed/by-species/Nan", nil, 404,
			"invalid breed ID", true},
		{"GetSpeciesHandler_SqlError",
			"GET", "/api/v1/breed/by-species/2", nil, 500,
			"get by id failed", true},

		// SearchCityHandler tests
		{"SearchCityHandler_OK",
			"GET", "/api/v1/city?q=test", nil, 200,
			`{"total":1,"items":[{"id":1,"name":"city","zip":"000"}]}` + "\n", false},
		{"SearchCityHandler_OK_empty_query",
			"GET", "/api/v1/city", nil, 200,
			`{"total":1,"items":[{"id":1,"name":"city","zip":"000"}]}` + "\n", false},
		{"SearchCityHandler_SqlError",
			"GET", "/api/v1/city?q=fail", nil, 500,
			"search city failed", true},

		// SearchStreetByCityHandler tests
		{"SearchStreetByCityHandler_OK",
			"GET", "/api/v1/street/by-city/1?q=test", nil, 200,
			`{"total":1,"items":[{"id":1,"name":"street","zip":"000"}]}` + "\n", false},
		{"SearchStreetByCityHandler_OK_empty_query",
			"GET", "/api/v1/street/by-city/1", nil, 200,
			`{"total":1,"items":[{"id":1,"name":"street","zip":"000"}]}` + "\n", false},
		{"SearchStreetByCityHandler_BadParam",
			"GET", "/api/v1/street/by-city/Nan", nil, 404,
			"invalid city ID", true},
		{"SearchStreetByCityHandler_SqlError",
			"GET", "/api/v1/street/by-city/1?q=fail", nil, 500,
			"search street failed", true},

		// GetTagHandler tests
		{"GetTagHandler_OK",
			"GET", "/api/v1/tag/1", nil, 200,
			`{"id":1,"version":0,"creator":"","created":"0001-01-01T00:00:00Z","modifier":"","modified":"0001-01-01T00:00:00Z","type":"RFID","value":"007","data":""}` + "\n", false},
		// CreateTagHandler tests
		{"CreateTagHandler_OK",
			"POST", "/api/v1/tag",
			strings.NewReader(`{"patientId":2,"type":"RFID","value":"007", "data": "QUJD"}`),
			200, "42", false},
		{"CreateTagHandler_BadJSON",
			"POST", "/api/v1/tag",
			strings.NewReader(`:-)`),
			400, "json decode error", true},
		{"CreateTagHandler_NoBody",
			"POST", "/api/v1/tag", strings.NewReader(""), 400,
			`json decode error: EOF`, true},

		// UpdateTagHandler tests
		{"UpdateTagHandler_OK",
			"PUT", "/api/v1/tag/1",
			strings.NewReader(`{"version":2, "value":"007", "data": "QUJD"}`),
			200, "", false},
	}

	handler := newHttpHandler()
	for _, tt := range tests {
		req, err := syshttp.NewRequest(tt.reqMethod, tt.reqURL, tt.reqBody)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Authorization", "Bearer: test-token")

		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)

		if tt.expCode != resp.Code {
			t.Fatalf("%s failed. Expected return code %d but was %d",
				tt.name, tt.expCode, resp.Code)
		}

		actBody := resp.Body.String()
		var equals bool

		if tt.useBodyContains {
			equals = strings.Contains(actBody, tt.expBody)
		} else {
			equals = tt.expBody == actBody
		}

		if !equals {
			t.Fatalf("%s failed. Expected \n--%s--\n but was \n--%s--\n",
				tt.name, tt.expBody, actBody)
		}
	}
}
