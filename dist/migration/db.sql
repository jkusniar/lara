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

-- USER TABLE
CREATE TABLE "user" (
     id SERIAL PRIMARY KEY,
     login character varying(20) NOT NULL UNIQUE,
     pass_salt bytea NOT NULL,
     pass_hash bytea NOT NULL
);

CREATE TABLE permission (
     id SERIAL PRIMARY KEY,
     name character varying(20) NOT NULL UNIQUE
);

CREATE TABLE user_permission (
     user_id integer NOT NULL REFERENCES "user",
     permission_id integer NOT NULL REFERENCES permission,
     UNIQUE(user_id, permission_id)
);

-- PATIENTS NAME REQUIRED
ALTER TABLE patient ALTER COLUMN "name" SET NOT NULL;

-- VERSION
ALTER TABLE client ADD COLUMN version integer NOT NULL DEFAULT 0;
ALTER TABLE patient ADD COLUMN version integer NOT NULL DEFAULT 0;
ALTER TABLE record ADD COLUMN version integer NOT NULL DEFAULT 0;

-- CREATOR/MODIFIER
-- Creator of legacy records is "tara"
ALTER TABLE client ADD COLUMN creator character varying(20);
UPDATE client set creator = 'tara';
ALTER TABLE client ALTER COLUMN "creator" SET NOT NULL;
-- TODO: PROPER "CREATED" MIGRATION BASED ON OLDEST RECORD OF CLIENT
ALTER TABLE client ADD COLUMN created TIMESTAMP;
UPDATE client set created = current_timestamp;
ALTER TABLE client ALTER COLUMN "created" SET NOT NULL;
ALTER TABLE client ADD COLUMN modifier character varying(20);
ALTER TABLE client ADD COLUMN modified TIMESTAMP;

ALTER TABLE patient ADD COLUMN creator character varying(20);
UPDATE patient set creator = 'tara';
ALTER TABLE patient ALTER COLUMN "creator" SET NOT NULL;
-- TODO: PROPER "CREATED" MIGRATION BASED ON OLDEST RECORD OF PATIENT
ALTER TABLE patient ADD COLUMN created TIMESTAMP;
UPDATE patient set created = current_timestamp;
ALTER TABLE patient ALTER COLUMN "created" SET NOT NULL;
ALTER TABLE patient ADD COLUMN modifier character varying(20);
ALTER TABLE patient ADD COLUMN modified TIMESTAMP;

ALTER TABLE record ADD COLUMN creator character varying(20);
UPDATE record set creator = 'tara';
ALTER TABLE record ALTER COLUMN "creator" SET NOT NULL;
-- TODO: PROPER "CREATED" MIGRATION BASED ON REC_DATE
ALTER TABLE record ADD COLUMN created TIMESTAMP;
UPDATE record set created = current_timestamp;
ALTER TABLE record ALTER COLUMN "created" SET NOT NULL;
ALTER TABLE record ADD COLUMN modifier character varying(20);
ALTER TABLE record ADD COLUMN modified TIMESTAMP;

ALTER TABLE client RENAME TO owner;
ALTER SEQUENCE client_id_seq RENAME TO owner_id_seq;
ALTER TABLE patient RENAME COLUMN client_id TO owner_id;

-- Rename SEX to GENDER
ALTER TABLE lov_sex RENAME TO lov_gender;
ALTER SEQUENCE lov_sex_id_seq RENAME TO lov_gender_id_seq;
ALTER TABLE patient RENAME COLUMN sex_id TO gender_id;

-- unique constraints to LOVs
--
delete FROM lov_street
WHERE  ctid NOT IN (
  SELECT min(ctid)
  FROM   lov_street
  GROUP  BY city_id, street, psc); -- removes duplicate street entries
-- fix duplicate products
UPDATE lov_product SET name = name || ' XXX' WHERE  ctid NOT IN (SELECT max(ctid) FROM lov_product GROUP  BY name, unit_id);

ALTER TABLE lov_title ADD CONSTRAINT lov_title_uq UNIQUE (name);
ALTER TABLE lov_species ADD CONSTRAINT lov_species_uq UNIQUE (name);
ALTER TABLE lov_breed ADD CONSTRAINT lov_breed_uq UNIQUE (lov_species_id, name);
ALTER TABLE lov_city ADD CONSTRAINT lov_city_uq UNIQUE (city, psc);
ALTER TABLE lov_street ADD CONSTRAINT lov_street_uq UNIQUE (city_id, street, psc);
ALTER TABLE lov_gender ADD CONSTRAINT lov_gender_uq UNIQUE (name);
ALTER TABLE lov_unit ADD CONSTRAINT lov_unit_uq UNIQUE (name);
ALTER TABLE lov_product ADD CONSTRAINT lov_product_uq UNIQUE (name, unit_id);
ALTER TABLE lov_phrase ADD CONSTRAINT lov_phrase_uq UNIQUE (name);

-- TAGS
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
  creator character varying(20) NOT NULL,
  created TIMESTAMP NOT NULL,
  modifier character varying(20),
  modified TIMESTAMP,
  version integer NOT NULL DEFAULT 0
);
CREATE INDEX "idx_tag$patient_id" ON tag USING btree (patient_id);

