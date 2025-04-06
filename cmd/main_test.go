package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	csvInput := strings.NewReader(`given_name,family_name,birthday
John,Doe,1990-01-01
Jane,Smith,1985-05-15
Alice,Johnson,1992-07-20
Bob,Brown,1988-11-30
Charlie,Wilson,1995-03-10`)

	expectedJSONOutput := `[{"birthday":"1990-01-01","family_name":"Doe","given_name":"John"},{"birthday":"1985-05-15","family_name":"Smith","given_name":"Jane"},{"birthday":"1992-07-20","family_name":"Johnson","given_name":"Alice"},{"birthday":"1988-11-30","family_name":"Brown","given_name":"Bob"},{"birthday":"1995-03-10","family_name":"Wilson","given_name":"Charlie"}]`

	actualJSONOutput := &strings.Builder{}

	StartConversion(csvInput, actualJSONOutput, 2)

	require.JSONEq(t, expectedJSONOutput, actualJSONOutput.String(), "The JSON output does not match the expected output")
}
