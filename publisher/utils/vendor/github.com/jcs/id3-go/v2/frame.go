// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/jcs/id3-go/encodedbytes"
)

const (
	FrameHeaderSize = 10
)

// FrameType holds frame id metadata and constructor method
// A set number of these are created in the version specific files
type FrameType struct {
	id          string
	description string
	constructor func(FrameHead, []byte) Framer
}

// Framer provides a generic interface for frames
// This is the default type returned when creating frames
type Framer interface {
	Id() string
	Size() uint
	StatusFlags() byte
	FormatFlags() byte
	String() string
	Bytes() []byte
	setOwner(*Tag)
}

// FrameHead represents the header of each frame
// Additional metadata is kept through the embedded frame type
// These do not usually need to be manually created
type FrameHead struct {
	FrameType
	statusFlags byte
	formatFlags byte
	size        uint32
	owner       *Tag
}

func (ft FrameType) Id() string {
	return ft.id
}

func (h FrameHead) Size() uint {
	return uint(h.size)
}

func (h *FrameHead) changeSize(diff int) {
	if diff >= 0 {
		h.size += uint32(diff)
	} else {
		h.size -= uint32(-diff)
	}

	if h.owner != nil {
		h.owner.changeSize(diff)
	}
}

func (h FrameHead) StatusFlags() byte {
	return h.statusFlags
}

func (h FrameHead) FormatFlags() byte {
	return h.formatFlags
}

func (h *FrameHead) setOwner(t *Tag) {
	h.owner = t
}

// DataFrame is the default frame for binary data
type DataFrame struct {
	FrameHead
	data []byte
}

func NewDataFrame(ft FrameType, data []byte) *DataFrame {
	head := FrameHead{
		FrameType: ft,
		size:      uint32(len(data)),
	}

	return &DataFrame{head, data}
}

func ParseDataFrame(head FrameHead, data []byte) Framer {
	return &DataFrame{head, data}
}

func (f DataFrame) Data() []byte {
	return f.data
}

func (f *DataFrame) SetData(b []byte) {
	diff := len(b) - len(f.data)
	f.changeSize(diff)
	f.data = b
}

func (f DataFrame) String() string {
	return "<binary data>"
}

func (f DataFrame) Bytes() []byte {
	return f.data
}

// IdFrame represents identification tags
type IdFrame struct {
	FrameHead
	ownerIdentifier string
	identifier      []byte
}

func NewIdFrame(ft FrameType, ownerId string, id []byte) *IdFrame {
	head := FrameHead{
		FrameType: ft,
		size:      uint32(1 + len(ownerId) + len(id)),
	}

	return &IdFrame{
		FrameHead:       head,
		ownerIdentifier: ownerId,
		identifier:      id,
	}
}

func ParseIdFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &IdFrame{FrameHead: head}
	rd := encodedbytes.NewReader(data)

	if f.ownerIdentifier, err = rd.ReadNullTermString(encodedbytes.NativeEncoding); err != nil {
		return nil
	}

	if f.identifier, err = rd.ReadRest(); len(f.identifier) > 64 || err != nil {
		return nil
	}

	return f
}

func (f IdFrame) OwnerIdentifier() string {
	return f.ownerIdentifier
}

func (f *IdFrame) SetOwnerIdentifier(ownerId string) {
	f.changeSize(len(ownerId) - len(f.ownerIdentifier))
	f.ownerIdentifier = ownerId
}

func (f IdFrame) Identifier() []byte {
	return f.identifier
}

func (f *IdFrame) SetIdentifier(id []byte) error {
	if len(id) > 64 {
		return errors.New("identifier: identifier too long")
	}

	f.changeSize(len(id) - len(f.identifier))
	f.identifier = id

	return nil
}

func (f IdFrame) String() string {
	return fmt.Sprintf("%s: %v", f.ownerIdentifier, f.identifier)
}

func (f IdFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := encodedbytes.NewWriter(bytes)

	if err = wr.WriteString(f.ownerIdentifier, encodedbytes.NativeEncoding); err != nil {
		return bytes
	}

	if _, err = wr.Write(f.identifier); err != nil {
		return bytes
	}

	return bytes
}

// TextFramer represents frames that contain encoded text
type TextFramer interface {
	Framer
	Encoding() string
	SetEncoding(string) error
	Text() string
	SetText(string) error
}

// TextFrame represents frames that contain encoded text
type TextFrame struct {
	FrameHead
	encoding byte
	text     string
}

