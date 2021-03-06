// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/infobloxopen/feature-flag-service/pkg/pb/service.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	query "github.com/infobloxopen/atlas-app-toolkit/query"
	_ "github.com/lyft/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
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

type VersionResponse struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{0}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type FeatureFlag struct {
	FeatureName          string   `protobuf:"bytes,1,opt,name=feature_name,json=featureName,proto3" json:"feature_name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Origin               string   `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeatureFlag) Reset()         { *m = FeatureFlag{} }
func (m *FeatureFlag) String() string { return proto.CompactTextString(m) }
func (*FeatureFlag) ProtoMessage()    {}
func (*FeatureFlag) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{1}
}

func (m *FeatureFlag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureFlag.Unmarshal(m, b)
}
func (m *FeatureFlag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureFlag.Marshal(b, m, deterministic)
}
func (m *FeatureFlag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureFlag.Merge(m, src)
}
func (m *FeatureFlag) XXX_Size() int {
	return xxx_messageInfo_FeatureFlag.Size(m)
}
func (m *FeatureFlag) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureFlag.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureFlag proto.InternalMessageInfo

func (m *FeatureFlag) GetFeatureName() string {
	if m != nil {
		return m.FeatureName
	}
	return ""
}

func (m *FeatureFlag) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *FeatureFlag) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

type ListFeatureFlagsRequest struct {
	Filter               *query.Filtering  `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	Labels               map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ListFeatureFlagsRequest) Reset()         { *m = ListFeatureFlagsRequest{} }
func (m *ListFeatureFlagsRequest) String() string { return proto.CompactTextString(m) }
func (*ListFeatureFlagsRequest) ProtoMessage()    {}
func (*ListFeatureFlagsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{2}
}

func (m *ListFeatureFlagsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListFeatureFlagsRequest.Unmarshal(m, b)
}
func (m *ListFeatureFlagsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListFeatureFlagsRequest.Marshal(b, m, deterministic)
}
func (m *ListFeatureFlagsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListFeatureFlagsRequest.Merge(m, src)
}
func (m *ListFeatureFlagsRequest) XXX_Size() int {
	return xxx_messageInfo_ListFeatureFlagsRequest.Size(m)
}
func (m *ListFeatureFlagsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListFeatureFlagsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListFeatureFlagsRequest proto.InternalMessageInfo

func (m *ListFeatureFlagsRequest) GetFilter() *query.Filtering {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *ListFeatureFlagsRequest) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type ListFeatureFlagsResponse struct {
	Results              []*FeatureFlag  `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	Page                 *query.PageInfo `protobuf:"bytes,2,opt,name=page,proto3" json:"page,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ListFeatureFlagsResponse) Reset()         { *m = ListFeatureFlagsResponse{} }
func (m *ListFeatureFlagsResponse) String() string { return proto.CompactTextString(m) }
func (*ListFeatureFlagsResponse) ProtoMessage()    {}
func (*ListFeatureFlagsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{3}
}

func (m *ListFeatureFlagsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListFeatureFlagsResponse.Unmarshal(m, b)
}
func (m *ListFeatureFlagsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListFeatureFlagsResponse.Marshal(b, m, deterministic)
}
func (m *ListFeatureFlagsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListFeatureFlagsResponse.Merge(m, src)
}
func (m *ListFeatureFlagsResponse) XXX_Size() int {
	return xxx_messageInfo_ListFeatureFlagsResponse.Size(m)
}
func (m *ListFeatureFlagsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListFeatureFlagsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListFeatureFlagsResponse proto.InternalMessageInfo

func (m *ListFeatureFlagsResponse) GetResults() []*FeatureFlag {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *ListFeatureFlagsResponse) GetPage() *query.PageInfo {
	if m != nil {
		return m.Page
	}
	return nil
}

type ReadFeatureFlagRequest struct {
	FeatureName          string                `protobuf:"bytes,1,opt,name=feature_name,json=featureName,proto3" json:"feature_name,omitempty"`
	Fields               *query.FieldSelection `protobuf:"bytes,2,opt,name=fields,proto3" json:"fields,omitempty"`
	Labels               map[string]string     `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ReadFeatureFlagRequest) Reset()         { *m = ReadFeatureFlagRequest{} }
