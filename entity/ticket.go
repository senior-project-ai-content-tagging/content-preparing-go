package entity

type Ticket struct {
	Id        int64  `db:"id"`
	ContentId int64  `db:"content_id"`
	Status    string `db:"status"`
}

type TicketSubmitWeblinkPubSub struct {
	TicketId int64  `json:"ticketId"`
	Url      string `json:"url"`
}

type TicketContentPubSub struct {
	TicketId int64 `json:"ticketId"`
}

type TicketPubSub struct {
	TicketId int64 `json:"ticketId"`
}
