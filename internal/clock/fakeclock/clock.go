// Package fakeclock 擬似的な時刻を返す clock.Clock 実装。
package fakeclock

import (
	"sync"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/clock"
)

var _ clock.Clock = (*Clock)(nil)

// Clock 固定の時刻を返す clock.Clock 実装。
type Clock struct {
	mu  sync.RWMutex
	now time.Time
}

// New 指定した時刻で固定された Clock を生成して返す。
func New(t time.Time) *Clock {
	return &Clock{
		now: t,
	}
}

// Now 固定された時刻を返す。
func (s *Clock) Now() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.now
}

// SetTime 時刻を指定した値に変更する。
func (s *Clock) SetTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.now = t
}
