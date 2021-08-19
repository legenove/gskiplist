package gskiplist

const (
    ExCur = 1 // 是否除了当前值 0 为>= | 1 为 >
    Inf   = 2 // 是否是最大值
)

type RangeSpec struct {
    Min, Max     float64
    MinEx, MaxEx uint8 // 是否除了最大(小)值
}

func (spec *RangeSpec) ValueGteMin(value float64) bool {
    if spec.MinEx & Inf > 0 {
        return true
    }
    if spec.MinEx & ExCur > 0 {
        return value > spec.Min
    }
    return value >= spec.Min
}

func (spec *RangeSpec) ValueLteMax(value float64) bool {
    if spec.MaxEx & Inf > 0 {
        return true
    }
    if spec.MaxEx & ExCur > 0 {
        return value < spec.Max
    }
    return value <= spec.Max
}

func (spec *RangeSpec) OpenRange() bool {
    if spec.MaxEx & ExCur > 0 || spec.MinEx & ExCur > 0 {
        return true
    }
    return false
}