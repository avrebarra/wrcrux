package busway_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/akyoto/assert"
	"github.com/avrebarra/busway"
)

// writerWithError errors the Write call after `successfulWrites` writes.
type writerWithError struct {
	io.Writer

	countWrites      int
	successfulWrites int
}

func (writer *writerWithError) Write(buffer []byte) (int, error) {
	return 0, errors.New("Artificial error")
}

// zeroWriter always writes zero bytes.
type zeroWriter struct {
	io.Writer
}

func (writer *zeroWriter) Write(buffer []byte) (int, error) {
	return 0, nil
}

func TestBusway(t *testing.T) {
	var err error
	datastr := fmt.Sprintf("Info message %d %f %f %s", 1, float32(3.14), 3.14, "some text")

	fileWriter := busway.File("hello.log")
	defer os.Remove("hello.log")
	defer fileWriter.Close()

	errorWriter := &writerWithError{
		Writer: fileWriter,
	}

	zero := &zeroWriter{
		Writer: fileWriter,
	}

	t.Run("ok", func(t *testing.T) {
		bus := busway.NewWux(busway.ConfigWux{})
		bus.AddWriter(fileWriter)

		// untagged
		if _, err = bus.WriteRich(busway.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(busway.BImmediate, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// writer
		if _, err = bus.Write([]byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// close
		bus.Close()

		// write after close
		if _, err = bus.Write([]byte(datastr)); err == nil {
			fmt.Println(err)
			t.Fatal("case failed")
		}
	})

	t.Run("error write", func(t *testing.T) {
		bus := busway.NewWux(busway.ConfigWux{})
		bus.AddWriter(errorWriter)

		// untagged
		if _, err = bus.WriteRich(busway.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(busway.BImmediate, []byte(datastr)); err == nil {
			t.Fatal("case failed")
		}

		// writer
		if _, err = bus.Write([]byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// close
		bus.Close()

		// write after close
		if _, err = bus.Write([]byte(datastr)); err == nil {
			t.Fatal("case failed")
		}
	})

	t.Run("error incomplete write", func(t *testing.T) {
		bus := busway.NewWux(busway.ConfigWux{})
		bus.AddWriter(zero)

		// untagged
		if _, err = bus.WriteRich(busway.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(busway.BImmediate, []byte(datastr)); err == nil {
			t.Fatal("case failed")
		}

		// writer
		if _, err = bus.Write([]byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// close
		bus.Close()

		// write after close
		if _, err = bus.Write([]byte(datastr)); err == nil {
			t.Fatal("case failed")
		}
	})
}

func TestInvalidFilePath(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()

	busway.File("")
}

func BenchmarkWriter(b *testing.B) {
	defer os.Remove("hello.log")

	hello := busway.NewWux(busway.ConfigWux{})
	hello.AddWriter(busway.File("hello.log"))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hello.Write([]byte("Hello World"))
		}
	})
}

func BenchmarkWriteRich(b *testing.B) {
	defer os.Remove("hello.log")

	hello := busway.NewWux(busway.ConfigWux{})
	hello.AddWriter(busway.File("hello.log"))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hello.WriteRich(0, []byte("Hello World"))
		}
	})
}

func BenchmarkWriteRichImmediate(b *testing.B) {
	defer os.Remove("hello.log")

	hello := busway.NewWux(busway.ConfigWux{})
	hello.AddWriter(busway.File("hello.log"))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hello.WriteRich(busway.BImmediate, []byte("Hello World"))
		}
	})
}
