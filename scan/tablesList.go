package scan

import (
	"context"
	"slices"
	"sync"
)

type TableList struct {
	TableNames []string
}

func (tl *TableList) Add(tableName string) {
	if !slices.Contains(tl.TableNames, tableName) {
		tl.TableNames = append(tl.TableNames, tableName)
	}

	slices.Sort(tl.TableNames)
}

func (tl *TableList) Describe(client TableClient) ([]TableInfo, error) {
	res := make([]TableInfo, 0)

	errCh := make(chan error)
	tfCh := make(chan *TableInfo)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}
	for _, n := range tl.TableNames {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ti, err := client.GetTableInfo(context.Background(), n)
			if err != nil {
				errCh <- err
				return
			}

			tfCh <- ti
		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return nil, err
		case ti := <-tfCh:
			res = append(res, *ti)
		case <-doneCh:
			return res, nil
		}
	}
}
