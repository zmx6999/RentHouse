// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_User

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

type GetImageCptRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=Uuid,proto3" json:"Uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetImageCptRequest) Reset()         { *m = GetImageCptRequest{} }
func (m *GetImageCptRequest) String() string { return proto.CompactTextString(m) }
func (*GetImageCptRequest) ProtoMessage()    {}
func (*GetImageCptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{1}
}

func (m *GetImageCptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetImageCptRequest.Unmarshal(m, b)
}
func (m *GetImageCptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetImageCptRequest.Marshal(b, m, deterministic)
}
func (m *GetImageCptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetImageCptRequest.Merge(m, src)
}
func (m *GetImageCptRequest) XXX_Size() int {
	return xxx_messageInfo_GetImageCptRequest.Size(m)
}
func (m *GetImageCptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetImageCptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetImageCptRequest proto.InternalMessageInfo

func (m *GetImageCptRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type GetImageCptResponse struct {
	ErrCode              string                     `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string                     `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Pix                  []byte                     `protobuf:"bytes,3,opt,name=Pix,proto3" json:"Pix,omitempty"`
	Stride               int64                      `protobuf:"varint,4,opt,name=Stride,proto3" json:"Stride,omitempty"`
	Max                  *GetImageCptResponse_Point `protobuf:"bytes,5,opt,name=Max,proto3" json:"Max,omitempty"`
	Min                  *GetImageCptResponse_Point `protobuf:"bytes,6,opt,name=Min,proto3" json:"Min,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *GetImageCptResponse) Reset()         { *m = GetImageCptResponse{} }
func (m *GetImageCptResponse) String() string { return proto.CompactTextString(m) }
func (*GetImageCptResponse) ProtoMessage()    {}
func (*GetImageCptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{2}
}

func (m *GetImageCptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetImageCptResponse.Unmarshal(m, b)
}
func (m *GetImageCptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetImageCptResponse.Marshal(b, m, deterministic)
}
func (m *GetImageCptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetImageCptResponse.Merge(m, src)
}
func (m *GetImageCptResponse) XXX_Size() int {
	return xxx_messageInfo_GetImageCptResponse.Size(m)
}
func (m *GetImageCptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetImageCptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetImageCptResponse proto.InternalMessageInfo

func (m *GetImageCptResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetImageCptResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *GetImageCptResponse) GetPix() []byte {
	if m != nil {
		return m.Pix
	}
	return nil
}

func (m *GetImageCptResponse) GetStride() int64 {
	if m != nil {
		return m.Stride
	}
	return 0
}

func (m *GetImageCptResponse) GetMax() *GetImageCptResponse_Point {
	if m != nil {
		return m.Max
	}
	return nil
}

func (m *GetImageCptResponse) GetMin() *GetImageCptResponse_Point {
	if m != nil {
		return m.Min
	}
	return nil
}

type GetImageCptResponse_Point struct {
	X                    int64    `protobuf:"varint,1,opt,name=X,proto3" json:"X,omitempty"`
	Y                    int64    `protobuf:"varint,2,opt,name=Y,proto3" json:"Y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetImageCptResponse_Point) Reset()         { *m = GetImageCptResponse_Point{} }
func (m *GetImageCptResponse_Point) String() string { return proto.CompactTextString(m) }
func (*GetImageCptResponse_Point) ProtoMessage()    {}
func (*GetImageCptResponse_Point) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{2, 0}
}

func (m *GetImageCptResponse_Point) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetImageCptResponse_Point.Unmarshal(m, b)
}
func (m *GetImageCptResponse_Point) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetImageCptResponse_Point.Marshal(b, m, deterministic)
}
func (m *GetImageCptResponse_Point) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetImageCptResponse_Point.Merge(m, src)
}
func (m *GetImageCptResponse_Point) XXX_Size() int {
	return xxx_messageInfo_GetImageCptResponse_Point.Size(m)
}
func (m *GetImageCptResponse_Point) XXX_DiscardUnknown() {
	xxx_messageInfo_GetImageCptResponse_Point.DiscardUnknown(m)
}

