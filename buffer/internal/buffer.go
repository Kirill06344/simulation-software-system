package internal

import (
	"stewie.com/buffer/internal/types"
)

type Buffer struct {
	requests     []*types.Request
	queue        []int32
	size         int32
	writePointer int32
}

func NewBuffer(capacity int32) *Buffer {
	return &Buffer{
		requests: make([]*types.Request, capacity),
	}
}

func (buf *Buffer) AddRequest(request *types.Request) int32 {
	if buf.IsFilled() {
		return -1
	}

	inx := buf.getFreePlace(buf.writePointer)
	buf.requests[inx] = request
	buf.size++
	buf.queue = push(buf.queue, inx)
	buf.writePointer = (inx + 1) % int32(len(buf.requests))
	return inx
}

func (buf *Buffer) Poll() *types.Request {
	if len(buf.queue) == 0 {
		return nil
	}

	readIndex := peek(buf.queue)
	buf.queue = poll(buf.queue)
	request := buf.requests[readIndex]
	buf.requests[readIndex] = nil
	buf.size--
	return request
}

func (buf *Buffer) Peek() types.Request {
	readIndex := peek(buf.queue)
	return *buf.requests[readIndex]
}

func (buf *Buffer) RefuseRequest(request *types.Request) *types.Request {
	if len(buf.queue) == 0 {
		return nil
	}

	readIndex := peek(buf.queue)
	buf.queue = poll(buf.queue)
	refusedRequest := buf.requests[readIndex]
	buf.requests[readIndex] = request
	buf.queue = push(buf.queue, readIndex)
	return refusedRequest
}

func (buf *Buffer) getFreePlace(startIndex int32) int32 {
	i := startIndex
	for buf.requests[i] != nil {
		i = (i + 1) % int32(len(buf.requests))
	}
	return i
}

func (buf *Buffer) IsFilled() bool {
	return buf.size == int32(len(buf.requests))
}

func (buf *Buffer) IsEmpty() bool {
	return buf.size == 0
}
