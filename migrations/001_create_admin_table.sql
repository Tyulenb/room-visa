--13.09.2025 create_admin.sql
BEGIN;

CREATE TABLE admin (
    login varchar(40) PRIMARY KEY,
    hashpass varchar(255) NOT NULL
);

--DEFAULT admin record
INSERT INTO admin VALUES('seal', '$2a$10$4B6sLxtzFbdJSLBIAPMOwu9AqHZDzNrz.EF52eNGZ0BXzhugx4IUW');

COMMIT;
