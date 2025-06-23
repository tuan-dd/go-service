package dtos

import "time"

type (
	PaginationReq struct {
		PageIndex int32
		PageSize  int32
	}

	PaginationRes struct {
		PageIndex int32
		PageSize  int32
		Total     int32
	}

	ModifiedEntity struct {
		CreatedBy string
		UpdatedBy string
		CreatedAt *time.Time
		UpdatedAt *time.Time
	}
)
