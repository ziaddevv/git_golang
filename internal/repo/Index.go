package repo

import (
	"fmt"
	"mygit/internal/index"
	"mygit/internal/utils"
	"strings"
)

func UpdateIndexCacheInfo(cacheInfo string) error {
	info := strings.Split(cacheInfo, ",")

	if len(info) != 3 {
		return fmt.Errorf("invalid cacheinfo format")
	}
	mode, hash, path := info[0], info[1], info[2]

	entry, err := index.NewEntryFromCacheInfo(mode, hash, path)
	if err != nil {
		return err
	}

	indexContent, err := utils.ReadFile(".mygit/INDEX")
	var indexEntries *index.Index
	if err != nil {
		// first time — index doesn't exist yet
		indexEntries = &index.Index{}
	} else {
		indexEntries = index.UnpackIndex(indexContent)
	}
	if err != nil {
		return err
	}

	indexEntries.Add(entry)

	// save to disk

	return Save(indexEntries)
}

func UpdateIndexAdd(path string) error {

	// read the file fro working directory
	// hash the content
	// make sure object exists
	// read index file --> bytes
	// check the index for path --> if exists change with the new hash
	// if doesn't exists adds a new entry

	file, err := utils.ReadFile(path)

	if err != nil {
		return err
	}
	hash, err := WriteObject("blob", file)
	if err != nil {
		return err
	}

	indexEntries := ReadIndex()
	// indexEntries := index.UnpackIndex(indexContent)

	// the entry

	entry, err := index.NewEntryFromFile(path, hash)
	// fmt.Println(entry)
	if err != nil {
		return err
	}

	indexEntries.Add(entry)

	// save to disk

	return Save(indexEntries)
}

func ReadIndex() *index.Index {
	indexContent, err := utils.ReadFile(".mygit/INDEX")
	var indexEntries *index.Index
	if err != nil {
		// first time — index doesn't exist yet
		indexEntries = &index.Index{}
	} else {
		indexEntries = index.UnpackIndex(indexContent)
	}
	return indexEntries
}
func UpdateIndexRemove(path string) error {
	indexContent, err := utils.ReadFile(".mygit/INDEX")
	if err != nil {
		return fmt.Errorf("index not found, nothing to remove")
	}

	idx := index.UnpackIndex(indexContent)
	idx.Remove(path)

	return Save(idx)
}

func Save(idxEntries *index.Index) error {

	var path = ".mygit/INDEX"

	content := idxEntries.Pack()
	return utils.WriteFile(path, content)
}
