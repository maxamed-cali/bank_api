package dtos

import "time"

type UserWithRoleDTO struct {
	ID          uint      `json:"id"`
    Name      string    `json:"name"`
    Phone     string    `json:"phone"`
    Role      string    `json:"role"`
    Status    bool      `json:"status"`
    CreatedAt time.Time `json:"created"`
}
