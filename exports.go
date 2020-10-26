package main

import "C"
import (
	"unsafe"

	"agent/singleton"
	"agent/types"
)

//X265Nal callback
//export onX265Nal
func onX265Nal(nalType, nalSizeBytes C.uint, nalPayload *C.uchar) {
	singleton.X265Queue <- &types.X265Nal{
		Type:    int(nalType),
		Size:    int(nalSizeBytes),
		Payload: C.GoBytes(unsafe.Pointer(nalPayload), C.int(nalSizeBytes)),
	}
}