var xxx_messageInfo_GetImageCptResponse_Point proto.InternalMessageInfo

func (m *GetImageCptResponse_Point) GetX() int64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *GetImageCptResponse_Point) GetY() int64 {
	if m != nil {
		return m.Y
	}
	return 0
}

type GetSmsCptRequest struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=Uuid,proto3" json:"Uuid,omitempty"`
	Mobile               string   `protobuf:"bytes,2,opt,name=Mobile,proto3" json:"Mobile,omitempty"`
	Text                 string   `protobuf:"bytes,3,opt,name=Text,proto3" json:"Text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSmsCptRequest) Reset()         { *m = GetSmsCptRequest{} }
func (m *GetSmsCptRequest) String() string { return proto.CompactTextString(m) }
func (*GetSmsCptRequest) ProtoMessage()    {}
func (*GetSmsCptRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{3}
}

func (m *GetSmsCptRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSmsCptRequest.Unmarshal(m, b)
}
func (m *GetSmsCptRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSmsCptRequest.Marshal(b, m, deterministic)
}
func (m *GetSmsCptRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSmsCptRequest.Merge(m, src)
}
func (m *GetSmsCptRequest) XXX_Size() int {
	return xxx_messageInfo_GetSmsCptRequest.Size(m)
}
func (m *GetSmsCptRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSmsCptRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSmsCptRequest proto.InternalMessageInfo

func (m *GetSmsCptRequest) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *GetSmsCptRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *GetSmsCptRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type GetSmsCptResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSmsCptResponse) Reset()         { *m = GetSmsCptResponse{} }
func (m *GetSmsCptResponse) String() string { return proto.CompactTextString(m) }
func (*GetSmsCptResponse) ProtoMessage()    {}
func (*GetSmsCptResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{4}
}

func (m *GetSmsCptResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSmsCptResponse.Unmarshal(m, b)
}
func (m *GetSmsCptResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSmsCptResponse.Marshal(b, m, deterministic)
}
func (m *GetSmsCptResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSmsCptResponse.Merge(m, src)
}
func (m *GetSmsCptResponse) XXX_Size() int {
	return xxx_messageInfo_GetSmsCptResponse.Size(m)
}
func (m *GetSmsCptResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSmsCptResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSmsCptResponse proto.InternalMessageInfo

func (m *GetSmsCptResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetSmsCptResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type RegisterRequest struct {
	Mobile               string   `protobuf:"bytes,1,opt,name=Mobile,proto3" json:"Mobile,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	Text                 string   `protobuf:"bytes,3,opt,name=Text,proto3" json:"Text,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{5}
}

func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(m, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *RegisterRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *RegisterRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type RegisterResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	SessionId            string   `protobuf:"bytes,3,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{6}
}

func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (m *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(m, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *RegisterResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *RegisterResponse) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type LoginRequest struct {
	Mobile               string   `protobuf:"bytes,1,opt,name=Mobile,proto3" json:"Mobile,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=Password,proto3" json:"Password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{7}
}

