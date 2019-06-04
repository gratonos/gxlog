package logger

import (
	"reflect"

	"github.com/gratonos/gxlog/iface"
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
	Level: iface.Off,
}

func (log *Logger) Slot(index SlotIndex) Slot {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index]
}

func (log *Logger) SetSlot(index SlotIndex, slot Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index] = slot
	log.updateEquivalents()
}

func (log *Logger) UpdateSlot(index SlotIndex, fn func(Slot) Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index] = fn(log.slots[index])
	log.updateEquivalents()
}

func (log *Logger) ResetSlot(index SlotIndex) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index] = nullSlot
	log.updateEquivalents()
}

func (log *Logger) ResetAllSlots() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for i := range log.slots {
		log.slots[i] = nullSlot
	}
	log.updateEquivalents()
}

func (log *Logger) CopySlot(dst, src SlotIndex) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[dst] = log.slots[src]
	log.updateEquivalents()
}

func (log *Logger) MoveSlot(to, from SlotIndex) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[to] = log.slots[from]
	log.slots[from] = nullSlot
	log.updateEquivalents()
}

func (log *Logger) SwapSlot(left, right SlotIndex) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[left], log.slots[right] = log.slots[right], log.slots[left]
	log.updateEquivalents()
}

func (log *Logger) SlotFormatter(index SlotIndex) iface.Formatter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index].Formatter
}

func (log *Logger) SetSlotFormatter(index SlotIndex, formatter iface.Formatter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index].Formatter = formatter
	log.updateEquivalents()
}

func (log *Logger) SlotWriter(index SlotIndex) iface.Writer {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index].Writer
}

func (log *Logger) SetSlotWriter(index SlotIndex, writer iface.Writer) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index].Writer = writer
}

func (log *Logger) SlotLevel(index SlotIndex) iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index].Level
}

func (log *Logger) SetSlotLevel(index SlotIndex, level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index].Level = level
}

func (log *Logger) SlotFilter(index SlotIndex) Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index].Filter
}

func (log *Logger) SetSlotFilter(index SlotIndex, filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index].Filter = filter
}

func (log *Logger) SlotErrorHandler(index SlotIndex) ErrorHandler {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[index].ErrorHandler
}

func (log *Logger) SetSlotErrorHandler(index SlotIndex, handler ErrorHandler) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[index].ErrorHandler = handler
}

func (log *Logger) updateEquivalents() {
	for i := 0; i < MaxSlot; i++ {
		log.equivalents[i] = log.equivalents[i][:0]
		if log.slots[i].Formatter == nil ||
			!reflect.TypeOf(log.slots[i].Formatter).Comparable() {
			continue
		}
		for j := i + 1; j < MaxSlot; j++ {
			if log.slots[j].Formatter == nil ||
				!reflect.TypeOf(log.slots[j].Formatter).Comparable() ||
				log.slots[i].Formatter != log.slots[j].Formatter {
				continue
			}
			log.equivalents[i] = append(log.equivalents[i], j)
		}
	}
}

func initSlots() []Slot {
	var slots []Slot
	for i := 0; i < MaxSlot; i++ {
		slots = append(slots, nullSlot)
	}
	return slots
}
