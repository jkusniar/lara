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

type TagService struct {
	GetFn      func(id uint64) (*lara.GetTag, error)
	GetInvoked bool

	UpdateFn      func(id uint64, t *lara.UpdateTag) error
	UpdateInvoked bool

	CreateFn      func(t *lara.CreateTag) (uint64, error)
	CreateInvoked bool

	GetPatientByTagFn      func(tagValue string) (*lara.PatientByTag, error)
	GetPatientByTagInvoked bool
}

func (s *TagService) GetPatientByTag(ctx context.Context, tagValue string) (*lara.PatientByTag, error) {
	s.GetPatientByTagInvoked = true
	return s.GetPatientByTagFn(tagValue)
}

func (s *TagService) Get(ctx context.Context, id uint64) (*lara.GetTag, error) {
	s.GetInvoked = true
	return s.GetFn(id)
}

func (s *TagService) Update(ctx context.Context, id uint64, t *lara.UpdateTag) error {
	s.UpdateInvoked = true
	return s.UpdateFn(id, t)
}

func (s *TagService) Create(ctx context.Context, t *lara.CreateTag) (uint64, error) {
	s.CreateInvoked = true
	return s.CreateFn(t)
}
