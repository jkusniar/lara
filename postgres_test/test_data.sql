-- test user
INSERT INTO "user" (login, pass_salt, pass_hash)
VALUES ('test', decode('11', 'hex'), decode('12', 'hex'));

INSERT INTO permission (name) VALUES ('EditRecord');
INSERT INTO permission (name) VALUES ('ViewReports');

-- permissinons to test user
INSERT INTO user_permission (user_id, permission_id) VALUES (1, 1);
INSERT INTO user_permission (user_id, permission_id) VALUES (1, 2);

-- LOVs
-----------------------------------------------------------------------------
--id=1
INSERT INTO lov_title (name) VALUES ('Ing.');

--id=1
INSERT INTO lov_unit (name) VALUES ('ml.');
--id=2
INSERT INTO lov_unit (name) VALUES ('tbl.');

--id=1
INSERT INTO lov_product (NAME, UNIT_ID, PRICE, VALID_TO, PLU)
VALUES ('Vyšetrenie uší a čistenie', 1, 7.00, to_date('31 Mar 2015', 'DD Mon YYYY'), '10');
--id=2
INSERT INTO lov_product (NAME, UNIT_ID, PRICE, VALID_TO)
VALUES ('Vyšetrenie 2', 1, 5.00, current_date);
--id=3
INSERT INTO lov_product (NAME, UNIT_ID, PRICE, PLU)
VALUES ('Vystavenie potvrdenia o zdravotnom stave psa', 2, 4.00, 42);

-- id=1
INSERT INTO lov_city (city,
                      district,
                      province,
                      psc,
                      post_office) VALUES ('test city', 'test discrict', 'TP', '99900', 'test post office');

-- id =2
INSERT INTO lov_city (city,
                      district,
                      province,
                      psc,
                      post_office) VALUES ('test city 2', 'test discrict', 'TP', '11100', 'test post office');

-- id =3
INSERT INTO lov_city (city,
                      district,
                      province,
                      psc,
                      post_office) VALUES ('not to be found', 'test discrict', 'TP', '99900', 'test post office');

-- id=1
INSERT INTO lov_street (city_id,
                        street,
                        psc,
                        post_office) VALUES (1, 'test street', '88800', 'street post office');

-- id=2
INSERT INTO lov_street (city_id,
                        street,
                        psc,
                        post_office) VALUES (1, 'another street', '88800', 'street post office');

-- id=3
INSERT INTO lov_street (city_id,
                        street,
                        psc,
                        post_office) VALUES (1, 'not found str', '88800', 'street post office');

-- id=1
INSERT INTO lov_species (name) VALUES ('dog');
-- id=2
INSERT INTO lov_species (name) VALUES ('cat');

-- id=1
INSERT INTO lov_breed (lov_species_id, name) VALUES (1, 'german shepard');
-- id=2
INSERT INTO lov_breed (lov_species_id, name) VALUES (1, 'german boxer');
-- id=3
INSERT INTO lov_breed (lov_species_id, name) VALUES (2, 'black cat');
-- id=4
INSERT INTO lov_breed (lov_species_id, name) VALUES (2, 'white cat');

-- id=1
INSERT INTO lov_gender (name) VALUES ('male');
-- id=2
INSERT INTO lov_gender (name) VALUES ('female');

-- id=1
INSERT INTO tag_type (name) VALUES ('LyssaVirus');
-- id=2
INSERT INTO tag_type (name) VALUES ('RFID');

-- GENERIC DATA
-----------------------------------------------------------------------------

-- OWNERS
-- id=1
INSERT INTO owner (first_name, last_name, creator, created) VALUES ('Get', 'Owner', 'testuser', current_timestamp);
-- id=2
INSERT INTO owner (first_name, last_name, version, creator, created)
VALUES ('To', 'Update', 42, 'testuser', current_timestamp);
-- id=3
INSERT INTO owner (first_name, last_name, creator, created)
VALUES ('Search', 'Test only', 'testuser', current_timestamp);
-- id=4
INSERT INTO owner (last_name, creator, created) VALUES ('Only Last Name', 'testuser', current_timestamp);

