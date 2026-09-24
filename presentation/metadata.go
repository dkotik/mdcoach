package presentation

import "time"

type Metadata struct {
	ID          string
	Title       string
	Description string
	Author      string
	Created     time.Time
}
