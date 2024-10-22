// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: streams.proto

package pb

// import (
// 	fmt "fmt"
// 	proto "github.com/gogo/protobuf/proto"
// 	io "io"
// 	math "math"
// 	math_bits "math/bits"
// )

// // Reference imports to suppress errors if they are not otherwise used.
// var _ = proto.Marshal
// var _ = fmt.Errorf
// var _ = math.Inf

// // This is a compile-time assertion to ensure that this generated file
// // is compatible with the proto package it is being compiled against.
// // A compilation error at this line likely means your copy of the
// // proto package needs to be updated.
// const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// // BasicMessage 定义了基本的消息结构
// type BasicMessage struct {
// 	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
// 	Sender               string   `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
// 	Receiver             string   `protobuf:"bytes,3,opt,name=receiver,proto3" json:"receiver,omitempty"`
// 	XXX_NoUnkeyedLiteral struct{} `json:"-"`
// 	XXX_unrecognized     []byte   `json:"-"`
// 	XXX_sizecache        int32    `json:"-"`
// }

// func (m *BasicMessage) Reset()         { *m = BasicMessage{} }
// func (m *BasicMessage) String() string { return proto.CompactTextString(m) }
// func (*BasicMessage) ProtoMessage()    {}
// func (*BasicMessage) Descriptor() ([]byte, []int) {
// 	return fileDescriptor_c6bbf8af0ec331d6, []int{0}
// }
// func (m *BasicMessage) XXX_Unmarshal(b []byte) error {
// 	return m.Unmarshal(b)
// }
// func (m *BasicMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
// 	if deterministic {
// 		return xxx_messageInfo_BasicMessage.Marshal(b, m, deterministic)
// 	} else {
// 		b = b[:cap(b)]
// 		n, err := m.MarshalToSizedBuffer(b)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return b[:n], nil
// 	}
// }
// func (m *BasicMessage) XXX_Merge(src proto.Message) {
// 	xxx_messageInfo_BasicMessage.Merge(m, src)
// }
// func (m *BasicMessage) XXX_Size() int {
// 	return m.Size()
// }
// func (m *BasicMessage) XXX_DiscardUnknown() {
// 	xxx_messageInfo_BasicMessage.DiscardUnknown(m)
// }

// var xxx_messageInfo_BasicMessage proto.InternalMessageInfo

// func (m *BasicMessage) GetType() string {
// 	if m != nil {
// 		return m.Type
// 	}
// 	return ""
// }

// func (m *BasicMessage) GetSender() string {
// 	if m != nil {
// 		return m.Sender
// 	}
// 	return ""
// }

// func (m *BasicMessage) GetReceiver() string {
// 	if m != nil {
// 		return m.Receiver
// 	}
// 	return ""
// }

// // StreamRequestMessage 定义了请求消息的结构
// type StreamRequestMessage struct {
// 	Message              *BasicMessage `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
// 	Payload              []byte        `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
// 	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
// 	XXX_unrecognized     []byte        `json:"-"`
// 	XXX_sizecache        int32         `json:"-"`
// }

// func (m *StreamRequestMessage) Reset()         { *m = StreamRequestMessage{} }
// func (m *StreamRequestMessage) String() string { return proto.CompactTextString(m) }
// func (*StreamRequestMessage) ProtoMessage()    {}
// func (*StreamRequestMessage) Descriptor() ([]byte, []int) {
// 	return fileDescriptor_c6bbf8af0ec331d6, []int{1}
// }
// func (m *StreamRequestMessage) XXX_Unmarshal(b []byte) error {
// 	return m.Unmarshal(b)
// }
// func (m *StreamRequestMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
// 	if deterministic {
// 		return xxx_messageInfo_StreamRequestMessage.Marshal(b, m, deterministic)
// 	} else {
// 		b = b[:cap(b)]
// 		n, err := m.MarshalToSizedBuffer(b)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return b[:n], nil
// 	}
// }
// func (m *StreamRequestMessage) XXX_Merge(src proto.Message) {
// 	xxx_messageInfo_StreamRequestMessage.Merge(m, src)
// }
// func (m *StreamRequestMessage) XXX_Size() int {
// 	return m.Size()
// }
// func (m *StreamRequestMessage) XXX_DiscardUnknown() {
// 	xxx_messageInfo_StreamRequestMessage.DiscardUnknown(m)
// }

