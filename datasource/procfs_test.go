package datasource

import (
	"testing"
)

func TestGetCurrentIRQStat(t *testing.T) {
	stat, err := GetCurrentIRQStat()
	if err != nil {
		t.Fatalf("%v\n", err)
		return
	}
	t.Logf("%v\n", stat)
}
