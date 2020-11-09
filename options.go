package gosharding

type Options struct {
	PipelineCount  int32
	BufferSize     int
	PrepareHandler func(int32, interface{}, chan interface{})
	Handler        func(int32, interface{})
}

func NewOptions() *Options {
	return &Options{
		PipelineCount:  32,
		BufferSize:     8192,
		PrepareHandler: nil,
		Handler: func(int32, interface{}) {
		},
	}
}
