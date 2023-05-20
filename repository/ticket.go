package repository

import (
	"github.com/senior-project-ai-content-tagging/content-preparing/entity"
)

type TicketRepository interface {
	UpdateTicketStatus(ticket entity.Ticket) (string, *entity.Ticket)
	FindTicketById(id int64) (string, int64)
}

type ticketRepositorySqlx struct {
}

func (r *ticketRepositorySqlx) UpdateTicketStatus(ticket entity.Ticket) (string, *entity.Ticket) {
	return "UPDATE tickets SET status = :status WHERE id = :id", &ticket
}

func (r *ticketRepositorySqlx) FindTicketById(id int64) (string, int64) {
	return "SELECT id, content_id, status FROM tickets WHERE id=$1", id
}

func NewTicketRepository() TicketRepository {
	return &ticketRepositorySqlx{}
}
