-- +goose Up

-- +goose StatementBegin
CREATE FUNCTION snapshot_has_ancestor(
    start_snapshot_id BIGINT,
    target_snapshot_id BIGINT
)
RETURNS BOOLEAN AS $func$
DECLARE
    current_snapshot_id BIGINT := start_snapshot_id;
BEGIN
    WHILE current_snapshot_id IS NOT NULL LOOP
        SELECT parent_snapshot_id INTO current_snapshot_id
        FROM snapshots
        WHERE snapshot_id = current_snapshot_id;

        IF current_snapshot_id = target_snapshot_id THEN
            RETURN TRUE;
        END IF;
    END LOOP;

    RETURN FALSE;
END;
$func$ LANGUAGE plpgsql STABLE STRICT;
-- +goose StatementEnd

-- +goose Down

DROP FUNCTION IF EXISTS snapshot_has_ancestor(BIGINT, BIGINT);
