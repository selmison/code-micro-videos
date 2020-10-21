package testdata

import (
	"github.com/bxcodec/faker/v3"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
)

var (
	FakeExistentCastMember, _ = cast_member.NewCastMember(
		faker.UUIDHyphenated(),
		cast_member.NewCastMemberDTO{
			Name: faker.FirstName(),
		},
	)
)
