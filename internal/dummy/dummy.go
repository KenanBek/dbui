/*
	Dummy implements dummy data provider.
*/
package dummy

import "fmt"

type Dummy struct {
}

func (Dummy) ListSchemas() []string {
	return []string{
		"omni",
		"omega",
		"beta",
	}
}

func (Dummy) ListTables(schema string) []string {
	return []string{
		fmt.Sprintf("%s_table1", schema),
		fmt.Sprintf("%s_table2", schema),
		fmt.Sprintf("%s_table3", schema),
		fmt.Sprintf("%s_table4", schema),
		fmt.Sprintf("%s_table5", schema),
	}
}

func (Dummy) PreviewTable(schema, table string) [][]string {
	return [][]string{
		{"abc", "adc"},
		{"bbc", "bdc"},
	}
}

func (Dummy) DescribeTable(schema, table string) [][]string {
	return [][]string{
	}
}

func (Dummy) Query(schema string) [][]string {
	return [][]string{
		{"qabc", "qadc"},
		{"qbbc", "qbdc"},
	}
}
