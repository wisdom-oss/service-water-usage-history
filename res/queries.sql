/* Queries related to the consumers */

-- Get a consumer by it's uuid
-- name: get-consumer-by-uuid
SELECT
    id,
    name,
    ST_AsGeoJSON(location) as location,
    usage_type,
    additional_properties
FROM
    consumers.consumers
WHERE
    id = $1::uuid;


/* Queries related to the usage types */
-- Get all usage types
-- name: get-all-usage-types
SELECT
    id,
    name,
    description,
    external_identifier
FROM
    water_usage.usage_types;

-- Get a usage type by it's UUID
-- name: get-single-usage-type-by-id
SELECT
    id,
    name,
    description,
    external_identifier
FROM
    water_usage.usage_types
WHERE
    id = $1::uuid;

-- Get a usage type by it's external identifier
-- name: get-single-usage-type-by-external-identifier
SELECT
    id,
    name,
    description,
    external_identifier
FROM
    water_usage.usage_types
WHERE
    external_identifier = $1;

-- Update a current usage type's name
-- name: update-consumer-type-name
UPDATE
    water_usage.usage_types
SET
    name = $1
WHERE
    id = $2;

-- Update a current usage type's description
-- name: update-usage-type-description
UPDATE
    water_usage.usage_types
SET
    description = $1
WHERE
    id = $2;

-- Update a current usage type's external identifier
-- name: update-usage-type-identifier
UPDATE
    water_usage.usage_types
SET
    external_identifier = $1
WHERE
    id = $2;

-- Insert a new usage type into the database
-- name: add-usage-type
INSERT INTO
    water_usage.usage_types (name, description, external_identifier)
VALUES
    ($1, $2, $3)
RETURNING
    id;

-- Delete a usage type
-- name: delete-usage-type
DELETE FROM
    water_usage.usage_types
WHERE
    id = $1;

/* Queries related to usage records */
-- Get all water usage records
-- name: get-all-usages
SELECT
    id,
    municipality,
    date,
    consumer,
    usage_type,
    created_at,
    amount
FROM
    water_usage.usages
WHERE id BETWEEN $1 AND $2;
-- Get water usages by consumer
-- name: get-consumers-usages
SELECT
    id,
    municipality,
    date,
    consumer,
    usage_type,
    created_at,
    amount
FROM
    water_usage.usages
WHERE
    consumer = $1::uuid;

-- Get water usages by type
-- name: get-type-based-usages
SELECT
    id,
    municipality,
    date,
    consumer,
    usage_type,
    created_at,
    amount
FROM
    water_usage.usages
WHERE
    usage_type = $1::uuid;

-- Get water usages by time range
-- name: get-usages-between-times
SELECT
    id,
    municipality,
    date,
    consumer,
    usage_type,
    created_at,
    amount
FROM
    water_usage.usages
WHERE
    date < $1
AND
    date > $2;

-- Get a single usage record by it's internal id
-- name: get-single-record
SELECT
    id,
    municipality,
    date,
    consumer,
    usage_type,
    created_at,
    amount
FROM
    water_usage.usages
WHERE
    id = $1;

-- Add a new usage record
-- name: add-usage-record
INSERT INTO
    water_usage.usages(municipality, date, consumer, usage_type, amount)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    id;

-- Delete a usage record
-- name: delete-usage-record
DELETE FROM
    water_usage.usages
WHERE
    id = $1;
