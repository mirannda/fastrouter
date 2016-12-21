package fastrouter

import (
	"bytes"
	"strings"
)

// cleanPath cleans the path without allocation.
func cleanPath(path string, buffer *bytes.Buffer) string {
	if path == "" {
		return "/"
	}

	buffer.Reset()
	if path[0] != Slash {
		buffer.WriteByte(Slash)
	}

	var (
		j  int
		k  int
		bs []byte
	)

	n := len(path)
	for i := 0; i < n; {
		if path[i] == Dot {
			for j = i + 1; j < n && path[j] == Dot; j++ {
			}

			if j == n || path[j] == Slash {
				if j-i != 1 {
					bs = buffer.Bytes()

					for k = len(bs) - 2; k > 0 && bs[k] != Slash; k-- {
					}
					if k < 0 {
						k = 0
					}
					buffer.Truncate(k + 1)
				}

				i = j + 1
				continue
			}
		}

		j = strings.IndexByte(path[i:], Slash)
		if j == -1 {
			j = n
		} else {
			j += i + 1
		}

		if !(path[i] == Slash && buffer.Len() > 0) {
			buffer.WriteString(path[i:j])
		}
		i = j
	}

	return buffer.String()
}
