-- 13.09.2025

BEGIN;

CREATE TYPE request_status AS ENUM ('Approved', 'Rejected', 'Awaiting');

CREATE TABLE request (
    id UUID PRIMARY KEY,
    status request_status NOT NULL DEFAULT 'Awaiting',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (created_at <= updated_at)
);

CREATE TYPE sex_type AS ENUM ('male', 'female');

CREATE TABLE request_data (
    id SERIAL PRIMARY KEY,
    request_id UUID NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    sex sex_type NOT NULL,
    ethnicity VARCHAR(100) NOT NULL,
    citizenship VARCHAR(100) NOT NULL,
    purpose text NOT NULL,
    photopath text NOT NULL,
    CONSTRAINT fk_request_id 
        FOREIGN KEY(request_id)
        REFERENCES request(id)
        ON DELETE CASCADE
);

CREATE TABLE visa (
    id UUID PRIMARY KEY,
    request_id UUID NOT NULL UNIQUE,
    token VARCHAR(511) NOT NULL,
    CONSTRAINT fk_request_id 
        FOREIGN KEY(request_id)
        REFERENCES request(id)
);

COMMIT;
