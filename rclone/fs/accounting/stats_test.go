package accounting

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/fserrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestETA(t *testing.T) {
	for _, test := range []struct {
		size, total int64
		rate        float64
		wantETA     time.Duration
		wantOK      bool
		wantString  string
	}{
		// Custom String Cases
		{size: 0, total: 365 * 86400, rate: 1.0, wantETA: 365 * 86400 * time.Second, wantOK: true, wantString: "1y"},
		{size: 0, total: 7 * 86400, rate: 1.0, wantETA: 7 * 86400 * time.Second, wantOK: true, wantString: "1w"},
		{size: 0, total: 1 * 86400, rate: 1.0, wantETA: 1 * 86400 * time.Second, wantOK: true, wantString: "1d"},
		{size: 0, total: 1110 * 86400, rate: 1.0, wantETA: 1110 * 86400 * time.Second, wantOK: true, wantString: "3y2w1d"},
		{size: 0, total: 15 * 86400, rate: 1.0, wantETA: 15 * 86400 * time.Second, wantOK: true, wantString: "2w1d"},
		// Composite Custom String Cases
		{size: 0, total: 1.5 * 86400, rate: 1.0, wantETA: 1.5 * 86400 * time.Second, wantOK: true, wantString: "1d12h"},
		{size: 0, total: 95000, rate: 1.0, wantETA: 95000 * time.Second, wantOK: true, wantString: "1d2h23m20s"},
		// Standard Duration String Cases
		{size: 0, total: 100, rate: 1.0, wantETA: 100 * time.Second, wantOK: true, wantString: "1m40s"},
		{size: 50, total: 100, rate: 1.0, wantETA: 50 * time.Second, wantOK: true, wantString: "50s"},
		{size: 100, total: 100, rate: 1.0, wantETA: 0 * time.Second, wantOK: true, wantString: "0s"},
		// No String Cases
		{size: -1, total: 100, rate: 1.0, wantETA: 0, wantOK: false, wantString: "-"},
		{size: 200, total: 100, rate: 1.0, wantETA: 0, wantOK: false, wantString: "-"},
		{size: 10, total: -1, rate: 1.0, wantETA: 0, wantOK: false, wantString: "-"},
		{size: 10, total: 20, rate: 0.0, wantETA: 0, wantOK: false, wantString: "-"},
		{size: 10, total: 20, rate: -1.0, wantETA: 0, wantOK: false, wantString: "-"},
		{size: 0, total: 0, rate: 1.0, wantETA: 0, wantOK: false, wantString: "-"},
	} {
		t.Run(fmt.Sprintf("size=%d/total=%d/rate=%f", test.size, test.total, test.rate), func(t *testing.T) {
			gotETA, gotOK := eta(test.size, test.total, test.rate)
			assert.Equal(t, test.wantETA, gotETA)
			assert.Equal(t, test.wantOK, gotOK)
			gotString := etaString(test.size, test.total, test.rate)
			assert.Equal(t, test.wantString, gotString)
		})
	}
}

func TestPercentage(t *testing.T) {
	assert.Equal(t, percent(0, 1000), "0%")
	assert.Equal(t, percent(1, 1000), "0%")
	assert.Equal(t, percent(9, 1000), "1%")
	assert.Equal(t, percent(500, 1000), "50%")
	assert.Equal(t, percent(1000, 1000), "100%")
	assert.Equal(t, percent(1e8, 1e9), "10%")
	assert.Equal(t, percent(1e8, 1e9), "10%")
	assert.Equal(t, percent(0, 0), "-")
	assert.Equal(t, percent(100, -100), "-")
	assert.Equal(t, percent(-100, 100), "-")
	assert.Equal(t, percent(-100, -100), "-")
}

