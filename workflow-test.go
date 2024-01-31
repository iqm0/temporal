package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessfulTicketWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	type CustomerSupportTicket struct {
		TicketID      int
		CustomerID    string
		CreatedOn     string
		CustomerName  string
		CustomerEmail string
		Subject       string
		Description   string
		Status        string
		AssignedTo    string
		Priority      string
	}

	testDetails := CustomerSupportTicket{
		TicketID:      1,
		CustomerID:    "CUST-001",
		CreatedOn:     "2021-08-31 15:47:00",
		CustomerName:  "John Doe",
		CustomerEmail: "john.doe@example.com",
		Subject:       "Issue with my account",
		Description:   "My account is not working properly.",
		Status:        "Open",
		AssignedTo:    "",
		Priority:      "Medium",
	}

	// Mock activity implementation
	env.OnActivity(Withdraw, mock.Anything, testDetails).Return("", nil)
	env.OnActivity(Deposit, mock.Anything, testDetails).Return("", nil)

	env.ExecuteWorkflow(MoneyTransfer, testDetails)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