// var xxx_messageInfo_StreamRequestMessage proto.InternalMessageInfo

// func (m *StreamRequestMessage) GetMessage() *BasicMessage {
// 	if m != nil {
// 		return m.Message
// 	}
// 	return nil
// }

// func (m *StreamRequestMessage) GetPayload() []byte {
// 	if m != nil {
// 		return m.Payload
// 	}
// 	return nil
// }

// // StreamResponseMessage 定义了响应消息的结构
// type StreamResponseMessage struct {
// 	Message              *BasicMessage `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
// 	Code                 int32         `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
// 	Msg                  string        `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
// 	Data                 []byte        `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
// 	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
// 	XXX_unrecognized     []byte        `json:"-"`
// 	XXX_sizecache        int32         `json:"-"`
// }

// func (m *StreamResponseMessage) Reset()         { *m = StreamResponseMessage{} }
// func (m *StreamResponseMessage) String() string { return proto.CompactTextString(m) }
// func (*StreamResponseMessage) ProtoMessage()    {}
// func (*StreamResponseMessage) Descriptor() ([]byte, []int) {
// 	return fileDescriptor_c6bbf8af0ec331d6, []int{2}
// }
// func (m *StreamResponseMessage) XXX_Unmarshal(b []byte) error {
// 	return m.Unmarshal(b)
// }
// func (m *StreamResponseMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
// 	if deterministic {
// 		return xxx_messageInfo_StreamResponseMessage.Marshal(b, m, deterministic)
// 	} else {
// 		b = b[:cap(b)]
// 		n, err := m.MarshalToSizedBuffer(b)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return b[:n], nil
// 	}
// }
// func (m *StreamResponseMessage) XXX_Merge(src proto.Message) {
// 	xxx_messageInfo_StreamResponseMessage.Merge(m, src)
// }
// func (m *StreamResponseMessage) XXX_Size() int {
// 	return m.Size()
// }
// func (m *StreamResponseMessage) XXX_DiscardUnknown() {
// 	xxx_messageInfo_StreamResponseMessage.DiscardUnknown(m)
// }

// var xxx_messageInfo_StreamResponseMessage proto.InternalMessageInfo

// func (m *StreamResponseMessage) GetMessage() *BasicMessage {
// 	if m != nil {
// 		return m.Message
// 	}
// 	return nil
// }

// func (m *StreamResponseMessage) GetCode() int32 {
// 	if m != nil {
// 		return m.Code
// 	}
// 	return 0
// }

// func (m *StreamResponseMessage) GetMsg() string {
// 	if m != nil {
// 		return m.Msg
// 	}
// 	return ""
// }

// func (m *StreamResponseMessage) GetData() []byte {
// 	if m != nil {
// 		return m.Data
// 	}
// 	return nil
// }

// func init() {
// 	proto.RegisterType((*BasicMessage)(nil), "pb.BasicMessage")
// 	proto.RegisterType((*StreamRequestMessage)(nil), "pb.StreamRequestMessage")
// 	proto.RegisterType((*StreamResponseMessage)(nil), "pb.StreamResponseMessage")
// }

// func init() { proto.RegisterFile("streams.proto", fileDescriptor_c6bbf8af0ec331d6) }

