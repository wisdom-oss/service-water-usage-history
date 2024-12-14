-- name: get-paginated
SELECT *
FROM timeseries.water_usage
LIMIT $1 OFFSET $2;