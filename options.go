package gosharding

type Options struct {
	PipelineCount  int32
	BufferSize     int
	PrepareHandler func(int32, interface{}, chan interface{})
	Handler        func(int32, interface{})
}

func NewOptions() *Options {
	return &Options{
		PipelineCount: 32,
		BufferSize:    8192,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			c <- data
		},
		Handler: func(int32, interface{}) {
		},
	}
}
