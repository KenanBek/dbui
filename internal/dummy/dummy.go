// Package dummy implements a dummy data source used for the demo purposes.
package dummy

import (
	"errors"
	"fmt"
)

func sptr(s string) *string {
	return &s
}

// Dummy exported.
type Dummy struct {
}

// Ping exported.
func (Dummy) Ping() error {
	return nil
}

// ListSchemas exported.
func (Dummy) ListSchemas() ([]string, error) {
	return []string{
		"omni",
		"errored",
		"beta",
	}, nil
}

// ListTables exported.
func (Dummy) ListTables(schema string) ([]string, error) {
	if schema == "errored" {
		return nil, errors.New("failed to load the list of tables")
	}

	return []string{
		fmt.Sprintf("%s_table1", schema),
		fmt.Sprintf("%s_table2", schema),
		fmt.Sprintf("%s_table3", schema),
		fmt.Sprintf("%s_table4", schema),
		fmt.Sprintf("%s_table5", schema),
	}, nil
}

// PreviewTable exported.
func (Dummy) PreviewTable(schema, table string) ([][]*string, error) {
	return [][]*string{
		{sptr("abc"), sptr("adc")},
		{sptr("bbc"), sptr("bdc")},
	}, nil
}

// DescribeTable exported.
func (Dummy) DescribeTable(schema, table string) ([][]*string, error) {
	return [][]*string{
		{sptr("header1"), sptr("header2")},
		{sptr("val1"), sptr("val2")},
	}, nil
}

// Query exported.
func (Dummy) Query(schema, query string) ([][]*string, error) {
	return [][]*string{
		{sptr("header1"), sptr("header2")},
		{sptr("val1"), sptr("val2")},
	}, nil
}