func NewTextFrame(ft FrameType, text string) *TextFrame {
	head := FrameHead{
		FrameType: ft,
		size:      uint32(1 + len(text)),
	}

	return &TextFrame{
		FrameHead: head,
		text:      text,
	}
}

func ParseTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := &TextFrame{FrameHead: head}
	rd := encodedbytes.NewReader(data)

	if f.encoding, err = rd.ReadByte(); err != nil {
		return nil
	}

	if f.text, err = rd.ReadRestString(f.encoding); err != nil {
		return nil
	}

	return f
}

func (f TextFrame) Encoding() string {
	return encodedbytes.EncodingForIndex(f.encoding)
}

func (f *TextFrame) SetEncoding(encoding string) error {
	i := byte(encodedbytes.IndexForEncoding(encoding))
	if i < 0 {
		return errors.New("encoding: invalid encoding")
	}

	diff, err := encodedbytes.EncodedDiff(i, f.text, f.encoding, f.text)
	if err != nil {
		return err
	}

	f.changeSize(diff)
	f.encoding = i
	return nil
}

func (f TextFrame) Text() string {
	return f.text
}

func (f *TextFrame) SetText(text string) error {
	diff, err := encodedbytes.EncodedDiff(f.encoding, text, f.encoding, f.text)
	if err != nil {
		return err
	}

	f.changeSize(diff)
	f.text = text
	return nil
}

func (f TextFrame) String() string {
	return f.text
}

