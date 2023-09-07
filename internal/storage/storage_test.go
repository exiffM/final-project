package storage

import (
	types "final-project/internal/statistics"
	"testing"

	"github.com/stretchr/testify/require"
)

var storage *Storage

func init() {
	storage = NewStorage()
}

func TestAppend(t *testing.T) {
	st := types.Statistic{}
	storage.Append(st)
	require.Equal(t, 1, storage.Len(), "Length is not 1 after appending 1 element")
}

func TestPullOut(t *testing.T) {
	storage.Clear()
	// Append 5 elems in storage
	for i := 0; i < 5; i++ {
		st := &types.Statistic{
			ASLStat:  types.NewASLS(),
			ACLStat:  types.NewACLS(),
			DIStat:   make(types.DiskInfoStats),
			FSDIStat: types.NewFSDIS(0),
			TTStat:   &types.TopTalkersStat{},
			NStat:    types.NewNetStat(0),
		}
		st.ASLStat.OneMinLoad = 1 + float64(i)
		st.ASLStat.FiveMinLoad = 5 + float64(i)
		st.ASLStat.QuaterLoad = 15 + float64(i)
		storage.Append(*st)
	}
	result := storage.PullOut(5)
	require.ElementsMatch(t, result, storage.stats, "Elements missmatch!")
	// Append 5 more
	tempTail := make([]types.Statistic, 0)
	for i := 5; i < 10; i++ {
		st := &types.Statistic{
			ASLStat:  types.NewASLS(),
			ACLStat:  types.NewACLS(),
			DIStat:   make(types.DiskInfoStats),
			FSDIStat: types.NewFSDIS(0),
			TTStat:   &types.TopTalkersStat{},
			NStat:    types.NewNetStat(0),
		}
		st.ASLStat.OneMinLoad = 1 + float64(i)
		st.ASLStat.FiveMinLoad = 5 + float64(i)
		st.ASLStat.QuaterLoad = 15 + float64(i)
		tempTail = append(tempTail, *st)
		storage.Append(*st)
	}
	result = storage.PullOut(5)
	require.ElementsMatch(t, result, tempTail, "Elements missmatch!")
	// Pull more than in storage (now we have 10 in storage, pulling 15, should got 10)
	result = storage.PullOut(15)
	require.ElementsMatch(t, result, storage.stats, "Elements missmatch!")
}

func TestClear(t *testing.T) {
	storage.Clear()
	require.Zero(t, storage.Len(), "Length of storage is not zero!")
}