-- PATIENTS
-- id=1
INSERT INTO patient (owner_id, name, creator, created) VALUES (1, 'test-pet', 'testuser', current_timestamp);
-- id=1
INSERT INTO record (patient_id, rec_date, billed, creator, created, data) VALUES
  (1, to_timestamp('21 Apr 2003 23:58:00', 'DD Mon YYYY HH24:MI:SS'), true, 'testuser', current_timestamp, 'RECORD');
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (1, 1, 1.0, 3.14, 3.14, 1);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (1, 2, 1.0, 3.01, 3.01, 1);
-- id=1
INSERT INTO tag (value, patient_id, tag_type_id, creator, created) VALUES ('2017-SK-0007', 1, 1, 'testuser', current_timestamp);
-- id=2
INSERT INTO record (patient_id, rec_date, billed, creator, created, data) VALUES
  (1, to_timestamp('21 Apr 2003 23:50:00', 'DD Mon YYYY HH24:MI:SS'), true, 'testuser', current_timestamp, 'FOR-UPDATE');
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (2, 1, 1.0, 3.14, 3.14, 1);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (2, 1, 1.0, 3.01, 3.01, 1);
-- id=3
INSERT INTO record (patient_id, rec_date, billed, creator, created)
VALUES (1, to_timestamp('22 Apr 2003 00:01:00', 'DD Mon YYYY HH24:MI:SS'), true, 'testuser', current_timestamp);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (3, 1, 1.0, 3.14, 3.14, 0);
-- id=4
INSERT INTO record (patient_id, rec_date, billed, creator, created)
VALUES (1, to_timestamp('05 Jun 2010 12:00:00', 'DD Mon YYYY HH24:MI:SS'), true, 'testuser', current_timestamp);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (4, 1, 1.0, 3.14, 3.14, 0);
-- id=5
INSERT INTO record (patient_id, rec_date, billed, creator, created)
VALUES (1, to_timestamp('04 Feb 2017 23:58:00', 'DD Mon YYYY HH24:MI:SS'), false, 'testuser', current_timestamp);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (5, 1, 1.0, 3.14, 3.14, 0);
-- id=6
INSERT INTO record (patient_id, rec_date, billed, creator, created)
VALUES (1, to_timestamp('05 Feb 2017 00:01:00', 'DD Mon YYYY HH24:MI:SS'), false, 'testuser', current_timestamp);
INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
VALUES (6, 1, 1.0, 3.14, 3.14, 0);
-- id=7
INSERT INTO record (patient_id, rec_date, billed, creator, created)
VALUES (1, to_timestamp('05 Apr 2017 00:01:00', 'DD Mon YYYY HH24:MI:SS'), false, 'testuser', current_timestamp);

-- DATA FOR CORE TESTS
-----------------------------------------------------------------------------


-- GetOwner
-- id=5
INSERT INTO owner (first_name, last_name,
                   title_id,
                   phone_1,
                   phone_2,
                   email,
                   city_id,
                   street_id,
                   house_no,
                   note,
                   ic,
                   dic,
                   icdph,
                   creator,
                   created,
                   modifier,
                   modified,
                   version)
VALUES ('Test', 'GetOwner', 1, '000111', '000222', 'test@test.com', 1, 1, '1B', 'test note', '0001', '0002', '0003',
        'testuser', current_timestamp, 'testmodifier', current_timestamp, 3);
-- id=2
INSERT INTO patient (owner_id,
                     name,
                     birth_date,
                     species_id,
                     breed_id,
                     gender_id,
                     note,
                     dead,
                     creator,
                     created,
                     modifier,
                     modified,
                     version) VALUES
  (5, 'get-owner-pet', current_date, 1, 1, 1, 'test-note', false, 'testuser', current_timestamp, 'testmodifier',
   current_timestamp, 5);

-- CreateTag/UpdateTag
-- id=6
INSERT INTO owner (first_name, last_name,
                   title_id,
                   phone_1,
                   phone_2,
                   email,
                   city_id,
                   street_id,
                   house_no,
                   note,
                   ic,
                   dic,
                   icdph,
                   creator,
                   created,
                   modifier,
                   modified,
                   version)
VALUES ('Test', 'CreateTag', 1, '000111', '000222', 'test@test.com', 1, 1, '1B', 'test note', '0001', '0002', '0003',
        'testuser', current_timestamp, 'testmodifier', current_timestamp, 3);
-- id=3
INSERT INTO patient (owner_id,
                     name,
                     birth_date,
                     species_id,
                     breed_id,
                     gender_id,
                     note,
                     dead,
                     creator,
                     created,
                     modifier,
                     modified,
                     version) VALUES
  (6, 'create-tag-pet', current_date, 1, 1, 1, 'test-note', false, 'testuser', current_timestamp, 'testmodifier',
   current_timestamp, 5);

-- id=2
INSERT INTO tag (value, patient_id, tag_type_id, creator, created, version) VALUES ('tag-id', 3,2,'testuser',current_timestamp, 2);