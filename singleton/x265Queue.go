package singleton

import "openOEP/types"

//X265Queue x256 data queue
var X265Queue = make(chan *types.X265Nal, 256)
