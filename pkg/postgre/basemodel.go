package postgre

import (
	"fmt"
	"strings"
	"time"
)

type BaseModel struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Args map[string]interface{}

func (a Args) String() string {
	var (
		sb    strings.Builder
		count int
	)

	for key, val := range a {
		sb.WriteString(fmt.Sprintf("%s = %s", key, val))
		if count < len(a) {
			sb.WriteString(", ")
		}

		count++
	}

	return sb.String()
}

var BaseModelFields = []string{"id", "created_at", "updated_at"}
