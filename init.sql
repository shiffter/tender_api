CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
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
