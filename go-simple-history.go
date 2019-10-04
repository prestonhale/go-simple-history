package goSimpleHistory

import (
	"fmt"
	sqlx "github.com/jmoiron/sqlx"
)

// ========================= SAFETY NOTICE =============================
// The tableName param to CreateAndRecordHistory MUST be safe (e.g. not '; DROP TABLE users')
// It is not escaped. Do not send random crap through that param.
// ============================ by Preston ==============================

const ERROR_NO_ROWS_IN_RESULT  = `sql: no rows in result set`


func CreateAndRecordHistory(dbConn *sqlx.DB, tableName string){
	tx, err := dbConn.Beginx()
	if err != nil {
		//logger
		return
	}

	// 1. Check if history table already exists
	// TODO: Check if base table has been updated and update history table to match
	historyTableExists := true
	dynamicTableCheck := fmt.Sprintf(dynamicCheckHistoryStmt, tableName)
	checkHistroyTableSql := dbConn.Rebind(dynamicTableCheck)
	_, err = tx.Exec(checkHistroyTableSql)
	if err != nil {
		if err.Error() == ERROR_NO_ROWS_IN_RESULT{
			historyTableExists = false
		} else {
			//logger
			return
		}
	}
	if historyTableExists {
		//logger
		return
	}

	// 2. Create history table
	createHistoryTableStmt := fmt.Sprintf(dynamicCreateHistoryTableStmt, tableName)
	createHistoryTableSql := dbConn.Rebind(createHistoryTableStmt)
	_, err = tx.Exec(createHistoryTableSql)
	if err != nil {
		//logger
		return
	}

	// 3. Create function to run on trigger
	createTriggerFuncStmt := fmt.Sprintf(dynamicCreateTriggerFuncStmt, tableName)
	createTriggerFuncSql := dbConn.Rebind(createTriggerFuncStmt)
	_, err = tx.Exec(createTriggerFuncSql)
	if err != nil {
		//logger
		return
	}

	// 4. Create trigger
	createTriggerStmt := fmt.Sprintf(dynamicCreateTriggerStmt, tableName)
	createTriggerSql := dbConn.Rebind(createTriggerStmt)
	_, err = tx.Exec(createTriggerSql)
	if err != nil {
		//logger
		return
	}
}
