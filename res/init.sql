-- Create the database schema into which the water usages are written
-- name: create-schema
CREATE SCHEMA IF NOT EXISTS waterUsage;

-- Create a table which will store the different usage types for the recorded
-- usages
-- name: create-usage-types-table
CREATE TABLE IF NOT EXISTS waterUsage.usageTypes (
    id serial primary key,
    name varchar not null,
    description text,
    externalIdentifier varchar not null
);

-- Create a table which will hold the water usages and allows the association
-- of a entry with a consumer
-- name: create-usage-table
CREATE TABLE IF NOT EXISTS waterUsage.usages (
    id serial primary key,
    municipality varchar(12) not null,
    date timestamp not null,
    consumer uuid,
    usageType int REFERENCES waterUsage.usageTypes(id) ON DELETE SET NULL ON UPDATE CASCADE,
    createdAt timestamp not null DEFAULT NOW(),
    amount double precision not null
);

-- Create some example usage types to allow users to get started right away
-- name: create-initial-usage-types
INSERT INTO waterUsage.usageTypes  (name, description, externalIdentifier)
VALUES
('Private Households', 'This usage type contains usages recorded for residential buildings like flats and houses',
 'privateHouseholds'),
('Businesses', 'This usage type contains usages recorded for business buildings like factories and grocery stores',
 'businesses'),
('Agriculture, Forestry, Fishery', 'This usage type contains usages recorded for businesses handling agriculture or forestry',
 'agriculture'),
('Public Institutions', 'This usage type contains usages recorded for public institutions like governmental buildings',
 'publicInstitutions');

