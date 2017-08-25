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

package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jkusniar/lara"
)

//go:generate stringer -type=objectType -output obj_string.go
//requires golang.org/x/tools/cmd/stringer installed locally
//if new object type added to enum, run "go generate"
type objectType int

const (
	owner objectType = iota
	patient
	record
	title
	unit
	gender
	city
	street
	species
	breed
	tag
)

func parseID(r *http.Request) (uint64, error) {
	return strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
}

// searchHandler searches owners and pets by name. Parameter is called "q"
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	resp, err := s.SearchService.Search(r.Context(), q)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getOwnerHandler returns JSON formatted GetOwner data by ID
func (s *Server) getOwnerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, owner, err)
		return
	}

	resp, err := s.OwnerService.Get(r.Context(), id)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// createOwnerHandler creates new owner from JSON encoded body of request.
// New record's ID is returned in response body as text
func (s *Server) createOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var o lara.CreateOwner
	if err := render.DecodeJSON(r.Body, &o); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := s.OwnerService.Create(r.Context(), &o)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.PlainText(w, r, fmt.Sprintf("%d", id))
}

// updateOwnerHandler updates existing owner identified by id param. Owner to
// update is JSON encodes in request's body. Result is indicated by response
// status only (204/4xx/5xx).
func (s *Server) updateOwnerHandler(w http.ResponseWriter, r *http.Request) {
	var o lara.UpdateOwner
	if err := render.DecodeJSON(r.Body, &o); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, owner, err)
		return
	}

	if err := s.OwnerService.Update(r.Context(), id, &o); err != nil {
		renderError(w, r, err)
	}
}

// searchProductHandler searches products valid to specified date by name
func (s *Server) searchProductHandler(w http.ResponseWriter, r *http.Request) {
	var p lara.ProductSearchRequest
	if err := render.DecodeJSON(r.Body, &p); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	resp, err := s.ProductService.Search(r.Context(), &p)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getIncomeStatisticsHandler counts records and income for specified time period
func (s *Server) getIncomeStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	var rr lara.ReportRequest
	if err := render.DecodeJSON(r.Body, &rr); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	resp, err := s.ReportService.GetIncomeStatistics(r.Context(), &rr)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getPatientHandler returns JSON formatted GetPatient data by ID
func (s *Server) getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, patient, err)
		return
	}

	resp, err := s.PatientSevice.Get(r.Context(), id)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// searchPatientByTagHandler returns JSON formatted PatientByTag data by tag value
func (s *Server) searchPatientByTagHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "tag")

	resp, err := s.TagService.GetPatientByTag(r.Context(), t)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getRecordHandler returns JSON formatted GetRecord data by ID
func (s *Server) getRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, record, err)
		return
	}

	resp, err := s.RecordService.Get(r.Context(), id)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// createRecordHandler creates new record from JSON encoded body of request.
// New record's ID is returned in response body as text
func (s *Server) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var rec lara.CreateRecord
	if err := render.DecodeJSON(r.Body, &rec); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := s.RecordService.Create(r.Context(), &rec)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.PlainText(w, r, fmt.Sprintf("%d", id))
}

// createPatientHandler creates new patient from JSON encoded body of request.
// New patient's ID is returned in response body as text
func (s *Server) createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var p lara.CreatePatient
	if err := render.DecodeJSON(r.Body, &p); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := s.PatientSevice.Create(r.Context(), &p)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.PlainText(w, r, fmt.Sprintf("%d", id))
}

// updatePatientHandler updates existing patient identified by id param. Patient to
// update is JSON encoded in request's body. Result is indicated by response
// status only (204/4xx/5xx).
func (s *Server) updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	var p lara.UpdatePatient
	if err := render.DecodeJSON(r.Body, &p); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, patient, err)
		return
	}

	if err := s.PatientSevice.Update(r.Context(), id, &p); err != nil {
		renderError(w, r, err)
	}
}

// updateRecordHandler updates existing record identified by id param. Record to
// update is JSON encoded in request's body. Result is indicated by response
// status only (204/4xx/5xx).
func (s *Server) updateRecordHandler(w http.ResponseWriter, r *http.Request) {
	var rec lara.UpdateRecord
	if err := render.DecodeJSON(r.Body, &rec); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, record, err)
		return
	}

	if err := s.RecordService.Update(r.Context(), id, &rec); err != nil {
		renderError(w, r, err)
	}
}

// getAllTitlesHandler returns JSON formatted list of all titles
func (s *Server) getAllTitlesHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.TitleService.GetAllTitles(r.Context())
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getAllTitlesHandler returns JSON formatted list of all titles
func (s *Server) getAllUnitsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.UnitService.GetAllUnits(r.Context())
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getAllGendersHandler returns JSON formatted list of all titles
func (s *Server) getAllGendersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.GenderService.GetAllGenders(r.Context())
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getAllSpeciesHandler returns JSON formatted list of all titles
func (s *Server) getAllSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := s.SpeciesService.GetAllSpecies(r.Context())
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getAllBreedsBySpeciesHandler returns JSON formatted breed data by species ID
func (s *Server) getAllBreedsBySpeciesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, breed, err)
		return
	}

	resp, err := s.BreedService.GetAllBreedsBySpecies(r.Context(), id)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// searchCityHandler returns JSON formated cities queried by name
func (s *Server) searchCityHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	resp, err := s.AddressService.SearchCity(r.Context(), q)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// searchStreetByCityHandler returns JSON formatted street data for city queried by street name
func (s *Server) searchStreetByCityHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	cityID, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, city, err)
		return
	}

	resp, err := s.AddressService.SearchStreetForCity(r.Context(), cityID, q)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// getTagHandler returns JSON formatted GetTag data by ID
func (s *Server) getTagHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, tag, err)
		return
	}

	resp, err := s.TagService.Get(r.Context(), id)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}

// createTagHandler creates new patient's tag from JSON encoded body of request.
// New tag's ID is returned in response body as text
func (s *Server) createTagHandler(w http.ResponseWriter, r *http.Request) {
	var t lara.CreateTag
	if err := render.DecodeJSON(r.Body, &t); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := s.TagService.Create(r.Context(), &t)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.PlainText(w, r, fmt.Sprintf("%d", id))
}

// updateTagHandler updates existing patient's tag identified by id param. Tag to
// update is JSON encoded in request's body. Result is indicated by response
// status only (204/4xx/5xx).
func (s *Server) updateTagHandler(w http.ResponseWriter, r *http.Request) {
	var t lara.UpdateTag
	if err := render.DecodeJSON(r.Body, &t); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	id, err := parseID(r)
	if err != nil {
		renderNotFoundError(w, r, tag, err)
		return
	}

	if err := s.TagService.Update(r.Context(), id, &t); err != nil {
		renderError(w, r, err)
	}
}
