package gosharding

import (
	"encoding/json"
	"testing"
)

type Payload struct {
	String   string   `json:"string,omitempty"`
	Number   float64  `json:"number,omitempty"`
	Elements []string `json:"elements,omitempty"`
}

var payload []byte

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

func BenchmarkWithoutPipeline(b *testing.B) {

	output := make(chan interface{}, b.N)

	b.RunParallel(func(b *testing.PB) {

		for b.Next() {
			var result map[string]interface{}
			json.Unmarshal(payload, &result)
			output <- result
		}
	})

	b.RunParallel(func(b *testing.PB) {

		for b.Next() {
			<-output
		}
	})
}

func BenchmarkPipeline32(b *testing.B) {

	output := make(chan interface{}, 10240)

	// Create Options object
	options := &Options{
		PipelineCount: 32,
		BufferSize:    102400,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			var result map[string]interface{}
			json.Unmarshal(data.([]byte), &result)
			c <- result
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

func BenchmarkPipeline128(b *testing.B) {

	output := make(chan interface{}, 10240)

	// Create Options object
	options := &Options{
		PipelineCount: 128,
		BufferSize:    10240,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			var result map[string]interface{}
			json.Unmarshal(data.([]byte), &result)
			c <- result
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

func BenchmarkPipelines512(b *testing.B) {

	output := make(chan interface{}, 10240)

	// Create Options object
	options := &Options{
		PipelineCount: 512,
		BufferSize:    10240,
		PrepareHandler: func(id int32, data interface{}, c chan interface{}) {
			var result map[string]interface{}
			json.Unmarshal(data.([]byte), &result)
			c <- result
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
