package gosharding

type Pipeline struct {
	id             int32
	bufferSize     int
	input          chan interface{}
	output         chan interface{}
	prepareHandler func(int32, interface{}, chan interface{})
	handler        func(int32, interface{})
}

func (pipeline *Pipeline) initialize() {

	pipeline.output = make(chan interface{}, pipeline.bufferSize)

	if pipeline.prepareHandler != nil {
		pipeline.input = make(chan interface{}, pipeline.bufferSize)
		go pipeline.dispatcher()
	} else {
		go pipeline.outputHandler()
	}
}

func (pipeline *Pipeline) close() {
	if pipeline.prepareHandler != nil {
		close(pipeline.input)
	}

	close(pipeline.output)
}

func (pipeline *Pipeline) push(data interface{}) {

	if pipeline.prepareHandler != nil {
		pipeline.input <- data
		return
	}

	pipeline.output <- data
}

func (pipeline *Pipeline) dispatcher() {

	for {
		select {
		case data, ok := <-pipeline.input:

			if !ok {
				// closed
				return
			}

			pipeline.prepareHandler(pipeline.id, data, pipeline.output)
		case data, ok := <-pipeline.output:

			if !ok {
				// closed
				return
			}

			pipeline.handler(pipeline.id, data)
		}
	}
}

func (pipeline *Pipeline) outputHandler() {

	for data := range pipeline.output {
		pipeline.handler(pipeline.id, data)
	}
}
