SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;
SET default_with_oids = false;

CREATE TABLE lov_title (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE lov_species (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE lov_breed (
    id SERIAL PRIMARY KEY,
    lov_species_id integer NOT NULL REFERENCES lov_species,
    name TEXT NOT NULL UNIQUE,
    UNIQUE(lov_species_id, name)
);

CREATE TABLE lov_city (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    district TEXT NOT NULL,
    province TEXT NOT NULL,
    psc TEXT,
    post_office TEXT,
    UNIQUE(city, psc)
);

CREATE TABLE lov_street (
  id SERIAL PRIMARY KEY,
  city_id integer NOT NULL REFERENCES lov_city,
  street TEXT NOT NULL,
  psc TEXT NOT NULL,
  post_office TEXT NOT NULL,
  UNIQUE(city_id, street, psc)
);

CREATE TABLE lov_gender (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE lov_unit (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE lov_product (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  unit_id integer NOT NULL REFERENCES lov_unit,
  price numeric(8,2) NOT NULL,
  valid_to date,
  plu integer,
  UNIQUE(name, unit_id)
);

CREATE TABLE lov_phrase (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  phrase_text text
);

CREATE TABLE owner (
  id SERIAL PRIMARY KEY,
  first_name TEXT,
  last_name TEXT NOT NULL,
  title_id integer REFERENCES lov_title,
  phone_1 TEXT,
  phone_2 TEXT,
  email TEXT,
  city_id integer REFERENCES lov_city,
  street_id integer REFERENCES lov_street,
  house_no TEXT,
  note TEXT,
  ic TEXT,
  dic TEXT,
  icdph TEXT,
  creator TEXT CHECK (length(creator) <= 20) NOT NULL,
  created TIMESTAMP NOT NULL,
  modifier TEXT CHECK (length(modifier) <= 20),
  modified TIMESTAMP,
  version integer NOT NULL DEFAULT 0
);

CREATE TABLE patient (
    id SERIAL PRIMARY KEY,
    owner_id integer NOT NULL REFERENCES owner,
    name TEXT NOT NULL,
    birth_date date,
    species_id integer REFERENCES lov_species,
    breed_id integer REFERENCES lov_breed,
    gender_id integer REFERENCES lov_gender,
    note TEXT,
    dead BOOLEAN NOT NULL DEFAULT FALSE,
    creator TEXT CHECK (length(creator) <= 20) NOT NULL,
    created TIMESTAMP NOT NULL,
    modifier TEXT CHECK (length(modifier) <= 20),
    modified TIMESTAMP,
    version integer NOT NULL DEFAULT 0
);

CREATE TABLE tag_type (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE tag (
  id SERIAL PRIMARY KEY,
  patient_id integer NOT NULL REFERENCES patient,
  tag_type_id integer NOT NULL REFERENCES tag_type,
  value TEXT NOT NULL UNIQUE,
  data bytea,
  creator TEXT CHECK (length(creator) <= 20) NOT NULL,
  created TIMESTAMP NOT NULL,
  modifier TEXT CHECK (length(modifier) <= 20),
  modified TIMESTAMP,
  version integer NOT NULL DEFAULT 0
);

CREATE TABLE record (
    id SERIAL PRIMARY KEY,
    patient_id integer NOT NULL REFERENCES patient,
    rec_date timestamp without time zone NOT NULL,
    data text,
    invoice_id character varying(10),
    payed date,
    inv_delivery_date date,
    inv_create_date date,
    inv_payment_date date,
    billed BOOLEAN NOT NULL,
    creator TEXT CHECK (length(creator) <= 20) NOT NULL,
    created TIMESTAMP NOT NULL,
    modifier TEXT CHECK (length(modifier) <= 20),
    modified TIMESTAMP,
    version integer NOT NULL DEFAULT 0
);

CREATE TABLE record_item (
    id SERIAL PRIMARY KEY,
    record_id integer NOT NULL REFERENCES record,
    prod_id integer NOT NULL REFERENCES lov_product,
    amount numeric(10,4) NOT NULL,
    item_price numeric(8,2) NOT NULL,
    prod_price numeric(8,2) NOT NULL,
    item_type integer NOT NULL
);

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE CHECK (length(login) <= 20),
    pass_salt bytea NOT NULL,
    pass_hash bytea NOT NULL
);

CREATE TABLE permission (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE user_permission (
    user_id integer NOT NULL REFERENCES "user",
    permission_id integer NOT NULL REFERENCES permission,
    UNIQUE(user_id, permission_id)
);

CREATE INDEX "idx_owner$last_name" ON owner USING btree (last_name);
CREATE INDEX "idx_lov_breed$lov_species_id" ON lov_breed USING btree (lov_species_id);
CREATE INDEX "idx_lov_street$city_id" ON lov_street USING btree (city_id);
CREATE INDEX "idx_patient$owner_id" ON patient USING btree (owner_id);
CREATE INDEX "idx_record$invoice_id" ON record USING btree (invoice_id);
CREATE INDEX "idx_record$patient_id" ON record USING btree (patient_id);
CREATE INDEX "idx_record_item$record_id" ON record_item USING btree (record_id);
CREATE INDEX "idx_tag$patient_id" ON tag USING btree (patient_id);

