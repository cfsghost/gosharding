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

	if pipeline.prepareHandler != nil {
		pipeline.inputClosed = make(chan struct{})
		pipeline.input = make(chan interface{}, pipeline.bufferSize)
		go pipeline.inputReceivier()
	}

	pipeline.outputClosed = make(chan struct{})
	pipeline.output = make(chan interface{}, pipeline.bufferSize)
	go pipeline.outputHandler()
}

func (pipeline *Pipeline) close() {
	if pipeline.prepareHandler != nil {
		pipeline.inputClosed <- struct{}{}
	}

	pipeline.outputClosed <- struct{}{}
}

func (pipeline *Pipeline) push(data interface{}) {

	if pipeline.prepareHandler != nil {
		pipeline.input <- data
		return
	}

	pipeline.output <- data
}

func (pipeline *Pipeline) inputReceivier() {

	defer close(pipeline.inputClosed)

	for {
		select {
		case data := <-pipeline.input:
			pipeline.prepareHandler(pipeline.id, data, pipeline.output)
		case <-pipeline.inputClosed:
			return
		}
	}
}

func (pipeline *Pipeline) outputHandler() {

	defer close(pipeline.outputClosed)

	for {
		select {
		case data := <-pipeline.output:
			pipeline.handler(pipeline.id, data)
		case <-pipeline.outputClosed:
			return
		}
	}
}
