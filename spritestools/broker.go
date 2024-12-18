package spritestools

// initial reference material
// https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel
// https://stackoverflow.com/a/49877632

// A data broker that can be published to, and any entity that has subscribed will get a copy of the message.
// Any subscriber should make sure to unsubscribe before deleting the reference.
// An independent go routine is started to receive published messages and to distribute them to all subscribers.
// The go routine stops once Stop() is called and all subscribers have unsubscribed.
// Ideally this is a 1-to-n communication system. If there are multiple publishers, it would be hard to know when to call Stop.
type Broker[T any] struct {
	bufferSize    int
	stopCh        chan struct{}
	publishCh     chan T
	subscribeCh   chan chan T
	unsubscribeCh chan chan T
}

// Creates a new Broker of any given type and start it running.
func NewBroker[T any](bufferSize int) *Broker[T] {
	ret := &Broker[T]{
		bufferSize:    bufferSize,
		stopCh:        make(chan struct{}),
		publishCh:     make(chan T, bufferSize),
		subscribeCh:   make(chan chan T, bufferSize),
		unsubscribeCh: make(chan chan T, bufferSize),
	}

	go ret.start()
	return ret
}

// Should be started as a go routine
func (b *Broker[T]) start() {
	on := true
	subs := map[chan T]struct{}{}
	for {
		select {
		case msg := <-b.publishCh:
			if on {
				for msgCh := range subs {
					// msgCh is buffered, use non-blocking send to protect the broker:
					select {
					case msgCh <- msg:
						// The message has been sent
					default:
						// The channel is full, and has been skipped.
					}
				}
			}
		case msgCh := <-b.subscribeCh:
			if on {
				subs[msgCh] = struct{}{}
			}
		case msgCh := <-b.unsubscribeCh:
			if _, found := subs[msgCh]; !found {
				continue
			}
			if on {
				close(msgCh)
			}
			delete(subs, msgCh)
			if !on && len(subs) == 0 {
				return
			}
		case <-b.stopCh:
			on = false
			for msgCh := range subs {
				close(msgCh)
			}
			if len(subs) == 0 {
				return
			}
		}
	}
}

// Stop the data broker. Only call this once. All subscriptions will be closed
func (b *Broker[T]) Stop() {
	b.stopCh <- struct{}{}
}

// Subscribe to the data broker. Messages can be received on the returned channel.
func (b *Broker[T]) Subscribe() chan T {
	msgCh := make(chan T, b.bufferSize)
	b.subscribeCh <- msgCh
	return msgCh
}

// Unsubscribe from the data broker. The reference to the channel can be safely discarded after calling this.
func (b *Broker[T]) Unsubscribe(msgCh chan T) {
	b.unsubscribeCh <- msgCh
}

// Send a message to every current subscriber.
func (b *Broker[T]) Publish(msg T) {
	b.publishCh <- msg
}
