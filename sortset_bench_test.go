package gskiplist

import (
    "fmt"
    "strconv"
    "sync"
    "testing"
    "time"
)

const nnn = 100000

func TestNewSortSet1(t *testing.T) {
    s := NewSortSet()
    wg := sync.WaitGroup{}
    for i := 0; i < nnn; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            s.Add(float64(i), MockElement(fmt.Sprintf("%d", i)))
        }(i)
    }
    for i := 0; i < nnn; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            time.Sleep(5 * time.Microsecond)
            s.Del(MockElement(fmt.Sprintf("%d", i)))
        }(i)
    }
    //for i := 0; i < nnn; i++ {
    //   wg.Add(1)
    //   go func() {
    //       defer wg.Done()
    //       s.RangeByRank(0, 1, 10000)
    //   }()
    //}
    for i := 0; i < nnn; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            s.RangeByScore(&RangeSpec{
                Min:   0,
                Max:   19999,
                MinEx: false,
                MaxEx: false,
            }, 10000)
        }()
    }
    wg.Wait()
}

func BenchmarkNewSortSet_Add_Go(b *testing.B) {
    s := NewSortSet()
    wg := sync.WaitGroup{}
    for i := 0; i < b.N; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            s.Add(float64(i), MockElement(strconv.Itoa(i)))
        }(i)
    }
    wg.Wait()
}

func BenchmarkNewSortSet_Add(b *testing.B) {
    s := NewSortSet()
    for i := 0; i < b.N; i++ {
        s.Add(float64(i), MockElement(strconv.Itoa(i)))
    }
}

func BenchmarkNewSkipList_Insert(b *testing.B) {
    s := NewSkipList()
    for i := 0; i < b.N; i++ {
        s.Insert(float64(i), MockElement(strconv.Itoa(i)))
    }
}

func BenchmarkNewSkipList_Update(b *testing.B) {
    s := NewSkipList()
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        s.Insert(float64(i), MockElement(strconv.Itoa(i)))
    }
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        s.Insert(1000, MockElement(strconv.Itoa(i)))
    }
}

func BenchmarkNewSortSet_Del(b *testing.B) {
    s := NewSortSet()
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        s.Add(float64(i), MockElement(strconv.Itoa(i)))
    }
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        s.Del(MockElement(strconv.Itoa(i)))
    }
}

func TestNil_String(t *testing.T) {
    fmt.Println(strconv.FormatUint(9223372036854775808, 10))
    fmt.Println(strconv.FormatUint(9223372036854775808, 16))
    fmt.Println(strconv.FormatUint(9223372036854775808, 32))
    fmt.Println(strconv.FormatUint(9223372036854775808, 2))
}

func Benchmark111(t *testing.B) {
    for i := 0; i < t.N ; i++ {
        strconv.FormatUint(9223372036854775808, 10)
    }
}


func Benchmark112(t *testing.B) {
    for i := 0; i < t.N ; i++ {
        strconv.FormatUint(9223372036854775808, 2)
    }
}

func Benchmark113(t *testing.B) {
    for i := 0; i < t.N ; i++ {
        strconv.FormatUint(9223372036854775808, 32)
    }
}

func Benchmark114(t *testing.B) {
    for i := 0; i < t.N ; i++ {
        strconv.FormatUint(9223372036854775808, 16)
    }
}
