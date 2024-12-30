package bchan

type BChan[T any] struct{}

func (_ BChan[T]) Unbounded() (chan T, chan T) {
	wch := make(chan T, 8)
	rch := make(chan T, 8)

	bufs := make(chan chan T, 64)

	go func() {
		buf := make(chan T, 1)
		bufs <- buf
		var v T
		for {
			v = <-wch
			select {
			case buf <- v:
			default:
				buf = make(chan T, cap(buf)*2)
				bufs <- buf
				buf <- v
			}
		}
	}()

	go func() {
		buf := <-bufs
		var b chan T
		var v T
		for {
			select {
			case b = <-bufs:
			loop:
				for {
					select {
					case v = <-buf:
						rch <- v
					default:
						buf = b
						break loop
					}
				}
			case v = <-buf:
				rch <- v
			}
		}
	}()

	return wch, rch
}

func (_ BChan[T]) Bounded(l int) (chan T, chan T) {
	wch := make(chan T, 8)
	rch := make(chan T, 8)

	bufs := make(chan chan T, 64)

	go func() {
		buf := make(chan T, 1)
		bufs <- buf
		var v T
	loop:
		for {
			v = <-wch
			select {
			case buf <- v:
			default:
				if len(buf)*2 > l {
					break loop
				}
				buf = make(chan T, cap(buf)*2)
				bufs <- buf
				buf <- v
			}
		}

		buf = make(chan T, l)
		bufs <- buf
		buf <- v

		for {
			buf <- <-wch
		}

	}()

	go func() {
		buf := <-bufs
		var b chan T
		var v T
		for {
			select {
			case b = <-bufs:
			loop:
				for {
					select {
					case v = <-buf:
						rch <- v
					default:
						buf = b
						break loop
					}
				}
			case v = <-buf:
				rch <- v
			}
		}
	}()

	return wch, rch
}
