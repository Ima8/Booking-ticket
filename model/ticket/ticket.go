package ticket

type TicketRemain struct {
	Seats                  []string `json:"seats"`
	UnconfimedTicketsCount int      `json:"unconfimedTicketsCount"`
}

func GetRemainTicket(round int) {

}

func IsTicketFinish(round int) bool {
	var isFinish = true

	return isFinish
}
