package http

import (
	"bytes"
	"io"
	"reflect"
	"unsafe"
)

type RepeatReadCloser struct {
	Reader *bytes.Buffer
}

// Read reset point after each read
//
//	@Description:
//	@receiver p
//	@param val
//	@return n
//	@return err
func (p *RepeatReadCloser) Read(val []byte) (n int, err error) {
	if p.Reader.Len() == 0 {
		// reset offset
		p.resetBufferOffset()

		return 0, io.EOF
	}

	n, err = p.Reader.Read(val)

	return
}

// reset offset and lastRead in buffer
func (p *RepeatReadCloser) resetBufferOffset() {
	r := reflect.ValueOf(p.Reader)
	buffer := r.Elem()

	// set buffer.off = 0
	offValue := buffer.FieldByName("off")
	offValue = reflect.NewAt(offValue.Type(), unsafe.Pointer(offValue.UnsafeAddr())).Elem()
	offValue.SetInt(0)

	// sync set buffer.lastRead = opInvalid
	lastReadValue := buffer.FieldByName("lastRead")
	lastReadValue = reflect.NewAt(lastReadValue.Type(), unsafe.Pointer(lastReadValue.UnsafeAddr())).Elem()
	lastReadValue.SetInt(0)
}

func (p *RepeatReadCloser) Close() error {

	return nil
}

// TransformReadCloser
//
//	@Description: quickly transform aio.Reader into RepeatReadCloser
//	@param r
func TransformReadCloser(r io.Reader) *RepeatReadCloser {

	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return &RepeatReadCloser{Reader: buf}
}
