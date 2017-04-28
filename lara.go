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

package lara

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// COMMON STRUCTURES

type Versioned struct {
	ID      uint64 `json:"id"`
	Version uint64 `json:"version"`
}

type CreatorModifier struct {
	Creator  string    `json:"creator"`
	Created  time.Time `json:"created"`
	Modifier string    `json:"modifier"`
	Modified time.Time `json:"modified"`
}

// -----------------------------------------------------------------------------
// OWNER MANAGEMENT SERVICE

// Owner is JSON encoded updatable owner fields
type Owner struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	TitleID   uint64 `json:"titleId"`
	CityID    uint64 `json:"cityId"`
	StreetID  uint64 `json:"streetId"`
	HouseNo   string `json:"houseNo"`
	Phone1    string `json:"phone1"`
	Phone2    string `json:"phone2"`
	Email     string `json:"email"`
	Note      string `json:"note"`
	IC        string `json:"IC"`
	DIC       string `json:"DIC"`
	ICDPH     string `json:"ICDPH"`
}

// GetOwner is JSON encoded owner data
type GetOwner struct {
	Versioned
	Owner
	Title  string `json:"title"`
	City   string `json:"city"`
	Street string `json:"street"`
	CreatorModifier
	Patients []OwnersPatient `json:"patients"`
}

type OwnersPatient struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Breed   string `json:"breed"`
	Gender  string `json:"gender"`
	Dead    bool   `json:"dead"`
}

// UpdateOwner is JSON encoded updatable owner data
type UpdateOwner struct {
	Version uint64 `json:"version"`
	Owner
}

type CreateOwner struct {
	Owner
	Patient NewPatient `json:"patient"`
}

type NewPatient struct {
	Patient
	Record NewRecord `json:"record"`
}

type NewRecord struct {
	Text   string       `json:"text"`
	Billed bool         `json:"billed"`
	Items  []RecordItem `json:"items"`
}

type RecordItem struct {
	ProductID    uint64         `json:"productId"`
	ProductPrice string         `json:"productPrice"`
	Amount       string         `json:"amount"`
	ItemPrice    string         `json:"itemPrice"`
	ItemType     RecordItemType `json:"itemType"`
}

// OwnerService manages owners
type OwnerService interface {
	Get(ctx context.Context, id uint64) (*GetOwner, error)
	Update(ctx context.Context, id uint64, o *UpdateOwner) error
	Create(ctx context.Context, o *CreateOwner) (uint64, error)
}

// -----------------------------------------------------------------------------
// PATIENT MANAGEMENT SERVICE

// Patient updatable fields
type Patient struct {
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birthDate"`
	SpeciesID uint64    `json:"speciesId"`
	BreedID   uint64    `json:"breedId"`
	GenderID  uint64    `json:"genderId"`
	Note      string    `json:"note"`
}

type GetPatient struct {
	Versioned
	CreatorModifier
	Patient
	Dead    bool             `json:"dead"`
	Species string           `json:"species"`
	Breed   string           `json:"breed"`
	Gender  string           `json:"gender"`
	Records []PatientsRecord `json:"records"`
	Tags    []PatientsTag    `json:"tags"`
}

type PatientsRecord struct {
	ID   uint64    `json:"id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

type PatientsTag struct {
	ID    uint64 `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CreatePatient struct {
	OwnerID uint64 `json:"ownerId"`
	NewPatient
}

type UpdatePatient struct {
	Version uint64 `json:"version"`
	Patient
	Dead bool `json:"dead"`
}

type PatientService interface {
	Get(ctx context.Context, id uint64) (*GetPatient, error)
	Update(ctx context.Context, id uint64, p *UpdatePatient) error
	Create(ctx context.Context, p *CreatePatient) (uint64, error)
}

// -----------------------------------------------------------------------------
// RECORD MANAGEMENT SERVICE

// RecordItemType is enum defining system permissions
//go:generate stringer -type=RecordItemType -output recorditem_string.go
//requires golang.org/x/tools/cmd/stringer installed locally
type RecordItemType int

// RecordItemType enum
const (
	Labor RecordItemType = iota
	Material
)

func (i RecordItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *RecordItemType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "labor":
		*i = Labor
	case "material":
		*i = Material
	default:
		return fmt.Errorf("bad RecordItemType: '%s'", s)
	}

	return nil

}