// var fileDescriptor_c6bbf8af0ec331d6 = []byte{
// 	// 232 bytes of a gzipped FileDescriptorProto
// 	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0xb1, 0x4a, 0x04, 0x31,
// 	0x10, 0x86, 0xc9, 0xdd, 0x7a, 0xa7, 0xe3, 0x0a, 0xc7, 0xa0, 0x12, 0x2c, 0x16, 0xd9, 0x4a, 0x2c,
// 	0xb6, 0xd0, 0x37, 0xb8, 0xde, 0x26, 0x82, 0x95, 0x4d, 0x76, 0x33, 0x1c, 0x07, 0xee, 0x25, 0x66,
// 	0xa2, 0x70, 0x85, 0xef, 0x67, 0xe9, 0x23, 0xc8, 0x3e, 0x89, 0xec, 0xb8, 0x11, 0xeb, 0xeb, 0xbe,
// 	0x3f, 0x19, 0xbe, 0x7f, 0x12, 0x38, 0xe3, 0x14, 0xc9, 0xf6, 0xdc, 0x84, 0xe8, 0x93, 0xc7, 0x59,
// 	0x68, 0xeb, 0x27, 0x28, 0xd7, 0x96, 0xb7, 0xdd, 0x03, 0x31, 0xdb, 0x0d, 0x21, 0x42, 0x91, 0xf6,
// 	0x81, 0xb4, 0xba, 0x56, 0x37, 0x27, 0x46, 0x18, 0x2f, 0x61, 0xc1, 0xb4, 0x73, 0x14, 0xf5, 0x4c,
// 	0x4e, 0xa7, 0x84, 0x57, 0x70, 0x1c, 0xa9, 0xa3, 0xed, 0x3b, 0x45, 0x3d, 0x97, 0x9b, 0xbf, 0x5c,
// 	0x3f, 0xc3, 0xf9, 0xa3, 0x94, 0x19, 0x7a, 0x7d, 0x23, 0x4e, 0xd9, 0x7f, 0x0b, 0xcb, 0xfe, 0x17,
// 	0xa5, 0xe2, 0xf4, 0x6e, 0xd5, 0x84, 0xb6, 0xf9, 0xbf, 0x82, 0xc9, 0x03, 0xa8, 0x61, 0x19, 0xec,
// 	0xfe, 0xc5, 0x5b, 0x27, 0xc5, 0xa5, 0xc9, 0xb1, 0xfe, 0x80, 0x8b, 0x6c, 0xe7, 0xe0, 0x77, 0x4c,
// 	0x87, 0xe8, 0x11, 0x8a, 0xce, 0x3b, 0x12, 0xf7, 0x91, 0x11, 0xc6, 0x15, 0xcc, 0x7b, 0xde, 0x4c,
// 	0xaf, 0x19, 0x71, 0x9c, 0x72, 0x36, 0x59, 0x5d, 0xc8, 0x06, 0xc2, 0xeb, 0xf2, 0x73, 0xa8, 0xd4,
// 	0xd7, 0x50, 0xa9, 0xef, 0xa1, 0x52, 0xed, 0x42, 0x7e, 0xf3, 0xfe, 0x27, 0x00, 0x00, 0xff, 0xff,
// 	0xe2, 0x32, 0x13, 0xa7, 0x5e, 0x01, 0x00, 0x00,
// }

// func (m *BasicMessage) Marshal() (dAtA []byte, err error) {
// 	size := m.Size()
// 	dAtA = make([]byte, size)
// 	n, err := m.MarshalToSizedBuffer(dAtA[:size])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dAtA[:n], nil
// }

// func (m *BasicMessage) MarshalTo(dAtA []byte) (int, error) {
// 	size := m.Size()
// 	return m.MarshalToSizedBuffer(dAtA[:size])
// }

