package gskiplist

import (
    "errors"
    "math/rand"
    "strings"
)

const (
    maxLevel = 32
    p        = 0.25
)

type Element interface {
    String() string
}

type SkipList struct {
    // 上层使用时候加，这里先不加了
    // mu           sync.RWMutex
    header, tail *Node
    length       int64
    level        int
}

type Nil struct {
}

func (n *Nil) String() string {
    return ""
}

type Level struct {
    forward *Node
    span    int64
}

type Node struct {
    score    float64
    element  Element
    level    []Level
    backward *Node
}

func (n *Node) IsNil() bool {
    if _, ok := n.element.(*Nil); ok {
        return true
    }
    return false
}

func createNode(level int, score float64, ele Element) *Node {
    if ele == nil {
        ele = &Nil{}
    }
    return &Node{
        score:   score,
        element: ele,
        level:   make([]Level, level),
    }
}

func NewSkipList() *SkipList {
    return &SkipList{
        header: createNode(maxLevel, 0, nil),
        tail:   nil,
        length: 0,
        level:  1,
    }
}

func randomLevel() int {
    var level = 1
    for level < maxLevel && float64(rand.Int()&0xFFFF) < float64(0xFFFF)*p {
        level += 1
    }
    return level
}

func (sl *SkipList) Insert(score float64, element Element) (*Node, error) {
    if sl == nil {
        return nil, errors.New("skip list is nil")
    }
    if element.String() == "" {
        return nil, errors.New("element is empty")
    }
    //sl.mu.Lock()
    //defer sl.mu.Unlock()
    update := make([]*Node, maxLevel)
    rank := make([]int64, maxLevel)
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        if sl.level-1 == i {
            rank[i] = 0
        } else {
            rank[i] = rank[i+1]
        }
        for x.level[i].forward != nil &&
            (x.level[i].forward.score < score ||
                (x.level[i].forward.score == score &&
                    strings.Compare(x.level[i].forward.element.String(), element.String()) < 0)) {
            rank[i] += x.level[i].span
            x = x.level[i].forward
        }
        update[i] = x
    }
    level := randomLevel()
    if level > sl.level {
        for i := sl.level; i < level; i++ {
            rank[i] = 0
            update[i] = sl.header
            update[i].level[i].span = sl.length
        }
        sl.level = level
    }
    x = createNode(level, score, element)
    for i := 0; i < level; i++ {
        x.level[i].forward = update[i].level[i].forward
        update[i].level[i].forward = x

        x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
        update[i].level[i].span = (rank[0] - rank[i]) + 1
    }

    for i := level; i < sl.level; i++ {
        update[i].level[i].span++
    }
    if update[0] != sl.header {
        x.backward = update[0]
    }
    if x.level[0].forward != nil {
        x.level[0].forward.backward = x
    } else {
        sl.tail = x
    }
    sl.length++
    return x, nil
}

func (sl *SkipList) deleteNode(node *Node, update []*Node) {
    for i := 0; i < sl.level; i++ {
        if update[i].level[i].forward == node {
            update[i].level[i].span += node.level[i].span - 1
            update[i].level[i].forward = node.level[i].forward
        } else {
            update[i].level[i].span -= 1
        }
    }
    if node.level[0].forward != nil {
        node.level[0].forward.backward = node.backward
    } else {
        sl.tail = node.backward
    }
    for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
        sl.level--
    }
    sl.length--
}

func (sl *SkipList) Delete(score float64, element Element) (int64, error) {
    if sl == nil {
        return 0, errors.New("skip list is nil")
    }
    if element.String() == "" {
        return 0, errors.New("element is empty")
    }
    update := make([]*Node, maxLevel)
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            (x.level[i].forward.score < score ||
                (x.level[i].forward.score == score &&
                    strings.Compare(x.level[i].forward.element.String(), element.String()) < 0)) {
            x = x.level[i].forward
        }
        update[i] = x
    }
    x = x.level[0].forward
    if x != nil && score == x.score &&
        strings.Compare(element.String(), x.element.String()) == 0 {
        sl.deleteNode(x, update)
        return 1, nil
    }
    return 0, nil
}