func (m *LoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginRequest.Unmarshal(m, b)
}
func (m *LoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginRequest.Marshal(b, m, deterministic)
}
func (m *LoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginRequest.Merge(m, src)
}
func (m *LoginRequest) XXX_Size() int {
	return xxx_messageInfo_LoginRequest.Size(m)
}
func (m *LoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *LoginRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type LoginResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	SessionId            string   `protobuf:"bytes,3,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{8}
}

func (m *LoginResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResponse.Unmarshal(m, b)
}
func (m *LoginResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResponse.Marshal(b, m, deterministic)
}
func (m *LoginResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResponse.Merge(m, src)
}
func (m *LoginResponse) XXX_Size() int {
	return xxx_messageInfo_LoginResponse.Size(m)
}
func (m *LoginResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResponse proto.InternalMessageInfo

func (m *LoginResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *LoginResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *LoginResponse) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type LogoutRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutRequest) Reset()         { *m = LogoutRequest{} }
func (m *LogoutRequest) String() string { return proto.CompactTextString(m) }
func (*LogoutRequest) ProtoMessage()    {}
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{9}
}

func (m *LogoutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutRequest.Unmarshal(m, b)
}
func (m *LogoutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutRequest.Marshal(b, m, deterministic)
}
func (m *LogoutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutRequest.Merge(m, src)
}
func (m *LogoutRequest) XXX_Size() int {
	return xxx_messageInfo_LogoutRequest.Size(m)
}
func (m *LogoutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutRequest proto.InternalMessageInfo

func (m *LogoutRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type LogoutResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutResponse) Reset()         { *m = LogoutResponse{} }
func (m *LogoutResponse) String() string { return proto.CompactTextString(m) }
func (*LogoutResponse) ProtoMessage()    {}
func (*LogoutResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{10}
}

func (m *LogoutResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutResponse.Unmarshal(m, b)
}
func (m *LogoutResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutResponse.Marshal(b, m, deterministic)
}
func (m *LogoutResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutResponse.Merge(m, src)
}
func (m *LogoutResponse) XXX_Size() int {
	return xxx_messageInfo_LogoutResponse.Size(m)
}
func (m *LogoutResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutResponse proto.InternalMessageInfo

func (m *LogoutResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *LogoutResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type GetUserInfoRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserInfoRequest) Reset()         { *m = GetUserInfoRequest{} }
func (m *GetUserInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetUserInfoRequest) ProtoMessage()    {}
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{11}
}

func (m *GetUserInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserInfoRequest.Unmarshal(m, b)
}
func (m *GetUserInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserInfoRequest.Marshal(b, m, deterministic)
}
func (m *GetUserInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserInfoRequest.Merge(m, src)
}
func (m *GetUserInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetUserInfoRequest.Size(m)
}
func (m *GetUserInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserInfoRequest proto.InternalMessageInfo

func (m *GetUserInfoRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type GetUserInfoResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserInfoResponse) Reset()         { *m = GetUserInfoResponse{} }
func (m *GetUserInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetUserInfoResponse) ProtoMessage()    {}
func (*GetUserInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{12}
}

func (m *GetUserInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserInfoResponse.Unmarshal(m, b)
}
func (m *GetUserInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserInfoResponse.Marshal(b, m, deterministic)
}
func (m *GetUserInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserInfoResponse.Merge(m, src)
}
func (m *GetUserInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetUserInfoResponse.Size(m)
}
func (m *GetUserInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserInfoResponse proto.InternalMessageInfo

func (m *GetUserInfoResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetUserInfoResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *GetUserInfoResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type RenameRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	NewName              string   `protobuf:"bytes,2,opt,name=NewName,proto3" json:"NewName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RenameRequest) Reset()         { *m = RenameRequest{} }
func (m *RenameRequest) String() string { return proto.CompactTextString(m) }
func (*RenameRequest) ProtoMessage()    {}
func (*RenameRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{13}
}

func (m *RenameRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RenameRequest.Unmarshal(m, b)
}
func (m *RenameRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RenameRequest.Marshal(b, m, deterministic)
}
func (m *RenameRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RenameRequest.Merge(m, src)
}
func (m *RenameRequest) XXX_Size() int {
	return xxx_messageInfo_RenameRequest.Size(m)
}
func (m *RenameRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RenameRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RenameRequest proto.InternalMessageInfo

func (m *RenameRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *RenameRequest) GetNewName() string {
	if m != nil {
		return m.NewName
	}
	return ""
}

type RenameResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RenameResponse) Reset()         { *m = RenameResponse{} }
func (m *RenameResponse) String() string { return proto.CompactTextString(m) }
func (*RenameResponse) ProtoMessage()    {}
func (*RenameResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{14}
}

func (m *RenameResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RenameResponse.Unmarshal(m, b)
}
func (m *RenameResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RenameResponse.Marshal(b, m, deterministic)
}
func (m *RenameResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RenameResponse.Merge(m, src)
}
func (m *RenameResponse) XXX_Size() int {
	return xxx_messageInfo_RenameResponse.Size(m)
}
func (m *RenameResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RenameResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RenameResponse proto.InternalMessageInfo

func (m *RenameResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *RenameResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type AuthRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	RealName             string   `protobuf:"bytes,2,opt,name=RealName,proto3" json:"RealName,omitempty"`
	IdCard               string   `protobuf:"bytes,3,opt,name=IdCard,proto3" json:"IdCard,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthRequest) Reset()         { *m = AuthRequest{} }
func (m *AuthRequest) String() string { return proto.CompactTextString(m) }
func (*AuthRequest) ProtoMessage()    {}
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{15}
}

func (m *AuthRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthRequest.Unmarshal(m, b)
}
func (m *AuthRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthRequest.Marshal(b, m, deterministic)
}
func (m *AuthRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthRequest.Merge(m, src)
}
func (m *AuthRequest) XXX_Size() int {
	return xxx_messageInfo_AuthRequest.Size(m)
}
func (m *AuthRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AuthRequest proto.InternalMessageInfo

func (m *AuthRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *AuthRequest) GetRealName() string {
	if m != nil {
		return m.RealName
	}
	return ""
}

func (m *AuthRequest) GetIdCard() string {
	if m != nil {
		return m.IdCard
	}
	return ""
}

type AuthResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthResponse) Reset()         { *m = AuthResponse{} }
func (m *AuthResponse) String() string { return proto.CompactTextString(m) }
func (*AuthResponse) ProtoMessage()    {}
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{16}
}

func (m *AuthResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthResponse.Unmarshal(m, b)
}
func (m *AuthResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthResponse.Marshal(b, m, deterministic)
}
func (m *AuthResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthResponse.Merge(m, src)
}
func (m *AuthResponse) XXX_Size() int {
	return xxx_messageInfo_AuthResponse.Size(m)
}
func (m *AuthResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AuthResponse proto.InternalMessageInfo

func (m *AuthResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *AuthResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type UploadAvatarRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	FileName             string   `protobuf:"bytes,3,opt,name=FileName,proto3" json:"FileName,omitempty"`
	FileSize             int64    `protobuf:"varint,4,opt,name=FileSize,proto3" json:"FileSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadAvatarRequest) Reset()         { *m = UploadAvatarRequest{} }
func (m *UploadAvatarRequest) String() string { return proto.CompactTextString(m) }
func (*UploadAvatarRequest) ProtoMessage()    {}
func (*UploadAvatarRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{17}
}

func (m *UploadAvatarRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadAvatarRequest.Unmarshal(m, b)
}
func (m *UploadAvatarRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadAvatarRequest.Marshal(b, m, deterministic)
}
func (m *UploadAvatarRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadAvatarRequest.Merge(m, src)
}
func (m *UploadAvatarRequest) XXX_Size() int {
	return xxx_messageInfo_UploadAvatarRequest.Size(m)
}
func (m *UploadAvatarRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadAvatarRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UploadAvatarRequest proto.InternalMessageInfo

func (m *UploadAvatarRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *UploadAvatarRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *UploadAvatarRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *UploadAvatarRequest) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
	}
	return 0
}

type UploadAvatarResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadAvatarResponse) Reset()         { *m = UploadAvatarResponse{} }
func (m *UploadAvatarResponse) String() string { return proto.CompactTextString(m) }
func (*UploadAvatarResponse) ProtoMessage()    {}
func (*UploadAvatarResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{18}
}

func (m *UploadAvatarResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadAvatarResponse.Unmarshal(m, b)
}
func (m *UploadAvatarResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadAvatarResponse.Marshal(b, m, deterministic)
}
func (m *UploadAvatarResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadAvatarResponse.Merge(m, src)
}
func (m *UploadAvatarResponse) XXX_Size() int {
	return xxx_messageInfo_UploadAvatarResponse.Size(m)
}
func (m *UploadAvatarResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadAvatarResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UploadAvatarResponse proto.InternalMessageInfo

func (m *UploadAvatarResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *UploadAvatarResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.User.Message")
	proto.RegisterType((*GetImageCptRequest)(nil), "go.micro.srv.User.GetImageCptRequest")
	proto.RegisterType((*GetImageCptResponse)(nil), "go.micro.srv.User.GetImageCptResponse")
	proto.RegisterType((*GetImageCptResponse_Point)(nil), "go.micro.srv.User.GetImageCptResponse.Point")
	proto.RegisterType((*GetSmsCptRequest)(nil), "go.micro.srv.User.GetSmsCptRequest")
	proto.RegisterType((*GetSmsCptResponse)(nil), "go.micro.srv.User.GetSmsCptResponse")
	proto.RegisterType((*RegisterRequest)(nil), "go.micro.srv.User.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "go.micro.srv.User.RegisterResponse")
	proto.RegisterType((*LoginRequest)(nil), "go.micro.srv.User.LoginRequest")
	proto.RegisterType((*LoginResponse)(nil), "go.micro.srv.User.LoginResponse")
	proto.RegisterType((*LogoutRequest)(nil), "go.micro.srv.User.LogoutRequest")
	proto.RegisterType((*LogoutResponse)(nil), "go.micro.srv.User.LogoutResponse")
	proto.RegisterType((*GetUserInfoRequest)(nil), "go.micro.srv.User.GetUserInfoRequest")
	proto.RegisterType((*GetUserInfoResponse)(nil), "go.micro.srv.User.GetUserInfoResponse")
	proto.RegisterType((*RenameRequest)(nil), "go.micro.srv.User.RenameRequest")
	proto.RegisterType((*RenameResponse)(nil), "go.micro.srv.User.RenameResponse")
	proto.RegisterType((*AuthRequest)(nil), "go.micro.srv.User.AuthRequest")
	proto.RegisterType((*AuthResponse)(nil), "go.micro.srv.User.AuthResponse")
	proto.RegisterType((*UploadAvatarRequest)(nil), "go.micro.srv.User.UploadAvatarRequest")
	proto.RegisterType((*UploadAvatarResponse)(nil), "go.micro.srv.User.UploadAvatarResponse")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 701 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdf, 0x4f, 0x13, 0x4f,
	0x10, 0xe7, 0x68, 0x29, 0x74, 0x28, 0xdf, 0x2f, 0x2c, 0xc6, 0x5c, 0x0e, 0x23, 0x75, 0xf1, 0x47,
	0x1f, 0xf4, 0x4c, 0xf0, 0xcd, 0x18, 0x0d, 0x20, 0x92, 0x26, 0x14, 0xc9, 0xd5, 0x26, 0x6d, 0x62,
	0x34, 0x8b, 0xb7, 0x9e, 0x9b, 0xb4, 0xb7, 0xf5, 0x76, 0x0b, 0xd5, 0x17, 0xff, 0x05, 0x1f, 0xfd,
	0x73, 0xcd, 0xee, 0xdd, 0x5e, 0x7b, 0xf5, 0x5a, 0x4a, 0xf5, 0x89, 0x9d, 0xbd, 0xcf, 0xcc, 0x7c,
	0x66, 0x66, 0xe7, 0x43, 0x61, 0xa7, 0x1f, 0x71, 0xc9, 0x9f, 0xd2, 0x21, 0xe9, 0xf5, 0xbb, 0xd4,
	0xfc, 0x75, 0xf5, 0x2d, 0xda, 0x0a, 0xb8, 0xdb, 0x63, 0x9f, 0x22, 0xee, 0x8a, 0xe8, 0xd2, 0x6d,
	0x09, 0x1a, 0xe1, 0x1d, 0x58, 0x6d, 0x50, 0x21, 0x48, 0x40, 0xd1, 0x26, 0x14, 0x04, 0xf9, 0x66,
	0x5b, 0x55, 0xab, 0x56, 0xf6, 0xd4, 0x11, 0xd7, 0x00, 0x9d, 0x50, 0x59, 0xef, 0x91, 0x80, 0x1e,
	0xf5, 0xa5, 0x47, 0xbf, 0x0e, 0xa8, 0x90, 0x08, 0x41, 0xb1, 0x35, 0x60, 0x7e, 0x02, 0xd4, 0x67,
	0xfc, 0x73, 0x19, 0xb6, 0x33, 0x50, 0xd1, 0xe7, 0xa1, 0xa0, 0xc8, 0x86, 0xd5, 0xe3, 0x28, 0x3a,
	0xe2, 0x3e, 0x4d, 0xe0, 0xc6, 0x54, 0xd9, 0x1a, 0x22, 0xb0, 0x97, 0xe3, 0x6c, 0x0d, 0x11, 0xa8,
	0x9b, 0x73, 0x36, 0xb4, 0x0b, 0x55, 0xab, 0x56, 0xf1, 0xd4, 0x11, 0xdd, 0x86, 0x52, 0x53, 0x46,
	0xcc, 0xa7, 0x76, 0xb1, 0x6a, 0xd5, 0x0a, 0x5e, 0x62, 0xa1, 0x97, 0x50, 0x68, 0x90, 0xa1, 0xbd,
	0x52, 0xb5, 0x6a, 0xeb, 0xfb, 0x8f, 0xdd, 0x3f, 0xaa, 0x72, 0x73, 0xa8, 0xb8, 0xe7, 0x9c, 0x85,
	0xd2, 0x53, 0x8e, 0xda, 0x9f, 0x85, 0x76, 0x69, 0x21, 0x7f, 0x16, 0x3a, 0x7b, 0xb0, 0xa2, 0x2d,
	0x54, 0x01, 0xab, 0xad, 0x0b, 0x2b, 0x78, 0x56, 0x5b, 0x59, 0x1d, 0x5d, 0x50, 0xc1, 0xb3, 0x3a,
	0xd8, 0x83, 0xcd, 0x13, 0x2a, 0x9b, 0x3d, 0x31, 0xbb, 0x75, 0xaa, 0xc8, 0x06, 0xbf, 0x60, 0x5d,
	0x9a, 0xf4, 0x22, 0xb1, 0x14, 0xf6, 0x1d, 0x1d, 0x4a, 0xdd, 0x8f, 0xb2, 0xa7, 0xcf, 0xf8, 0x15,
	0x6c, 0x8d, 0xc5, 0xbc, 0x79, 0x8f, 0x71, 0x07, 0xfe, 0xf7, 0x68, 0xc0, 0x84, 0xa4, 0x91, 0xe1,
	0x34, 0xca, 0x6f, 0x65, 0xf2, 0x3b, 0xb0, 0x76, 0x4e, 0x84, 0xb8, 0xe2, 0x91, 0x9f, 0x44, 0x48,
	0xed, 0x5c, 0x6e, 0xef, 0x61, 0x73, 0x14, 0x7a, 0x81, 0xf1, 0xdf, 0x81, 0x72, 0x93, 0x0a, 0xc1,
	0x78, 0x58, 0xf7, 0x93, 0xc0, 0xa3, 0x0b, 0x7c, 0x08, 0x95, 0x53, 0x1e, 0xb0, 0xf0, 0x2f, 0x58,
	0xe3, 0x0e, 0x6c, 0x24, 0x31, 0xfe, 0x39, 0xbd, 0x27, 0x3a, 0x34, 0x1f, 0xa4, 0x93, 0xce, 0xc0,
	0xad, 0x49, 0xf8, 0x0b, 0xf8, 0xcf, 0xc0, 0x17, 0x18, 0xe2, 0xbe, 0x5e, 0x4b, 0xf5, 0x50, 0xeb,
	0xe1, 0x67, 0x3e, 0x5f, 0xc6, 0x96, 0xde, 0xcf, 0x91, 0xcf, 0x02, 0x1d, 0x40, 0x50, 0x7c, 0x4d,
	0x24, 0x49, 0x16, 0x54, 0x9f, 0xf1, 0x09, 0x6c, 0x78, 0x34, 0x24, 0x3d, 0x3a, 0x17, 0x0b, 0x95,
	0xee, 0x8c, 0x5e, 0x9d, 0x91, 0x9e, 0x79, 0xec, 0xc6, 0x54, 0x1d, 0x31, 0x81, 0x16, 0xe8, 0xc8,
	0x47, 0x58, 0x3f, 0x18, 0xc8, 0x2f, 0xf3, 0x91, 0x70, 0x60, 0xcd, 0xa3, 0xa4, 0x3b, 0xc6, 0x22,
	0xb5, 0xd5, 0xb3, 0xaa, 0xfb, 0x47, 0x24, 0x32, 0x23, 0x4e, 0x2c, 0xfc, 0x1c, 0x2a, 0x71, 0x82,
	0x05, 0xc8, 0xfd, 0x80, 0xed, 0x56, 0xbf, 0xcb, 0x89, 0x7f, 0x70, 0x49, 0x24, 0x89, 0xe6, 0x23,
	0x69, 0x9a, 0xbd, 0x3c, 0x6a, 0xb6, 0x22, 0xfe, 0x86, 0x75, 0xa9, 0x26, 0x1e, 0xd3, 0x4b, 0x6d,
	0xf3, 0xad, 0xc9, 0xbe, 0x1b, 0xb1, 0x4c, 0x6d, 0x7c, 0x08, 0xb7, 0xb2, 0x04, 0x6e, 0x5e, 0xc4,
	0xfe, 0xaf, 0x12, 0xac, 0x1e, 0xc7, 0xff, 0x4c, 0xd0, 0x07, 0x58, 0x1f, 0x13, 0x48, 0xf4, 0xe0,
	0x3a, 0x01, 0xd5, 0xf5, 0x3a, 0x0f, 0xe7, 0xd3, 0x59, 0xbc, 0x84, 0xda, 0x50, 0x4e, 0x55, 0x0e,
	0xed, 0xe5, 0xbb, 0x65, 0x74, 0xd5, 0xb9, 0x3f, 0x1b, 0x94, 0x46, 0x6e, 0xa9, 0xd1, 0xc7, 0x1a,
	0x85, 0x70, 0x8e, 0xcf, 0x84, 0x36, 0x3a, 0x7b, 0x33, 0x31, 0x69, 0xd8, 0x53, 0x58, 0xd1, 0xc2,
	0x82, 0x76, 0x73, 0xf0, 0xe3, 0xb2, 0xe5, 0x54, 0xa7, 0x03, 0xd2, 0x68, 0x6f, 0xa1, 0x14, 0x8b,
	0x03, 0x9a, 0x82, 0x1e, 0xc9, 0x8c, 0x73, 0x6f, 0x06, 0x22, 0x0d, 0x18, 0xcf, 0xcb, 0xec, 0xfe,
	0xb4, 0x79, 0x4d, 0xe8, 0xc9, 0xb4, 0x79, 0x4d, 0x4a, 0x48, 0x4c, 0x38, 0xde, 0xdd, 0x5c, 0xc2,
	0x19, 0x7d, 0xc8, 0x25, 0x9c, 0x5d, 0x7c, 0xbc, 0x84, 0xea, 0x50, 0x54, 0xdb, 0x86, 0xee, 0xe6,
	0x80, 0xc7, 0xf6, 0xdc, 0xd9, 0x9d, 0xfa, 0x3d, 0x0d, 0x45, 0xa0, 0x32, 0xfe, 0xf6, 0x51, 0x5e,
	0x55, 0x39, 0xdb, 0xe9, 0x3c, 0xba, 0x16, 0x67, 0x52, 0x5c, 0x94, 0xf4, 0x8f, 0xab, 0x67, 0xbf,
	0x03, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x49, 0x6c, 0x73, 0x7b, 0x09, 0x00, 0x00,
}