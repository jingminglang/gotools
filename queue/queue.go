package queue

import (
	"fmt"
	"sync"
)

var ROBOT_QUEUE = NewQueue(10)
var AGENT_QUEUE = NewQueue(10)

var mutex sync.Mutex


type Node struct {
	Value Message
}

func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]*Node, size),
		size:  size,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	lock sync.Mutex
	nodes []*Node
	size  int
	head  int
	tail  int
	count int
}

// Push adds a node to the queue.
func (q *Queue) Push(n *Node) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Node, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() *Node {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.count <= 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}


func (q *Queue) Count() int  {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.count
}
