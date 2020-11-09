package gosharding

import (
	"encoding/json"
	"testing"
)

func init() {

	// Prepare json
	data := Payload{
		String:   "string",
		Number:   99999,
		Elements: make([]string, 0, 1000),
	}

	for i := 0; i < 1000; i++ {
		data.Elements = append(data.Elements, "sample")
	}

	payload, _ = json.Marshal(&data)
}

func BenchmarkPH_Pipeline32(b *testing.B) {

	output := make(chan interface{}, 1024000)

	// Create Options object
	options := &Options{
		PipelineCount: 32,
		BufferSize:    1024000,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			c <- data
		},
		Handler: func(id int32, data interface{}) {
			output <- data
		},
	}

	// Create shard with options
	shard := NewShard(options)

	b.RunParallel(func(b *testing.PB) {
		var count uint64 = 0
		for b.Next() {
			shard.Push(count, payload)
			count++
		}
	})

	b.RunParallel(func(b *testing.PB) {

		for b.Next() {
			<-output
		}
	})
}

func BenchmarkPH_Nil_Pipeline32(b *testing.B) {

	output := make(chan interface{}, 1024000)

	// Create Options object
	options := &Options{
		PipelineCount:  32,
		BufferSize:     1024000,
		PrepareHandler: nil,
		Handler: func(id int32, data interface{}) {
			output <- data
		},
	}

	// Create shard with options
	shard := NewShard(options)

	b.RunParallel(func(b *testing.PB) {
		var count uint64 = 0
		for b.Next() {
			shard.Push(count, payload)
			count++
		}
	})

	b.RunParallel(func(b *testing.PB) {

		for b.Next() {
			<-output
		}
	})
}
