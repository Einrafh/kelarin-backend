package models

import "time"

// Card represents a card in a Kanban list.
type Card struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	ListID      uint       `gorm:"not null" json:"list_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// The list this card belongs to.
	List        BoardList        `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE" json:"-"`
	Subtasks    []Subtask        `gorm:"foreignKey:CardID" json:"subtasks"`
	Assignees   []CardAssignee   `gorm:"foreignKey:CardID" json:"assignees"`
	Attachments []CardAttachment `gorm:"foreignKey:CardID" json:"attachments"`
	Labels      []CardLabel      `gorm:"foreignKey:CardID" json:"labels"`
	Comments    []CardComment    `gorm:"foreignKey:CardID" json:"comments"`
}
