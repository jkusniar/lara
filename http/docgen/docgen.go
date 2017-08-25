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

package main

import (
	"io/ioutil"

	"github.com/go-chi/docgen"
	"github.com/jkusniar/lara/http"
)

// generates http routes doc - routes.md
// intended to be run by "go generate"
func main() {
	srv := &http.Server{
		WWWRoot: "static",
	}

	// Markdown docs
	ioutil.WriteFile("routes.md", []byte(docgen.MarkdownRoutesDoc(srv.Router(),
		docgen.MarkdownOpts{
			ProjectPath:        "github.com/jkusniar/lara",
			Intro:              "LARA REST API.",
			ForceRelativeLinks: true,
		})), 0666)

	// JSON docs
	ioutil.WriteFile("routes.json",
		[]byte(docgen.JSONRoutesDoc(srv.Router())), 0666)
}
