package lru

import "sync"

type Node[V any] struct {
	// TODO: попробовать заменить пустой интерфейс на тип V
	Value V
	Next  *Node[V]
	Prev  *Node[V]
}

func newNode[V any](value V) *Node[V] {
	return &Node[V]{Value: value}
}

type LRUCache[K comparable, V any] struct {
	mu               *sync.Mutex
	length, capacity int
	head, tail       *Node[V]
	lookup           map[K]*Node[V]

	// нужна толькоко в одном месте что бы узнать ключ по ноде
	// но может лучше хранить ключи в ноде? или запоминать ключ хвоста?
	reverseLookup map[*Node[V]]K
}

func newLRUCache[K comparable, V any](capacity int) LRUCache[K, V] {
	lru := LRUCache[K, V]{
		mu:            &sync.Mutex{},
		length:        0,
		capacity:      capacity,
		lookup:        make(map[K]*Node[V]),
		reverseLookup: make(map[*Node[V]]K),
	}

	return lru
}

func (l *LRUCache[K, V]) Reset(capacity int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.length = 0
	l.capacity = capacity
	l.head = nil
	l.tail = nil
	l.lookup = make(map[K]*Node[V])
	l.reverseLookup = make(map[*Node[V]]K)
}

func (l *LRUCache[K, V]) Update(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.update(key, value)
}

func (l *LRUCache[K, V]) Get(key K) (value V, exists bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.get(key)
}

// TODO Fetch
func (l *LRUCache[K, V]) Fetch(key K, f func(K) (V, bool)) (V, bool) {
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

func (l *LRUCache[K, V]) update(key K, value V) {
	if node, ok := l.lookup[key]; ok {
		l.detach(node)
		l.prepend(node)
		node.Value = value
	} else {
		node = newNode(value)
		l.prepend(node)
		l.length++
		l.trimCache()
		l.lookup[key] = node
		l.reverseLookup[node] = key
	}
}

func (l *LRUCache[K, V]) get(key K) (V, bool) {
	node, ok := l.lookup[key]
	if !ok {
		var defalultValue V
		return defalultValue, true
	}

	l.detach(node)
	l.prepend(node)
	return node.Value, true
}

func (l *LRUCache[K, V]) trimCache() {
	if l.length <= l.capacity {
		return
	}
	tail := l.tail
	l.detach(tail)
	key := l.reverseLookup[tail]
	delete(l.lookup, key)
	delete(l.reverseLookup, tail)
	l.length--
}

func (l *LRUCache[K, V]) detach(node *Node[V]) {
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

func (l *LRUCache[K, V]) prepend(node *Node[V]) {
	if l.head == nil {
		l.tail = node
		l.head = node
		return
	}
	node.Next = l.head
	l.head.Prev = node
	l.head = node
}
