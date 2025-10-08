package models

import (
	"time"
)

type Student struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Age       int       `json:"age"`
    Grade     string    `json:"grade"`
    CreatedAt time.Time `json:"created_at"`
}