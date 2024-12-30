package bchan

import (
	"testing"
)

const n = 1000000
const m = 100

func TestUnbounded(t *testing.T) {
	in, out := BChan[int]{}.Unbounded()

	go func() {
		for i := 0; i < n; i++ {
			in <- i
		}
	}()

	for i := 0; i < n; i++ {
		if v := <-out; v != i {
			t.Errorf("Expected %d, got %d", i, v)
		}
	}

	go func() {
		for i := 0; i < n; i++ {
			if v := <-out; v != i {
				t.Errorf("Expected %d, got %d", i, v)
			}
		}
	}()

	for i := 0; i < n; i++ {
		in <- i
	}
}

func TestBounded(t *testing.T) {
	in, out := BChan[int]{}.Bounded(n)

	go func() {
		for i := 0; i < n; i++ {
			in <- i
		}
	}()
	for i := 0; i < n; i++ {
		if v := <-out; v != i {
			t.Errorf("Expected %d, got %d", i, v)
		}
	}

	afterPause := false

	go func() {
		for i := 0; i < n; i++ {
			if v := <-out; v != i {
				t.Errorf("Expected %d, got %d", i, v)
			}
		}
		for i := 0; i < n; i++ {
			if !afterPause {
				t.Errorf("Bounded behaved like unbounded")
			}
			if v := <-out; v != i {
				t.Errorf("Expected %d, got %d", i, v)
			}
		}
	}()

	for i := 0; i < n; i++ {
		in <- i
	}
	afterPause = true
	for i := 0; i < n; i++ {
		in <- i
	}
}

func BenchmarkUnbounded(b *testing.B) {
	in, out := BChan[int]{}.Unbounded()

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			in <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-out
			_ = a
		}
	}
}

func BenchmarkBounded(b *testing.B) {
	in, out := BChan[int]{}.Bounded(b.N)

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			in <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-out
			_ = a
		}
	}
}

func BenchmarkCompareDouble(b *testing.B) {
	in, out := make(chan int, b.N), make(chan int, b.N)

	go func() {
		for {
			out <- <-in
		}
	}()

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			in <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-out
			_ = a
		}
	}
}

func BenchmarkCompareTriple(b *testing.B) {
	in, out, buf := make(chan int, b.N), make(chan int, b.N), make(chan int, b.N)

	go func() {
		for {
			buf <- <-in
		}
	}()

	go func() {
		for {
			out <- <-buf
		}
	}()

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			in <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-out
			_ = a
		}
	}
}

func BenchmarkCompareSelect(b *testing.B) {
	in, out, buf := make(chan int, b.N), make(chan int, b.N), make(chan int, b.N)

	go func() {
		var v int
		for {
			select {
			case v = <-in:
				buf <- v
			}
		}
	}()

	go func() {
		var v int
		for {
			select {
			case v = <-buf:
				out <- v
			}
		}
	}()

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			in <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-out
			_ = a
		}
	}
}

func BenchmarkBaseline(b *testing.B) {
	ch := make(chan int, b.N)

	for j := 0; j < m; j++ {
		for i := 0; i < b.N; i++ {
			ch <- i
		}

		for i := 0; i < b.N; i++ {
			a := <-ch
			_ = a
		}
	}
}
