package goSimpleHistory


const dynamicCheckHistoryStmt = `
SELECT * FROM information_schema.tables WHERE table_schema='public' AND table_name=%[1]s;
`

const dynamicCreateHistoryTableStmt = `
DO
$do$
BEGIN

EXECUTE (
	WITH table_columns AS (
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema='public' AND table_name='%[1]s'
	)
	CREATE TABLE %[1]s_history (%s)
;
`

// TODO: Potentially slow because we're fetching column data from info_schema constantly
// TODO: The escaping in this is just horrendous
//  %% are values that plpgsql will use that need to pass through the golang sprintf formatter first
const dynamicCreateTriggerFuncStmt = `
CREATE OR REPLACE FUNCTION record_player_currency_table_change()
RETURNS trigger AS
$BODY$
BEGIN

	SELECT string_agg(column_name, ',')
	INTO tmp_table_columns TEMP UNLOGGED
	FROM information_schema.columns
	WHERE table_schema = 'public' AND table_name='%[1]s'
	;
	
	SELECT string_agg(NEW.*, ',')
	INTO tmp_table_values TEMP UNLOGGED
	;

	EXECUTE format(
		'
		INSERT INTO %[1]s_history
			(changed_at, change_type, %%1$s)
		SELECT
			NOW()
			, CASE TG_OP
				WHEN 'INSERT' THEN '+'
				WHEN 'UPDATE' THEN '~'
				WHEN 'DELETE' THEN '-'
			END
			, %%2$s
		', tmp_table_columns, tmp_table_values
	);

RETURN NEW;
END;
$BODY$

LANGUAGE plpgsql VOLATILE
;
`

const dynamicCreateTriggerStmt = `
CREATE TRIGGER %[1]s_changes
BEFORE UPDATE OR INSERT OR DELETE
ON %[1]s
FOR EACH ROW
EXECUTE PROCEDURE record_%[1]s_table_change();
`

