package lru

import "sync"

type NodeT[K comparable, V any] struct {
	Value V
	Key   K
	Next  *NodeT[K, V]
	Prev  *NodeT[K, V]
}

func newNodeT[K comparable, V any](key K, value V) *NodeT[K, V] {
	return &NodeT[K, V]{Key: key, Value: value}
}

type LRUCacheT[K comparable, V any] struct {
	mu               *sync.Mutex
	length, capacity int
	head, tail       *NodeT[K, V]
	lookup           map[K]*NodeT[K, V]
}

func newLRUCacheT[K comparable, V any](capacity int) LRUCacheT[K, V] {
	lru := LRUCacheT[K, V]{
		mu:       &sync.Mutex{},
		length:   0,
		capacity: capacity,
		lookup:   make(map[K]*NodeT[K, V]),
	}

	return lru
}

func (l *LRUCacheT[K, V]) Reset(capacity int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.length = 0
	l.capacity = capacity
	l.head = nil
	l.tail = nil
	l.lookup = make(map[K]*NodeT[K, V])
}

func (l *LRUCacheT[K, V]) Update(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.update(key, value)
}

func (l *LRUCacheT[K, V]) Get(key K) (value V, exists bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.get(key)
}

// TODO Fetch
func (l *LRUCacheT[K, V]) Fetch(key K, f func(K) (V, bool)) (V, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	value, ok := l.get(key)
	if ok {
		return value, ok
	}

	if value, ok := f(key); ok {
		l.update(key, value)
		return value, true
	}

	return value, false
}

func (l *LRUCacheT[K, V]) update(key K, value V) {
	if node, ok := l.lookup[key]; ok {
		l.detach(node)
		l.prepend(node)
		node.Value = value
	} else {
		node = newNodeT(key, value)
		l.prepend(node)
		l.length++
		l.trimCache()
		l.lookup[key] = node
	}
}

func (l *LRUCacheT[K, V]) get(key K) (V, bool) {
	node, ok := l.lookup[key]
	if !ok {
		var defalultValue V
		return defalultValue, true
	}

	l.detach(node)
	l.prepend(node)
	return node.Value, true
}

func (l *LRUCacheT[K, V]) trimCache() {
	if l.length <= l.capacity {
		return
	}
	tail := l.tail
	l.detach(tail)
	key := tail.Key
	delete(l.lookup, key)
	l.length--
}

func (l *LRUCacheT[K, V]) detach(node *NodeT[K, V]) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if node == l.head {
		l.head = l.head.Next
	}
	if node == l.tail {
		l.tail = l.tail.Prev
	}
	node.Prev = nil
	node.Next = nil
}

func (l *LRUCacheT[K, V]) prepend(node *NodeT[K, V]) {
	if l.head == nil {
		l.tail = node
		l.head = node
		return
	}
	node.Next = l.head
	l.head.Prev = node
	l.head = node
}
