package scan

import "slices"

type TableList struct {
	TableNames []string
}

func (tl *TableList) Add(tableName string) {
	if !slices.Contains(tl.TableNames, tableName) {
		tl.TableNames = append(tl.TableNames, tableName)
	}

	slices.Sort(tl.TableNames)
}
