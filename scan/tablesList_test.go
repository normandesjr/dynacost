package scan_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/normandesjr/dynacost/scan"
)

func TestTablesList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expNames []string
	}{
		{"Ordered", []string{"t0", "t1", "t2"}, []string{"t0", "t1", "t2"}},
		{"UnOrdered", []string{"t2", "t1", "t0"}, []string{"t0", "t1", "t2"}},
		{"Duplicated", []string{"t1", "t1", "t0"}, []string{"t0", "t1"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tl := &scan.TableList{}
			for _, n := range tc.input {
				tl.Add(n)
			}

			if slices.Compare(tc.expNames, tl.TableNames) != 0 {
				t.Errorf("Expected %q, got %q\n", tc.expNames, tl.TableNames)
			}
		})
	}
}

type mockTableClient struct {
	expError error
}

func (m *mockTableClient) GetTableInfo(context context.Context, tableName string) (*scan.TableInfo, error) {
	if m.expError != nil {
		return nil, m.expError
	}

	return &scan.TableInfo{
		Name: tableName,
	}, nil
}

func TestTableInfo(t *testing.T) {
	testCases := []struct {
		name   string
		expErr error
	}{
		{name: "ListAllTableInfo"},
		{name: "ExpError", expErr: errors.New("Some error")},
	}

	tl := &scan.TableList{}
	tl.Add("t1")
	tl.Add("t2")
	tl.Add("t3")

	for _, tc := range testCases {
		ti, err := tl.Describe(&mockTableClient{expError: tc.expErr})

		if tc.expErr != nil {
			if err == nil {
				t.Error("Expected error here, got nil")
			}
			return
		}

		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}

		for _, n := range ti {
			if !strings.Contains("t1t2t3", n.Name) {
				t.Errorf("Expected name be in t1, t2 or t3 but got %s", n.Name)
			}
		}
	}
}
