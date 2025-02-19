package interpreter

type Frame struct {
	slots []Value
	ip    int
}

func (f *Frame) Slot(idx int) (Value, bool) {
	if len(f.slots) <= idx {
		return nil, false
	}
	val := f.slots[idx]
	if val == nil {
		return nil, false
	}
	return val, true
}

func (f *Frame) SetSlot(idx int, val Value) {
	if len(f.slots) <= idx {
		slots := make([]Value, (idx+1)*2)
		copy(slots, f.slots)
		f.slots = slots
	}
	f.slots[idx] = val
}
