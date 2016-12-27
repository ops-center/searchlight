package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient()
	assert.Nil(t, err)
}
