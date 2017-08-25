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
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/render"
	"github.com/jkusniar/lara"
)

// Server server REST API and client web application
type Server struct {
	srv *http.Server

	// Services
	OwnerService   lara.OwnerService
	PatientSevice  lara.PatientService
	RecordService  lara.RecordService
	SearchService  lara.SearchService
	UserService    lara.UserService
	ProductService lara.ProductService
	ReportService  lara.ReportService
	TitleService   lara.TitleService
	UnitService    lara.UnitService
	GenderService  lara.GenderService
	SpeciesService lara.SpeciesService
	BreedService   lara.BreedService
	AddressService lara.AddressService
	TagService     lara.TagService

	// Auth
	Token AuthToken

	// WWW root
	WWWRoot string
}

//go:generate go run docgen/docgen.go
// Run this every time routes are changed, to regenerate docs.

// Router configures and returns URL router for http server.
// Method is not intended to be called directly except for unit tests.
func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Consider using middleware.Throttle to limit maximum concurrent requests (DDOS).
	// At least for public routes.

	// serve web app
	fileServer(r, "/", http.Dir(s.WWWRoot))

	//heartbeat
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, ".")
	})

	r.Post("/login", s.authenticationHandler)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(s.requireAuthorizedUser)

		// owner
		r.Route("/owner", func(r chi.Router) {
			r.With(requirePermission(lara.EditRecord)).Post("/", s.createOwnerHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.With(requirePermission(lara.ViewRecord)).Get("/", s.getOwnerHandler)
				r.With(requirePermission(lara.EditRecord)).Put("/", s.updateOwnerHandler)
			})
		})

		// patient
		r.Route("/patient", func(r chi.Router) {
			r.With(requirePermission(lara.EditRecord)).Post("/", s.createPatientHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.With(requirePermission(lara.ViewRecord)).Get("/", s.getPatientHandler)
				r.With(requirePermission(lara.EditRecord)).Put("/", s.updatePatientHandler)
			})
		})

		// record
		r.Route("/record", func(r chi.Router) {
			r.With(requirePermission(lara.EditRecord)).Post("/", s.createRecordHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.With(requirePermission(lara.ViewRecord)).Get("/", s.getRecordHandler)
				r.With(requirePermission(lara.EditRecord)).Put("/", s.updateRecordHandler)
			})
		})

		// tags
		r.Route("/tag", func(r chi.Router) {
			r.With(requirePermission(lara.EditRecord)).Post("/", s.createTagHandler)
			r.Route("/{id}", func(r chi.Router) {
				r.With(requirePermission(lara.ViewRecord)).Get("/", s.getTagHandler)
				r.With(requirePermission(lara.EditRecord)).Put("/", s.updateTagHandler)
			})
		})

		// List Of Values
		r.With(requirePermission(lara.ViewRecord)).Get("/title", s.getAllTitlesHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/unit", s.getAllUnitsHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/gender", s.getAllGendersHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/species", s.getAllSpeciesHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/breed/by-species/{id}",
			s.getAllBreedsBySpeciesHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/city", s.searchCityHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/street/by-city/{id}",
			s.searchStreetByCityHandler)

		// TODO: move to owner/patient subroutes as "../search"
		// search
		r.With(requirePermission(lara.ViewRecord)).Get("/search", s.searchHandler)
		r.With(requirePermission(lara.ViewRecord)).Get("/search/patient-by-tag/{tag}",
			s.searchPatientByTagHandler)

		// reports
		r.With(requirePermission(lara.ViewReports)).Post("/report/income", s.getIncomeStatisticsHandler)

		// products
		r.With(requirePermission(lara.ViewRecord)).Post("/productsearch", s.searchProductHandler)
	})

	return r
}

// Serve starts TLS server.
// Method blocks until server stopped from another goroutine or error occurs.
func (s *Server) Serve(Hostname string, Port uint, CertFile string, KeyFile string) error {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", Hostname, strconv.FormatUint(uint64(Port), 10)),
		Handler: s.Router()}

	if err := s.srv.ListenAndServeTLS(CertFile, KeyFile); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown stops server gracefully
func (s *Server) Shutdown() error {
	log.Println("Server going down...")
	return s.srv.Shutdown(context.Background())
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// Stolen from go-chi _examples folder
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
