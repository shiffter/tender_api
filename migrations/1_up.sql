CREATE TABLE schema_migrations;

CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
);

CREATE TYPE propos_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCEL'
);

CREATE TYPE service_type AS ENUM (
    'CONSTRUCTION',
    'DELIVERY',
    'MANUFACTURE',
);

CREATE TABLE tender
(
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY ,
    name             VARCHAR(100),
    description      VARCHAR(200),
    status           tender_status,
    creator_username VARCHAR(50) REFERENCES employee (username),
    organization_id  uuid REFERENCES organization,
    service_type     service_type,
    version          integer   default 1,
    created_at       timestamp default now()
);

CREATE TABLE proposals (
                           id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
                           tender_id uuid NOT NULL,
                           author_id uuid NOT NULL,
                           organization_id uuid NOT NULL,
                           name varchar(100),
                           status propos_status NOT NULL DEFAULT 'CREATED',
                           version INT NOT NULL DEFAULT 1,
                           created_at TIMESTAMP DEFAULT NOW(),
                           updated_at TIMESTAMP DEFAULT NOW(),
                           description TEXT,
                           CONSTRAINT fk_tender FOREIGN KEY (tender_id) REFERENCES tender(id),
                           CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES employee(id),
                           CONSTRAINT fk_organization FOREIGN KEY (organization_id) REFERENCES organization(id),
                           UNIQUE (tender_id, author_id)
);

CREATE OR REPLACE FUNCTION increment_version()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.version := OLD.version + 1;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_version
    BEFORE UPDATE ON tender
    FOR EACH ROW
    EXECUTE FUNCTION increment_version();

CREATE TRIGGER update_version
    BEFORE UPDATE ON proposals
    FOR EACH ROW
    EXECUTE FUNCTION increment_version();

