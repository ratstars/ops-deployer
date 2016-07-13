package commons

import (
    "testing"
)

func TestResultOutput(t *testing.T) {
	outputSet := make([]ResultOutput, 0)
	std := NewStdoutOutput("OUTPUT")
	outputSet = append(outputSet, std)
	if outputSet[0].String() != "OUTPUT" || outputSet[0].Type() != "INFO" {
		t.Error("StdoutOutput not expect one")
		return
	}
	err := NewStderrOutput("ERROROUTPUT")
	outputSet = append(outputSet, err)
	if outputSet[1].String() != "ERROROUTPUT" || outputSet[1].Type() != "ERROR" {
		t.Error("StderrOutput not expect one")
		return
	}
	custom := NewCustomOutput("DISPLAY", "DISPLAY_CONTENT")
	outputSet = append(outputSet, custom)
	if outputSet[2].String() != "DISPLAY_CONTENT" || outputSet[2].Type() != "DISPLAY" {
		t.Error("CustomOutput not expect one")
		return
	}
	
	result := Merge(outputSet[1], outputSet[2])
	if result.String() != "ERROROUTPUTDISPLAY_CONTENT" || result.Type() != "ERROR" {
		t.Error("Merge result not expect one")
		return
	}

}