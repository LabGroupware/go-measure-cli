// Package clock 時刻の固定を行えるようにするためのパッケージ。
package clock

import "time"

// Clock 現在時刻を扱うインターフェース。
// 時刻に基づいた処理を行う際にこのインターフェースを使用することで、
// テストなどでの時刻の偽装が容易になる。
type Clock interface {
	// Now 現在時刻を返す。
	Now() time.Time
}

var clk = &RealClock{}

// New 実際の時刻を扱う Clock を返す。
func New() Clock {
	return clk
}

// RealClock time.Now() による Clock 実装。
type RealClock struct{}

// Now 現在時刻を返す。
func (s *RealClock) Now() time.Time {
	return time.Now()
}