func (sl *SkipList) UpdateScore(score float64, element Element, newscore float64) (*Node, error) {
    update := make([]*Node, maxLevel)
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            (x.level[i].forward.score < score ||
                (x.level[i].forward.score == score &&
                    strings.Compare(x.level[i].forward.element.String(), element.String()) < 0)) {
            x = x.level[i].forward
        }
        update[i] = x
    }
    x = x.level[0].forward
    if x == nil || x.score != score || strings.Compare(x.element.String(), element.String()) != 0 {
        return nil, errors.New("element not found")
    }

    if (x.backward == nil || x.backward.score < newscore) &&
        (x.level[0].forward == nil || x.level[0].forward.score > newscore) {
        x.score = newscore
        x.element = element
        return x, nil
    }
    sl.deleteNode(x, update)
    newNode, err := sl.Insert(newscore, x.element)
    x.element = nil
    return newNode, err
}

func (sl *SkipList) IsInRank(index int64) bool {
    if index >= 0 && index < int64(sl.length) {
        return true
    }
    return false
}

func (sl *SkipList) IsInRange(spec *RangeSpec) bool {
    if spec.Min > spec.Max ||
        (spec.Min == spec.Max && spec.OpenRange()) {
        return false
    }
    x := sl.tail
    if x == nil || !spec.ValueGteMin(x.score) {
        return false
    }
    x = sl.header.level[0].forward
    if x == nil || !spec.ValueLteMax(x.score) {
        return false
    }
    return true
}

func (sl *SkipList) FirstInRange(spec *RangeSpec) (*Node, int64) {
    if !sl.IsInRange(spec) {
        return nil, -1
    }

    var traversed int64
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            !spec.ValueGteMin(x.level[i].forward.score) {
            traversed += x.level[i].span
            x = x.level[i].forward
        }
    }

    traversed++
    x = x.level[0].forward
    if x == nil {
        // 应该不为空
        return nil, -1
    }

    if !spec.ValueLteMax(x.score) {
        return nil, -1
    }
    return x, traversed
}

func (sl *SkipList) LastInRange(spec *RangeSpec) (*Node, int64) {
    if !sl.IsInRange(spec) {
        return nil, -1
    }
    var traversed int64
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            !spec.ValueLteMax(x.level[i].forward.score) {
            traversed += x.level[i].span
            x = x.level[i].forward
        }
    }
    if x == nil {
        // 应该不为空
        return nil, -1
    }
    if !spec.ValueGteMin(x.score) {
        return nil, -1
    }
    return x, traversed
}

func (sl *SkipList) GetByRank(index int64) *Node {
    if !sl.IsInRank(index) {
        return nil
    }
    var traversed int64
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            traversed+x.level[i].span <= index {
            traversed += x.level[i].span
            x = x.level[i].forward
        }
    }
    traversed++
    x = x.level[0].forward
    return x
}

func (sl *SkipList) DeleteRangeByScore(spce *RangeSpec, dict map[string]Element) (int, error) {
    removed := 0
    update := make([]*Node, maxLevel)
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            spce.ValueGteMin(x.level[i].forward.score) {
            x = x.level[i].forward
        }
        update[i] = x
    }
    x = x.level[0].forward
    for x != nil && spce.ValueLteMax(x.score) {
        next := x.level[0].forward
        sl.deleteNode(x, update)
        delete(dict, x.element.String())
        removed++
        x = next
    }
    return removed, nil
}

func (sl *SkipList) DeleteRangeByRank(start, end int64, dict map[string]Element) (int, error) {
    var traversed int64
    var removed int
    update := make([]*Node, maxLevel)
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            traversed+x.level[i].span < start {
            traversed += x.level[i].span
            x = x.level[i].forward
        }
        update[i] = x
    }

    traversed++
    x = x.level[0].forward

    for x != nil && traversed <= end {
        next := x.level[0].forward
        sl.deleteNode(x, update)
        delete(dict, x.element.String())
        removed++
        traversed++
        x = next
    }
    return removed, nil
}

func (sl *SkipList) GetRank(score float64, element Element) int64 {
    var rank int64
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil &&
            (x.level[i].forward.score < score ||
                (x.level[i].forward.score == score &&
                    strings.Compare(x.level[i].forward.element.String(), element.String()) <= 0)) {
            rank += x.level[i].span
            x = x.level[i].forward
        }
        if !x.IsNil() && x.score == score && strings.Compare(x.element.String(), element.String()) == 0 {
            return rank
        }
    }
    return 0
}

func (sl *SkipList) GetElementByRank(rank int64) Element {
    var traversed int64
    x := sl.header
    for i := sl.level - 1; i >= 0; i-- {
        for x.level[i].forward != nil && (traversed+x.level[i].span) <= rank {
            traversed += x.level[i].span
            x = x.level[i].forward
        }
        if traversed == rank {
            return x.element
        }
    }
    return nil
}
