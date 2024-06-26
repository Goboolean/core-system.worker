package chanutil

// DummyChannelConsumer consumes data from the 'in' channel until it is closed.
// If DummyChannelConsumer is called with a closed channel, it exits immediately.
func DummyChannelConsumer[T any](in chan T) {
	for range in {
	}
}
