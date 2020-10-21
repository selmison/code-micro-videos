package testdata

import (
	"context"

	"github.com/bxcodec/faker/v3"

	"github.com/selmison/code-micro-videos/pkg/category"
)

var (
	FakeCtx                   = context.Background()
	FakeId                    = faker.UUIDHyphenated()
	FakeName                  = faker.FirstName()
	FakeDesc                  = faker.Sentence()
	FakeExistentCategory      = category.Category{Id: faker.UUIDHyphenated(), Name: faker.FirstName()}
	FakeExistentCategoryId    = FakeExistentCategory.Id
	FakeNonExistentCategory   = category.Category{Id: faker.UUIDHyphenated(), Name: faker.FirstName()}
	FakeNonExistentCategoryId = FakeNonExistentCategory.Id
)
