package porcelain

import (
	"container/heap"
	"fmt"
	"mygit/internal/object"
	"mygit/internal/repo"
	"mygit/internal/revision"
	"time"
)

/*
ok we need to follow commits to the last one with no parent(which is the first commit)

where is the last commit , it is in the HEAD file

ok so we can a coupl of funcs to do that

1 that reads the head and gets us the last commit --> resolve ref

2 returns the parent commit hash

3






implement the pq to take commit objects then it sortes based on timestamp then

when printing it used the function commitToLog by passing the commit object to it

*/

func Log() error {

	// get last commit

	_, lastCommit, err := repo.ResolveRef("HEAD")

	if err != nil {
		return err
	}

	if lastCommit == "" {
		fmt.Println("fatal: your current branch 'master' does not have any commits yet")
		return nil
	}
	commit, err := object.GetCommitbyHash(lastCommit)
	if err != nil {
		return err
	}
	traverse(commit)

	return nil
}

func traverse(lastCommit *object.Commit) {
	pq := &revision.CommitPQ{}

	heap.Init(pq)

	heap.Push(pq, lastCommit)

	for pq.Len() > 0 {
		commit := heap.Pop(pq).(*object.Commit)

		fmt.Println(ToLogCommit(commit))

		for _, parent := range commit.ParentCommits {

			parentCommit, err := object.GetCommitbyHash(parent)

			if err != nil {
				fmt.Println("issue happened while fetching commit ", parentCommit)
				return
			}
			heap.Push(pq, parentCommit)
		}
	}
}

func ToLogCommit(commit *object.Commit) string {
	commitTime := time.Unix(commit.Date, 0)

	dateStr := commitTime.Format("Mon Jan 02 15:04:05 2006 -0700")

	return fmt.Sprintf("commit %s\nAuthor: %s\nDate:   %s\n\n    %s\n\n",
		commit.Hash,
		commit.Author,
		dateStr,
		commit.Message,
	)
}
