package main

import (
	"bufio"
	"fmt"
	L "log"
	"os"
	log_v1 "server-transactions-commit-log/api/v1"
	"server-transactions-commit-log/log"
	"strings"
)

func main() {
	config := log.Config{}
	config.Segment.MaxStoreBytes = 1024
	config.Segment.MaxIndexBytes = 1024

	dir := "./log_data"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			L.Fatalf("Error while creating the directory %s: %v", dir, err)
		}
	}

	lg, err := log.NewLog(dir, config)
	if err != nil {
		L.Fatalf("Error while creating the log: %v", err)
	}
	defer lg.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Add a new record")
		fmt.Println("2. Read a specific record")
		fmt.Println("3. View lower and highest offsets")
		fmt.Println("4. View all records")
		fmt.Println("S. Exit")

		fmt.Print("Option: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Print("Register value: ")
			value, _ := reader.ReadString('\n')
			value = strings.TrimSpace(value)

			record := &log_v1.Record{Value: []byte(value)}
			offset, err := lg.Append(record)
			if err != nil {
				L.Printf("Error while adding the record: %v", err)
			} else {
				fmt.Printf("Record added at offset %d\n", offset)
			}

		case "2":
			fmt.Print("Offset to search: ")
			var offset uint64
			fmt.Scanf("%d", &offset)

			record, err := lg.Read(offset)
			if err != nil {
				L.Printf("Error while reading the record: %v", err)
			} else {
				fmt.Printf("Record at offset %d: %s\n", offset, record.Value)
			}

		case "3":
			lowestOffset, err := lg.LowestOffset()
			if err != nil {
				L.Printf("Error while obtaining the lowest offset: %v", err)
			} else {
				fmt.Printf("Lowest offset: %d\n", lowestOffset)
			}

			highestOffset, err := lg.HighestOffset()
			if err != nil {
				L.Printf("Error while obtaining the highest offset: %v", err)
			} else {
				fmt.Printf("Highest offset: %d\n", highestOffset)
			}

		case "4":
			lowestOffset, err := lg.LowestOffset()
			if err != nil {
				L.Printf("Error while obtaining the lowest offset: %v", err)
				break
			}

			highestOffset, err := lg.HighestOffset()
			if err != nil {
				L.Printf("Error while obtaining the highest offset: %v", err)
				break
			}

			for i := lowestOffset; i <= highestOffset; i++ {
				record, err := lg.Read(i)
				if err != nil {
					L.Printf("Error while reading the record at offset %d: %v", i, err)
					continue
				}
				fmt.Printf("Record at offset %d: %s\n", i, record.Value)
			}

		case "S":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid option")
		}
	}
}
