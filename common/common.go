package common

import U "github.com/daifiyum/cat-box/common/utils"

var (
	Port             string          = "3000"
	IsCoreRunning    *U.BoolState    = &U.BoolState{}
	IsTun            *U.BoolState    = &U.BoolState{}
	Box              *U.Box          = &U.Box{}
	PrevCrc32        uint32          = 0
	Broadcaster      *U.BroadcastHub = U.NewBroadcaster()
	DefaultUserAgent string          = "sing-box"
)
