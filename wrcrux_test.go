package wrcrux_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/avrebarra/wrcrux"
)

// filewriter provides a simple interface to create a log file.
// The given file path must be writable, otherwise it will panic.
func filewriter(path string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	return file
}

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

	fileWriter := filewriter("hello.log")
	defer os.Remove("hello.log")
	defer fileWriter.Close()

	errorWriter := &writerWithError{
		Writer: fileWriter,
	}

	zero := &zeroWriter{
		Writer: fileWriter,
	}

	t.Run("ok", func(t *testing.T) {
		bus := wrcrux.NewWux(wrcrux.ConfigWux{})
		bus.AddWriter(fileWriter)

		// untagged
		if _, err = bus.WriteRich(wrcrux.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(wrcrux.BImmediate, []byte(datastr)); err != nil {
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
		bus := wrcrux.NewWux(wrcrux.ConfigWux{})
		bus.AddWriter(errorWriter)

		// untagged
		if _, err = bus.WriteRich(wrcrux.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(wrcrux.BImmediate, []byte(datastr)); err == nil {
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
		bus := wrcrux.NewWux(wrcrux.ConfigWux{})
		bus.AddWriter(zero)

		// untagged
		if _, err = bus.WriteRich(wrcrux.BNormal, []byte(datastr)); err != nil {
			t.Fatal("case failed")
		}

		// immediate
		if _, err = bus.WriteRich(wrcrux.BImmediate, []byte(datastr)); err == nil {
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
		if err == nil {
			panic(err)
		}
	}()

	filewriter("")
}

func BenchmarkWriter(b *testing.B) {
	defer os.Remove("hello.log")

	hello := wrcrux.NewWux(wrcrux.ConfigWux{})
	hello.AddWriter(filewriter("hello.log"))

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

	hello := wrcrux.NewWux(wrcrux.ConfigWux{})
	hello.AddWriter(filewriter("hello.log"))

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

	hello := wrcrux.NewWux(wrcrux.ConfigWux{})
	hello.AddWriter(filewriter("hello.log"))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hello.WriteRich(wrcrux.BImmediate, []byte("Hello World"))
		}
	})
}
