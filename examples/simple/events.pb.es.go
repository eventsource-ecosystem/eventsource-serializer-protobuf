// Code generated by eventsource-protobuf. DO NOT EDIT.
// source: events.proto

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/eventsource-ecosystem/eventsource"
	"github.com/gogo/protobuf/proto"
)

type serializer struct {
}

func (s *serializer) MarshalEvent(event eventsource.Event) (eventsource.Record, error) {
	data, err := MarshalEvent(event)
	if err != nil {
		return eventsource.Record{}, err
	}

	return eventsource.Record{
		Version: event.EventVersion(),
		Data:    data,
	}, nil
}

func (s *serializer) UnmarshalEvent(record eventsource.Record) (eventsource.Event, error) {
	return UnmarshalEvent(record.Data)
}

func NewSerializer() eventsource.Serializer {
	return &serializer{}
}

func (m *A) AggregateID() string { return m.Id }
func (m *A) EventVersion() int   { return int(m.Version) }
func (m *A) EventAt() time.Time  { return time.Unix(m.At, 0) }

func (m *B) AggregateID() string { return m.ID }
func (m *B) EventVersion() int   { return int(m.Version) }
func (m *B) EventAt() time.Time  { return time.Unix(m.At, 0) }


func MarshalEvent(event eventsource.Event) ([]byte, error) {
	container := &EventContainer{}

	switch v := event.(type) {

	case *A:
		container.Type = 2
		container.Ma = v

	case *B:
		container.Type = 3
		container.Mb = v

	default:
		return nil, fmt.Errorf("Unhandled type, %v", event)
	}

	data, err := proto.Marshal(container)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalEvent(data []byte) (eventsource.Event, error) {
	container := &EventContainer{};
	err := proto.Unmarshal(data, container)
	if err != nil {
		return nil, err
	}

	var event interface{}
	switch container.Type {

	case 2:
		event = container.Ma

	case 3:
		event = container.Mb

	default:
		return nil, fmt.Errorf("Unhandled type, %v", container.Type)
	}

	return event.(eventsource.Event), nil
}

type Encoder struct{
	w io.Writer
}

func (e *Encoder) WriteEvent(event eventsource.Event) (int, error) {
	data, err := MarshalEvent(event)
	if err != nil {
		return 0, err
	}

	// Write the length of the marshaled event as uint64
	//
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, uint64(len(data)))
	if _, err := e.w.Write(buffer); err != nil {
		return 0, err
	}

	n, err := e.w.Write(data)
	if err != nil {
		return 0, err
	}

	return n + 8, nil
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

type Decoder struct {
	r       *bufio.Reader
	scratch *bytes.Buffer
}

func (d *Decoder) readN(n uint64) ([]byte, error) {
	d.scratch.Reset()
	for i := uint64(0); i < n; i++ {
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if err := d.scratch.WriteByte(b); err != nil {
			return nil, err
		}
	}
	return d.scratch.Bytes(), nil
}

func (d *Decoder) ReadEvent() (eventsource.Event, error) {
	data, err := d.readN(8)
	if err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint64(data)

	data, err = d.readN(length)
	if err != nil {
		return nil, err
	}

	event, err := UnmarshalEvent(data)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder {
		r:       bufio.NewReader(r),
		scratch: bytes.NewBuffer(nil),
	}
}

type Builder struct {
	id      string
	version int
	Events  []eventsource.Event
}

func NewBuilder(id string, version int) *Builder {
	return &Builder {
		id:      id,
		version: version,
	}
}

func (b *Builder) nextVersion() int32 {
	b.version++
	return int32(b.version)
}


func (b *Builder) A() {
	event := &A{
		Id:      b.id,
		Version: b.nextVersion(),
		At:      time.Now().Unix(),

	}
	b.Events = append(b.Events, event)
}

func (b *Builder) B(name string, ) {
	event := &B{
		ID:      b.id,
		Version: b.nextVersion(),
		At:      time.Now().Unix(),
	Name: name,

	}
	b.Events = append(b.Events, event)
}

