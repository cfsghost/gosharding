package gosharding

type Pipeline struct {
	id             int32
	bufferSize     int
	input          chan interface{}
	inputClosed    chan struct{}
	output         chan interface{}
	outputClosed   chan struct{}
	prepareHandler func(int32, interface{}, chan interface{})
	handler        func(int32, interface{})
}

func (pipeline *Pipeline) initialize() {

	pipeline.input = make(chan interface{}, pipeline.bufferSize)
	pipeline.output = make(chan interface{}, pipeline.bufferSize)

	go pipeline.inputReceivier()
	go pipeline.outputHandler()
}

func (pipeline *Pipeline) close() {
	pipeline.inputClosed <- struct{}{}
	pipeline.outputClosed <- struct{}{}
}

func (pipeline *Pipeline) inputReceivier() {

	go func() {

		for {
			select {
			case data := <-pipeline.input:
				pipeline.prepareHandler(pipeline.id, data, pipeline.output)
			case <-pipeline.inputClosed:
				return
			}
		}
	}()
}

func (pipeline *Pipeline) outputHandler() {

	go func() {

		for {
			select {
			case data := <-pipeline.output:
				pipeline.handler(pipeline.id, data)
			case <-pipeline.outputClosed:
				return
			}
		}
	}()
}
