package state

import (
	"sync"
	"ws-trial/config/types"
)

var WorkerMappings map[int]chan types.Payload_t = map[int]chan types.Payload_t{}
var Jobs types.Jobs_t = types.Jobs_t{List: make([]types.Payload_t, 0), Mut: sync.RWMutex{}}
