package gskiplist

import (
    "fmt"
    "testing"
)

type MockElement string

func (me MockElement) String() string {
    return string(me)
}

func TestNewSortSet(t *testing.T) {
    s := NewSortSet()
    s.Add(1, MockElement("abc1"))
    s.Add(2, MockElement("abc2"))
    s.Add(3, MockElement("abc3"))
    s.Add(4, MockElement("abc4"))
    s.Add(5, MockElement("abc5"))
    s.Add(6, MockElement("abc6"))
    s.Add(7, MockElement("abc7"))
    fmt.Println(s.RangeByRank(0, 5, 10000))
    fmt.Println(s.Del(MockElement("abc7")))
    fmt.Println("------------------")
    fmt.Println(s.RangeByRank(0, -1, 10000))
    fmt.Println(3, 5,s.RangeByRank(3, 5, 10000))
    fmt.Println(0, 1,s.RangeByRankWithScores(0, 1, 10000))
    fmt.Println(s.RangeByScore(&RangeSpec{
        Min:   4,
        Max:   4,
        MinEx: false,
        MaxEx: false,
    }, 10000))
    fmt.Println(s.Add(3, MockElement("abc6")))
    fmt.Println(s.Add(3, MockElement("abc5")))
    fmt.Println("------------------")
    fmt.Println(s.RangeByRank(0, -1, 10000))
    fmt.Println(2, 4,s.RangeByRank(2, 4, 10000))
    fmt.Println(1, 6,s.RangeByRankWithScores(1, 6, 10000))
    fmt.Println(s.RangeByScore(&RangeSpec{
        Min:   2,
        Max:   10000,
        MinEx: false,
        MaxEx: false,
    }, 10000))
    fmt.Println(s.Del(MockElement("abc5")))
    fmt.Println("------------------")
    fmt.Println(s.RangeByRank(0, -1, 10000))
    fmt.Println(2, 5,s.RangeByRank(2, 5, 10000))
    fmt.Println(0, -1, s.RangeByRankWithScores(0, -1, 10000))
    fmt.Println(s.RangeByScore(&RangeSpec{
        Min:   1,
        Max:   3,
        MinEx: false,
        MaxEx: false,
    }, 10000))
    fmt.Println(s.Score(MockElement("abc5")))
    fmt.Println(s.Score(MockElement("abc6")))
    fmt.Println(s.Card())
}
