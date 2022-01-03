package main

import (
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {

	// setting env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPath := os.Getenv("DBPATH")
	token := os.Getenv("API_TOKEN")
	organization := os.Getenv("ORGANIZATION")
	bucket := os.Getenv("BUCKET")

	fmt.Println("InfluxDB + Golang Monitoring Test")

	client := influxdb2.NewClient(dbPath, token)
	// get non-blocking write client
	writeAPI := client.WriteAPI(organization, bucket)

	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45.0},
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	// Or write directly line protocol
	// line := fmt.Sprintf("stat,unit=temperature avg=%f,max=%f", 23.5, 45.0)
	// writeAPI.WriteRecord(context.Background(), line)

	// Query
	// queryAPI := client.QueryAPI(organization)
	// result, err := queryAPI.Query(context.Background(), `from(bucket: "test")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "stat")`)
	// if err != nil {
	// 	log.Println(err)
	// }

	// // Use Next() to iterate over query result lines
	// for result.Next() {
	// 	// Observe when there is new grouping key producing new table
	// 	if result.TableChanged() {
	// 		fmt.Printf("table: %s\n", result.TableMetadata().String())
	// 	}
	// 	// read result
	// 	fmt.Printf("row: %s\n", result.Record().String())
	// }
	// if result.Err() != nil {
	// 	fmt.Printf("Query error: %s\n", result.Err().Error())
	// }

	client.Close()
}
