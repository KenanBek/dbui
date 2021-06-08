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
		"demo_omni",
		"demo_errored",
		"demo_beta",
	}, nil
}

// ListTables exported.
func (Dummy) ListTables(schema string) ([]string, error) {
	if schema == "demo_errored" {
		return nil, errors.New("failed to load the list of tables")
	}

	return []string{
		fmt.Sprintf("demo_%s_table1", schema),
		fmt.Sprintf("demo_%s_table2", schema),
		fmt.Sprintf("demo_%s_table3", schema),
		fmt.Sprintf("demo_%s_table4", schema),
		fmt.Sprintf("demo_%s_table5", schema),
		fmt.Sprintf("demo_%s_table6", schema),
		fmt.Sprintf("demo_%s_table7", schema),
		fmt.Sprintf("demo_%s_table8", schema),
		fmt.Sprintf("demo_%s_table9", schema),
		fmt.Sprintf("demo_%s_table10", schema),
		fmt.Sprintf("demo_%s_table11", schema),
	}, nil
}

// PreviewTable exported.
func (Dummy) PreviewTable(schema, table string) ([][]*string, error) {
	return [][]*string{
		{sptr("Name"), sptr("Surname"), sptr("Department"), sptr("Position")},
		{sptr("Alex"), sptr("Doe"), sptr("IT"), sptr("Cool")},
		{sptr("Bob"), sptr("Excellent"), sptr("Finance"), sptr("Cool")},
		{sptr("Joe"), sptr("Cool"), sptr("Growth"), sptr("Cool")},
	}, nil
}

// DescribeTable exported.
func (Dummy) DescribeTable(schema, table string) ([][]*string, error) {
	return [][]*string{
		{sptr("Column Name"), sptr("Column Type"), sptr("Size")},
		{sptr("Name"), sptr("string"), sptr("12")},
		{sptr("Surname"), sptr("string"), sptr("12")},
		{sptr("Department"), sptr("string"), sptr("12")},
		{sptr("Position"), sptr("string"), sptr("12")},
	}, nil
}

// Query exported.
func (Dummy) Query(schema, query string) ([][]*string, error) {
	return [][]*string{
		{sptr("header1"), sptr("header2")},
		{sptr("val1"), sptr("val2")},
	}, nil
}
