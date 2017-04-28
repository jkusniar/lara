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

package mock

import (
	"context"

	"github.com/jkusniar/lara"
)

// ProductService mock implementation
type ProductService struct {
	SearchFn func(p *lara.ProductSearchRequest) (*lara.ProductSearchResult,
		error)
	SearchInvoked bool
}

// Search mock implementation
func (s *ProductService) Search(ctx context.Context, p *lara.ProductSearchRequest) (*lara.ProductSearchResult,
	error) {
	s.SearchInvoked = true
	return s.SearchFn(p)
}
