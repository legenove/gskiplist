package gskiplist

import (
    "fmt"
    "sync"
)

type SortSet struct {
    mu       sync.RWMutex
    Mapper   map[string]*ElementWithScores
    skipList *SkipList
}

type ElementWithScores struct {
    Score float64
    Ele   Element
}

func (es *ElementWithScores) String() string {
    return fmt.Sprintf("%f, %s", es.Score, es.Ele.String())
}

func NewSortSet() *SortSet {
    return &SortSet{
        Mapper:   map[string]*ElementWithScores{},
        skipList: NewSkipList(),
    }
}

func (ss *SortSet) Add(score float64, element Element) (int64, error) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    var err error
    var num int64
    if oldE, ok := ss.Mapper[element.String()]; ok {
        _, err = ss.skipList.UpdateScore(oldE.Score, element, score)
    } else {
        num += 1
        _, err = ss.skipList.Insert(score, element)
    }
    if err != nil {
        return 0, err
    }
    ss.Mapper[element.String()] = &ElementWithScores{Score: score, Ele: element}
    return num, nil
}

func (ss *SortSet) Del(element Element) (int64, error) {
    ss.mu.Lock()
    defer ss.mu.Unlock()
    if oldE, ok := ss.Mapper[element.String()]; ok {
        i, err := ss.skipList.Delete(oldE.Score, oldE.Ele)
        delete(ss.Mapper, element.String())
        return i, err
    }
    return 0, nil
}

func (ss *SortSet) Score(element Element) (float64, bool) {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    e, ok := ss.Mapper[element.String()]
    if ok {
        return e.Score, ok
    }
    return 0, false
}

func (ss *SortSet) Card() int64 {
    return ss.skipList.length
}

func (ss *SortSet) RangeByRank(start, end, count int64) []Element {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    if count <= 0 {
        return nil
    }
    if end == -1 {
        end = ss.Card()
    } else if end >= ss.Card() {
        end = ss.Card() - 1
    }
    node := ss.skipList.GetByRank(start)
    if node == nil {
        return nil
    }
    l := end - start + 1
    if count < l {
        l = count
    }
    res := make([]Element, 0, l)
    for node != nil && count > 0 && start <= end {
        res = append(res, node.element)
        count--
        start++
        node = node.level[0].forward
    }
    return res
}

func (ss *SortSet) RangeByRankWithScores(start, end, count int64) []*ElementWithScores {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    if count <= 0 {
        return nil
    }
    if end == -1 {
        end = ss.Card()
    } else if end >= ss.Card() {
        end = ss.Card() - 1
    }
    node := ss.skipList.GetByRank(start)
    if node == nil {
        return nil
    }
    l := end - start + 1
    if count < l {
        l = count
    }
    res := make([]*ElementWithScores, 0, l)
    for node != nil && count > 0 && start <= end {
        res = append(res, ss.Mapper[node.element.String()])
        count--
        start++
        node = node.level[0].forward
    }
    return res
}

func (ss *SortSet) RangeByScore(spec *RangeSpec, count int64) []Element {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    if count <= 0 {
        return nil
    }
    node, k := ss.skipList.FirstInRange(spec)
    if node == nil {
        return nil
    }

    l := count
    if count > ss.Card()-k+1 {
        l = ss.Card() - k + 1
    }
    res := make([]Element, 0, l)
    for node != nil && spec.ValueLteMax(node.score) && count > 0 {
        res = append(res, node.element)
        count--
        node = node.level[0].forward
    }
    return res
}

func (ss *SortSet) RangeByScoreWithScores(spec *RangeSpec, count int64) []*ElementWithScores {
    ss.mu.RLock()
    defer ss.mu.RUnlock()
    if count <= 0 {
        return nil
    }
    node, k := ss.skipList.FirstInRange(spec)
    if node == nil {
        return nil
    }

    l := count
    if count > ss.Card()-k+2 {
        l = ss.Card() - k + 2
    }
    res := make([]*ElementWithScores, 0, l)
    for node != nil && spec.ValueLteMax(node.score) && count > 0 {
        res = append(res, ss.Mapper[node.element.String()])
        count--
        node = node.level[0].forward
    }
    return res
}
