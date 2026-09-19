package revision

import (
	"container/heap"
	"mygit/internal/object"
)

type CommitPQ []*object.Commit

var _ heap.Interface = (*CommitPQ)(nil)

func (h CommitPQ) Len() int {
	return len(h)
}

func (h CommitPQ) Less(i, j int) bool {
	return h[i].Date < h[j].Date
}

func (h CommitPQ) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *CommitPQ) Push(x any) {
	commit := x.(*object.Commit)

	*h = append(*h, commit)
}

func (h *CommitPQ) Pop() any {
	old := *h

	n := len(old)

	commit := old[n-1]

	*h = old[:n-1]

	return commit
}