// func (m *BasicMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
// 	i := len(dAtA)
// 	_ = i
// 	var l int
// 	_ = l
// 	if m.XXX_unrecognized != nil {
// 		i -= len(m.XXX_unrecognized)
// 		copy(dAtA[i:], m.XXX_unrecognized)
// 	}
// 	if len(m.Receiver) > 0 {
// 		i -= len(m.Receiver)
// 		copy(dAtA[i:], m.Receiver)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Receiver)))
// 		i--
// 		dAtA[i] = 0x1a
// 	}
// 	if len(m.Sender) > 0 {
// 		i -= len(m.Sender)
// 		copy(dAtA[i:], m.Sender)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Sender)))
// 		i--
// 		dAtA[i] = 0x12
// 	}
// 	if len(m.Type) > 0 {
// 		i -= len(m.Type)
// 		copy(dAtA[i:], m.Type)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Type)))
// 		i--
// 		dAtA[i] = 0xa
// 	}
// 	return len(dAtA) - i, nil
// }

// func (m *StreamRequestMessage) Marshal() (dAtA []byte, err error) {
// 	size := m.Size()
// 	dAtA = make([]byte, size)
// 	n, err := m.MarshalToSizedBuffer(dAtA[:size])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dAtA[:n], nil
// }

// func (m *StreamRequestMessage) MarshalTo(dAtA []byte) (int, error) {
// 	size := m.Size()
// 	return m.MarshalToSizedBuffer(dAtA[:size])
// }

// func (m *StreamRequestMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
// 	i := len(dAtA)
// 	_ = i
// 	var l int
// 	_ = l
// 	if m.XXX_unrecognized != nil {
// 		i -= len(m.XXX_unrecognized)
// 		copy(dAtA[i:], m.XXX_unrecognized)
// 	}
// 	if len(m.Payload) > 0 {
// 		i -= len(m.Payload)
// 		copy(dAtA[i:], m.Payload)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Payload)))
// 		i--
// 		dAtA[i] = 0x12
// 	}
// 	if m.Message != nil {
// 		{
// 			size, err := m.Message.MarshalToSizedBuffer(dAtA[:i])
// 			if err != nil {
// 				return 0, err
// 			}
// 			i -= size
// 			i = encodeVarintStreams(dAtA, i, uint64(size))
// 		}
// 		i--
// 		dAtA[i] = 0xa
// 	}
// 	return len(dAtA) - i, nil
// }

// func (m *StreamResponseMessage) Marshal() (dAtA []byte, err error) {
// 	size := m.Size()
// 	dAtA = make([]byte, size)
// 	n, err := m.MarshalToSizedBuffer(dAtA[:size])
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dAtA[:n], nil
// }

// func (m *StreamResponseMessage) MarshalTo(dAtA []byte) (int, error) {
// 	size := m.Size()
// 	return m.MarshalToSizedBuffer(dAtA[:size])
// }

// func (m *StreamResponseMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
// 	i := len(dAtA)
// 	_ = i
// 	var l int
// 	_ = l
// 	if m.XXX_unrecognized != nil {
// 		i -= len(m.XXX_unrecognized)
// 		copy(dAtA[i:], m.XXX_unrecognized)
// 	}
// 	if len(m.Data) > 0 {
// 		i -= len(m.Data)
// 		copy(dAtA[i:], m.Data)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Data)))
// 		i--
// 		dAtA[i] = 0x22
// 	}
// 	if len(m.Msg) > 0 {
// 		i -= len(m.Msg)
// 		copy(dAtA[i:], m.Msg)
// 		i = encodeVarintStreams(dAtA, i, uint64(len(m.Msg)))
// 		i--
// 		dAtA[i] = 0x1a
// 	}
// 	if m.Code != 0 {
// 		i = encodeVarintStreams(dAtA, i, uint64(m.Code))
// 		i--
// 		dAtA[i] = 0x10
// 	}
// 	if m.Message != nil {
// 		{
// 			size, err := m.Message.MarshalToSizedBuffer(dAtA[:i])
// 			if err != nil {
// 				return 0, err
// 			}
// 			i -= size
// 			i = encodeVarintStreams(dAtA, i, uint64(size))
// 		}
// 		i--
// 		dAtA[i] = 0xa
// 	}
// 	return len(dAtA) - i, nil
// }

