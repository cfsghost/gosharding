package gosharding

import (
	"sync/atomic"

	jump "github.com/lithammer/go-jump-consistent-hash"
)

type Shard struct {
	options   *Options
	pipelines []*Pipeline
	counter   int32
}

func NewShard(opts *Options) *Shard {

	// Initialize piplines
	pipelines := make([]*Pipeline, 0, opts.PipelineCount)
	for i := int32(0); i < opts.PipelineCount; i++ {

		pipeline := &Pipeline{
			id:             i,
			bufferSize:     opts.BufferSize,
			prepareHandler: opts.PrepareHandler,
			handler:        opts.Handler,
		}

		pipeline.initialize()

		pipelines = append(pipelines, pipeline)
	}

	return &Shard{
		options:   opts,
		pipelines: pipelines,
		counter:   0,
	}
}

func (shard *Shard) dispatch(data interface{}) {

	// Push data to pipeline
	shard.pipelines[shard.counter].push(data)
	// Update counter
	counter := atomic.AddInt32((*int32)(&shard.counter), 1)
	if counter == shard.options.PipelineCount {
		atomic.StoreInt32((*int32)(&shard.counter), 0)
	}
}

// ComputePipelineByString takes string as key to genereate a pipeline ID.
func (shard *Shard) computePipelineByString(key string) int32 {
	if len(key) == 0 {
		return -1
	}

	return jump.HashString(key, shard.options.PipelineCount, jump.NewCRC64())
}

// Push will put data to the pipelines by 64 bit key.
func (shard *Shard) Push(key uint64, data interface{}) {

	// Figure out pipeline for this data
	pipelineID := jump.Hash(key, shard.options.PipelineCount)

	// Push data to pipeline
	shard.pipelines[pipelineID].push(data)
}

// PushKV will put data to the pipelines by string key
func (shard *Shard) PushKV(key string, data interface{}) {

	pipelineID := shard.computePipelineByString(key)

	if pipelineID == -1 {
		shard.dispatch(data)
		return
	}

	// Push data to pipeline
	shard.pipelines[pipelineID].push(data)
}
