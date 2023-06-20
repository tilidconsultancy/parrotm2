package domain

import "github.com/google/uuid"

type (
	PercentageLabelMeaning struct {
		LabelMeaning LabelMeaning
		Percentage   uint8
	}
	LabelMeaning struct {
		OwnerId     uuid.UUID
		Label       string
		Description string
	}
)
