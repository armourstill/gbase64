package main

import (
	"encoding/base64"
	"io"
)

// selectEncoding returns a base64 encoding according to url and padding flags.
func selectEncoding(url bool, noPadding bool) *base64.Encoding {
	var enc *base64.Encoding
	if url {
		enc = base64.URLEncoding
	} else {
		enc = base64.StdEncoding
	}
	if noPadding {
		enc = enc.WithPadding(base64.NoPadding)
	}
	return enc
}

// lineWriter wraps an io.Writer and inserts a newline after every 'limit' bytes.
// It also tracks total bytes written, so callers can decide whether to emit a final newline.
type lineWriter struct {
	w     io.Writer
	col   int
	limit int
	total int64
}

func newLineWriter(w io.Writer, limit int) *lineWriter {
	return &lineWriter{w: w, limit: limit}
}

func (lw *lineWriter) Write(p []byte) (int, error) {
	if lw.limit <= 0 {
		n, err := lw.w.Write(p)
		lw.total += int64(n)
		return n, err
	}
	written := 0
	for len(p) > 0 {
		if lw.col == lw.limit {
			if _, err := lw.w.Write([]byte{'\n'}); err != nil {
				return written, err
			}
			lw.col = 0
		}
		k, err := lw.w.Write(p[:min(len(p), lw.limit-lw.col)])
		if err != nil {
			return written, err
		}
		lw.col += k
		lw.total += int64(k)
		written += k
		p = p[k:]
	}
	return written, nil
}

// writeFinalNewlineIfNeeded emits a trailing newline if any byte has been written.
// For wrapping mode, it ensures the output ends with a newline even if the last line is exactly full.
func (lw *lineWriter) writeFinalNewlineIfNeeded() error {
	if lw.total == 0 {
		return nil
	}
	// Always ensure a trailing newline for non-empty output
	_, err := lw.w.Write([]byte{'\n'})
	return err
}

// filterReader filters characters on the fly.
// - If dropCRLF is true, it removes '\r' and '\n'.
// - If dropPadding is true, it removes '=' used as base64 padding.
// This lets us decode both padded and unpadded inputs using a NoPadding encoding.
type filterReader struct {
	r           io.Reader
	dropCRLF    bool
	dropPadding bool
}

func (fr *filterReader) Read(p []byte) (int, error) {
	n, err := fr.r.Read(p)
	if n == 0 {
		return n, err
	}
	j := 0
	for i := range n {
		b := p[i]
		if fr.dropCRLF && (b == '\r' || b == '\n') {
			continue
		}
		if fr.dropPadding && b == '=' {
			continue
		}
		p[j] = b
		j++
	}
	return j, err
}

// countingWriter wraps a writer and counts the number of bytes written through it.
type countingWriter struct {
	w io.Writer
	n int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
	k, err := cw.w.Write(p)
	cw.n += int64(k)
	return k, err
}

// EncodeStream encodes from r to w using streaming with optional line wrapping.
func EncodeStream(url, noPadding bool, wrap int, r io.Reader, w io.Writer) error {
	enc := selectEncoding(url, noPadding)
	lw := newLineWriter(w, wrap)
	ew := base64.NewEncoder(enc, lw)
	if _, err := io.Copy(ew, r); err != nil {
		_ = ew.Close()
		return err
	}
	if err := ew.Close(); err != nil {
		return err
	}
	return lw.writeFinalNewlineIfNeeded()
}

// DecodeStream decodes from r to w using streaming.
// It tolerates both padded and unpadded inputs by stripping '=' and using a no-padding encoding.
// If ignoreNewlines is true, it ignores CR and LF in the input.
func DecodeStream(url bool, ignoreNewlines bool, r io.Reader, w io.Writer) error {
	// Use no-padding encoding; strip '=' from the input stream for compatibility with padded input.
	enc := selectEncoding(url, true)
	fr := &filterReader{r: r, dropCRLF: ignoreNewlines, dropPadding: true}
	dr := base64.NewDecoder(enc, fr)
	cw := &countingWriter{w: w}
	_, err := io.Copy(cw, dr)
	return err
}