func (m *ReadFeatureFlagRequest) String() string { return proto.CompactTextString(m) }
func (*ReadFeatureFlagRequest) ProtoMessage()    {}
func (*ReadFeatureFlagRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{4}
}

func (m *ReadFeatureFlagRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadFeatureFlagRequest.Unmarshal(m, b)
}
func (m *ReadFeatureFlagRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadFeatureFlagRequest.Marshal(b, m, deterministic)
}
func (m *ReadFeatureFlagRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadFeatureFlagRequest.Merge(m, src)
}
func (m *ReadFeatureFlagRequest) XXX_Size() int {
	return xxx_messageInfo_ReadFeatureFlagRequest.Size(m)
}
func (m *ReadFeatureFlagRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadFeatureFlagRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReadFeatureFlagRequest proto.InternalMessageInfo

func (m *ReadFeatureFlagRequest) GetFeatureName() string {
	if m != nil {
		return m.FeatureName
	}
	return ""
}

func (m *ReadFeatureFlagRequest) GetFields() *query.FieldSelection {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *ReadFeatureFlagRequest) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type ReadFeatureFlagResponse struct {
	Result               *FeatureFlag `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ReadFeatureFlagResponse) Reset()         { *m = ReadFeatureFlagResponse{} }
func (m *ReadFeatureFlagResponse) String() string { return proto.CompactTextString(m) }
func (*ReadFeatureFlagResponse) ProtoMessage()    {}
func (*ReadFeatureFlagResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6544829dd69549ee, []int{5}
}

func (m *ReadFeatureFlagResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadFeatureFlagResponse.Unmarshal(m, b)
}
func (m *ReadFeatureFlagResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadFeatureFlagResponse.Marshal(b, m, deterministic)
}
func (m *ReadFeatureFlagResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadFeatureFlagResponse.Merge(m, src)
}
func (m *ReadFeatureFlagResponse) XXX_Size() int {
	return xxx_messageInfo_ReadFeatureFlagResponse.Size(m)
}
func (m *ReadFeatureFlagResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadFeatureFlagResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReadFeatureFlagResponse proto.InternalMessageInfo

func (m *ReadFeatureFlagResponse) GetResult() *FeatureFlag {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*VersionResponse)(nil), "service.VersionResponse")
	proto.RegisterType((*FeatureFlag)(nil), "service.FeatureFlag")
	proto.RegisterType((*ListFeatureFlagsRequest)(nil), "service.ListFeatureFlagsRequest")
	proto.RegisterMapType((map[string]string)(nil), "service.ListFeatureFlagsRequest.LabelsEntry")
	proto.RegisterType((*ListFeatureFlagsResponse)(nil), "service.ListFeatureFlagsResponse")
	proto.RegisterType((*ReadFeatureFlagRequest)(nil), "service.ReadFeatureFlagRequest")
	proto.RegisterMapType((map[string]string)(nil), "service.ReadFeatureFlagRequest.LabelsEntry")
	proto.RegisterType((*ReadFeatureFlagResponse)(nil), "service.ReadFeatureFlagResponse")
}

func init() {
	proto.RegisterFile("github.com/infobloxopen/feature-flag-service/pkg/pb/service.proto", fileDescriptor_6544829dd69549ee)
}

var fileDescriptor_6544829dd69549ee = []byte{
	// 758 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x55, 0xc1, 0x6e, 0x23, 0x35,
	0x18, 0xd6, 0x24, 0x21, 0x65, 0x1d, 0x24, 0x8a, 0xb5, 0x4a, 0x86, 0xb0, 0xd2, 0x7a, 0x73, 0x2a,
	0xbb, 0x8d, 0xbd, 0x1b, 0x7a, 0x00, 0x7a, 0x40, 0xa5, 0xb4, 0x15, 0xa8, 0x82, 0x6a, 0x5a, 0x71,
	0xe0, 0xd0, 0xca, 0x93, 0xf9, 0x33, 0x35, 0x75, 0xc6, 0xae, 0xed, 0x19, 0x88, 0x10, 0x17, 0x6e,
	0x5c, 0xe1, 0x45, 0xfa, 0x02, 0x3c, 0x01, 0x47, 0x5e, 0x81, 0x17, 0xe0, 0xc4, 0x15, 0xcd, 0x8c,
	0x13, 0x0d, 0x6d, 0x4a, 0x41, 0xe2, 0x36, 0xf6, 0xff, 0xf9, 0xfb, 0xfc, 0x7f, 0xfe, 0xec, 0x41,
	0x1f, 0xa5, 0xc2, 0x5d, 0xe6, 0x31, 0x9d, 0xaa, 0x39, 0xfb, 0x34, 0x9b, 0xa9, 0x58, 0xaa, 0x6f,
	0xc7, 0xfb, 0x67, 0x5f, 0x30, 0xee, 0x24, 0xb7, 0x74, 0x06, 0xdc, 0xe5, 0x06, 0xe8, 0x4c, 0xf2,
	0x94, 0xe9, 0xab, 0x94, 0xe9, 0x98, 0x59, 0x30, 0x85, 0x98, 0x02, 0xd5, 0x46, 0x39, 0x85, 0x37,
	0xfc, 0x70, 0xf8, 0x4e, 0xaa, 0x54, 0x2a, 0x81, 0x55, 0xd3, 0x71, 0x3e, 0x63, 0x30, 0xd7, 0x6e,
	0x51, 0xa3, 0x86, 0x4f, 0x7c, 0x91, 0x6b, 0xc1, 0x78, 0x96, 0x29, 0xc7, 0x9d, 0x50, 0x99, 0xf5,
	0xd5, 0xdd, 0xc6, 0x26, 0xe4, 0x62, 0xe6, 0x6a, 0x8e, 0xe9, 0x38, 0x85, 0x6c, 0x5c, 0x70, 0x29,
	0x12, 0xee, 0x80, 0xdd, 0xf9, 0xf0, 0x8b, 0xb7, 0x1b, 0x60, 0xfb, 0x0d, 0x4f, 0x53, 0x30, 0x4c,
	0xe9, 0x8a, 0x7e, 0x8d, 0xd4, 0x67, 0x0d, 0x29, 0xe1, 0xfb, 0x55, 0x1a, 0xb2, 0xba, 0xdf, 0x31,
	0xd7, 0x7a, 0xec, 0x94, 0x92, 0x57, 0xc2, 0xb1, 0xeb, 0x1c, 0xcc, 0x82, 0x4d, 0x95, 0x94, 0x30,
	0x2d, 0x29, 0x2e, 0x94, 0x06, 0xc3, 0x9d, 0x32, 0x9e, 0x6b, 0xf4, 0x02, 0xbd, 0xf9, 0x25, 0x18,
	0x2b, 0x54, 0x16, 0x81, 0xd5, 0x2a, 0xb3, 0x80, 0x43, 0xb4, 0x51, 0xd4, 0x53, 0x61, 0x40, 0x82,
	0xad, 0x47, 0xd1, 0x72, 0x38, 0x3a, 0x47, 0xbd, 0xc3, 0xda, 0xcc, 0x43, 0xc9, 0x53, 0xfc, 0x0c,
	0xbd, 0xe1, 0xbd, 0xbd, 0xc8, 0xf8, 0x1c, 0x3c, 0xba, 0xe7, 0xe7, 0x3e, 0xe7, 0x73, 0xc0, 0x8f,
	0xd1, 0x6b, 0x05, 0x97, 0x39, 0x84, 0xad, 0xaa, 0x56, 0x0f, 0x70, 0x1f, 0x75, 0x95, 0x11, 0xa9,
	0xc8, 0xc2, 0x76, 0x35, 0xed, 0x47, 0xa3, 0x5f, 0x03, 0x34, 0x38, 0x16, 0xd6, 0x35, 0x44, 0x6c,
	0x04, 0xd7, 0x39, 0x58, 0x87, 0x19, 0xea, 0xce, 0x84, 0x74, 0x60, 0x2a, 0x99, 0xde, 0x64, 0x40,
	0x97, 0xad, 0x53, 0xae, 0x05, 0x3d, 0xac, 0x6a, 0x22, 0x4b, 0x23, 0x0f, 0xc3, 0x9f, 0xa0, 0xae,
	0xe4, 0x31, 0x48, 0x1b, 0xb6, 0x48, 0x7b, 0xab, 0x37, 0xd9, 0xa6, 0xcb, 0x43, 0xbf, 0x47, 0x82,
	0x1e, 0x57, 0xf0, 0x83, 0xcc, 0x99, 0x45, 0xe4, 0xd7, 0x0e, 0x3f, 0x40, 0xbd, 0xc6, 0x34, 0xde,
	0x44, 0xed, 0x2b, 0x58, 0xf8, 0x4e, 0xcb, 0xcf, 0xf5, 0x1d, 0x7e, 0xd8, 0x7a, 0x3f, 0x18, 0x15,
	0x28, 0xbc, 0xab, 0xe4, 0x3d, 0xa6, 0x68, 0xc3, 0x80, 0xcd, 0xa5, 0xb3, 0x61, 0x50, 0xed, 0xee,
	0xf1, 0x6a, 0x77, 0x0d, 0x7c, 0xb4, 0x04, 0xe1, 0xe7, 0xa8, 0xa3, 0x79, 0x5a, 0x8b, 0xf4, 0x26,
	0xfd, 0xbf, 0xf7, 0x7e, 0xc2, 0x53, 0x28, 0x73, 0x1f, 0x55, 0x98, 0xd1, 0x1f, 0x01, 0xea, 0x47,
	0xc0, 0x93, 0x26, 0x91, 0x37, 0xf1, 0x5f, 0x9c, 0xd8, 0x4e, 0xe9, 0x33, 0xc8, 0xc4, 0x7a, 0xad,
	0x27, 0xb7, 0x7d, 0x06, 0x99, 0x9c, 0x82, 0x8f, 0x53, 0xe4, 0xb1, 0x78, 0x7f, 0x65, 0x76, 0xbb,
	0x6a, 0xe7, 0xc5, 0xaa, 0x9d, 0xf5, 0x3b, 0xf9, 0xbf, 0xbd, 0x3e, 0x42, 0x83, 0x3b, 0x42, 0xde,
	0xea, 0x6d, 0xd4, 0xad, 0x5d, 0xf4, 0xc1, 0x59, 0xef, 0xb4, 0xc7, 0x4c, 0x6e, 0x5a, 0x68, 0x73,
	0xaf, 0xbc, 0x46, 0xcd, 0xa0, 0x9f, 0x20, 0x74, 0x04, 0xce, 0xdf, 0x13, 0xdc, 0xa7, 0xf5, 0x43,
	0x40, 0x97, 0xaf, 0x04, 0x3d, 0x28, 0x5f, 0x89, 0x61, 0xb8, 0x22, 0xbe, 0x75, 0xa3, 0x46, 0x9b,
	0x3f, 0xfc, 0xf6, 0xfb, 0xcf, 0x2d, 0x84, 0x5f, 0x67, 0xfe, 0x26, 0xe1, 0x73, 0xd4, 0x29, 0xb3,
	0x81, 0xc9, 0x43, 0xa1, 0x1c, 0x3e, 0xfb, 0x07, 0x84, 0xa7, 0x7f, 0xab, 0xa2, 0xef, 0xe1, 0x47,
	0xcc, 0x1f, 0xa4, 0xc5, 0x5f, 0xa3, 0x4e, 0xe9, 0x07, 0x7e, 0xfa, 0xc0, 0x39, 0x0c, 0xc9, 0xfd,
	0x00, 0xcf, 0xfe, 0xb4, 0x62, 0x7f, 0x1b, 0x0f, 0x96, 0xec, 0xec, 0xbb, 0x66, 0x86, 0xbe, 0xff,
	0xf8, 0xcf, 0xe0, 0xa7, 0xbd, 0x5f, 0x02, 0x7c, 0x13, 0xa0, 0x61, 0xe5, 0x1c, 0xf1, 0x2c, 0xa4,
	0xa4, 0x21, 0xa7, 0x35, 0x3d, 0xfe, 0x31, 0x58, 0x53, 0xf4, 0xda, 0x64, 0xeb, 0x0c, 0xcc, 0x5c,
	0x64, 0xb9, 0x7d, 0x97, 0x68, 0xa3, 0x0a, 0x91, 0x80, 0x25, 0xee, 0x12, 0xc8, 0x94, 0x6b, 0x1e,
	0x0b, 0x29, 0xdc, 0x82, 0x38, 0x45, 0x44, 0xe6, 0x8c, 0x4a, 0xf2, 0x29, 0x2c, 0xd7, 0xd9, 0x55,
	0xc3, 0x65, 0x4d, 0x95, 0x6b, 0x93, 0xbc, 0x0a, 0x26, 0x51, 0x19, 0xe1, 0xc4, 0x19, 0x2e, 0x24,
	0x51, 0x86, 0x48, 0x31, 0x17, 0x0e, 0x12, 0x12, 0x73, 0x2b, 0x2c, 0x9d, 0x74, 0x8b, 0x57, 0xf4,
	0x25, 0x7d, 0x39, 0xea, 0xb0, 0xe2, 0x15, 0x7b, 0xde, 0x0a, 0x5a, 0x5f, 0xed, 0xfc, 0xe7, 0x5f,
	0xc8, 0xae, 0x8e, 0xe3, 0x6e, 0x95, 0x80, 0xf7, 0xfe, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x9a,
	0x6a, 0xf2, 0x81, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AtlasFeatureFlagClient is the client API for AtlasFeatureFlag service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AtlasFeatureFlagClient interface {
	GetVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	List(ctx context.Context, in *ListFeatureFlagsRequest, opts ...grpc.CallOption) (*ListFeatureFlagsResponse, error)
	Read(ctx context.Context, in *ReadFeatureFlagRequest, opts ...grpc.CallOption) (*ReadFeatureFlagResponse, error)
}

type atlasFeatureFlagClient struct {
	cc *grpc.ClientConn
}

func NewAtlasFeatureFlagClient(cc *grpc.ClientConn) AtlasFeatureFlagClient {
	return &atlasFeatureFlagClient{cc}
}

func (c *atlasFeatureFlagClient) GetVersion(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/service.AtlasFeatureFlag/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atlasFeatureFlagClient) List(ctx context.Context, in *ListFeatureFlagsRequest, opts ...grpc.CallOption) (*ListFeatureFlagsResponse, error) {
	out := new(ListFeatureFlagsResponse)
	err := c.cc.Invoke(ctx, "/service.AtlasFeatureFlag/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *atlasFeatureFlagClient) Read(ctx context.Context, in *ReadFeatureFlagRequest, opts ...grpc.CallOption) (*ReadFeatureFlagResponse, error) {
	out := new(ReadFeatureFlagResponse)
	err := c.cc.Invoke(ctx, "/service.AtlasFeatureFlag/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AtlasFeatureFlagServer is the server API for AtlasFeatureFlag service.
type AtlasFeatureFlagServer interface {
	GetVersion(context.Context, *empty.Empty) (*VersionResponse, error)
	List(context.Context, *ListFeatureFlagsRequest) (*ListFeatureFlagsResponse, error)
	Read(context.Context, *ReadFeatureFlagRequest) (*ReadFeatureFlagResponse, error)
}

func RegisterAtlasFeatureFlagServer(s *grpc.Server, srv AtlasFeatureFlagServer) {
	s.RegisterService(&_AtlasFeatureFlag_serviceDesc, srv)
}

func _AtlasFeatureFlag_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtlasFeatureFlagServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.AtlasFeatureFlag/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtlasFeatureFlagServer).GetVersion(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtlasFeatureFlag_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFeatureFlagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtlasFeatureFlagServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.AtlasFeatureFlag/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtlasFeatureFlagServer).List(ctx, req.(*ListFeatureFlagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AtlasFeatureFlag_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadFeatureFlagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AtlasFeatureFlagServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.AtlasFeatureFlag/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AtlasFeatureFlagServer).Read(ctx, req.(*ReadFeatureFlagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AtlasFeatureFlag_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.AtlasFeatureFlag",
	HandlerType: (*AtlasFeatureFlagServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _AtlasFeatureFlag_GetVersion_Handler,
		},
		{
			MethodName: "List",
			Handler:    _AtlasFeatureFlag_List_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _AtlasFeatureFlag_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/infobloxopen/feature-flag-service/pkg/pb/service.proto",
}
