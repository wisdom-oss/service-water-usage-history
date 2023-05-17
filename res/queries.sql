-- Insert a new consumer type into the database
-- name: add-consumer-type
INSERT INTO water_usage.usage_types (name, description, external_identifier)
VALUES ($1, $2, $3)
RETURNING id;

-- Update a current consumer type's name
-- name: update-consumer-type-name
UPDATE water_usage.usage_types
SET name = $1
WHERE id = $2;

-- Update a current consumer type's description
-- name: update-consumer-type-description
UPDATE water_usage.usage_types
SET description = $1
WHERE id = $2;

-- Update a current consumer type's external identifier
-- name: update-consumer-type-identifier
UPDATE water_usage.usage_types
SET external_identifier = $1
WHERE id = $2;

-- Delete a consumer type
-- name: delete-consumer-type
DELETE FROM water_usage.usage_types
WHERE id = $1;

-- Get all water usage records
-- name: get-all-usages
SELECT *
FROM water_usage.usages;

-- Get water usages by consumer
-- name: get-consumers-usages
SELECT *
FROM water_usage.usages
WHERE consumer = $1::uuid;

-- Get water usages by type
-- name: get-type-based-usages
SELECT *
FROM water_usage.usages
WHERE usage_type = (
    SELECT id
    FROM water_usage.usage_types
    WHERE external_identifier = $1
);

-- Get water usages by time range
-- name: get-usages-between-times
SELECT *
FROM water_usage.usages
WHERE date < $1 AND date > $2;

-- Add a new usage record
-- name: add-usage-record
INSERT INTO water_usage.usages(municipality, date, consumer, usage_type, amount)
VALUES ($1, $2, $3, $4, $5);

-- Delete a usage record
-- name: delete-usage-record
DELETE FROM water_usage.usages
WHERE id = $1;