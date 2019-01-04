// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_UserAuth

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

type Request struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	RealName             string   `protobuf:"bytes,2,opt,name=Real_name,json=RealName,proto3" json:"Real_name,omitempty"`
	IdCard               string   `protobuf:"bytes,3,opt,name=Id_card,json=IdCard,proto3" json:"Id_card,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{1}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *Request) GetRealName() string {
	if m != nil {
		return m.RealName
	}
	return ""
}

func (m *Request) GetIdCard() string {
	if m != nil {
		return m.IdCard
	}
	return ""
}

type Response struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *Response) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.UserAuth.Message")
	proto.RegisterType((*Request)(nil), "go.micro.srv.UserAuth.Request")
	proto.RegisterType((*Response)(nil), "go.micro.srv.UserAuth.Response")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0xad, 0x85, 0x6c, 0x33, 0x78, 0x90, 0x01, 0x35, 0x18, 0x51, 0xc9, 0xc9, 0xd3, 0x0a,
	0x7a, 0xf5, 0x22, 0x25, 0x87, 0x1c, 0x2a, 0x1a, 0xf1, 0xa6, 0x94, 0xb5, 0x3b, 0xc4, 0x42, 0x93,
	0x89, 0x33, 0xa9, 0xe8, 0xbf, 0x97, 0xa6, 0x1b, 0xbc, 0xd8, 0xd3, 0xee, 0x7b, 0x6f, 0x78, 0xcc,
	0x37, 0x90, 0xb6, 0xc2, 0x1d, 0x5f, 0xd3, 0xb7, 0xab, 0xdb, 0x15, 0x0d, 0xaf, 0xed, 0x5d, 0x3c,
	0xaa, 0xd8, 0xd6, 0xcb, 0x85, 0xb0, 0x55, 0xf9, 0xb2, 0x2f, 0x4a, 0x72, 0xbf, 0xee, 0x3e, 0xb2,
	0x14, 0xcc, 0x8c, 0x54, 0x5d, 0x45, 0x78, 0x08, 0x63, 0x75, 0x3f, 0xc9, 0xe8, 0x72, 0x74, 0x15,
	0x97, 0x9b, 0x6f, 0xf6, 0x06, 0xa6, 0xa4, 0xcf, 0x35, 0x69, 0x87, 0x67, 0x10, 0x3f, 0x93, 0xea,
	0x92, 0x9b, 0xc2, 0x87, 0x91, 0x3f, 0x03, 0x53, 0x88, 0x4b, 0x72, 0xab, 0x79, 0xe3, 0x6a, 0x4a,
	0xf6, 0xfb, 0x74, 0xb2, 0x31, 0x1e, 0x5c, 0x4d, 0x78, 0x02, 0xa6, 0xf0, 0xf3, 0x85, 0x13, 0x9f,
	0x8c, 0xfb, 0x28, 0x2a, 0xfc, 0xd4, 0x89, 0xcf, 0xee, 0x60, 0x52, 0x92, 0xb6, 0xdc, 0x28, 0x61,
	0x02, 0x26, 0x17, 0x99, 0xb2, 0xa7, 0xd0, 0x3e, 0x48, 0x3c, 0x86, 0x28, 0x17, 0x99, 0x69, 0x15,
	0x8a, 0x83, 0xba, 0x79, 0x05, 0x93, 0x6f, 0x09, 0xf1, 0x09, 0x0e, 0x1e, 0x59, 0xbb, 0x01, 0x0a,
	0xcf, 0xed, 0xbf, 0xb0, 0x36, 0xc0, 0x9c, 0x5e, 0xec, 0xcc, 0xb7, 0xdb, 0x64, 0x7b, 0xef, 0x51,
	0x7f, 0xb5, 0xdb, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x91, 0x91, 0x34, 0x54, 0x01, 0x00,
	0x00,
}
