package server

import "testing"

func TestExport2Excel(t *testing.T) {
	InitDB("test.db")
	Export2Excel("bec5da0c-aa51-4844-956b-ab4c35df0d8a")
}
