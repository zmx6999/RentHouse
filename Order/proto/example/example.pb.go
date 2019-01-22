// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_Order

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

type AddRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}
func (*AddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{1}
}

func (m *AddRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddRequest.Unmarshal(m, b)
}
func (m *AddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddRequest.Marshal(b, m, deterministic)
}
func (m *AddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddRequest.Merge(m, src)
}
func (m *AddRequest) XXX_Size() int {
	return xxx_messageInfo_AddRequest.Size(m)
}
func (m *AddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddRequest proto.InternalMessageInfo

func (m *AddRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *AddRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type AddResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddResponse) Reset()         { *m = AddResponse{} }
func (m *AddResponse) String() string { return proto.CompactTextString(m) }
func (*AddResponse) ProtoMessage()    {}
func (*AddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{2}
}

func (m *AddResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddResponse.Unmarshal(m, b)
}
func (m *AddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddResponse.Marshal(b, m, deterministic)
}
func (m *AddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddResponse.Merge(m, src)
}
func (m *AddResponse) XXX_Size() int {
	return xxx_messageInfo_AddResponse.Size(m)
}
func (m *AddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddResponse proto.InternalMessageInfo

func (m *AddResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *AddResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type GetOrdersRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	Role                 string   `protobuf:"bytes,2,opt,name=Role,proto3" json:"Role,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetOrdersRequest) Reset()         { *m = GetOrdersRequest{} }
func (m *GetOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*GetOrdersRequest) ProtoMessage()    {}
func (*GetOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{3}
}

func (m *GetOrdersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetOrdersRequest.Unmarshal(m, b)
}
func (m *GetOrdersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetOrdersRequest.Marshal(b, m, deterministic)
}
func (m *GetOrdersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetOrdersRequest.Merge(m, src)
}
func (m *GetOrdersRequest) XXX_Size() int {
	return xxx_messageInfo_GetOrdersRequest.Size(m)
}
func (m *GetOrdersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetOrdersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetOrdersRequest proto.InternalMessageInfo

func (m *GetOrdersRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *GetOrdersRequest) GetRole() string {
	if m != nil {
		return m.Role
	}
	return ""
}

type GetOrdersResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetOrdersResponse) Reset()         { *m = GetOrdersResponse{} }
func (m *GetOrdersResponse) String() string { return proto.CompactTextString(m) }
func (*GetOrdersResponse) ProtoMessage()    {}
func (*GetOrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{4}
}

func (m *GetOrdersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetOrdersResponse.Unmarshal(m, b)
}
func (m *GetOrdersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetOrdersResponse.Marshal(b, m, deterministic)
}
func (m *GetOrdersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetOrdersResponse.Merge(m, src)
}
func (m *GetOrdersResponse) XXX_Size() int {
	return xxx_messageInfo_GetOrdersResponse.Size(m)
}
func (m *GetOrdersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetOrdersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetOrdersResponse proto.InternalMessageInfo

func (m *GetOrdersResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetOrdersResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *GetOrdersResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type HandleRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	OrderId              int64    `protobuf:"varint,2,opt,name=OrderId,proto3" json:"OrderId,omitempty"`
	Action               string   `protobuf:"bytes,3,opt,name=Action,proto3" json:"Action,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HandleRequest) Reset()         { *m = HandleRequest{} }
func (m *HandleRequest) String() string { return proto.CompactTextString(m) }
func (*HandleRequest) ProtoMessage()    {}
func (*HandleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{5}
}

func (m *HandleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HandleRequest.Unmarshal(m, b)
}
func (m *HandleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HandleRequest.Marshal(b, m, deterministic)
}
func (m *HandleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HandleRequest.Merge(m, src)
}
func (m *HandleRequest) XXX_Size() int {
	return xxx_messageInfo_HandleRequest.Size(m)
}
func (m *HandleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HandleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HandleRequest proto.InternalMessageInfo

func (m *HandleRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *HandleRequest) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *HandleRequest) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

type HandleResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HandleResponse) Reset()         { *m = HandleResponse{} }
func (m *HandleResponse) String() string { return proto.CompactTextString(m) }
func (*HandleResponse) ProtoMessage()    {}
func (*HandleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{6}
}

func (m *HandleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HandleResponse.Unmarshal(m, b)
}
func (m *HandleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HandleResponse.Marshal(b, m, deterministic)
}
func (m *HandleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HandleResponse.Merge(m, src)
}
func (m *HandleResponse) XXX_Size() int {
	return xxx_messageInfo_HandleResponse.Size(m)
}
func (m *HandleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HandleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HandleResponse proto.InternalMessageInfo

func (m *HandleResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *HandleResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type CommentRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	OrderId              int64    `protobuf:"varint,2,opt,name=OrderId,proto3" json:"OrderId,omitempty"`
	Comment              string   `protobuf:"bytes,3,opt,name=Comment,proto3" json:"Comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommentRequest) Reset()         { *m = CommentRequest{} }
func (m *CommentRequest) String() string { return proto.CompactTextString(m) }
func (*CommentRequest) ProtoMessage()    {}
func (*CommentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{7}
}

func (m *CommentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommentRequest.Unmarshal(m, b)
}
func (m *CommentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommentRequest.Marshal(b, m, deterministic)
}
func (m *CommentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommentRequest.Merge(m, src)
}
func (m *CommentRequest) XXX_Size() int {
	return xxx_messageInfo_CommentRequest.Size(m)
}
func (m *CommentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommentRequest proto.InternalMessageInfo

func (m *CommentRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *CommentRequest) GetOrderId() int64 {
	if m != nil {
		return m.OrderId
	}
	return 0
}

func (m *CommentRequest) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

type CommentResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommentResponse) Reset()         { *m = CommentResponse{} }
func (m *CommentResponse) String() string { return proto.CompactTextString(m) }
func (*CommentResponse) ProtoMessage()    {}
func (*CommentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{8}
}

func (m *CommentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommentResponse.Unmarshal(m, b)
}
func (m *CommentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommentResponse.Marshal(b, m, deterministic)
}
func (m *CommentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommentResponse.Merge(m, src)
}
func (m *CommentResponse) XXX_Size() int {
	return xxx_messageInfo_CommentResponse.Size(m)
}
func (m *CommentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CommentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CommentResponse proto.InternalMessageInfo

func (m *CommentResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *CommentResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.Order.Message")
	proto.RegisterType((*AddRequest)(nil), "go.micro.srv.Order.AddRequest")
	proto.RegisterType((*AddResponse)(nil), "go.micro.srv.Order.AddResponse")
	proto.RegisterType((*GetOrdersRequest)(nil), "go.micro.srv.Order.GetOrdersRequest")
	proto.RegisterType((*GetOrdersResponse)(nil), "go.micro.srv.Order.GetOrdersResponse")
	proto.RegisterType((*HandleRequest)(nil), "go.micro.srv.Order.HandleRequest")
	proto.RegisterType((*HandleResponse)(nil), "go.micro.srv.Order.HandleResponse")
	proto.RegisterType((*CommentRequest)(nil), "go.micro.srv.Order.CommentRequest")
	proto.RegisterType((*CommentResponse)(nil), "go.micro.srv.Order.CommentResponse")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 369 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x41, 0x4f, 0xf2, 0x40,
	0x10, 0xfd, 0xa0, 0x5f, 0x68, 0x3a, 0x2a, 0xe2, 0x1c, 0x4c, 0x53, 0x8c, 0xe2, 0xaa, 0x09, 0xa7,
	0x9a, 0xe8, 0x5d, 0x83, 0x40, 0x84, 0x03, 0x31, 0x56, 0x2f, 0x7a, 0x31, 0x85, 0x9d, 0x10, 0x12,
	0xda, 0xc5, 0xdd, 0x6a, 0xf4, 0x9f, 0xfa, 0x73, 0x0c, 0xcb, 0xb6, 0x80, 0x22, 0x12, 0x3c, 0x31,
	0xb3, 0xbc, 0xf7, 0xe6, 0xed, 0xce, 0x4b, 0xa1, 0x3c, 0x92, 0x22, 0x11, 0xa7, 0xf4, 0x16, 0x46,
	0xa3, 0x21, 0xa5, 0xbf, 0xbe, 0x3e, 0x45, 0xec, 0x0b, 0x3f, 0x1a, 0xf4, 0xa4, 0xf0, 0x95, 0x7c,
	0xf5, 0x6f, 0x24, 0x27, 0xc9, 0xca, 0x60, 0x77, 0x48, 0xa9, 0xb0, 0x4f, 0x58, 0x02, 0x4b, 0x85,
	0xef, 0x6e, 0xae, 0x92, 0xab, 0x3a, 0xc1, 0xb8, 0x64, 0x17, 0x00, 0x35, 0xce, 0x03, 0x7a, 0x7e,
	0x21, 0x95, 0xe0, 0x1e, 0x38, 0x77, 0xa4, 0xd4, 0x40, 0xc4, 0x6d, 0x6e, 0x50, 0xd3, 0x03, 0x44,
	0xf8, 0xdf, 0x08, 0x93, 0xd0, 0xcd, 0x57, 0x72, 0xd5, 0xcd, 0x40, 0xd7, 0xec, 0x12, 0x36, 0x34,
	0x5f, 0x8d, 0x44, 0xac, 0x08, 0x5d, 0xb0, 0x9b, 0x52, 0xd6, 0x05, 0x27, 0x43, 0x4f, 0x5b, 0xdc,
	0x85, 0x42, 0x53, 0xca, 0x8e, 0xea, 0x6b, 0xba, 0x13, 0x98, 0x8e, 0x35, 0xa0, 0x74, 0x4d, 0x89,
	0x76, 0xaa, 0x56, 0xb6, 0x11, 0x88, 0x21, 0x19, 0x1d, 0x5d, 0xb3, 0x07, 0xd8, 0x99, 0x51, 0x59,
	0xd7, 0x4c, 0x76, 0x43, 0x6b, 0xe6, 0x86, 0x4f, 0xb0, 0xd5, 0x0a, 0x63, 0x3e, 0xa4, 0xd5, 0xdc,
	0xb9, 0x60, 0x6b, 0x1b, 0x6d, 0xae, 0xb5, 0xad, 0x20, 0x6d, 0xc7, 0x43, 0x6b, 0xbd, 0x64, 0x20,
	0x62, 0x2d, 0xef, 0x04, 0xa6, 0x63, 0x57, 0x50, 0x4c, 0x07, 0xac, 0xfd, 0x8a, 0x5d, 0x28, 0xd6,
	0x45, 0x14, 0x51, 0x9c, 0xfc, 0xd5, 0xa5, 0x0b, 0xb6, 0x51, 0x32, 0x36, 0xd3, 0x96, 0xd5, 0x61,
	0x3b, 0x9b, 0xb1, 0xae, 0xd1, 0xb3, 0x8f, 0x3c, 0xd8, 0xcd, 0x49, 0x64, 0xb1, 0x05, 0x56, 0x8d,
	0x73, 0xdc, 0xf7, 0xbf, 0x87, 0xd6, 0x9f, 0x86, 0xd2, 0x3b, 0xf8, 0xf1, 0xff, 0x89, 0x0b, 0xf6,
	0x0f, 0x1f, 0xc1, 0xc9, 0xd6, 0x8f, 0xc7, 0x8b, 0xf0, 0x5f, 0x33, 0xe6, 0x9d, 0xfc, 0x82, 0xca,
	0xb4, 0x6f, 0xa1, 0x30, 0x59, 0x0f, 0x1e, 0x2e, 0xa2, 0xcc, 0x65, 0xc3, 0x63, 0xcb, 0x20, 0x99,
	0xe4, 0x7d, 0xf6, 0xc6, 0xb8, 0x90, 0x30, 0xbf, 0x4a, 0xef, 0x68, 0x29, 0x26, 0x55, 0xed, 0x16,
	0xf4, 0x27, 0xe0, 0xfc, 0x33, 0x00, 0x00, 0xff, 0xff, 0x48, 0xe7, 0xc1, 0x1f, 0x21, 0x04, 0x00,
	0x00,
}
