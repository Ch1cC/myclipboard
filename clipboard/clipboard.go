package clipboard

type Clipboard struct {
	UnixMicro int64  `json:"unixMicro"`
	Msg       []byte `json:"msg"`
}
