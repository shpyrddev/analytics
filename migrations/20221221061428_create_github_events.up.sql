CREATE TABLE 'github_events' (
    PushID Int64,
    Head LowCardinality(String),
    Ref String,
    -- Size Int8
    -- Commits Nested(
    --     Message String,
    --     Author Nested (
    --        Name String,
    --        Login String 
    --     ),
    --     URL String
    -- ),
    Before String,
    -- DistinctSize Int8

    After String,
    Created Boolean,
    Deleted Boolean,
    Forced Boolean,
    BaseRef String,
    Repo Nested(
        Name String
    )
) ENGINE = MergeTree ORDER BY (PushID);