func (f TextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := encodedbytes.NewWriter(bytes)

	if err = wr.WriteByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

type DescTextFrame struct {
	TextFrame
	description string
}

func NewDescTextFrame(ft FrameType, desc, text string) *DescTextFrame {
	f := NewTextFrame(ft, text)
	nullLength := encodedbytes.EncodingNullLengthForIndex(f.encoding)
	f.size += uint32(len(desc) + nullLength)

	return &DescTextFrame{
		TextFrame:   *f,
		description: desc,
	}
}

// DescTextFrame represents frames that contain encoded text and descriptions
func ParseDescTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := new(DescTextFrame)
	f.FrameHead = head
	rd := encodedbytes.NewReader(data)

	if f.encoding, err = rd.ReadByte(); err != nil {
		return nil
	}
	f.size = uint32(1)

	if f.description, err = rd.ReadNullTermString(f.encoding); err != nil {
		return nil
	}
	l, err := encodedbytes.EncodedNullTermStringBytes(f.description, f.encoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	if f.text, err = rd.ReadRestString(f.encoding); err != nil {
		return nil
	}
	l, err = encodedbytes.EncodedStringBytes(f.text, f.encoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	return f
}

func (f DescTextFrame) Description() string {
	return f.description
}

func (f *DescTextFrame) SetDescription(description string) error {
	diff, err := encodedbytes.EncodedDiff(f.encoding, description, f.encoding, f.description)
	if err != nil {
		return err
	}

	f.changeSize(diff)
	f.description = description
	return nil
}

func (f *DescTextFrame) SetEncoding(encoding string) error {
	i := byte(encodedbytes.IndexForEncoding(encoding))
	if i < 0 {
		return errors.New("encoding: invalid encoding")
	}

	descDiff, err := encodedbytes.EncodedDiff(i, f.text, f.encoding, f.text)
	if err != nil {
		return err
	}

	newNullLength := encodedbytes.EncodingNullLengthForIndex(i)
	oldNullLength := encodedbytes.EncodingNullLengthForIndex(f.encoding)
	nullDiff := newNullLength - oldNullLength

	textDiff, err := encodedbytes.EncodedDiff(i, f.description, f.encoding, f.description)
	if err != nil {
		return err
	}

	f.changeSize(descDiff + nullDiff + textDiff)
	f.encoding = i
	return nil
}

func (f DescTextFrame) String() string {
	return fmt.Sprintf("%s: %s", f.description, f.text)
}

func (f DescTextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := encodedbytes.NewWriter(bytes)

	if err = wr.WriteByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteNullTermString(f.description, f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteNullTermString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

// UnsynchTextFrame represents frames that contain unsynchronized text
type UnsynchTextFrame struct {
	DescTextFrame
	language string
}

func NewUnsynchTextFrame(ft FrameType, desc, text string) *UnsynchTextFrame {
	f := NewDescTextFrame(ft, desc, text)
	f.size += uint32(3)

	return &UnsynchTextFrame{
		DescTextFrame: *f,
		language:      "eng",
	}
}

func ParseUnsynchTextFrame(head FrameHead, data []byte) Framer {
	var err error
	f := new(UnsynchTextFrame)
	f.FrameHead = head
	rd := encodedbytes.NewReader(data)

	if f.encoding, err = rd.ReadByte(); err != nil {
		return nil
	}
	f.size = uint32(1)

	if f.language, err = rd.ReadNumBytesString(3); err != nil {
		return nil
	}
	f.size += uint32(3)

	if f.description, err = rd.ReadNullTermString(f.encoding); err != nil {
		return nil
	}
	l, err := encodedbytes.EncodedNullTermStringBytes(f.description, f.encoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	if f.text, err = rd.ReadRestString(f.encoding); err != nil {
		return nil
	}
	l, err = encodedbytes.EncodedStringBytes(f.text, f.encoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	return f
}

func (f UnsynchTextFrame) Language() string {
	return f.language
}

func (f *UnsynchTextFrame) SetLanguage(language string) error {
	if len(language) != 3 {
		return errors.New("language: invalid language string")
	}

	f.language = language
	f.changeSize(0)
	return nil
}

func (f UnsynchTextFrame) String() string {
	return fmt.Sprintf("%s\t%s:\n%s", f.language, f.description, f.text)
}

func (f UnsynchTextFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := encodedbytes.NewWriter(bytes)

	if err = wr.WriteByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteString(f.language, encodedbytes.NativeEncoding); err != nil {
		return bytes
	}

	if err = wr.WriteNullTermString(f.description, f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteString(f.text, f.encoding); err != nil {
		return bytes
	}

	return bytes
}

// ImageFrame represent frames that have media attached
type ImageFrame struct {
	DataFrame
	encoding    byte
	mimeType    string
	pictureType byte
	description string
}

func ParseImageFrame(head FrameHead, data []byte) Framer {
	var err error
	f := new(ImageFrame)
	f.FrameHead = head
	rd := encodedbytes.NewReader(data)

	if f.encoding, err = rd.ReadByte(); err != nil {
		return nil
	}
	f.size = uint32(1)

	if f.mimeType, err = rd.ReadNullTermString(encodedbytes.NativeEncoding); err != nil {
		return nil
	}
	l, err := encodedbytes.EncodedNullTermStringBytes(f.mimeType, encodedbytes.NativeEncoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	if f.pictureType, err = rd.ReadByte(); err != nil {
		return nil
	}
	f.size += uint32(1)

	if f.description, err = rd.ReadNullTermString(f.encoding); err != nil {
		return nil
	}
	l, err = encodedbytes.EncodedNullTermStringBytes(f.description, f.encoding)
	if err != nil {
		return nil
	}
	f.size += uint32(len(l))

	if f.data, err = rd.ReadRest(); err != nil {
		return nil
	}
	f.size += uint32(len(f.data))

	return f
}

func (f ImageFrame) Encoding() string {
	return encodedbytes.EncodingForIndex(f.encoding)
}

func (f *ImageFrame) SetEncoding(encoding string) error {
	i := byte(encodedbytes.IndexForEncoding(encoding))
	if i < 0 {
		return errors.New("encoding: invalid encoding")
	}

	diff, err := encodedbytes.EncodedDiff(i, f.description, f.encoding, f.description)
	if err != nil {
		return err
	}

	f.changeSize(diff)
	f.encoding = i
	return nil
}

func (f ImageFrame) MIMEType() string {
	return f.mimeType
}

func (f *ImageFrame) SetMIMEType(mimeType string) {
	diff := len(mimeType) - len(f.mimeType)
	if mimeType[len(mimeType)-1] != 0 {
		nullTermBytes := append([]byte(mimeType), 0x00)
		f.mimeType = string(nullTermBytes)
		diff += 1
	} else {
		f.mimeType = mimeType
	}

	f.changeSize(diff)
}

func (f ImageFrame) String() string {
	return fmt.Sprintf("%s\t%s: <binary data>", f.mimeType, f.description)
}

func (f ImageFrame) Bytes() []byte {
	var err error
	bytes := make([]byte, f.Size())
	wr := encodedbytes.NewWriter(bytes)

	if err = wr.WriteByte(f.encoding); err != nil {
		return bytes
	}

	if err = wr.WriteNullTermString(f.mimeType, encodedbytes.NativeEncoding); err != nil {
		return bytes
	}

	if err = wr.WriteByte(f.pictureType); err != nil {
		return bytes
	}

	if err = wr.WriteNullTermString(f.description, f.encoding); err != nil {
		return bytes
	}

	if n, err := wr.Write(f.data); n < len(f.data) || err != nil {
		return bytes
	}

	return bytes
}

// ChapterFrame represents chapter frames
type ChapterFrame struct {
	FrameHead
	Element    string
	StartTime  uint32
	EndTime    uint32
	StartByte  uint32
	EndByte    uint32
	UseTime    bool
	titleFrame Framer
	linkFrame  Framer
}

func NewChapterFrame(ft FrameType, element string, startTime uint32, endTime uint32, startByte uint32, endByte uint32, useTime bool, title string, link string, linkTitle string) *ChapterFrame {
	var titleFrame Framer
	var linkFrame Framer

	if title != "" {
		ft := V23FrameTypeMap["TIT2"]
		titleFrame = NewTextFrame(ft, title)
		titleFrame.(*TextFrame).SetEncoding("UTF-8")
	}

	if link != "" {
		ft := V23FrameTypeMap["WXXX"]
		linkFrame = NewDescTextFrame(ft, linkTitle, link)
		linkFrame.(*DescTextFrame).SetEncoding("UTF-8")
	}

	head := FrameHead{
		FrameType: ft,
	}

	cf := &ChapterFrame{head, element, startTime, endTime, startByte, endByte, useTime, titleFrame, linkFrame}
	cf.size = uint32(len(cf.Bytes()))

	return cf
}

func ParseChapterFrame(head FrameHead, data []byte) Framer {
	var err error
	var d []byte
	var empty uint32
	f := new(ChapterFrame)
	f.FrameHead = head
	rd := encodedbytes.NewReader(data)

	// http://id3.org/id3v2-chapters-1.0

	empty = binary.BigEndian.Uint32([]byte{0xff, 0xff, 0xff, 0xff})

	if f.Element, err = rd.ReadNullTermString(encodedbytes.NativeEncoding); err != nil {
		return nil
	}

	if d, err = rd.ReadNumBytes(encodedbytes.BytesPerInt); err != nil {
		return nil
	}
	f.StartTime = binary.BigEndian.Uint32(d)

	if d, err = rd.ReadNumBytes(encodedbytes.BytesPerInt); err != nil {
		return nil
	}
	f.EndTime = binary.BigEndian.Uint32(d)

	if d, err = rd.ReadNumBytes(encodedbytes.BytesPerInt); err != nil {
		return nil
	}
	f.StartByte = binary.BigEndian.Uint32(d)

	if d, err = rd.ReadNumBytes(encodedbytes.BytesPerInt); err != nil {
		return nil
	}
	f.EndByte = binary.BigEndian.Uint32(d)

	if f.StartTime == empty && f.EndTime == empty {
		f.StartTime = 0
		f.EndTime = 0
		f.UseTime = false
	} else if f.StartByte == empty && f.EndByte == empty {
		f.StartByte = 0
		f.EndByte = 0
		f.UseTime = true
	} else {
		return nil
	}

	f.size = uint32(len(f.Element) + 1 + (4 * 4))

	if d, err = rd.ReadRest(); err != nil {
		return nil
	}

	// individual TIT2 labels will be subframes which are just normal frames
	// but contained within the CHAP frame's size
	if d != nil {
		var frame Framer
		dsize := len(d)
		pos := 0
		for pos < dsize {
			reader := bytes.NewReader(d[pos:])
			if frame = ParseV23Frame(reader); frame == nil {
				break
			}

			switch frame.Id() {
			case "TIT1", "TIT2", "TIT3":
				f.titleFrame = frame
			case "WXXX":
				f.linkFrame = frame
			}

			fsize := int(frame.Size()) + FrameHeaderSize
			pos += fsize
			f.size += uint32(fsize)
		}
	}

	return f
}

func (f ChapterFrame) String() string {
	if f.UseTime {
		return fmt.Sprintf("chapter: %d ms to %d ms: %v", f.StartTime, f.EndTime, f.Title())
	} else {
		return fmt.Sprintf("chapter: byte %d to %d: %v", f.StartByte, f.EndByte, f.Title())
	}
}

func (f ChapterFrame) Link() string {
	if f.linkFrame != nil {
		return f.linkFrame.(*DescTextFrame).Text()
	}
	return ""
}

func (f ChapterFrame) Title() string {
	if f.titleFrame != nil {
		return f.titleFrame.(*TextFrame).String()
	}
	return ""
}

func (f *ChapterFrame) Bytes() []byte {
	f.size = uint32(len(f.Element) + 1 + 4 + 4 + 4 + 4)

	var titleBytes []byte
	if f.titleFrame != nil {
		titleBytes = V23Bytes(f.titleFrame)
		f.size += uint32(len(titleBytes))
	}

	var linkBytes []byte
	if f.linkFrame != nil {
		linkBytes = V23Bytes(f.linkFrame)
		f.size += uint32(len(linkBytes))
	}

	bs := make([]byte, f.size)
	wr := encodedbytes.NewWriter(bs)

	if err := wr.WriteNullTermString(f.Element, encodedbytes.NativeEncoding); err != nil {
		return bs
	}

	b4 := make([]byte, 4)
	if f.UseTime {
		binary.BigEndian.PutUint32(b4, f.StartTime)
		if _, err := wr.Write(b4); err != nil {
			return bs
		}

		binary.BigEndian.PutUint32(b4, f.EndTime)
		if _, err := wr.Write(b4); err != nil {
			return bs
		}

		if _, err := wr.Write(bytes.Repeat([]byte{0xff}, 8)); err != nil {
			return bs
		}
	} else {
		if _, err := wr.Write(bytes.Repeat([]byte{0xff}, 8)); err != nil {
			return bs
		}

		binary.BigEndian.PutUint32(b4, f.StartByte)
		if _, err := wr.Write(b4); err != nil {
			return bs
		}

		binary.BigEndian.PutUint32(b4, f.EndByte)
		if _, err := wr.Write(b4); err != nil {
			return bs
		}
	}

	if f.titleFrame != nil {
		wr.Write(titleBytes)
	}
	if f.linkFrame != nil {
		wr.Write(linkBytes)
	}

	return bs
}

// TOCFrame represents Table of Contents frames
type TOCFrame struct {
	FrameHead
	Element       string
	TopLevel      bool
	Ordered       bool
	ChildElements []string
}

func NewTOCFrame(ft FrameType, element string, topLevel bool, ordered bool, childElements []string) *TOCFrame {
	head := FrameHead{
		FrameType: ft,
	}

	tf := &TOCFrame{head, element, topLevel, ordered, childElements}
	tf.size = uint32(len(tf.Bytes()))

	return tf
}

func ParseTOCFrame(head FrameHead, data []byte) Framer {
	var err error
	f := new(TOCFrame)
	f.FrameHead = head
	rd := encodedbytes.NewReader(data)

	if f.Element, err = rd.ReadNullTermString(encodedbytes.NativeEncoding); err != nil {
		return nil
	}

	f.size = uint32(len(f.Element) + 1)

	b, err := rd.ReadByte()
	if err != nil {
		return nil
	}
	f.Ordered = (b&(1<<0) != 0)
	f.TopLevel = (b&(1<<1) != 0)
	f.size += 1

	b, err = rd.ReadByte()
	if err != nil {
		return nil
	}
	f.size += 1

	for i := 0; i < int(b); i++ {
		s, err := rd.ReadNullTermString(encodedbytes.NativeEncoding)
		if err != nil {
			return nil
		}

		f.size += uint32(len(s) + 1)
		f.ChildElements = append(f.ChildElements, s)
	}

	return f
}

func (f *TOCFrame) SetChildElements(elements []string) {
	f.ChildElements = elements
	old := int(f.size)
	now := len(f.Bytes())
	f.changeSize(now - old)
}

func (f TOCFrame) String() string {
	return fmt.Sprintf("<TOC %v>", f.ChildElements)
}

func (f *TOCFrame) Bytes() []byte {
	var err error

	size := uint32(len(f.Element) + 1 + 1 + 1)
	for _, e := range f.ChildElements {
		size += uint32(len(e) + 1)
	}

	bs := make([]byte, size)
	wr := encodedbytes.NewWriter(bs)

	if err := wr.WriteNullTermString(f.Element, encodedbytes.NativeEncoding); err != nil {
		return bs
	}

	flags := 0
	if f.Ordered {
		flags |= (1 << 0)
	}
	if f.TopLevel {
		flags |= (1 << 1)
	}

	if err = wr.WriteByte(byte(flags)); err != nil {
		return bs
	}

	if err = wr.WriteByte(byte(len(f.ChildElements))); err != nil {
		return bs
	}

	for _, e := range f.ChildElements {
		if err := wr.WriteNullTermString(e, encodedbytes.NativeEncoding); err != nil {
			return bs
		}
	}

	return bs
}
