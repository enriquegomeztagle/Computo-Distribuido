package main

import (
	"encoding/json"
	"fmt"
	"json-http-server-mux-rpc/internal/log"
	"net/rpc"
)

func main() {
	// Connect to RPC server
	client, err := rpc.Dial("tcp", "mux-rpc-log-server:3214")

	if err != nil {
		fmt.Println("Error connecting to the RPC server:", err)
		return
	}
	defer client.Close()

	// Read encoded value
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter base64 encoded value: ")
	// encodedValue, _ := reader.ReadString('\n')
	encodedValue := "enriquegomeztagle" // hardcoded for docker tests
	// encodedValue = encodedValue[:len(encodedValue)-1]

	appendArgs := log.AppendArgs{
		Record: log.Record{
			Value: []byte(encodedValue),
		},
	}

	var appendReply log.AppendResponse
	err = client.Call("LogRPC.Append", &appendArgs, &appendReply)
	if err != nil {
		fmt.Println("Error appending record:", err)
		return
	}
	fmt.Printf("Record appended at offset: %d\n", appendReply.Offset)

	// Change operation
	fmt.Println("Append operation complete, starting fetch operation...")

	fetchArgs := log.FetchRequest(appendReply)

	var fetchReply log.FetchResponse
	err = client.Call("LogRPC.Fetch", &fetchArgs, &fetchReply)
	if err != nil {
		fmt.Println("Error fetching record:", err)
		return
	}

	// Confirm new operation
	fmt.Println("Fetch operation started...")

	// Show response
	response := map[string]interface{}{
		"record": map[string]interface{}{
			"value":  string(fetchReply.Record.Value),
			"offset": fetchReply.Record.Offset,
		},
	}

	responseJSON, _ := json.Marshal(response)
	fmt.Println(string(responseJSON))
}