func TestStatsError(t *testing.T) {
	s := NewStats()
	assert.Equal(t, int64(0), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.False(t, s.HadRetryError())
	assert.Equal(t, time.Time{}, s.RetryAfter())
	assert.Equal(t, nil, s.GetLastError())
	assert.False(t, s.Errored())

	t0 := time.Now()
	t1 := t0.Add(time.Second)

	s.Error(nil)
	assert.Equal(t, int64(0), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.False(t, s.HadRetryError())
	assert.Equal(t, time.Time{}, s.RetryAfter())
	assert.Equal(t, nil, s.GetLastError())
	assert.False(t, s.Errored())

	s.Error(io.EOF)
	assert.Equal(t, int64(1), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.True(t, s.HadRetryError())
	assert.Equal(t, time.Time{}, s.RetryAfter())
	assert.Equal(t, io.EOF, s.GetLastError())
	assert.True(t, s.Errored())

	e := fserrors.ErrorRetryAfter(t0)
	s.Error(e)
	assert.Equal(t, int64(2), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.True(t, s.HadRetryError())
	assert.Equal(t, t0, s.RetryAfter())
	assert.Equal(t, e, s.GetLastError())

	err := errors.Wrap(fserrors.ErrorRetryAfter(t1), "potato")
	s.Error(err)
	assert.Equal(t, int64(3), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.True(t, s.HadRetryError())
	assert.Equal(t, t1, s.RetryAfter())
	assert.Equal(t, t1, fserrors.RetryAfterErrorTime(err))

	s.Error(fserrors.FatalError(io.EOF))
	assert.Equal(t, int64(4), s.GetErrors())
	assert.True(t, s.HadFatalError())
	assert.True(t, s.HadRetryError())
	assert.Equal(t, t1, s.RetryAfter())

	s.ResetErrors()
	assert.Equal(t, int64(0), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.False(t, s.HadRetryError())
	assert.Equal(t, time.Time{}, s.RetryAfter())
	assert.Equal(t, nil, s.GetLastError())
	assert.False(t, s.Errored())

	s.Error(fserrors.NoRetryError(io.EOF))
	assert.Equal(t, int64(1), s.GetErrors())
	assert.False(t, s.HadFatalError())
	assert.False(t, s.HadRetryError())
	assert.Equal(t, time.Time{}, s.RetryAfter())
}

func TestStatsTotalDuration(t *testing.T) {
	startTime := time.Now()
	time1 := startTime.Add(-40 * time.Second)
	time2 := time1.Add(10 * time.Second)
	time3 := time2.Add(10 * time.Second)
	time4 := time3.Add(10 * time.Second)

	t.Run("Single completed transfer", func(t *testing.T) {
		s := NewStats()
		tr1 := &Transfer{
			startedAt:   time1,
			completedAt: time2,
		}
		s.AddTransfer(tr1)

		s.mu.Lock()
		total := s.totalDuration()
		s.mu.Unlock()

		assert.Equal(t, 1, len(s.startedTransfers))
		assert.Equal(t, 10*time.Second, total)
		s.RemoveTransfer(tr1)
		assert.Equal(t, 10*time.Second, total)
		assert.Equal(t, 0, len(s.startedTransfers))
	})

	t.Run("Single uncompleted transfer", func(t *testing.T) {
		s := NewStats()
		tr1 := &Transfer{
			startedAt: time1,
		}
		s.AddTransfer(tr1)

		s.mu.Lock()
		total := s.totalDuration()
		s.mu.Unlock()

		assert.Equal(t, time.Since(time1)/time.Second, total/time.Second)
		s.RemoveTransfer(tr1)
		assert.Equal(t, time.Since(time1)/time.Second, total/time.Second)
	})

	t.Run("Overlapping without ending", func(t *testing.T) {
		s := NewStats()
		tr1 := &Transfer{
			startedAt:   time2,
			completedAt: time3,
		}
		s.AddTransfer(tr1)
		tr2 := &Transfer{
			startedAt:   time2,
			completedAt: time2.Add(time.Second),
		}
		s.AddTransfer(tr2)
		tr3 := &Transfer{
			startedAt:   time1,
			completedAt: time3,
		}
		s.AddTransfer(tr3)
		tr4 := &Transfer{
			startedAt:   time3,
			completedAt: time4,
		}
		s.AddTransfer(tr4)
		tr5 := &Transfer{
			startedAt: time.Now(),
		}
		s.AddTransfer(tr5)

		time.Sleep(time.Millisecond)

		s.mu.Lock()
		total := s.totalDuration()
		s.mu.Unlock()

		assert.Equal(t, time.Duration(30), total/time.Second)
		s.RemoveTransfer(tr1)
		assert.Equal(t, time.Duration(30), total/time.Second)
		s.RemoveTransfer(tr2)
		assert.Equal(t, time.Duration(30), total/time.Second)
		s.RemoveTransfer(tr3)
		assert.Equal(t, time.Duration(30), total/time.Second)
		s.RemoveTransfer(tr4)
		assert.Equal(t, time.Duration(30), total/time.Second)
	})

	t.Run("Mixed completed and uncompleted transfers", func(t *testing.T) {
		s := NewStats()
		s.AddTransfer(&Transfer{
			startedAt:   time1,
			completedAt: time2,
		})
		s.AddTransfer(&Transfer{
			startedAt: time2,
		})
		s.AddTransfer(&Transfer{
			startedAt: time3,
		})
		s.AddTransfer(&Transfer{
			startedAt: time3,
		})

		s.mu.Lock()
		total := s.totalDuration()
		s.mu.Unlock()

		assert.Equal(t, startTime.Sub(time1)/time.Second, total/time.Second)
	})
}

// make time ranges from string description for testing
func makeTimeRanges(t *testing.T, in []string) timeRanges {
	trs := make(timeRanges, len(in))
	for i, Range := range in {
		var start, end int64
		n, err := fmt.Sscanf(Range, "%d-%d", &start, &end)
		require.NoError(t, err)
		require.Equal(t, 2, n)
		trs[i] = timeRange{time.Unix(start, 0), time.Unix(end, 0)}
	}
	return trs
}

func (trs timeRanges) toStrings() (out []string) {
	out = []string{}
	for _, tr := range trs {
		out = append(out, fmt.Sprintf("%d-%d", tr.start.Unix(), tr.end.Unix()))
	}
	return out
}

func TestTimeRangeMerge(t *testing.T) {
	for _, test := range []struct {
		in   []string
		want []string
	}{{
		in:   []string{},
		want: []string{},
	}, {
		in:   []string{"1-2"},
		want: []string{"1-2"},
	}, {
		in:   []string{"1-4", "2-3"},
		want: []string{"1-4"},
	}, {
		in:   []string{"2-3", "1-4"},
		want: []string{"1-4"},
	}, {
		in:   []string{"1-3", "2-4"},
		want: []string{"1-4"},
	}, {
		in:   []string{"2-4", "1-3"},
		want: []string{"1-4"},
	}, {
		in:   []string{"1-2", "2-3"},
		want: []string{"1-3"},
	}, {
		in:   []string{"2-3", "1-2"},
		want: []string{"1-3"},
	}, {
		in:   []string{"1-2", "3-4"},
		want: []string{"1-2", "3-4"},
	}, {
		in:   []string{"1-3", "7-8", "4-6", "2-5", "7-8", "7-8"},
		want: []string{"1-6", "7-8"},
	}} {

		in := makeTimeRanges(t, test.in)
		in.merge()

		got := in.toStrings()
		assert.Equal(t, test.want, got)
	}
}

func TestTimeRangeCull(t *testing.T) {
	for _, test := range []struct {
		in           []string
		cutoff       int64
		want         []string
		wantDuration time.Duration
	}{{
		in:           []string{},
		cutoff:       1,
		want:         []string{},
		wantDuration: 0 * time.Second,
	}, {
		in:           []string{"1-2"},
		cutoff:       1,
		want:         []string{"1-2"},
		wantDuration: 0 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       1,
		want:         []string{"2-5", "7-9"},
		wantDuration: 0 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       4,
		want:         []string{"2-5", "7-9"},
		wantDuration: 0 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       5,
		want:         []string{"7-9"},
		wantDuration: 3 * time.Second,
	}, {
		in:           []string{"2-5", "7-9", "2-5", "2-5"},
		cutoff:       6,
		want:         []string{"7-9"},
		wantDuration: 9 * time.Second,
	}, {
		in:           []string{"7-9", "3-3", "2-5"},
		cutoff:       7,
		want:         []string{"7-9"},
		wantDuration: 3 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       8,
		want:         []string{"7-9"},
		wantDuration: 3 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       9,
		want:         []string{},
		wantDuration: 5 * time.Second,
	}, {
		in:           []string{"2-5", "7-9"},
		cutoff:       10,
		want:         []string{},
		wantDuration: 5 * time.Second,
	}} {

		in := makeTimeRanges(t, test.in)
		cutoff := time.Unix(test.cutoff, 0)
		gotDuration := in.cull(cutoff)

		what := fmt.Sprintf("in=%q, cutoff=%d", test.in, test.cutoff)
		got := in.toStrings()
		assert.Equal(t, test.want, got, what)
		assert.Equal(t, test.wantDuration, gotDuration, what)
	}
}

func TestTimeRangeDuration(t *testing.T) {
	assert.Equal(t, 0*time.Second, timeRanges{}.total())
	assert.Equal(t, 1*time.Second, makeTimeRanges(t, []string{"1-2"}).total())
	assert.Equal(t, 91*time.Second, makeTimeRanges(t, []string{"1-2", "10-100"}).total())
}

func TestPruneTransfers(t *testing.T) {
	max := maxCompletedTransfers + fs.Config.Transfers

	s := NewStats()
	for i := int64(1); i <= int64(max+100); i++ {
		s.AddTransfer(&Transfer{
			startedAt:   time.Unix(i, 0),
			completedAt: time.Unix(i+1, 0),
		})
	}

	s.mu.Lock()
	assert.Equal(t, time.Duration(max+100)*time.Second, s.totalDuration())
	assert.Equal(t, max+100, len(s.startedTransfers))
	s.mu.Unlock()

	for i := 0; i < 200; i++ {
		s.PruneTransfers()
	}

	s.mu.Lock()
	assert.Equal(t, time.Duration(max+100)*time.Second, s.totalDuration())
	assert.Equal(t, max, len(s.startedTransfers))
	s.mu.Unlock()

}