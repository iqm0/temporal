package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	// This workflow ID can be user defined or an ID from your database
	workflowID := "customerSupportWorkflow_1"
	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "customerSupport",
	}

	// Start a workflow execution
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, CustomerSupportWorkflow, "ticket123")
	if err != nil {
		panic(err)
	}

	// Wait for the workflow to complete
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get workflow result:", err)
	}
	// Print the workflow completion result
	log.Println("Workflow completed with result:", result)
	// Output: Workflow completed with result: ticket123

}
