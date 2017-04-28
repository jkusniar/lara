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
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/pressly/chi/render"
)

func renderBadJSONError(w http.ResponseWriter, r *http.Request, err error) {
	e := errors.Wrap(err, "json decode error")
	logErr(e)

	render.Status(r, http.StatusBadRequest)
	render.PlainText(w, r, e.Error())
}

func renderNotFoundError(w http.ResponseWriter, r *http.Request, objType fmt.Stringer, err error) {
	e := errors.Wrapf(err, "invalid %s ID", objType)
	logErr(e)

	render.Status(r, http.StatusNotFound)
	render.PlainText(w, r, e.Error())
}

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	logErr(err)

	render.Status(r, getErrorCode(err))
	render.PlainText(w, r, err.Error())
}

// getErrorCode extracts HTTP error code from err.
func getErrorCode(err error) int {
	type isCoded interface {
		Code() int
	}

	code := 500
	if ce, ok := err.(isCoded); ok {
		code = ce.Code()
	}
	return code
}

func logErr(err error) {
	log.Printf("ERROR: request failed: %+v\n", err)
}
