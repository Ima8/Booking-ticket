package ticket

type TicketRemain struct {
	Seats                  []string `json:"seats"`
	UnconfimedTicketsCount int      `json:"unconfimedTicketsCount"`
}