type GetRecord struct {
	Versioned
	CreatorModifier
	Date   time.Time       `json:"date"`
	Text   string          `json:"text"`
	Billed bool            `json:"billed"`
	Items  []GetRecordItem `json:"items"`
	Total  string          `json:"total"`
}

type GetRecordItem struct {
	ID uint64 `json:"id"`
	RecordItem
	Product string `json:"product"`
	Unit    string `json:"unit"`
	PLU     string `json:"plu"`
}

type CreateRecord struct {
	PatientID uint64 `json:"patientId"`
	NewRecord
}

type UpdateRecord struct {
	Version uint64       `json:"version"`
	Text    string       `json:"text"`
	Items   []RecordItem `json:"items"`
}

type RecordService interface {
	Get(ctx context.Context, id uint64) (*GetRecord, error)
	Update(ctx context.Context, id uint64, r *UpdateRecord) error
	Create(ctx context.Context, r *CreateRecord) (uint64, error)
}

// -----------------------------------------------------------------------------
// RECORD SEARCH SERVICE

// SearchResult is JSON encoded search result structure
type SearchResult struct {
	Total   int            `json:"total"`
	Records []SearchRecord `json:"records"`
}

// SearchRecord is JSON encoded search result
type SearchRecord struct {
	ID      uint64 `json:"id"`      // DB primary key
	Name    string `json:"name"`    // formatted owner name (first, last, title)
	Address string `json:"address"` // owner's address
}

// SearchService searches owners/pets
type SearchService interface {
	Search(ctx context.Context, q string) (*SearchResult, error)
}

// -----------------------------------------------------------------------------
// LIST OF VALUES MANAGEMENT SERVICES

type LOVItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type LOVItemList struct {
	Items []LOVItem `json:"items"`
}

type TitleService interface {
	GetAllTitles(ctx context.Context) (*LOVItemList, error)
}

type UnitService interface {
	GetAllUnits(ctx context.Context) (*LOVItemList, error)
}

type GenderService interface {
	GetAllGenders(ctx context.Context) (*LOVItemList, error)
}

type SpeciesService interface {
	GetAllSpecies(ctx context.Context) (*LOVItemList, error)
}

type BreedService interface {
	GetAllBreedsBySpecies(ctx context.Context, speciesID uint64) (*LOVItemList, error)
}

// City or street
type CityStreet struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	ZIP  string `json:"zip"`
}

// List of cities/streets
type CityStreetList struct {
	Total int          `json:"total"`
	Items []CityStreet `json:"items"`
}

type AddressService interface {
	SearchCity(ctx context.Context, query string) (*CityStreetList, error)
	SearchStreetForCity(ctx context.Context, cityID uint64, query string) (*CityStreetList, error)
}

// -----------------------------------------------------------------------------
// PRODUCT MANAGEMENT SERVICES

// Product is JSON encoded product structure
type Product struct {
	ID    uint64 `json:"id"`    // DB primary key
	Name  string `json:"name"`  // product's name
	Unit  string `json:"unit"`  // product's unit of measure
	Price string `json:"price"` // product's price (formatted decimal, precision: 8.2)
}

// ProductSearchResult is JSON encoded search Product result structure
type ProductSearchResult struct {
	Total    int       `json:"total"`
	Products []Product `json:"products"`
}

// ProductSearchRequest is JSON encoded search Product request structure
type ProductSearchRequest struct {
	ValidTo time.Time `json:"validTo"`
	Query   string    `json:"query"`
}

// ProductService manages products
type ProductService interface {
	Search(ctx context.Context, p *ProductSearchRequest) (*ProductSearchResult, error)
}

// -----------------------------------------------------------------------------
// REPORTING SERVICES

// ReportRequest is JSON encoded request structure for various reports
type ReportRequest struct {
	ValidFrom time.Time `json:"validFrom"`
	ValidTo   time.Time `json:"validTo"`
}

// IncomeStatistics is JSON encoded income statistics
type IncomeStatistics struct {
	Records         int    `json:"records"`         // count
	Income          string `json:"income"`          // currency (formatted decimal, precision: 8.2)
	IncomeBilled    string `json:"incomeBilled"`    // currency
	IncomeNotBilled string `json:"incomeNotBilled"` // currency
}

// ReportService generates data for reports
type ReportService interface {
	GetIncomeStatistics(ctx context.Context, r *ReportRequest) (*IncomeStatistics, error)
}
