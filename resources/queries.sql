-- name: get-paginated
SELECT
    *
FROM
    timeseries.water_usage
LIMIT
    $1
OFFSET
    $2;

-- name: consumer-exists
SELECT
    EXISTS (
        SELECT
            id
        FROM
            consumers.consumers
        WHERE
            id = $1
    );

-- name: consumer-usages
SELECT
    *
FROM
    timeseries.water_usage
WHERE
    consumer = $1
LIMIT
    $2
OFFSET
    $3;

-- name: municipal-usages
SELECT
    *
FROM
    timeseries.water_usage
WHERE
    municipality = $1
LIMIT
    $2
OFFSET
    $3;


-- name: typed-usages
SELECT
    *
FROM
    timeseries.water_usage
WHERE
    usage_type = $1
LIMIT
    $2
OFFSET
    $3;