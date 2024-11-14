package views

type size string

const (
	// Small size
	Small size = "small"
	// Medium size
	Medium size = "medium"
	// Large size
	Large size = "large"
)

type vSize struct {
	size size
	v    int
}

var wSizeBroadcast = NewBroadcaster[vSize]()
var hSizeBroadcast = NewBroadcaster[vSize]()

func broadcastWidthSize(s int) {
	switch {
	case s > 100:
		wSizeBroadcast.broadcast(vSize{size: Large, v: s})
	case s > 80:
		wSizeBroadcast.broadcast(vSize{size: Medium, v: s})
	default:
		wSizeBroadcast.broadcast(vSize{size: Small, v: s})
	}
}

func broadcastHeightSize(s int) {
	switch {
	case s > 25:
		hSizeBroadcast.broadcast(vSize{size: Large, v: s})
	case s > 20:
		hSizeBroadcast.broadcast(vSize{size: Medium, v: s})
	default:
		hSizeBroadcast.broadcast(vSize{size: Small, v: s})
	}
}

func subscribeWidthResize(consumer func(vSize)) {
	ch := wSizeBroadcast.subscribe()
	go func(ch chan vSize) {
		for value := range ch {
			consumer(value)
		}
	}(ch)
}

func subscribeHeightResize(consumer func(vSize)) {
	ch := hSizeBroadcast.subscribe()
	go func(ch chan vSize) {
		for value := range ch {
			consumer(value)
		}
	}(ch)
}
