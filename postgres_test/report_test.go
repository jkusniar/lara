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
	"testing"
	"time"

	"github.com/jkusniar/lara"
)

func TestGetIncomeStatisticsEmpty(t *testing.T) {
	report, err := reportService.GetIncomeStatistics(testCtx, &lara.ReportRequest{})
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if report.Records != 0 {
		t.Fatalf("expected record count 0 but was %d", report.Records)
	}

	if report.Income != "0.00" {
		t.Fatalf("expected Income 0 but was %s", report.Income)
	}

	if report.IncomeBilled != "0.00" {
		t.Fatalf("expected IncomeBilled billed 0 but was %s", report.IncomeBilled)
	}

	if report.IncomeNotBilled != "0.00" {
		t.Fatalf("expected IncomeNotBilled 0 but was %s", report.IncomeNotBilled)
	}
}

func TestGetIncomeStatistics(t *testing.T) {
	const timeFormat = "Jan 2, 2006 15:04:05 (MST)"
	from, _ := time.Parse(timeFormat, "Apr 21, 2003 22:00:00 (UTC)")
	to, _ := time.Parse(timeFormat, "Feb 4, 2017 22:59:59 (UTC)")

	report, err := reportService.GetIncomeStatistics(testCtx, &lara.ReportRequest{
		ValidFrom: from,
		ValidTo:   to,
	})

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if report.Records != 3 {
		t.Fatalf("expected record count 0 but was %d", report.Records)
	}

	if report.Income != "9.42" {
		t.Fatalf("expected Income 9.42 but was %s", report.Income)
	}

	if report.IncomeBilled != "6.28" {
		t.Fatalf("expected IncomeBilled billed 6.28 but was %s", report.IncomeBilled)
	}

	if report.IncomeNotBilled != "3.14" {
		t.Fatalf("expected IncomeNotBilled 3.14 but was %s", report.IncomeNotBilled)
	}
}
