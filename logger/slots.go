package logger

import (
	"reflect"

	"github.com/gratonos/gxlog/formatter"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/writer"
)

type Slot int

const (
	Slot0 Slot = iota
	Slot1
	Slot2
	Slot3
	Slot4
	Slot5
	Slot6
	Slot7
)

const MaxSlot = 8

type slotLink struct {
	Formatter    iface.Formatter
	Writer       iface.Writer
	Level        iface.Level
	Filter       Filter
	ErrorHandler ErrorHandler
}

var nullSlotLink = slotLink{
	Formatter: formatter.Null(),
	Writer:    writer.Null(),
	Level:     iface.Off,
}

func (log *Logger) Link(slot Slot, formatter iface.Formatter, writer iface.Writer,
	level iface.Level, filter Filter, handler ErrorHandler) {

	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = slotLink{
		Formatter:    formatter,
		Writer:       writer,
		Level:        level,
		Filter:       filter,
		ErrorHandler: handler,
	}
	log.updateEquivalents()
}

func (log *Logger) Unlink(slot Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot] = nullSlotLink
	log.updateEquivalents()
}

func (log *Logger) UnlinkAll() {
	log.lock.Lock()
	defer log.lock.Unlock()

	for i := range log.slots {
		log.slots[i] = nullSlotLink
	}
	log.updateEquivalents()
}

func (log *Logger) CopySlot(dst, src Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[dst] = log.slots[src]
	log.updateEquivalents()
}

func (log *Logger) MoveSlot(to, from Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[to] = log.slots[from]
	log.slots[from] = nullSlotLink
	log.updateEquivalents()
}

func (log *Logger) SwapSlot(left, right Slot) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[left], log.slots[right] = log.slots[right], log.slots[left]
	log.updateEquivalents()
}

func (log *Logger) SlotFormatter(slot Slot) iface.Formatter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Formatter
}

func (log *Logger) SetSlotFormatter(slot Slot, formatter iface.Formatter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Formatter = formatter
	log.updateEquivalents()
}

func (log *Logger) SlotWriter(slot Slot) iface.Writer {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Writer
}

func (log *Logger) SetSlotWriter(slot Slot, writer iface.Writer) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Writer = writer
}

func (log *Logger) SlotLevel(slot Slot) iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Level
}

func (log *Logger) SetSlotLevel(slot Slot, level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Level = level
}

func (log *Logger) SlotFilter(slot Slot) Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].Filter
}

func (log *Logger) SetSlotFilter(slot Slot, filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].Filter = filter
}

func (log *Logger) SlotErrorHandler(slot Slot) ErrorHandler {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.slots[slot].ErrorHandler
}

func (log *Logger) SetSlotErrorHandler(slot Slot, handler ErrorHandler) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.slots[slot].ErrorHandler = handler
}

func (log *Logger) updateEquivalents() {
	for i := 0; i < MaxSlot; i++ {
		log.equivalents[i] = log.equivalents[i][:0]
		if !reflect.TypeOf(log.slots[i].Formatter).Comparable() {
			continue
		}
		for j := i + 1; j < MaxSlot; j++ {
			if !reflect.TypeOf(log.slots[j].Formatter).Comparable() ||
				log.slots[i].Formatter != log.slots[j].Formatter {
				continue
			}
			log.equivalents[i] = append(log.equivalents[i], j)
		}
	}
}

func initSlots() []slotLink {
	var slots []slotLink
	for slot := 0; slot < MaxSlot; slot++ {
		slots = append(slots, nullSlotLink)
	}
	return slots
}
