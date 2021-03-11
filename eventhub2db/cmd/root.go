/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/ci4rail/ads-demo-environment/eventhub2db/internal/avro"
	e "github.com/ci4rail/ads-demo-environment/eventhub2db/internal/eventhub"
	tdb "github.com/ci4rail/ads-demo-environment/eventhub2db/internal/timescaledb"
	"github.com/spf13/cobra"
)

const (
	defaultDbConnStr = "postgresql://postgres:password@localhost:5432/postgres"
)

var (
	dbConnStr string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "eventhub2db",
	Short: "Cyclical read data from eventhub and store it in db",
	Long:  `Cyclical read data from eventhub and store it in db`,
	Run: func(cmd *cobra.Command, args []string) {

		// Connect to database
		dbConn, err := tdb.NewConnection(dbConnStr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer dbConn.Close()

		connStr, err := e.ReadConnectionStringFromEnv()
		if err != nil {
			fmt.Println(err)
			return
		}

		hub, err := eventhub.NewHubFromConnectionString(connStr)

		if err != nil {
			fmt.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		handler := func(c context.Context, event *eventhub.Event) error {
			// Output as direct access
			fmt.Println("Reading...")
			fmt.Println("Received data:")
			fmt.Println(event.Data)
			fmt.Println()

			avro, err := avro.NewAvroReader(event.Data)
			if err != nil {
				return err
			}

			m, err := avro.AvroToMap()
			if err != nil {
				return err
			}

			// Output as json
			jbytes, err := avro.AvroToByteString()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Decoded data as JSON:")
			fmt.Println(avro.AvroToJson())
			fmt.Println("Decoded single fields:")
			fmt.Println("device: " + m["device"].(string))
			fmt.Println("acqTime: ", time.Unix(int64(m["acqTime"].(int32)), 0))
			fmt.Println()

			// Remove the timestamp and the device id from json
			var i interface{}
			if err := json.Unmarshal([]byte(string(jbytes)), &i); err != nil {
				panic(err)
			}
			if m2, ok := i.(map[string]interface{}); ok {
				if _, ok := m2["device"].(string); ok {
					delete(m2, "device")
				}
				if _, ok := m2["acqTime"].(string); ok {
					delete(m2, "acqTime")
				}
			}

			data, err := json.Marshal(i)
			if err != nil {
				return err
			}
			fmt.Println("Removed timestamp and device id from json:")
			fmt.Println(string(data))
			fmt.Println()

			// Data to put into the database
			entry := tdb.TableEntry{
				Time:     time.Unix(int64(m["acqTime"].(int32)), 0),
				DeviceID: m["device"].(string),
				Data:     string(data),
			}

			err = dbConn.Write(entry)
			if err != nil {
				return err
			} else {
				fmt.Println("Successfully wrote into db:")
				fmt.Println(entry)
				fmt.Println()
			}

			return nil
		}

		// listen to each partition of the Event Hub
		runtimeInfo, err := hub.GetRuntimeInformation(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, partitionID := range runtimeInfo.PartitionIDs {
			listenerHandle, err := hub.Receive(ctx, partitionID, handler, eventhub.ReceiveWithLatestOffset())
			if err != nil {
				fmt.Println(err)
				return
			}
			listenerHandle.Done()
		}

		// Wait for a signal to quit:
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		<-signalChan

		err = hub.Close(context.Background())
		if err != nil {
			fmt.Println(err)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringVarP(&dbConnStr, "dbconnstr", "d", defaultDbConnStr, "use alternative db connection string")
}
