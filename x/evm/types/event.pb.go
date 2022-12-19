// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ethermint/evm/v1/event.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MergeAccountEvent struct {
	EthAddress          string `protobuf:"bytes,1,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty" yaml:"eth_address"`
	CosmosAddress       string `protobuf:"bytes,2,opt,name=cosmos_address,json=cosmosAddress,proto3" json:"cosmos_address,omitempty" yaml:"cosmos_address"`
	NewCosmosAccCreated bool   `protobuf:"varint,3,opt,name=new_cosmos_acc_created,json=newCosmosAccCreated,proto3" json:"new_cosmos_acc_created,omitempty" yaml:"new_cosmos_acc_created"`
}

func (m *MergeAccountEvent) Reset()         { *m = MergeAccountEvent{} }
func (m *MergeAccountEvent) String() string { return proto.CompactTextString(m) }
func (*MergeAccountEvent) ProtoMessage()    {}
func (*MergeAccountEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_1648308e5d6fa622, []int{0}
}
func (m *MergeAccountEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MergeAccountEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MergeAccountEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MergeAccountEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MergeAccountEvent.Merge(m, src)
}
func (m *MergeAccountEvent) XXX_Size() int {
	return m.Size()
}
func (m *MergeAccountEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_MergeAccountEvent.DiscardUnknown(m)
}

var xxx_messageInfo_MergeAccountEvent proto.InternalMessageInfo

func (m *MergeAccountEvent) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

func (m *MergeAccountEvent) GetCosmosAddress() string {
	if m != nil {
		return m.CosmosAddress
	}
	return ""
}

func (m *MergeAccountEvent) GetNewCosmosAccCreated() bool {
	if m != nil {
		return m.NewCosmosAccCreated
	}
	return false
}

func init() {
	proto.RegisterType((*MergeAccountEvent)(nil), "ethermint.evm.v1.MergeAccountEvent")
}

func init() { proto.RegisterFile("ethermint/evm/v1/event.proto", fileDescriptor_1648308e5d6fa622) }

var fileDescriptor_1648308e5d6fa622 = []byte{
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x49, 0x2d, 0xc9, 0x48,
	0x2d, 0xca, 0xcd, 0xcc, 0x2b, 0xd1, 0x4f, 0x2d, 0xcb, 0xd5, 0x2f, 0x33, 0xd4, 0x4f, 0x2d, 0x4b,
	0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x80, 0xcb, 0xea, 0xa5, 0x96, 0xe5,
	0xea, 0x95, 0x19, 0x4a, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x25, 0xf5, 0x41, 0x2c, 0x88, 0x3a,
	0xa5, 0x67, 0x8c, 0x5c, 0x82, 0xbe, 0xa9, 0x45, 0xe9, 0xa9, 0x8e, 0xc9, 0xc9, 0xf9, 0xa5, 0x79,
	0x25, 0xae, 0x20, 0x33, 0x84, 0xcc, 0xb9, 0xb8, 0x53, 0x4b, 0x32, 0xe2, 0x13, 0x53, 0x52, 0x8a,
	0x52, 0x8b, 0x8b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x9d, 0xc4, 0x3e, 0xdd, 0x93, 0x17, 0xaa,
	0x4c, 0xcc, 0xcd, 0xb1, 0x52, 0x42, 0x92, 0x54, 0x0a, 0xe2, 0x4a, 0x2d, 0xc9, 0x70, 0x84, 0x70,
	0x84, 0x1c, 0xb8, 0xf8, 0x92, 0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xe1, 0x7a, 0x99, 0xc0, 0x7a, 0x25,
	0x3f, 0xdd, 0x93, 0x17, 0x85, 0xe8, 0x45, 0x95, 0x57, 0x0a, 0xe2, 0x85, 0x08, 0xc0, 0x4c, 0x08,
	0xe3, 0x12, 0xcb, 0x4b, 0x2d, 0x8f, 0x87, 0xa9, 0x4a, 0x4e, 0x8e, 0x4f, 0x2e, 0x4a, 0x4d, 0x2c,
	0x49, 0x4d, 0x91, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x70, 0x52, 0xfc, 0x74, 0x4f, 0x5e, 0x16, 0x62,
	0x12, 0x76, 0x75, 0x4a, 0x41, 0xc2, 0x79, 0xa9, 0xe5, 0xce, 0x10, 0x43, 0x93, 0x93, 0x9d, 0x21,
	0xa2, 0x4e, 0x0e, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3,
	0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10, 0xa5, 0x96, 0x9e,
	0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0x0b, 0x0a, 0xc9, 0xfc, 0x62, 0x7d, 0x44, 0xc8,
	0x56, 0x80, 0xc3, 0xb6, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0x1c, 0x62, 0xc6, 0x80, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x80, 0x94, 0xe9, 0x20, 0x79, 0x01, 0x00, 0x00,
}

func (m *MergeAccountEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MergeAccountEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MergeAccountEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NewCosmosAccCreated {
		i--
		if m.NewCosmosAccCreated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.CosmosAddress) > 0 {
		i -= len(m.CosmosAddress)
		copy(dAtA[i:], m.CosmosAddress)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.CosmosAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.EthAddress) > 0 {
		i -= len(m.EthAddress)
		copy(dAtA[i:], m.EthAddress)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.EthAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MergeAccountEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EthAddress)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.CosmosAddress)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.NewCosmosAccCreated {
		n += 2
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MergeAccountEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MergeAccountEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MergeAccountEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CosmosAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CosmosAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewCosmosAccCreated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.NewCosmosAccCreated = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
