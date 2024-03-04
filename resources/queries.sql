-- name: create-data-table
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS time_series.usages
(
    time         timestamptz,
    amount       double precision,
    usage_type   uuid,
    consumer     uuid,
    municipality text
);
SELECT create_hypertable('usages', by_range('time'));
END TRANSACTION;

-- name: get-all
SELECT *
FROM time_series.usages
LIMIT $1 OFFSET $2;

-- name: get-from
SELECT *
FROM time_series.usages
WHERE time > $1
LIMIT $2 OFFSET $3;

-- name: get-until
SELECT *
FROM time_series.usages
WHERE time < $1
LIMIT $2 OFFSET $3;

-- name: get-range
SELECT *
FROM time_series.usages
WHERE time > $1
AND time < $2
LIMIT $3 OFFSET $4;

-- name: filter-consumer
WHERE consumer = $1::uuid

-- name: filter-municipality
WHERE municipality = $1