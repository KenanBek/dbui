// Package dummy implements a dummy data source used for the demo purposes.
package dummy

import "fmt"

// Dummy exported.
type Dummy struct {
}

// ListSchemas exported.
func (Dummy) ListSchemas() []string {
	return []string{
		"omni",
		"omega",
		"beta",
	}
}

// ListTables exported.
func (Dummy) ListTables(schema string) []string {
	return []string{
		fmt.Sprintf("%s_table1", schema),
		fmt.Sprintf("%s_table2", schema),
		fmt.Sprintf("%s_table3", schema),
		fmt.Sprintf("%s_table4", schema),
		fmt.Sprintf("%s_table5", schema),
	}
}

// PreviewTable exported.
func (Dummy) PreviewTable(schema, table string) [][]string {
	return [][]string{
		{"abc", "adc"},
		{"bbc", "bdc"},
	}
}

// DescribeTable exported.
func (Dummy) DescribeTable(schema, table string) [][]string {
	return [][]string{}
}

// Query exported.
func (Dummy) Query(schema string) [][]string {
	return [][]string{
		{"qabc", "qadc"},
		{"qbbc", "qbdc"},
	}
}
