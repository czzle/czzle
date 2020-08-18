package lrucache

import (
	"github.com/czzle/czzle/internal/pkg/cache"
	"sync"
)

type node struct {
	key   string
	value []byte
	prev  *node
	next  *node
}

type lrucache struct {
	m    map[string]*node
	cap  int
	size int
	head *node
	tail *node
	sync.Mutex
}

func (c *lrucache) Clear() {
	c.Lock()
	defer c.Unlock()
	for c.size > 0 {
		f := c.head.next
		c.size -= len(f.value)
		c.unlink(f)
		delete(c.m, f.key)
	}
}

func (c *lrucache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()
	if n, ok := c.m[key]; ok {
		c.unlink(n)
		c.link(n)
		return n.value, true
	}
	return nil, false
}

func (c *lrucache) Set(key string, value []byte) {
	c.Lock()
	defer c.Unlock()
	n, ok := c.m[key]
	if ok {
		c.size -= len(n.value)
		c.unlink(n)
		n.value = value
	} else {
		n = &node{
			key:   key,
			value: value,
		}
	}
	c.size += len(value)
	c.link(n)
	c.m[key] = n
	for c.size > c.cap {
		f := c.head.next
		c.size -= len(f.value)
		c.unlink(f)
		delete(c.m, f.key)
	}
}

func (c *lrucache) Remove(key string) {
	c.Lock()
	defer c.Unlock()
	n, ok := c.m[key]
	if !ok {
		return
	}
	c.unlink(n)
	delete(c.m, key)
}

func (c *lrucache) unlink(n *node) {
	n.prev.next, n.next.prev = n.next, n.prev
}

func (c *lrucache) link(n *node) {
	c.tail.prev.next,
		c.tail.prev,
		n.prev,
		n.next =
		n,
		n,
		c.tail.prev,
		c.tail
}

const (
	B  int = 1
	KB int = 1000 * B
	MB int = 1000 * KB
	GB int = 1000 * MB
)

func New(cap int) cache.Cache {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	return &lrucache{
		m:    make(map[string]*node),
		cap:  cap,
		head: head,
		tail: tail,
		size: 0,
	}
}
