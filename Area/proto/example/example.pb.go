// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_Area

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Message struct {
	Say                  string   `protobuf:"bytes,1,opt,name=say,proto3" json:"say,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetSay() string {
	if m != nil {
		return m.Say
	}
	return ""
}

type GetAreaRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAreaRequest) Reset()         { *m = GetAreaRequest{} }
func (m *GetAreaRequest) String() string { return proto.CompactTextString(m) }
func (*GetAreaRequest) ProtoMessage()    {}
func (*GetAreaRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{1}
}

func (m *GetAreaRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAreaRequest.Unmarshal(m, b)
}
func (m *GetAreaRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAreaRequest.Marshal(b, m, deterministic)
}
func (m *GetAreaRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAreaRequest.Merge(m, src)
}
func (m *GetAreaRequest) XXX_Size() int {
	return xxx_messageInfo_GetAreaRequest.Size(m)
}
func (m *GetAreaRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAreaRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAreaRequest proto.InternalMessageInfo

type GetAreaResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAreaResponse) Reset()         { *m = GetAreaResponse{} }
func (m *GetAreaResponse) String() string { return proto.CompactTextString(m) }
func (*GetAreaResponse) ProtoMessage()    {}
func (*GetAreaResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{2}
}

func (m *GetAreaResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAreaResponse.Unmarshal(m, b)
}
func (m *GetAreaResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAreaResponse.Marshal(b, m, deterministic)
}
func (m *GetAreaResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAreaResponse.Merge(m, src)
}
func (m *GetAreaResponse) XXX_Size() int {
	return xxx_messageInfo_GetAreaResponse.Size(m)
}
func (m *GetAreaResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAreaResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAreaResponse proto.InternalMessageInfo

func (m *GetAreaResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetAreaResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *GetAreaResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.Area.Message")
	proto.RegisterType((*GetAreaRequest)(nil), "go.micro.srv.Area.GetAreaRequest")
	proto.RegisterType((*GetAreaResponse)(nil), "go.micro.srv.Area.GetAreaResponse")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2e, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x85, 0xd1, 0x7a, 0x60, 0x51, 0x21, 0xc1,
	0xf4, 0x7c, 0xbd, 0xdc, 0xcc, 0xe4, 0xa2, 0x7c, 0xbd, 0xe2, 0xa2, 0x32, 0x3d, 0xc7, 0xa2, 0xd4,
	0x44, 0x25, 0x69, 0x2e, 0x76, 0xdf, 0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0x21, 0x01, 0x2e, 0xe6,
	0xe2, 0xc4, 0x4a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x10, 0x53, 0x49, 0x80, 0x8b, 0xcf,
	0x3d, 0xb5, 0x04, 0xa4, 0x2e, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x29, 0x90, 0x8b, 0x1f,
	0x2e, 0x52, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x2a, 0x24, 0xc1, 0xc5, 0xee, 0x5a, 0x54, 0xe4, 0x9c,
	0x9f, 0x92, 0x0a, 0xd5, 0x0a, 0xe3, 0x82, 0x0c, 0xf4, 0x2d, 0x4e, 0x97, 0x60, 0x82, 0x18, 0xe8,
	0x5b, 0x9c, 0x2e, 0x24, 0xc4, 0xc5, 0xe2, 0x92, 0x58, 0x92, 0x28, 0xc1, 0xac, 0xc0, 0xa8, 0xc1,
	0x13, 0x04, 0x66, 0x1b, 0xc5, 0x72, 0xb1, 0xbb, 0x42, 0x5c, 0x29, 0x14, 0xc4, 0xc5, 0x0e, 0x35,
	0x5d, 0x48, 0x51, 0x0f, 0xc3, 0xad, 0x7a, 0xa8, 0x6e, 0x91, 0x52, 0xc2, 0xa7, 0x04, 0xe2, 0x38,
	0x25, 0x86, 0x24, 0x36, 0xb0, 0xd7, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xfa, 0xc4, 0x4b,
	0x00, 0x19, 0x01, 0x00, 0x00,
}
