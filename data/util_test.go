package data_test

import (
	"fmt"
	"log"
	"testing"

	. "github.com/appscode/searchlight/data"
	"github.com/stretchr/testify/assert"
)

func TestIcingaData(t *testing.T) {
	ic, err := LoadIcingaData()
	if err != nil {
		log.Fatal(err)
	}
	assert.NotZero(t, len(ic.Command), "No check agent found")
	fmt.Println(ic.Command[0].Name)
	fmt.Println(ic.Command[0].Description)
}
