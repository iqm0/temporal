package app

const CustomerSupportQueue = "CUSTOMER_SUPPORT_QUEUE"

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
