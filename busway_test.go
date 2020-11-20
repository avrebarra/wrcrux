package busway_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

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
	if writer.countWrites == writer.successfulWrites {
		return 0, errors.New("Artificial error")
	}

	writer.countWrites++
	return writer.Write(buffer)
}

// zeroWriter always writes zero bytes.
type zeroWriter struct {
	io.Writer
}

func (writer *zeroWriter) Write(buffer []byte) (int, error) {
	return 0, nil
}

func TestBasic(t *testing.T) {
	var err error

	fileWriter := busway.File("hello.log")
	defer os.Remove("hello.log")
	defer fileWriter.Close()

	fileWriter2 := busway.File("hello2.log")
	defer os.Remove("hello2.log")
	defer fileWriter2.Close()

	fileWriter3 := busway.File("hello3.log")
	defer os.Remove("hello3.log")
	defer fileWriter3.Close()

	errorWriter := &writerWithError{
		Writer: fileWriter2,
	}

	zero := &zeroWriter{
		Writer: fileWriter3,
	}

	hello := busway.New(busway.Config{})

	hello.AddWriter(fileWriter)
	hello.AddWriter(errorWriter)
	hello.AddWriter(zero)

	// untagged
	_, err = hello.WriteRich(0, []byte(fmt.Sprintf(
		"Info message %d %f %f %s",
		1,
		float32(3.14),
		3.14,
		"some text",
	)))
	if err != nil {
		panic(err)
	}

	// immediate
	_, err = hello.WriteRich(busway.TImmedate, []byte(fmt.Sprintf(
		"Info message %d %f %f %s",
		1,
		float32(3.14),
		3.14,
		"some text",
	)))
	if err != nil {
		panic(err)
	}

	// writer
	_, err = hello.Write([]byte(fmt.Sprintf(
		"Info message %d %f %f %s",
		1,
		float32(3.14),
		3.14,
		"some text",
	)))
	if err != nil {
		panic(err)
	}

	time.Sleep(500 * time.Millisecond)
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

	hello := busway.New(busway.Config{})
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

	hello := busway.New(busway.Config{})
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

	hello := busway.New(busway.Config{})
	hello.AddWriter(busway.File("hello.log"))

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hello.WriteRich(busway.TImmedate, []byte("Hello World"))
		}
	})
}
