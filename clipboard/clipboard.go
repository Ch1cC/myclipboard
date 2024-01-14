package clipboard

type Clipboard struct {
	UnixMicro int64  `json:"unixMicro"`
	Msg       string `json:"msg"`
}
