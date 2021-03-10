/*
Copyright © 2021 Ci4Rail GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package timescaledb

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

const (
	tableName string = "adsData"
)

var (
	// colums to be created in the table tableName of the database
	colums = []string{"time TIMESTAMPTZ", "device_id VARCHAR(50)", "data JSONB"}
)

// TableEntry struct with data to insert into table
type TableEntry struct {
	Time     time.Time `db:"time"`
	DeviceID string    `db:"device_id"`
	Data     string    `db:"data"`
}

// Connection to a timescale db
type Connection struct {
	ctx  context.Context
	conn *pgx.Conn
}

// NewConnection create a connection to a timescale db by connection string c
// * Extend the database with TimescaleDB (if not happed before)
// * Create table if not exists
// * Insert columns if not exist
// * Convert the table created into a hypertable (if not exists)
// * Insert sample data into database
func NewConnection(connStr string) (*Connection, error) {
	// connect to database using a single connection
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v", err)
	}

	c := Connection{
		ctx:  ctx,
		conn: conn,
	}

	// Extend the database with TimescaleDB
	query := `CREATE EXTENSION IF NOT EXISTS timescaledb;`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Unable to extend the database with TimescaleDB: %v", err)
	}

	// Create table if not exists
	_, err = conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS "+tableName+"();")
	if err != nil {
		return nil, fmt.Errorf("Unable to create TABLE "+tableName+": %v", err)
	}

	// Insert columns if not exist
	for _, column := range colums {
		_, err = conn.Exec(ctx, "ALTER TABLE "+tableName+" ADD COLUMN IF NOT EXISTS "+column+";")
		if err != nil {
			return nil, fmt.Errorf("Unable to add column %s to table %s: %v", column, tableName, err)
		}
	}

	// Convert the table created into a hypertable
	// This needs to be executed after column "time" was added
	query = "SELECT create_hypertable('" + tableName + "', 'time', if_not_exists => TRUE);"
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert TABLE into a hypertable: %v", err)
	}

	return &c, nil
}

// Write write data into the database
func (c Connection) Write(t TableEntry) error {
	query := "INSERT INTO " + tableName + " (time, device_id, data) VALUES ($1, $2, $3);"

	_, err := c.conn.Exec(c.ctx, query, t.Time, t.DeviceID, t.Data)
	if err != nil {
		return fmt.Errorf("Unable to insert sample into Timescale: %v", err)
	}

	return nil
}

// Close connection to db
func (c Connection) Close() {
	c.conn.Close(c.ctx)
}
