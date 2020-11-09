/*
Package gosharding provides the ability to craete shard for process data in parallel.
Here is example to create a shard:

	// Create Options object
	options := &gosharding.Options{
		PipelineCount: 32,
		BufferSize:    8192,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			// Customized handler to process data before output
			c <- data
		},
		Handler: func(int32, interface{}) {
			// Customized handler to consume data
		},
	}

	// Create shard with options
	shard := gosharding.NewShard(options)

	// Push data into shard
	shard.Push(1, "test")
*/

package gosharding
