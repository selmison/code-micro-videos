package testdata

import (
	"github.com/bxcodec/faker/v3"

	"github.com/selmison/code-micro-videos/pkg/genre"
)

var (
	FakeExistentGenre      = genre.Genre{Id: faker.UUIDHyphenated(), Name: faker.FirstName()}
	FakeExistentGenreId    = FakeExistentGenre.Id
	FakeNonExistentGenre   = genre.Genre{Id: faker.UUIDHyphenated(), Name: faker.FirstName()}
	FakeNonExistentGenreId = FakeNonExistentGenre.Id
)
