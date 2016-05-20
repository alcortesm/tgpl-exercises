package main

import (
	"sort"
	"sync"
)

// FreqCountWithTags keeps track of the number of times a unique
// identifier is added to it, along with a list of the tags provided
// while adding.
type FreqCountWithTags struct {
	ids   map[string]countAndTags
	mutex sync.Mutex
}

type countAndTags struct {
	count int
	tags  map[string]bool
}

func NewFreqCountWithTags() *FreqCountWithTags {
	return &FreqCountWithTags{
		ids: make(map[string]countAndTags),
	}
}

func (f *FreqCountWithTags) Add(id string, tag string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isNewID(id) {
		f.addID(id, tag)
	} else {
		f.incrementID(id, tag)
	}
}

func (f *FreqCountWithTags) isNewID(id string) bool {
	_, ok := f.ids[id]

	return !ok
}

func (f *FreqCountWithTags) addID(id string, tag string) {
	tags := make(map[string]bool)
	tags[tag] = true

	f.ids[id] = countAndTags{
		count: 1,
		tags:  tags,
	}
}

func (f *FreqCountWithTags) incrementID(id string, tag string) {
	countAndTags, _ := f.ids[id]
	countAndTags.tags[tag] = true
	countAndTags.count++
	f.ids[id] = countAndTags
}

func (f *FreqCountWithTags) Get(id string) (count int, tags []string) {
	f.mutex.Lock()
	count, tagsMap := f.get(id)
	f.mutex.Unlock()

	asSlice := mapToSlice(tagsMap)
	sort.Strings(asSlice)

	return count, asSlice
}

func (f *FreqCountWithTags) get(id string) (count int, tags map[string]bool) {
	countAndTags, ok := f.ids[id]
	if !ok {
		return 0, nil
	}

	return countAndTags.count, countAndTags.tags
}

func mapToSlice(m map[string]bool) []string {
	s := make([]string, len(m))
	for k := range m {
		s = append(s, k)
	}

	return s
}

func (f *FreqCountWithTags) GetAll() (ids []string, counts []int, tagss [][]string) {
	ids = make([]string, 0, len(f.ids))
	counts = make([]int, 0, len(f.ids))
	tagss = make([][]string, 0, len(f.ids))

	f.mutex.Lock()
	for id := range f.ids {
		ids = append(ids, id)

		count, tags := f.get(id)
		counts = append(counts, count)

		asSlice := mapToSlice(tags)
		sort.Strings(asSlice)
		tagss = append(tagss, asSlice)
	}
	f.mutex.Unlock()

	return ids, counts, tagss
}
