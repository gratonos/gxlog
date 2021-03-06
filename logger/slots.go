package logger

import (
	"reflect"

	"github.com/gratonos/gxlog/formatter"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/writer"
)

type SlotIndex int

const (
	Slot0 SlotIndex = iota
	Slot1
	Slot2
	Slot3
	Slot4
	Slot5
	Slot6
	Slot7
)

const MaxSlot = 8

type Slot struct {
	Formatter    iface.Formatter
	Writer       iface.Writer
	Level        iface.Level
	Filter       Filter
	ErrorHandler ErrorHandler
}

var nullSlot = Slot{
	Formatter:    formatter.Null(),
	Writer:       writer.Null(),
	Level:        iface.Off,
	Filter:       nullFilter,
	ErrorHandler: nullErrorHandler,
}

func (this *Logger) Slot(index SlotIndex) Slot {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index]
}

func (this *Logger) SetSlot(index SlotIndex, slot Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index] = fillSlot(slot)
	this.updateEquivalents()
}

func (this *Logger) UpdateSlot(index SlotIndex, fn func(Slot) Slot) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index] = fillSlot(fn(this.slots[index]))
	this.updateEquivalents()
}

func (this *Logger) ResetSlot(index SlotIndex) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index] = nullSlot
	this.updateEquivalents()
}

func (this *Logger) ResetAllSlots() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for i := range this.slots {
		this.slots[i] = nullSlot
	}
	this.updateEquivalents()
}

func (this *Logger) CopySlot(dst, src SlotIndex) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[dst] = this.slots[src]
	this.updateEquivalents()
}

func (this *Logger) MoveSlot(to, from SlotIndex) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[to] = this.slots[from]
	this.slots[from] = nullSlot
	this.updateEquivalents()
}

func (this *Logger) SwapSlot(left, right SlotIndex) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[left], this.slots[right] = this.slots[right], this.slots[left]
	this.updateEquivalents()
}

func (this *Logger) SlotFormatter(index SlotIndex) iface.Formatter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index].Formatter
}

func (this *Logger) SetSlotFormatter(index SlotIndex, formatter iface.Formatter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index].Formatter = fillFormatter(formatter)
	this.updateEquivalents()
}

func (this *Logger) SlotWriter(index SlotIndex) iface.Writer {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index].Writer
}

func (this *Logger) SetSlotWriter(index SlotIndex, writer iface.Writer) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index].Writer = fillWriter(writer)
}

func (this *Logger) SlotLevel(index SlotIndex) iface.Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index].Level
}

func (this *Logger) SetSlotLevel(index SlotIndex, level iface.Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index].Level = level
}

func (this *Logger) SlotFilter(index SlotIndex) Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index].Filter
}

func (this *Logger) SetSlotFilter(index SlotIndex, filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index].Filter = fillFilter(filter)
}

func (this *Logger) SlotErrorHandler(index SlotIndex) ErrorHandler {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.slots[index].ErrorHandler
}

func (this *Logger) SetSlotErrorHandler(index SlotIndex, handler ErrorHandler) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.slots[index].ErrorHandler = fillErrorHandler(handler)
}

func (this *Logger) updateEquivalents() {
	for i := 0; i < MaxSlot; i++ {
		this.equivalents[i] = this.equivalents[i][:0]
		if reflect.TypeOf(this.slots[i].Formatter).Comparable() {
			for j := i + 1; j < MaxSlot; j++ {
				if reflect.TypeOf(this.slots[j].Formatter).Comparable() &&
					this.slots[i].Formatter == this.slots[j].Formatter {
					this.equivalents[i] = append(this.equivalents[i], j)
				}
			}
		}
	}
}

func initSlots() []Slot {
	slots := make([]Slot, MaxSlot)
	for i := range slots {
		slots[i] = nullSlot
	}
	return slots
}

func fillSlot(slot Slot) Slot {
	slot.Formatter = fillFormatter(slot.Formatter)
	slot.Writer = fillWriter(slot.Writer)
	slot.Filter = fillFilter(slot.Filter)
	slot.ErrorHandler = fillErrorHandler(slot.ErrorHandler)
	return slot
}

func fillFormatter(formatter iface.Formatter) iface.Formatter {
	if formatter == nil {
		formatter = nullSlot.Formatter
	}
	return formatter
}

func fillWriter(writer iface.Writer) iface.Writer {
	if writer == nil {
		writer = nullSlot.Writer
	}
	return writer
}
