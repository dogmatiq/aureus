package streamdiff

import (
	"bufio"
	"bytes"
	"io"
)

// Diff writes a diff of the content read from x and y to w.
//
// If x and y produce the same content no data is written to w and eq is true.
func Diff(
	w io.Writer,
	x string, rx io.Reader,
	y string, ry io.Reader,
	context int,
) (eq bool, err error) {
	readerX := bufio.NewReader(rx)
	readerY := bufio.NewReader(ry)

	ringX := newRing[string](context)
	ringY := newRing[string](context)
	var countX, countY int

	for {
		lineX, okX, err := readLine(readerX)
		if err != nil {
			return false, err
		}

		lineY, okY, err := readLine(readerY)
		if err != nil {
			return false, err
		}

		if okX {
			ringX.Push(string(lineX))
			countX++
		}

		if okY {
			ringY.Push(string(lineY))
			countY++
		}

		if !okX && !okY {
			return true, nil
		}

		if !okX || !okY || !bytes.Equal(lineX, lineY) {
			if _, err := io.WriteString(w, "--- "); err != nil {
				return false, err
			}
			if _, err := io.WriteString(w, x); err != nil {
				return false, err
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return false, err
			}
			if _, err := io.WriteString(w, "+++ "); err != nil {
				return false, err
			}
			if _, err := io.WriteString(w, y); err != nil {
				return false, err
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return false, err
			}
			break
		}
	}

	linesX := ringX.Slice()
	linesY := ringY.Slice()

	for i := 0; i < context; i++ {
		lineX, okX, err := readLine(readerX)
		if err != nil {
			return false, err
		}

		lineY, okY, err := readLine(readerY)
		if err != nil {
			return false, err
		}

		if !okX && !okY {
			break
		}

		if okX {
			linesX = append(linesX, string(lineX))
			countX++
		}

		if okY {
			linesY = append(linesY, string(lineY))
			countY++
		}
	}

	return false, diff(
		w,
		linesX,
		linesY,
		countX-len(linesX),
		countY-len(linesY),
	)
}

func readLine(r *bufio.Reader) ([]byte, bool, error) {
	line, err := r.ReadBytes('\n')
	if err == io.EOF {
		return line, len(line) > 0, nil
	}
	if err != nil {
		return nil, false, err
	}
	return line, true, nil
}
