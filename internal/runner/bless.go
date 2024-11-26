package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/dogmatiq/aureus/internal/test"
)

// BlessStrategy is a strategy for accepting failed test output as the new
// expectation, known as "blessing" the output.
type BlessStrategy int

const (
	// BlessAvailable is a [BlessStrategy] that instructs the user that blessing
	// may be activated by using the -aureus.bless flag on the command line.
	BlessAvailable BlessStrategy = iota

	// BlessEnabled is a [BlessStrategy] that explicitly enables blessing of
	// failed tests.
	BlessEnabled

	// BlessDisabled is a [BlessStrategy] that explicitly disables blessing of
	// failed tests.
	BlessDisabled
)

func bless(a test.Assertion, r *os.File) error {
	// TODO: Tests are loaded using an [fs.FS], so the file system is not
	// necessarily the host file system.
	//
	// Ultimately, we probably need to make it the loader's responsibility to
	// bless tests, since it is the loader that knows where the tests came from.

	w, err := os.OpenFile(a.Output.File, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file containing expected output: %w", err)
	}
	defer w.Close()

	if a.Output.IsEntireFile() {
		n, err := io.Copy(w, r)
		if err != nil {
			return err
		}
		return w.Truncate(n)
	}

	if err := resize(a, r, w); err != nil {
		return fmt.Errorf("unable to resize expected output: %w", err)
	}

	if _, err := w.Seek(a.Output.Begin, io.SeekStart); err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	return err
}

func resize(a test.Assertion, r, w *os.File) error {
	sizeExpected := a.Output.End - a.Output.Begin
	sizeActual, err := fileSize(r)
	if err != nil {
		return err
	}

	if sizeExpected == sizeActual {
		return nil
	}

	sizeBefore, err := fileSize(w)
	if err != nil {
		return err
	}

	sizeAfter := sizeBefore - sizeExpected + sizeActual

	op := shrink
	if sizeAfter > sizeBefore {
		op = grow
	}

	return op(
		w,
		a.Output.End,
		sizeBefore,
		sizeAfter,
	)
}

func shrink(w *os.File, endOfContent, fileLengthBefore, fileLengthAfter int64) error {
	moveDistance := fileLengthBefore - fileLengthAfter
	buf := make([]byte, 4096)

	for {
		n, err := w.ReadAt(buf, endOfContent)
		if err != nil && err != io.EOF {
			return err
		}

		if n > 0 {
			if _, err := w.WriteAt(buf[:n], endOfContent-moveDistance); err != nil {
				return err
			}

			endOfContent += int64(n)
		}

		if err == io.EOF {
			return w.Truncate(fileLengthAfter)
		}
	}
}

func grow(w *os.File, endOfContent, fileLengthBefore, fileLengthAfter int64) error {
	moveDistance := fileLengthAfter - fileLengthBefore
	moveLength := fileLengthBefore - endOfContent
	bufSize := min(moveLength, 4096)
	buf := make([]byte, bufSize)

	// (1) move the partial chunk that doesn't fill the entire buffer first.
	chunkSize := moveLength % bufSize
	if chunkSize == 0 {
		chunkSize = bufSize
	}

	for cursor := fileLengthBefore - chunkSize; cursor >= endOfContent; cursor -= chunkSize {
		_, err := w.ReadAt(buf[:chunkSize], cursor)
		if err != nil {
			return err
		}

		if _, err := w.WriteAt(buf[:chunkSize], cursor+moveDistance); err != nil {
			return err
		}

		// (2) move the rest of the content in chunks of the full buffer size.
		chunkSize = bufSize
	}

	return nil
}

func fileSize(f *os.File) (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("unable to determine file size of %s: %w", f.Name(), err)
	}
	return stat.Size(), nil
}