// func encodeVarintStreams(dAtA []byte, offset int, v uint64) int {
// 	offset -= sovStreams(v)
// 	base := offset
// 	for v >= 1<<7 {
// 		dAtA[offset] = uint8(v&0x7f | 0x80)
// 		v >>= 7
// 		offset++
// 	}
// 	dAtA[offset] = uint8(v)
// 	return base
// }
// func (m *BasicMessage) Size() (n int) {
// 	if m == nil {
// 		return 0
// 	}
// 	var l int
// 	_ = l
// 	l = len(m.Type)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	l = len(m.Sender)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	l = len(m.Receiver)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	if m.XXX_unrecognized != nil {
// 		n += len(m.XXX_unrecognized)
// 	}
// 	return n
// }

// func (m *StreamRequestMessage) Size() (n int) {
// 	if m == nil {
// 		return 0
// 	}
// 	var l int
// 	_ = l
// 	if m.Message != nil {
// 		l = m.Message.Size()
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	l = len(m.Payload)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	if m.XXX_unrecognized != nil {
// 		n += len(m.XXX_unrecognized)
// 	}
// 	return n
// }

// func (m *StreamResponseMessage) Size() (n int) {
// 	if m == nil {
// 		return 0
// 	}
// 	var l int
// 	_ = l
// 	if m.Message != nil {
// 		l = m.Message.Size()
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	if m.Code != 0 {
// 		n += 1 + sovStreams(uint64(m.Code))
// 	}
// 	l = len(m.Msg)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	l = len(m.Data)
// 	if l > 0 {
// 		n += 1 + l + sovStreams(uint64(l))
// 	}
// 	if m.XXX_unrecognized != nil {
// 		n += len(m.XXX_unrecognized)
// 	}
// 	return n
// }

// func sovStreams(x uint64) (n int) {
// 	return (math_bits.Len64(x|1) + 6) / 7
// }
// func sozStreams(x uint64) (n int) {
// 	return sovStreams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
// }
// func (m *BasicMessage) Unmarshal(dAtA []byte) error {
// 	l := len(dAtA)
// 	iNdEx := 0
// 	for iNdEx < l {
// 		preIndex := iNdEx
// 		var wire uint64
// 		for shift := uint(0); ; shift += 7 {
// 			if shift >= 64 {
// 				return ErrIntOverflowStreams
// 			}
// 			if iNdEx >= l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			b := dAtA[iNdEx]
// 			iNdEx++
// 			wire |= uint64(b&0x7F) << shift
// 			if b < 0x80 {
// 				break
// 			}
// 		}
// 		fieldNum := int32(wire >> 3)
// 		wireType := int(wire & 0x7)
// 		if wireType == 4 {
// 			return fmt.Errorf("proto: BasicMessage: wiretype end group for non-group")
// 		}
// 		if fieldNum <= 0 {
// 			return fmt.Errorf("proto: BasicMessage: illegal tag %d (wire type %d)", fieldNum, wire)
// 		}
// 		switch fieldNum {
// 		case 1:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
// 			}
// 			var stringLen uint64
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				stringLen |= uint64(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			intStringLen := int(stringLen)
// 			if intStringLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + intStringLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Type = string(dAtA[iNdEx:postIndex])
// 			iNdEx = postIndex
// 		case 2:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
// 			}
// 			var stringLen uint64
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				stringLen |= uint64(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			intStringLen := int(stringLen)
// 			if intStringLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + intStringLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Sender = string(dAtA[iNdEx:postIndex])
// 			iNdEx = postIndex
// 		case 3:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
// 			}
// 			var stringLen uint64
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				stringLen |= uint64(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			intStringLen := int(stringLen)
// 			if intStringLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + intStringLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Receiver = string(dAtA[iNdEx:postIndex])
// 			iNdEx = postIndex
// 		default:
// 			iNdEx = preIndex
// 			skippy, err := skipStreams(dAtA[iNdEx:])
// 			if err != nil {
// 				return err
// 			}
// 			if (skippy < 0) || (iNdEx+skippy) < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if (iNdEx + skippy) > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
// 			iNdEx += skippy
// 		}
// 	}

// 	if iNdEx > l {
// 		return io.ErrUnexpectedEOF
// 	}
// 	return nil
// }
// func (m *StreamRequestMessage) Unmarshal(dAtA []byte) error {
// 	l := len(dAtA)
// 	iNdEx := 0
// 	for iNdEx < l {
// 		preIndex := iNdEx
// 		var wire uint64
// 		for shift := uint(0); ; shift += 7 {
// 			if shift >= 64 {
// 				return ErrIntOverflowStreams
// 			}
// 			if iNdEx >= l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			b := dAtA[iNdEx]
// 			iNdEx++
// 			wire |= uint64(b&0x7F) << shift
// 			if b < 0x80 {
// 				break
// 			}
// 		}
// 		fieldNum := int32(wire >> 3)
// 		wireType := int(wire & 0x7)
// 		if wireType == 4 {
// 			return fmt.Errorf("proto: StreamRequestMessage: wiretype end group for non-group")
// 		}
// 		if fieldNum <= 0 {
// 			return fmt.Errorf("proto: StreamRequestMessage: illegal tag %d (wire type %d)", fieldNum, wire)
// 		}
// 		switch fieldNum {
// 		case 1:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
// 			}
// 			var msglen int
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				msglen |= int(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			if msglen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + msglen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			if m.Message == nil {
// 				m.Message = &BasicMessage{}
// 			}
// 			if err := m.Message.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
// 				return err
// 			}
// 			iNdEx = postIndex
// 		case 2:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
// 			}
// 			var byteLen int
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				byteLen |= int(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			if byteLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + byteLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
// 			if m.Payload == nil {
// 				m.Payload = []byte{}
// 			}
// 			iNdEx = postIndex
// 		default:
// 			iNdEx = preIndex
// 			skippy, err := skipStreams(dAtA[iNdEx:])
// 			if err != nil {
// 				return err
// 			}
// 			if (skippy < 0) || (iNdEx+skippy) < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if (iNdEx + skippy) > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
// 			iNdEx += skippy
// 		}
// 	}

// 	if iNdEx > l {
// 		return io.ErrUnexpectedEOF
// 	}
// 	return nil
// }
// func (m *StreamResponseMessage) Unmarshal(dAtA []byte) error {
// 	l := len(dAtA)
// 	iNdEx := 0
// 	for iNdEx < l {
// 		preIndex := iNdEx
// 		var wire uint64
// 		for shift := uint(0); ; shift += 7 {
// 			if shift >= 64 {
// 				return ErrIntOverflowStreams
// 			}
// 			if iNdEx >= l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			b := dAtA[iNdEx]
// 			iNdEx++
// 			wire |= uint64(b&0x7F) << shift
// 			if b < 0x80 {
// 				break
// 			}
// 		}
// 		fieldNum := int32(wire >> 3)
// 		wireType := int(wire & 0x7)
// 		if wireType == 4 {
// 			return fmt.Errorf("proto: StreamResponseMessage: wiretype end group for non-group")
// 		}
// 		if fieldNum <= 0 {
// 			return fmt.Errorf("proto: StreamResponseMessage: illegal tag %d (wire type %d)", fieldNum, wire)
// 		}
// 		switch fieldNum {
// 		case 1:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
// 			}
// 			var msglen int
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				msglen |= int(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			if msglen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + msglen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			if m.Message == nil {
// 				m.Message = &BasicMessage{}
// 			}
// 			if err := m.Message.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
// 				return err
// 			}
// 			iNdEx = postIndex
// 		case 2:
// 			if wireType != 0 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
// 			}
// 			m.Code = 0
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				m.Code |= int32(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 		case 3:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
// 			}
// 			var stringLen uint64
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				stringLen |= uint64(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			intStringLen := int(stringLen)
// 			if intStringLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + intStringLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Msg = string(dAtA[iNdEx:postIndex])
// 			iNdEx = postIndex
// 		case 4:
// 			if wireType != 2 {
// 				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
// 			}
// 			var byteLen int
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				byteLen |= int(b&0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			if byteLen < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			postIndex := iNdEx + byteLen
// 			if postIndex < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if postIndex > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
// 			if m.Data == nil {
// 				m.Data = []byte{}
// 			}
// 			iNdEx = postIndex
// 		default:
// 			iNdEx = preIndex
// 			skippy, err := skipStreams(dAtA[iNdEx:])
// 			if err != nil {
// 				return err
// 			}
// 			if (skippy < 0) || (iNdEx+skippy) < 0 {
// 				return ErrInvalidLengthStreams
// 			}
// 			if (iNdEx + skippy) > l {
// 				return io.ErrUnexpectedEOF
// 			}
// 			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
// 			iNdEx += skippy
// 		}
// 	}

// 	if iNdEx > l {
// 		return io.ErrUnexpectedEOF
// 	}
// 	return nil
// }
// func skipStreams(dAtA []byte) (n int, err error) {
// 	l := len(dAtA)
// 	iNdEx := 0
// 	depth := 0
// 	for iNdEx < l {
// 		var wire uint64
// 		for shift := uint(0); ; shift += 7 {
// 			if shift >= 64 {
// 				return 0, ErrIntOverflowStreams
// 			}
// 			if iNdEx >= l {
// 				return 0, io.ErrUnexpectedEOF
// 			}
// 			b := dAtA[iNdEx]
// 			iNdEx++
// 			wire |= (uint64(b) & 0x7F) << shift
// 			if b < 0x80 {
// 				break
// 			}
// 		}
// 		wireType := int(wire & 0x7)
// 		switch wireType {
// 		case 0:
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return 0, ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return 0, io.ErrUnexpectedEOF
// 				}
// 				iNdEx++
// 				if dAtA[iNdEx-1] < 0x80 {
// 					break
// 				}
// 			}
// 		case 1:
// 			iNdEx += 8
// 		case 2:
// 			var length int
// 			for shift := uint(0); ; shift += 7 {
// 				if shift >= 64 {
// 					return 0, ErrIntOverflowStreams
// 				}
// 				if iNdEx >= l {
// 					return 0, io.ErrUnexpectedEOF
// 				}
// 				b := dAtA[iNdEx]
// 				iNdEx++
// 				length |= (int(b) & 0x7F) << shift
// 				if b < 0x80 {
// 					break
// 				}
// 			}
// 			if length < 0 {
// 				return 0, ErrInvalidLengthStreams
// 			}
// 			iNdEx += length
// 		case 3:
// 			depth++
// 		case 4:
// 			if depth == 0 {
// 				return 0, ErrUnexpectedEndOfGroupStreams
// 			}
// 			depth--
// 		case 5:
// 			iNdEx += 4
// 		default:
// 			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
// 		}
// 		if iNdEx < 0 {
// 			return 0, ErrInvalidLengthStreams
// 		}
// 		if depth == 0 {
// 			return iNdEx, nil
// 		}
// 	}
// 	return 0, io.ErrUnexpectedEOF
// }

// var (
// 	ErrInvalidLengthStreams        = fmt.Errorf("proto: negative length found during unmarshaling")
// 	ErrIntOverflowStreams          = fmt.Errorf("proto: integer overflow")
// 	ErrUnexpectedEndOfGroupStreams = fmt.Errorf("proto: unexpected end of group")
// )
