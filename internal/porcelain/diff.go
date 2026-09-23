package porcelain

import (
	"fmt"
	"mygit/internal/index"
	"mygit/internal/object"
	"mygit/internal/utils"
	"mygit/internal/worktree"
	"strings"
)

type FileDiff struct {
	Path    string
	OldHash string
	NewHash string
}

type OperationType int

const (
	Equal OperationType = iota
	Delete
	Insert
)

type Operation struct {
	Type    OperationType
	OldLine string
	NewLine string
}

type Myers struct {
	oldLines []string
	newLines []string
	v        map[int]int
	traces   []map[int]int
}

func NewMyers(oldLines, newLines []string) *Myers {
	return &Myers{
		oldLines: oldLines,
		newLines: newLines,
		v:        make(map[int]int),
		traces:   make([]map[int]int, 0),
	}
}

func (m *Myers) FindShortestEditScript() []Operation {

	N := len(m.oldLines)
	M := len(m.newLines)

	Max := N + M

	m.v[1] = 0

	for d := 0; d <= Max; d++ {

		for k := -d; k <= d; k += 2 {
			x := 0
			if k == -d {
				// here i'm sure we moved down that means x didn't move only y moved down
				x = m.v[k+1]
			} else if k == d {
				// same here we are sure that x moved to right (increased)
				x = m.v[k-1] + 1
			} else {
				if m.v[k-1] < m.v[k+1] {
					x = m.v[k+1]
				} else {
					x = m.v[k-1] + 1
				}
			}

			y := x - k

			// then check if we can move diagonal

			for (x < N) && y < M && m.oldLines[x] == m.newLines[y] {
				x, y = x+1, y+1
			}

			m.v[k] = x

			// fmt.Printf("d=%d k=%d x=%d y=%d\n", d, k, x, y)

			if x >= N && y >= M {

				// for _, op := range operations {
				// 	fmt.Println(op)
				m.traces = append(m.traces, copyMap(m.v))
				// }
				return m.Backtrack(d, x, y)
			}
		}
		m.traces = append(m.traces, copyMap(m.v))

	}

	return nil
}
func copyMap(src map[int]int) map[int]int {
	dst := make(map[int]int)

	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func (m *Myers) Backtrack(d, x, y int) []Operation {

	var operations []Operation

	for d > 0 {
		k := x - y

		prev := m.traces[d-1]

		var prefK int

		if k == -d || (k != d && prev[k-1] < prev[k+1]) {
			// we simply came from k + 1
			// forward movement was down , so this is an insertion
			prefK = k + 1
		} else {
			//we came from k-1
			// deletion
			prefK = k - 1
		}

		prevX := prev[prefK]
		prevY := prevX - prefK

		// now we are maybe at the end of a snake we have to move back through all matching lines
		// exactly the reverse of what we did already in FindShortestEditScript

		for x > prevX && y > prevY {

			// all those are just equal or match ones
			operations = append(operations, Operation{
				Type:    Equal,
				OldLine: m.oldLines[x-1],
				NewLine: m.newLines[y-1],
			})

			x--
			y--
		}

		if x == prevX {
			// meaning we didn't change x ---> simply we moved down --> inserting operation
			operations = append(operations, Operation{
				Type:    Insert,
				NewLine: m.newLines[y-1],
			})

			y--
		} else {
			// Forward movement was RIGHT.
			// Therefore, a deletion happened.
			operations = append(operations, Operation{
				Type:    Delete,
				OldLine: m.oldLines[x-1],
			})

			x--
		}

		d--
	}

	// remaining snake from D=0 --> as loop finished before we loop back the last snake when d = 0
	// this happends when two strings start exatly the same for a one or more  chars
	for x > 0 && y > 0 {
		operations = append(operations, Operation{
			Type:    Equal,
			OldLine: m.oldLines[x-1],
			NewLine: m.newLines[y-1],
		})

		x--
		y--
	}

	// we built the path backwards, so reverse it
	for i, j := 0, len(operations)-1; i < j; i, j = i+1, j-1 {
		operations[i], operations[j] = operations[j], operations[i]
	}

	return operations

}

/*

k-1 ---> k move right

k + 1 --> k move down


*/

func Diff() error {
	idx, err := index.ReadIndex()
	if err != nil {
		return err
	}

	indexMap := ParseIndex(idx)

	workDirEntries, err := worktree.ReadWorkTree()
	if err != nil {
		return err
	}

	workDirMap := ParseWorkTree(workDirEntries)

	unstagedChanges := Compare(indexMap, workDirMap)

	filter := func() []FileDiff {
		var result []FileDiff
		for _, v := range unstagedChanges {
			if v.Type == Modified {
				newDiff := &FileDiff{Path: v.Path, OldHash: indexMap[v.Path], NewHash: workDirMap[v.Path]}
				result = append(result, *newDiff)
			}
		}
		return result
	}

	modifiedFiles := filter()

	_ = modifiedFiles

	// fmt.Println(ModifiedFiles)

	return CompareFiles(modifiedFiles)
}

func StagedDiff() error {

	return nil
}

func CompareFiles(ModifiedFiles []FileDiff) error {

	for _, file := range ModifiedFiles {

		fileContentA, err := object.Content(file.OldHash)
		if err != nil {
			return fmt.Errorf("Fatal: Couldn't fetch files")
		}
		fileContentB, err := utils.ReadFile(file.Path)

		if err != nil {
			return fmt.Errorf("Fatal: Couldn't fetch files")
		}

		oldLines := strings.Split(strings.TrimSuffix(string(fileContentA), "\n"), "\n")
		newLines := strings.Split(strings.TrimSuffix(string(fileContentB), "\n"), "\n")
		myers := NewMyers(oldLines, newLines)

		operations := myers.FindShortestEditScript()

		diffLines := FormatOperations(operations)

		for _, line := range diffLines {
			fmt.Println(line)
		}
	}
	return nil
}

func FormatOperations(operations []Operation) []string {
	var result []string

	for _, op := range operations {
		switch op.Type {
		case Equal:
			result = append(result, " "+op.OldLine)

		case Delete:
			result = append(result, "-"+op.OldLine)

		case Insert:
			result = append(result, "+"+op.NewLine)
		}
	}

	return result
}
