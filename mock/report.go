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

// ReportService mock implementation
type ReportService struct {
	GetIncomeStatisticsFn func(
		r *lara.ReportRequest) (*lara.IncomeStatistics, error)
	GetIncomeStatisticsInvoked bool
}

// GetIncomeStatistics mock implementation
func (s *ReportService) GetIncomeStatistics(ctx context.Context,
	r *lara.ReportRequest) (*lara.IncomeStatistics, error) {
	s.GetIncomeStatisticsInvoked = true
	return s.GetIncomeStatisticsFn(r)
}
