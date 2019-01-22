// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/example/example.proto

package go_micro_srv_House

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

type GetHousesRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetHousesRequest) Reset()         { *m = GetHousesRequest{} }
func (m *GetHousesRequest) String() string { return proto.CompactTextString(m) }
func (*GetHousesRequest) ProtoMessage()    {}
func (*GetHousesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{3}
}

func (m *GetHousesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHousesRequest.Unmarshal(m, b)
}
func (m *GetHousesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHousesRequest.Marshal(b, m, deterministic)
}
func (m *GetHousesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHousesRequest.Merge(m, src)
}
func (m *GetHousesRequest) XXX_Size() int {
	return xxx_messageInfo_GetHousesRequest.Size(m)
}
func (m *GetHousesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHousesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetHousesRequest proto.InternalMessageInfo

func (m *GetHousesRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type GetHousesResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetHousesResponse) Reset()         { *m = GetHousesResponse{} }
func (m *GetHousesResponse) String() string { return proto.CompactTextString(m) }
func (*GetHousesResponse) ProtoMessage()    {}
func (*GetHousesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{4}
}

func (m *GetHousesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHousesResponse.Unmarshal(m, b)
}
func (m *GetHousesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHousesResponse.Marshal(b, m, deterministic)
}
func (m *GetHousesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHousesResponse.Merge(m, src)
}
func (m *GetHousesResponse) XXX_Size() int {
	return xxx_messageInfo_GetHousesResponse.Size(m)
}
func (m *GetHousesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHousesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetHousesResponse proto.InternalMessageInfo

func (m *GetHousesResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetHousesResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *GetHousesResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type UploadImageRequest struct {
	SessionId            string   `protobuf:"bytes,1,opt,name=SessionId,proto3" json:"SessionId,omitempty"`
	HouseId              int64    `protobuf:"varint,2,opt,name=HouseId,proto3" json:"HouseId,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	FileSize             int64    `protobuf:"varint,4,opt,name=FileSize,proto3" json:"FileSize,omitempty"`
	FileName             string   `protobuf:"bytes,5,opt,name=FileName,proto3" json:"FileName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadImageRequest) Reset()         { *m = UploadImageRequest{} }
func (m *UploadImageRequest) String() string { return proto.CompactTextString(m) }
func (*UploadImageRequest) ProtoMessage()    {}
func (*UploadImageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{5}
}

func (m *UploadImageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageRequest.Unmarshal(m, b)
}
func (m *UploadImageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageRequest.Marshal(b, m, deterministic)
}
func (m *UploadImageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageRequest.Merge(m, src)
}
func (m *UploadImageRequest) XXX_Size() int {
	return xxx_messageInfo_UploadImageRequest.Size(m)
}
func (m *UploadImageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageRequest proto.InternalMessageInfo

func (m *UploadImageRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *UploadImageRequest) GetHouseId() int64 {
	if m != nil {
		return m.HouseId
	}
	return 0
}

func (m *UploadImageRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *UploadImageRequest) GetFileSize() int64 {
	if m != nil {
		return m.FileSize
	}
	return 0
}

func (m *UploadImageRequest) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

type UploadImageResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadImageResponse) Reset()         { *m = UploadImageResponse{} }
func (m *UploadImageResponse) String() string { return proto.CompactTextString(m) }
func (*UploadImageResponse) ProtoMessage()    {}
func (*UploadImageResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{6}
}

func (m *UploadImageResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadImageResponse.Unmarshal(m, b)
}
func (m *UploadImageResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadImageResponse.Marshal(b, m, deterministic)
}
func (m *UploadImageResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadImageResponse.Merge(m, src)
}
func (m *UploadImageResponse) XXX_Size() int {
	return xxx_messageInfo_UploadImageResponse.Size(m)
}
func (m *UploadImageResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadImageResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UploadImageResponse proto.InternalMessageInfo

func (m *UploadImageResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *UploadImageResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

type GetHouseDetailRequest struct {
	HouseId              int64    `protobuf:"varint,1,opt,name=HouseId,proto3" json:"HouseId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetHouseDetailRequest) Reset()         { *m = GetHouseDetailRequest{} }
func (m *GetHouseDetailRequest) String() string { return proto.CompactTextString(m) }
func (*GetHouseDetailRequest) ProtoMessage()    {}
func (*GetHouseDetailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{7}
}

func (m *GetHouseDetailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHouseDetailRequest.Unmarshal(m, b)
}
func (m *GetHouseDetailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHouseDetailRequest.Marshal(b, m, deterministic)
}
func (m *GetHouseDetailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHouseDetailRequest.Merge(m, src)
}
func (m *GetHouseDetailRequest) XXX_Size() int {
	return xxx_messageInfo_GetHouseDetailRequest.Size(m)
}
func (m *GetHouseDetailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHouseDetailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetHouseDetailRequest proto.InternalMessageInfo

func (m *GetHouseDetailRequest) GetHouseId() int64 {
	if m != nil {
		return m.HouseId
	}
	return 0
}

type GetHouseDetailResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetHouseDetailResponse) Reset()         { *m = GetHouseDetailResponse{} }
func (m *GetHouseDetailResponse) String() string { return proto.CompactTextString(m) }
func (*GetHouseDetailResponse) ProtoMessage()    {}
func (*GetHouseDetailResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{8}
}

func (m *GetHouseDetailResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetHouseDetailResponse.Unmarshal(m, b)
}
func (m *GetHouseDetailResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetHouseDetailResponse.Marshal(b, m, deterministic)
}
func (m *GetHouseDetailResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetHouseDetailResponse.Merge(m, src)
}
func (m *GetHouseDetailResponse) XXX_Size() int {
	return xxx_messageInfo_GetHouseDetailResponse.Size(m)
}
func (m *GetHouseDetailResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetHouseDetailResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetHouseDetailResponse proto.InternalMessageInfo

func (m *GetHouseDetailResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetHouseDetailResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *GetHouseDetailResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type GetIndexBannerRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetIndexBannerRequest) Reset()         { *m = GetIndexBannerRequest{} }
func (m *GetIndexBannerRequest) String() string { return proto.CompactTextString(m) }
func (*GetIndexBannerRequest) ProtoMessage()    {}
func (*GetIndexBannerRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{9}
}

func (m *GetIndexBannerRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetIndexBannerRequest.Unmarshal(m, b)
}
func (m *GetIndexBannerRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetIndexBannerRequest.Marshal(b, m, deterministic)
}
func (m *GetIndexBannerRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetIndexBannerRequest.Merge(m, src)
}
func (m *GetIndexBannerRequest) XXX_Size() int {
	return xxx_messageInfo_GetIndexBannerRequest.Size(m)
}
func (m *GetIndexBannerRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetIndexBannerRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetIndexBannerRequest proto.InternalMessageInfo

type GetIndexBannerResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetIndexBannerResponse) Reset()         { *m = GetIndexBannerResponse{} }
func (m *GetIndexBannerResponse) String() string { return proto.CompactTextString(m) }
func (*GetIndexBannerResponse) ProtoMessage()    {}
func (*GetIndexBannerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{10}
}

func (m *GetIndexBannerResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetIndexBannerResponse.Unmarshal(m, b)
}
func (m *GetIndexBannerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetIndexBannerResponse.Marshal(b, m, deterministic)
}
func (m *GetIndexBannerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetIndexBannerResponse.Merge(m, src)
}
func (m *GetIndexBannerResponse) XXX_Size() int {
	return xxx_messageInfo_GetIndexBannerResponse.Size(m)
}
func (m *GetIndexBannerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetIndexBannerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetIndexBannerResponse proto.InternalMessageInfo

func (m *GetIndexBannerResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *GetIndexBannerResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *GetIndexBannerResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type SearchRequest struct {
	AreaId               string   `protobuf:"bytes,1,opt,name=AreaId,proto3" json:"AreaId,omitempty"`
	StartDate            string   `protobuf:"bytes,2,opt,name=StartDate,proto3" json:"StartDate,omitempty"`
	EndDate              string   `protobuf:"bytes,3,opt,name=EndDate,proto3" json:"EndDate,omitempty"`
	Page                 string   `protobuf:"bytes,4,opt,name=Page,proto3" json:"Page,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchRequest) Reset()         { *m = SearchRequest{} }
func (m *SearchRequest) String() string { return proto.CompactTextString(m) }
func (*SearchRequest) ProtoMessage()    {}
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{11}
}

func (m *SearchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchRequest.Unmarshal(m, b)
}
func (m *SearchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchRequest.Marshal(b, m, deterministic)
}
func (m *SearchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchRequest.Merge(m, src)
}
func (m *SearchRequest) XXX_Size() int {
	return xxx_messageInfo_SearchRequest.Size(m)
}
func (m *SearchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SearchRequest proto.InternalMessageInfo

func (m *SearchRequest) GetAreaId() string {
	if m != nil {
		return m.AreaId
	}
	return ""
}

func (m *SearchRequest) GetStartDate() string {
	if m != nil {
		return m.StartDate
	}
	return ""
}

func (m *SearchRequest) GetEndDate() string {
	if m != nil {
		return m.EndDate
	}
	return ""
}

func (m *SearchRequest) GetPage() string {
	if m != nil {
		return m.Page
	}
	return ""
}

type SearchResponse struct {
	ErrCode              string   `protobuf:"bytes,1,opt,name=ErrCode,proto3" json:"ErrCode,omitempty"`
	ErrMsg               string   `protobuf:"bytes,2,opt,name=ErrMsg,proto3" json:"ErrMsg,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=Data,proto3" json:"Data,omitempty"`
	TotalPage            int64    `protobuf:"varint,4,opt,name=TotalPage,proto3" json:"TotalPage,omitempty"`
	Page                 int64    `protobuf:"varint,5,opt,name=Page,proto3" json:"Page,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchResponse) Reset()         { *m = SearchResponse{} }
func (m *SearchResponse) String() string { return proto.CompactTextString(m) }
func (*SearchResponse) ProtoMessage()    {}
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_097b3f5db5cf5789, []int{12}
}

func (m *SearchResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchResponse.Unmarshal(m, b)
}
func (m *SearchResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchResponse.Marshal(b, m, deterministic)
}
func (m *SearchResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchResponse.Merge(m, src)
}
func (m *SearchResponse) XXX_Size() int {
	return xxx_messageInfo_SearchResponse.Size(m)
}
func (m *SearchResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SearchResponse proto.InternalMessageInfo

func (m *SearchResponse) GetErrCode() string {
	if m != nil {
		return m.ErrCode
	}
	return ""
}

func (m *SearchResponse) GetErrMsg() string {
	if m != nil {
		return m.ErrMsg
	}
	return ""
}

func (m *SearchResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SearchResponse) GetTotalPage() int64 {
	if m != nil {
		return m.TotalPage
	}
	return 0
}

func (m *SearchResponse) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "go.micro.srv.House.Message")
	proto.RegisterType((*AddRequest)(nil), "go.micro.srv.House.AddRequest")
	proto.RegisterType((*AddResponse)(nil), "go.micro.srv.House.AddResponse")
	proto.RegisterType((*GetHousesRequest)(nil), "go.micro.srv.House.GetHousesRequest")
	proto.RegisterType((*GetHousesResponse)(nil), "go.micro.srv.House.GetHousesResponse")
	proto.RegisterType((*UploadImageRequest)(nil), "go.micro.srv.House.UploadImageRequest")
	proto.RegisterType((*UploadImageResponse)(nil), "go.micro.srv.House.UploadImageResponse")
	proto.RegisterType((*GetHouseDetailRequest)(nil), "go.micro.srv.House.GetHouseDetailRequest")
	proto.RegisterType((*GetHouseDetailResponse)(nil), "go.micro.srv.House.GetHouseDetailResponse")
	proto.RegisterType((*GetIndexBannerRequest)(nil), "go.micro.srv.House.GetIndexBannerRequest")
	proto.RegisterType((*GetIndexBannerResponse)(nil), "go.micro.srv.House.GetIndexBannerResponse")
	proto.RegisterType((*SearchRequest)(nil), "go.micro.srv.House.SearchRequest")
	proto.RegisterType((*SearchResponse)(nil), "go.micro.srv.House.SearchResponse")
}

func init() { proto.RegisterFile("proto/example/example.proto", fileDescriptor_097b3f5db5cf5789) }

var fileDescriptor_097b3f5db5cf5789 = []byte{
	// 523 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xdb, 0x6e, 0xd3, 0x4c,
	0x10, 0xae, 0x7f, 0xa7, 0xc9, 0xef, 0x29, 0x54, 0x65, 0x10, 0xc1, 0x72, 0x11, 0x94, 0x15, 0x87,
	0xc2, 0x85, 0x39, 0xdd, 0x83, 0x02, 0x09, 0x6d, 0x2e, 0x8a, 0xc0, 0x81, 0x0b, 0xb8, 0x40, 0x2c,
	0xdd, 0x91, 0xb1, 0xe4, 0x78, 0xc3, 0xae, 0x8b, 0x0a, 0x4f, 0xc0, 0x23, 0xf0, 0x2c, 0x3c, 0x1d,
	0xca, 0xfa, 0x10, 0x3b, 0x71, 0xda, 0xa8, 0xca, 0x55, 0x76, 0x4e, 0xfb, 0x7d, 0x33, 0x9e, 0x6f,
	0x03, 0xbb, 0x13, 0x25, 0x53, 0xf9, 0x88, 0x4e, 0xf9, 0x78, 0x12, 0x53, 0xf1, 0xeb, 0x1b, 0x2f,
	0x62, 0x28, 0xfd, 0x71, 0x74, 0xac, 0xa4, 0xaf, 0xd5, 0x0f, 0xff, 0x50, 0x9e, 0x68, 0x62, 0xbb,
	0xd0, 0x39, 0x22, 0xad, 0x79, 0x48, 0xb8, 0x03, 0xb6, 0xe6, 0x3f, 0x5d, 0x6b, 0xcf, 0xda, 0x77,
	0x82, 0xe9, 0x91, 0x3d, 0x07, 0xe8, 0x09, 0x11, 0xd0, 0xf7, 0x13, 0xd2, 0x29, 0xde, 0x00, 0x67,
	0x44, 0x5a, 0x47, 0x32, 0x19, 0x8a, 0x3c, 0x6b, 0xe6, 0x40, 0x84, 0x56, 0x9f, 0xa7, 0xdc, 0xfd,
	0x6f, 0xcf, 0xda, 0xbf, 0x14, 0x98, 0x33, 0x7b, 0x01, 0x5b, 0xa6, 0x5e, 0x4f, 0x64, 0xa2, 0x09,
	0x5d, 0xe8, 0x0c, 0x94, 0x7a, 0x25, 0x05, 0xe5, 0xe5, 0x85, 0x89, 0x5d, 0x68, 0x0f, 0x94, 0x3a,
	0xd2, 0xa1, 0x29, 0x77, 0x82, 0xdc, 0x62, 0x8f, 0x61, 0xe7, 0x80, 0x52, 0xc3, 0x54, 0xaf, 0x44,
	0x83, 0x7d, 0x84, 0x2b, 0x95, 0x8a, 0x8b, 0x02, 0x97, 0xdd, 0xd8, 0x95, 0x6e, 0xfe, 0x58, 0x80,
	0x1f, 0x26, 0xb1, 0xe4, 0x62, 0x38, 0xe6, 0x21, 0xad, 0x36, 0x16, 0x17, 0x3a, 0x86, 0xcc, 0x50,
	0x18, 0x04, 0x3b, 0x28, 0xcc, 0x26, 0x08, 0xf4, 0xe0, 0xff, 0xd7, 0x51, 0x4c, 0xa3, 0xe8, 0x17,
	0xb9, 0x2d, 0x93, 0x5e, 0xda, 0x45, 0xec, 0x0d, 0x1f, 0x93, 0xbb, 0x69, 0x60, 0x4a, 0x9b, 0x1d,
	0xc0, 0xd5, 0x1a, 0xb3, 0x0b, 0x0f, 0xfc, 0x09, 0x5c, 0x2b, 0xc6, 0xd7, 0xa7, 0x94, 0x47, 0x71,
	0xd1, 0x65, 0xa5, 0x0f, 0xab, 0xd6, 0x07, 0xfb, 0x0c, 0xdd, 0xf9, 0x92, 0xb5, 0x8e, 0xfd, 0xba,
	0xa1, 0x34, 0x4c, 0x04, 0x9d, 0xbe, 0xe4, 0x49, 0x42, 0x2a, 0xa7, 0x94, 0x03, 0xd7, 0x02, 0x6b,
	0x05, 0xd6, 0x70, 0x79, 0x44, 0x5c, 0x1d, 0x7f, 0x2b, 0x66, 0xd0, 0x85, 0x76, 0x4f, 0x11, 0x2f,
	0x3f, 0x73, 0x6e, 0x99, 0x0d, 0x48, 0xb9, 0x4a, 0xfb, 0x3c, 0xa5, 0xfc, 0xde, 0x99, 0xc3, 0x90,
	0x49, 0x84, 0x89, 0xd9, 0x39, 0x99, 0xcc, 0x9c, 0x82, 0xbe, 0xe5, 0x61, 0xf6, 0xa5, 0x9d, 0xc0,
	0x9c, 0xd9, 0x6f, 0x0b, 0xb6, 0x0b, 0xd4, 0x75, 0x76, 0x33, 0x25, 0xf9, 0x5e, 0xa6, 0x3c, 0x2e,
	0x11, 0xed, 0x60, 0xe6, 0x28, 0xa9, 0x6c, 0x9a, 0x80, 0x39, 0x3f, 0xfd, 0xdb, 0x82, 0xce, 0x20,
	0x7b, 0x40, 0xf0, 0x10, 0xec, 0x9e, 0x10, 0x78, 0xd3, 0x5f, 0x7c, 0x42, 0xfc, 0xd9, 0x13, 0xe1,
	0xdd, 0x5a, 0x1a, 0xcf, 0x7a, 0x61, 0x1b, 0xf8, 0x09, 0x9c, 0x52, 0xa0, 0x78, 0xa7, 0x29, 0x7f,
	0x5e, 0xf1, 0xde, 0xdd, 0x73, 0xb2, 0xca, 0xbb, 0xbf, 0xc0, 0x56, 0x45, 0x06, 0x78, 0xaf, 0xa9,
	0x6e, 0x51, 0xc1, 0xde, 0xfd, 0x73, 0xf3, 0x4a, 0x84, 0x08, 0xb6, 0xeb, 0xcb, 0x8e, 0x0f, 0xce,
	0x22, 0x57, 0xd3, 0x90, 0xf7, 0x70, 0x95, 0xd4, 0x39, 0xa8, 0xca, 0x7a, 0x2f, 0x85, 0x5a, 0xd4,
	0xc6, 0x52, 0xa8, 0x06, 0xb5, 0xb0, 0x0d, 0x7c, 0x07, 0xed, 0x6c, 0xe7, 0xf0, 0x76, 0x53, 0x5d,
	0x4d, 0x05, 0x1e, 0x3b, 0x2b, 0xa5, 0xb8, 0xf2, 0x6b, 0xdb, 0xfc, 0xe5, 0x3c, 0xfb, 0x17, 0x00,
	0x00, 0xff, 0xff, 0xa5, 0xe6, 0x13, 0xc8, 0x91, 0x06, 0x00, 0x00,
}
