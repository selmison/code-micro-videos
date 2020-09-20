package service

import "github.com/selmison/code-micro-videos/pkg/crud/domain"

type CastMemberDTO struct {
	Name string                `json:"name" validate:"not_blank"`
	Type domain.CastMemberType `json:"type"`
}
