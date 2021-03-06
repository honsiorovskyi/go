// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Lock-free stack.
// The following code runs only on g0 stack.

package runtime

import "unsafe"

func lfstackpush(head *uint64, node *lfnode) {
	node.pushcnt++
	new := lfstackPack(node, node.pushcnt)
	if node1, _ := lfstackUnpack(new); node1 != node {
		println("runtime: lfstackpush invalid packing: node=", node, " cnt=", hex(node.pushcnt), " packed=", hex(new), " -> node=", node1, "\n")
		throw("lfstackpush")
	}
	for {
		old := atomicload64(head)
		node.next = old
		if cas64(head, old, new) {
			break
		}
	}
}

func lfstackpop(head *uint64) unsafe.Pointer {
	for {
		old := atomicload64(head)
		if old == 0 {
			return nil
		}
		node, _ := lfstackUnpack(old)
		next := atomicload64(&node.next)
		if cas64(head, old, next) {
			return unsafe.Pointer(node)
		}
	}
}
