package gskiplist

type RangeSpec struct {
    Min, Max     float64
    MinEx, MaxEx bool // 是否除了最大(小)值
}

func (spec *RangeSpec) ValueGteMin(value float64) bool {
    if spec.MinEx {
        return value > spec.Min
    }
    return value >= spec.Min
}

func (spec *RangeSpec) ValueLteMax(value float64) bool {
    if spec.MaxEx {
        return value < spec.Max
    }
    return value <= spec.Max
}