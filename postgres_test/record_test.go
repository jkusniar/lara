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

	"github.com/jkusniar/lara"
)

func TestGetRecord(t *testing.T) {
	var r *lara.GetRecord
	var err error

	// get not existing
	r, err = recordService.Get(testCtx, 10000)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get OK
	r, err = recordService.Get(testCtx, 1)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if r == nil {
		t.Fatal("expected not nil result")
	}
	if r.Text != "RECORD" ||
		r.Total != "6.15" ||
		len(r.Items) != 2 ||
		r.Items[0].ItemType != lara.Material || r.Items[0].PLU != "10" {
		t.Fatalf("unexpected result %+v", r)
	}

	// Total == 0.00 (no items on record)
	r, err = recordService.Get(testCtx, 7)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if r == nil {
		t.Fatal("expected not nil result")
	}
	if r.Total != "0.00" {
		t.Fatalf("unexpected result %+v", r)
	}
}

func TestCreateRecord(t *testing.T) {
	var err error
	r := &lara.CreateRecord{
		PatientID: 1,
		NewRecord: lara.NewRecord{
			Text:   "test",
			Billed: true,
			Items: []lara.RecordItem{
				{ProductID: 3, Amount: "1.0000", ItemPrice: "2.00", ProductPrice: "2.00", ItemType: lara.Labor},
				{ProductID: 3, Amount: "2.0000", ItemPrice: "2.10", ProductPrice: "1.05", ItemType: lara.Material},
			},
		},
	}

	id, err := recordService.Create(testCtx, r)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if id == 0 {
		t.Fatal("incorrect ID returned")
	}
}

func TestUpdateRecord(t *testing.T) {
	var err error

	// update OK (with delete all items)
	u := &lara.UpdateRecord{Version: 0, Text: "updated-text"}
	err = recordService.Update(testCtx, 2, u)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	chk, err := recordService.Get(testCtx, 2)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if len(chk.Items) != 0 {
		t.Fatalf("expected zero items, but was %d", len(chk.Items))
	}

	// update OK (with one new item)
	u = &lara.UpdateRecord{Version: 1, Text: "updated-text2",
		Items: []lara.RecordItem{
			{ProductID: 3, Amount: "1.0000", ItemPrice: "2.00", ProductPrice: "2.00", ItemType: lara.Labor},
		}}
	err = recordService.Update(testCtx, 2, u)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	chk, err = recordService.Get(testCtx, 2)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if len(chk.Items) != 1 {
		t.Fatalf("expected zero items, but was %d", len(chk.Items))
	}

	// update not existing
	err = recordService.Update(testCtx, 1000, u)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// update without version upgrade
	err = recordService.Update(testCtx, 2, u)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 409); !ok {
		t.Fatalf("expected error code 409 but was %d, %+v", actual, err)
	}

	// update - required field missing
	err = recordService.Update(testCtx, 2, &lara.UpdateRecord{Version: 1, Text: "updated-text2",
		Items: []lara.RecordItem{{}}})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}
}
