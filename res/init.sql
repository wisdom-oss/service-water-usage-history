-- Create the database schema into which the water usages are written
-- name: create-schema
CREATE SCHEMA IF NOT EXISTS water_usage;

-- Create a table which will store the different usage types for the recorded
-- usages
-- name: create-usage-types-table
CREATE TABLE IF NOT EXISTS water_usage.usage_types(
    id serial primary key,
    name varchar not null unique ,
    description text,
    external_identifier varchar not null unique
);

-- Create a table which will hold the water usages and allows the association
-- of a entry with a consumer
-- name: create-usage-table
CREATE TABLE IF NOT EXISTS water_usage.usages(
    id serial primary key,
    municipality varchar(12) not null,
    date timestamp not null,
    consumer uuid,
    usage_type uuid REFERENCES water_usage.usage_types(id) ON DELETE SET NULL ON UPDATE CASCADE,
    createdAt timestamp not null DEFAULT NOW(),
    amount double precision not null
);

-- Create some example usage types to allow users to get started right away
-- name: create-initial-usage-types
INSERT INTO water_usage.usage_types  (name, description, external_identifier)
VALUES
('Private Households', 'This usage type contains usages recorded for residential buildings like flats and houses',
 'privateHouseholds'),
('Businesses', 'This usage type contains usages recorded for business buildings like factories and grocery stores',
 'businesses'),
('Agriculture, Forestry, Fishery', 'This usage type contains usages recorded for businesses handling agriculture or forestry',
 'agriculture'),
('Public Institutions', 'This usage type contains usages recorded for public institutions like governmental buildings',
 'publicInstitutions');

