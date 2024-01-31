package app

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Define a workflow.
func CustomerSupportWorkflow(ctx workflow.Context, ticketID string) error {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        0, // unlimited retries
		NonRetryableErrorTypes: []string{"VoidAccountError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Step 1: Create a new support ticket
	err := workflow.ExecuteActivity(ctx, CreateTicketActivity, ticketID).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Step 2: Assign the ticket to an available support agent
	err = workflow.ExecuteActivity(ctx, AssignTicketActivity, ticketID).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Step 3: Wait for the ticket to be resolved or escalate after a timeout
	selector := workflow.NewSelector(ctx)
	var resolved bool
	selector.AddFuture(workflow.NewTimer(ctx, 24*time.Hour), func(f workflow.Future) {
		// Escalate the ticket if not resolved in 24 hours
		workflow.ExecuteActivity(ctx, EscalateTicketActivity, ticketID)
	})
	selector.AddReceive(workflow.GetSignalChannel(ctx, "ticket_resolved"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &resolved)
	})

	selector.Select(ctx)

	if resolved {
		// Step 4: Ticket has been resolved
		err = workflow.ExecuteActivity(ctx, CloseTicketActivity, ticketID).Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	// Step 5: Follow up with the customer to ensure satisfaction
	err = workflow.ExecuteActivity(ctx, FollowUpCustomerActivity, ticketID).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
