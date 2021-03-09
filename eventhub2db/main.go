package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

const (
	connStr   string = "postgresql://postgres:password@127.0.0.1:5432/postgres"
	tableName string = "adsData"
)

var (
	columNames  = []string{"time", "device_id", "data"}
	columnTypes = []string{"TIMESTAMPTZ", "VARCHAR(50)", "JSONB"}
)

type tableEntry struct {
	time     time.Time `db:"time"`
	deviceID string    `db:"device_id"`
	data     string    `db:"data"`
}

func main() {
	// connect to database using a single connection
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Extend the database with TimescaleDB
	query := `CREATE EXTENSION IF NOT EXISTS timescaledb;`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to extend the database with TimescaleDB: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully extended the database with TimescaleDB")

	// Create table if not exists
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS "+tableName+"();")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create TABLE "+tableName+": %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully created TABLE", tableName)

	// Insert columns if not exist
	for id, column := range columNames {
		_, err = conn.Exec(ctx, "ALTER TABLE "+tableName+" ADD COLUMN IF NOT EXISTS "+column+" "+columnTypes[id]+";")
		if err != nil {
			fmt.Println("Unable to add column", column, "to table", tableName, ":", err)
			os.Exit(1)
		}
		fmt.Println("Successfully added column", column, "to table", tableName)
	}

	// Convert the table created into a hypertable
	// This needs to be executed after column "time" was added
	query = "SELECT create_hypertable('" + tableName + "', 'time', if_not_exists => TRUE);"
	_, err = conn.Exec(ctx, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to convert TABLE into a hypertable: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully converted TABLE into a hypertable")

	// Data to put into the database
	entry := tableEntry{
		time:     time.Now(),
		deviceID: "eval2",
		data: `{
			"name": "Chili",
			"ingredients": ["onion", "beef", "chili powder", "tomato paste"],
			"organic": true,
			"dimensions": {
				"weight": 1000.00
			}
		}`,
	}

	// Dynamically generate query to insert one row to the database
	columns := ""
	values := ""
	for idx, column := range columNames {
		columns += column
		values += "$" + strconv.Itoa(idx+1)

		if idx+1 < len(columNames) {
			columns += ", "
			values += ", "
		}
	}
	queryInsertTimeseriesData := "INSERT INTO " + tableName + " (" + columns + ") VALUES (" + values + ");"
	fmt.Println(queryInsertTimeseriesData)

	_, err = conn.Exec(ctx, queryInsertTimeseriesData, entry.time, entry.deviceID, entry.data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert sample into Timescale %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully inserted data into table", tableName)
}
