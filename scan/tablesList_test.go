package scan_test

import (
	"slices"
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