INSERT INTO tag_type (name) VALUES ('LyssaVirus');
UPDATE patient SET lysset = lysset || '-X' WHERE  ctid NOT IN (SELECT max(ctid) FROM patient GROUP  BY lysset);
DO language plpgsql $$
DECLARE
  lyss RECORD;
BEGIN
  FOR lyss IN select id, lysset FROM patient WHERE lysset IS NOT NULL
  LOOP
    EXECUTE 'INSERT INTO tag (value, patient_id, tag_type_id, creator, created) VALUES (''' || lyss.lysset ||  ''', ' || lyss.id || ', 1, ''tara'', current_timestamp)';
  END LOOP;
END
$$;

ALTER TABLE patient DROP COLUMN lysset;

-- VARCHAR -> TEXT conversion
ALTER TABLE lov_title ALTER COLUMN name TYPE text;
ALTER TABLE lov_species ALTER COLUMN name TYPE text;
ALTER TABLE lov_breed ALTER COLUMN name TYPE text;
ALTER TABLE lov_gender ALTER COLUMN name TYPE text;
ALTER TABLE lov_unit ALTER COLUMN name TYPE text;
ALTER TABLE lov_phrase ALTER COLUMN name TYPE text;
ALTER TABLE permission ALTER COLUMN name TYPE text;
ALTER TABLE "user" ALTER COLUMN login TYPE text;
ALTER TABLE "user" ADD CONSTRAINT login_chk CHECK (length(login) <= 20);
ALTER TABLE owner ALTER COLUMN creator TYPE text;
ALTER TABLE owner ADD CONSTRAINT owner_creator_chk CHECK (length(creator) <= 20);
ALTER TABLE owner ALTER COLUMN modifier TYPE text;
ALTER TABLE owner ADD CONSTRAINT owner_modifier_chk CHECK (length(modifier) <= 20);
ALTER TABLE patient ALTER COLUMN creator TYPE text;
ALTER TABLE patient ADD CONSTRAINT patient_creator_chk CHECK (length(creator) <= 20);
ALTER TABLE patient ALTER COLUMN modifier TYPE text;
ALTER TABLE patient ADD CONSTRAINT patient_modifier_chk CHECK (length(modifier) <= 20);
ALTER TABLE tag ALTER COLUMN creator TYPE text;
ALTER TABLE tag ADD CONSTRAINT tag_creator_chk CHECK (length(creator) <= 20);
ALTER TABLE tag ALTER COLUMN modifier TYPE text;
ALTER TABLE tag ADD CONSTRAINT tag_modifier_chk CHECK (length(modifier) <= 20);
ALTER TABLE record ALTER COLUMN creator TYPE text;
ALTER TABLE record ADD CONSTRAINT record_creator_chk CHECK (length(creator) <= 20);
ALTER TABLE record ALTER COLUMN modifier TYPE text;
ALTER TABLE record ADD CONSTRAINT record_modifier_chk CHECK (length(modifier) <= 20);
ALTER TABLE lov_product ALTER COLUMN name TYPE text;
ALTER TABLE owner ALTER COLUMN first_name TYPE text;
ALTER TABLE owner ALTER COLUMN last_name TYPE text;
ALTER TABLE owner ALTER COLUMN phone_1 TYPE text;
ALTER TABLE owner ALTER COLUMN phone_2 TYPE text;
ALTER TABLE owner ALTER COLUMN email TYPE text;
ALTER TABLE owner ALTER COLUMN house_no TYPE text;
ALTER TABLE owner ALTER COLUMN note TYPE text;
ALTER TABLE owner ALTER COLUMN ic TYPE text;
ALTER TABLE owner ALTER COLUMN dic TYPE text;
ALTER TABLE owner ALTER COLUMN icdph TYPE text;
ALTER TABLE patient ALTER COLUMN note TYPE text;
ALTER TABLE patient ALTER COLUMN name TYPE text;
ALTER TABLE lov_city ALTER COLUMN city TYPE text;
ALTER TABLE lov_city ALTER COLUMN district TYPE text;
ALTER TABLE lov_city ALTER COLUMN province TYPE text;
ALTER TABLE lov_city ALTER COLUMN psc TYPE text;
ALTER TABLE lov_city ALTER COLUMN post_office TYPE text;
ALTER TABLE lov_street ALTER COLUMN street TYPE text;
ALTER TABLE lov_street ALTER COLUMN psc TYPE text;
ALTER TABLE lov_street ALTER COLUMN post_office TYPE text;

-- Character(1) -> Boolean
ALTER TABLE patient ADD COLUMN dead BOOLEAN DEFAULT FALSE;
UPDATE patient SET dead = TRUE WHERE is_dead = '1';
ALTER TABLE patient ALTER COLUMN dead SET NOT NULL;
ALTER TABLE patient DROP COLUMN is_dead;

ALTER TABLE record ADD COLUMN billed_2 BOOLEAN;
UPDATE record SET billed_2 = TRUE WHERE billed = '1';
UPDATE record SET billed_2 = FALSE WHERE billed_2 IS NULL;
ALTER TABLE record ALTER COLUMN billed_2 SET NOT NULL;
ALTER TABLE record DROP COLUMN billed;
ALTER TABLE record RENAME COLUMN billed_2 TO billed;