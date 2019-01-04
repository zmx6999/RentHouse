// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_GetHouseInfo

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
	HouseId              int64    `protobuf:"varint,2,opt,name=HouseId,proto3" json:"HouseId,omitempty"`
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

func (m *Request) GetHouseId() int64 {
	if m != nil {
		return m.HouseId
	}
	return 0
}

type Response struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	UserId               int64    `protobuf:"varint,3,opt,name=UserId,proto3" json:"UserId,omitempty"`
	Data                 []byte   `protobuf:"bytes,4,opt,name=Data,proto3" json:"Data,omitempty"`
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

func (m *Response) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *Response) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.GetHouseInfo.Message")
	proto.RegisterType((*Request)(nil), "go.micro.srv.GetHouseInfo.Request")
	proto.RegisterType((*Response)(nil), "go.micro.srv.GetHouseInfo.Response")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 241 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0xb1, 0x4e, 0xc4, 0x30,
	0x10, 0x44, 0x09, 0x39, 0x5d, 0xc8, 0xea, 0x0a, 0xe4, 0x02, 0x19, 0x8e, 0x22, 0x32, 0x4d, 0x2a,
	0x23, 0xc1, 0x17, 0x20, 0x88, 0x20, 0xc5, 0x35, 0x46, 0x14, 0x94, 0x06, 0x2f, 0xe1, 0x24, 0x2e,
	0x1b, 0xbc, 0x09, 0x82, 0xbf, 0x47, 0x71, 0x1c, 0x41, 0x03, 0x95, 0x77, 0x46, 0xe3, 0x91, 0xde,
	0xc0, 0xba, 0xf3, 0xd4, 0xd3, 0x39, 0x7e, 0xda, 0x5d, 0xf7, 0x86, 0xf3, 0xab, 0x83, 0x2b, 0x8e,
	0x1b, 0xd2, 0xbb, 0xed, 0xb3, 0x27, 0xcd, 0xfe, 0x43, 0xdf, 0x62, 0x7f, 0x47, 0x03, 0x63, 0xdd,
	0xbe, 0x90, 0x5a, 0x43, 0xb6, 0x41, 0x66, 0xdb, 0xa0, 0x38, 0x84, 0x94, 0xed, 0x97, 0x4c, 0x8a,
	0xa4, 0xcc, 0xcd, 0x78, 0xaa, 0x2b, 0xc8, 0x0c, 0xbe, 0x0f, 0xc8, 0xbd, 0x38, 0x85, 0xfc, 0x1e,
	0x99, 0xb7, 0xd4, 0xd6, 0x2e, 0x46, 0x7e, 0x0c, 0x21, 0x21, 0x9b, 0x2a, 0x9d, 0xdc, 0x2f, 0x92,
	0x32, 0x35, 0xb3, 0x54, 0xaf, 0x70, 0x60, 0x90, 0x3b, 0x6a, 0x19, 0xc7, 0x54, 0xe5, 0xfd, 0x35,
	0x39, 0x8c, 0x0d, 0xb3, 0x14, 0x47, 0xb0, 0xac, 0xbc, 0xdf, 0x70, 0x13, 0xbe, 0xe7, 0x26, 0xaa,
	0xd1, 0x7f, 0x60, 0xf4, 0xb5, 0x93, 0x69, 0xa8, 0x8d, 0x4a, 0x08, 0x58, 0xdc, 0xd8, 0xde, 0xca,
	0x45, 0x91, 0x94, 0x2b, 0x13, 0xee, 0x0b, 0x07, 0x59, 0x35, 0x51, 0x8b, 0x47, 0x58, 0xfd, 0x86,
	0x14, 0x4a, 0xff, 0x39, 0x80, 0x8e, 0x80, 0x27, 0x67, 0xff, 0x66, 0x26, 0x02, 0xb5, 0xf7, 0xb4,
	0x0c, 0x8b, 0x5e, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x36, 0x43, 0xec, 0xe2, 0x70, 0x01, 0x00,
	0x00,
}