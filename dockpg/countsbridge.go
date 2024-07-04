package dockpg

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func pythonToFenAndBridge(db *sql.DB, data []FenCount) {
	const maxParams = 65535

	var countValues []string
	var countArgs []interface{}
	var bridgeValues []string
	var bridgeArgs []interface{}
	countIdx, bridgeIdx := 1, 1

	for _, entry := range data {
		if len(countArgs)+len(bridgeArgs)+len(entry.Urls)*2+2 > maxParams {
			// Execute batch
			executeBatch(db, countValues, countArgs, bridgeValues, bridgeArgs)
			countValues, countArgs, bridgeValues, bridgeArgs = nil, nil, nil, nil
			countIdx, bridgeIdx = 1, 1
		}

		// Append data for counts
		countValues = append(countValues, fmt.Sprintf("($%d, $%d)", countIdx, countIdx+1))
		countArgs = append(countArgs, entry.Fen, entry.Count)
		countIdx += 2

		// Handle multiple URLs for each FEN
		for _, url := range entry.Urls {
			bridgeValues = append(bridgeValues, fmt.Sprintf("($%d, $%d)", bridgeIdx, bridgeIdx+1))
			bridgeArgs = append(bridgeArgs, url, entry.Fen)
			bridgeIdx += 2
		}
	}

	if len(countArgs) > 0 || len(bridgeArgs) > 0 {
		executeBatch(db, countValues, countArgs, bridgeValues, bridgeArgs)
	}
}

func executeBatch(db *sql.DB, countValues []string, countArgs []interface{}, bridgeValues []string, bridgeArgs []interface{}) {

	countQuery := fmt.Sprintf("INSERT INTO counts (fen, count) VALUES %s ON CONFLICT (fen) DO UPDATE SET count = counts.count + excluded.count", strings.Join(countValues, ", "))
	bridgeQuery := fmt.Sprintf("INSERT INTO bridge (link, fen) VALUES %s", strings.Join(bridgeValues, ", "))

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("ON BEGIN")
		log.Fatal(err)
	}

	_, err = tx.Exec(countQuery, countArgs...)
	if err != nil {
		tx.Rollback()
		fmt.Println("ON COUNTS EXEC")
		log.Fatal(err)
	}

	_, err = tx.Exec(bridgeQuery, bridgeArgs...)
	if err != nil {
		tx.Rollback()
		fmt.Println("ON BRIDGE EXEC")
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		fmt.Println("ON COMMIT EXEC")
		log.Fatal(err)
	}
}
