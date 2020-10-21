package test

import (
	"testing"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/storage/inmem"
)

func SetupTestCase(t *testing.T, fakes interface{}) (func(t *testing.T), cast_member.Repository, error) {
	return inmem.SetupTestCase(t, fakes)
}
