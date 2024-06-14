package sse

import (
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

func Encode(w io.Writer, e *Event) (err error) {
	err = writeID(w, e.ID)
	if err != nil {
		return
	}
	err = writeEvent(w, e.Event)
	if err != nil {
		return
	}
	err = writeRetry(w, e.Retry)
	if err != nil {
		return
	}
	err = writeData(w, e.Data)
	if err != nil {
		return
	}
	return nil
}

func writeID(w io.Writer, id string) (err error) {
	if len(id) > 0 {
		_, err = w.Write([]byte("id:"))
		if err != nil {
			return
		}
		_, err = fieldReplacer.WriteString(w, id)
		if err != nil {
			return
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			return
		}
	}

	return
}

func writeEvent(w io.Writer, event string) (err error) {
	if len(event) > 0 {
		_, err = w.Write([]byte("event:"))
		if err != nil {
			return
		}
		_, err = fieldReplacer.WriteString(w, event)
		if err != nil {
			return
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return
		}
	}

	return
}

func writeRetry(w io.Writer, retry uint64) (err error) {
	if retry > 0 {
		_, err = w.Write([]byte("retry:"))
		if err != nil {
			return
		}
		_, err = w.Write(s2b(strconv.FormatUint(retry, 10)))
		if err != nil {
			return
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			return
		}
	}

	return
}

func writeData(w io.Writer, data []byte) (err error) {
	_, err = w.Write([]byte("data: "))
	if err != nil {
		return
	}

	_, err = dataReplacer.WriteString(w, b2s(data))
	if err != nil {
		return
	}

	_, err = w.Write([]byte("\n\n"))
	if err != nil {
		return
	}

	return nil
}

func b2s(in []byte) (out string) {
	op := (*reflect.StringHeader)(unsafe.Pointer(&out))
	ip := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	op.Data = ip.Data
	op.Len = ip.Len
	return
}

func s2b(in string) (out []byte) {
	op := (*reflect.SliceHeader)(unsafe.Pointer(&out))
	ip := (*reflect.StringHeader)(unsafe.Pointer(&in))
	op.Data = ip.Data
	op.Cap = ip.Len
	op.Len = ip.Len
	return
}
