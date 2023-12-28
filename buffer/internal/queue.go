package internal

func push(queue []int32, element int32) []int32 {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func poll(queue []int32) []int32 {
	return queue[1:] // Slice off the element once it is dequeued.
}

func peek(queue []int32) int32 {
	return queue[0]
}
