package orderedset

import (
	"sort"
	"fmt"
	"strconv"
)

/*
 * todo: optimize
 * 1. Score and key can put in one struct, do append once in add or remove.
 * 2. In "Remove", first find scores by sort.Search, then find key in items with these scores.
 */
type OrderedSet struct {
	scores []int
	keys   []interface{}
}

//func New() OrderedSet {
//	return OrderedSet{
//		scores: make([]int, 0),
//		keys:   make([]interface{}, 0),
//	}
//}

func (o *OrderedSet) Init() {
	o.scores = make([]int, 0)
	o.keys = make([]interface{}, 0)
}

func (o OrderedSet) Len() int {
	return len(o.scores)
}

func (o *OrderedSet) String() string {
	s := ""
	for i, v := range o.scores {
		s += fmt.Sprintf(" %+v:", i) + strconv.Itoa(v)
	}
	return s
}

func (o *OrderedSet) Add(k interface{}, score int) {
	i := sort.SearchInts(o.scores, score)
	rearScores := append([]int{}, o.scores[i:]...)
	o.scores = append(append(o.scores[:i], score), rearScores...)

	rearKeys := append([]interface{}{}, o.keys[i:]...)
	o.keys = append(append(o.keys[:i], k), rearKeys...)
}

func (o *OrderedSet) Remove(k interface{}) {
	for i, v := range o.keys {
		if v == k {
			o.scores = append(o.scores[:i], o.scores[i+1:]...)
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
		}
	}
}

func (o *OrderedSet) Update(k interface{}, score int) {
	o.Remove(k)
	o.Add(k, score)
}

func (o *OrderedSet) Front() interface{} {
	return o.keys[0]
}

func (o *OrderedSet) End() interface{} {
	return o.keys[len(o.keys)-1]
}
