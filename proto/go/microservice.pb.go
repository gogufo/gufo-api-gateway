// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: microservice.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Module     *string               `protobuf:"bytes,1,opt,name=Module,proto3,oneof" json:"Module,omitempty"`
	Param      *string               `protobuf:"bytes,2,opt,name=Param,proto3,oneof" json:"Param,omitempty"`
	ParamID    *string               `protobuf:"bytes,3,opt,name=ParamID,proto3,oneof" json:"ParamID,omitempty"`
	Action     *string               `protobuf:"bytes,4,opt,name=Action,proto3,oneof" json:"Action,omitempty"`
	Args       map[string]*anypb.Any `protobuf:"bytes,5,rep,name=Args,proto3" json:"Args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Path       *string               `protobuf:"bytes,6,opt,name=Path,proto3,oneof" json:"Path,omitempty"`
	Token      *string               `protobuf:"bytes,7,opt,name=Token,proto3,oneof" json:"Token,omitempty"`
	Sign       *string               `protobuf:"bytes,8,opt,name=Sign,proto3,oneof" json:"Sign,omitempty"`
	IP         *string               `protobuf:"bytes,9,opt,name=IP,proto3,oneof" json:"IP,omitempty"`
	UserAgent  *string               `protobuf:"bytes,10,opt,name=UserAgent,proto3,oneof" json:"UserAgent,omitempty"`
	TokenType  *string               `protobuf:"bytes,11,opt,name=TokenType,proto3,oneof" json:"TokenType,omitempty"`
	TimeStamp  *int32                `protobuf:"varint,12,opt,name=TimeStamp,proto3,oneof" json:"TimeStamp,omitempty"`
	Language   *string               `protobuf:"bytes,13,opt,name=Language,proto3,oneof" json:"Language,omitempty"`
	APIVersion *string               `protobuf:"bytes,14,opt,name=APIVersion,proto3,oneof" json:"APIVersion,omitempty"`
	Method     *string               `protobuf:"bytes,15,opt,name=Method,proto3,oneof" json:"Method,omitempty"`
	UID        *string               `protobuf:"bytes,16,opt,name=UID,proto3,oneof" json:"UID,omitempty"`
	IsAdmin    *int32                `protobuf:"varint,17,opt,name=IsAdmin,proto3,oneof" json:"IsAdmin,omitempty"`
	SessionEnd *int32                `protobuf:"varint,18,opt,name=SessionEnd,proto3,oneof" json:"SessionEnd,omitempty"`
	Completed  *int32                `protobuf:"varint,19,opt,name=Completed,proto3,oneof" json:"Completed,omitempty"`
	Readonly   *int32                `protobuf:"varint,20,opt,name=Readonly,proto3,oneof" json:"Readonly,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_microservice_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_microservice_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_microservice_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetModule() string {
	if x != nil && x.Module != nil {
		return *x.Module
	}
	return ""
}

func (x *Request) GetParam() string {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return ""
}

func (x *Request) GetParamID() string {
	if x != nil && x.ParamID != nil {
		return *x.ParamID
	}
	return ""
}

func (x *Request) GetAction() string {
	if x != nil && x.Action != nil {
		return *x.Action
	}
	return ""
}

func (x *Request) GetArgs() map[string]*anypb.Any {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *Request) GetPath() string {
	if x != nil && x.Path != nil {
		return *x.Path
	}
	return ""
}

func (x *Request) GetToken() string {
	if x != nil && x.Token != nil {
		return *x.Token
	}
	return ""
}

func (x *Request) GetSign() string {
	if x != nil && x.Sign != nil {
		return *x.Sign
	}
	return ""
}

func (x *Request) GetIP() string {
	if x != nil && x.IP != nil {
		return *x.IP
	}
	return ""
}

func (x *Request) GetUserAgent() string {
	if x != nil && x.UserAgent != nil {
		return *x.UserAgent
	}
	return ""
}

func (x *Request) GetTokenType() string {
	if x != nil && x.TokenType != nil {
		return *x.TokenType
	}
	return ""
}

func (x *Request) GetTimeStamp() int32 {
	if x != nil && x.TimeStamp != nil {
		return *x.TimeStamp
	}
	return 0
}

func (x *Request) GetLanguage() string {
	if x != nil && x.Language != nil {
		return *x.Language
	}
	return ""
}

func (x *Request) GetAPIVersion() string {
	if x != nil && x.APIVersion != nil {
		return *x.APIVersion
	}
	return ""
}

func (x *Request) GetMethod() string {
	if x != nil && x.Method != nil {
		return *x.Method
	}
	return ""
}

func (x *Request) GetUID() string {
	if x != nil && x.UID != nil {
		return *x.UID
	}
	return ""
}

func (x *Request) GetIsAdmin() int32 {
	if x != nil && x.IsAdmin != nil {
		return *x.IsAdmin
	}
	return 0
}

func (x *Request) GetSessionEnd() int32 {
	if x != nil && x.SessionEnd != nil {
		return *x.SessionEnd
	}
	return 0
}

func (x *Request) GetCompleted() int32 {
	if x != nil && x.Completed != nil {
		return *x.Completed
	}
	return 0
}

func (x *Request) GetReadonly() int32 {
	if x != nil && x.Readonly != nil {
		return *x.Readonly
	}
	return 0
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data        map[string]*anypb.Any `protobuf:"bytes,1,rep,name=Data,proto3" json:"Data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RequestBack *Request              `protobuf:"bytes,2,opt,name=RequestBack,proto3,oneof" json:"RequestBack,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_microservice_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_microservice_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_microservice_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetData() map[string]*anypb.Any {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Response) GetRequestBack() *Request {
	if x != nil {
		return x.RequestBack
	}
	return nil
}

var File_microservice_proto protoreflect.FileDescriptor

var file_microservice_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6d, 0x69, 0x63, 0x72, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x9f, 0x07, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x06, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x4d,
	0x6f, 0x64, 0x75, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x49, 0x44, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x07, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x49, 0x44, 0x88,
	0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x03, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12,
	0x26, 0x0a, 0x04, 0x41, 0x72, 0x67, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x72, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x04, 0x41, 0x72, 0x67, 0x73, 0x12, 0x17, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x04, 0x50, 0x61, 0x74, 0x68, 0x88, 0x01, 0x01,
	0x12, 0x19, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x05, 0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x53,
	0x69, 0x67, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x48, 0x06, 0x52, 0x04, 0x53, 0x69, 0x67,
	0x6e, 0x88, 0x01, 0x01, 0x12, 0x13, 0x0a, 0x02, 0x49, 0x50, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x07, 0x52, 0x02, 0x49, 0x50, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09, 0x55, 0x73, 0x65,
	0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x09,
	0x55, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01, 0x12, 0x21, 0x0a, 0x09,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x09, 0x52, 0x09, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x21, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x0a, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x88,
	0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52, 0x08, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0c, 0x52, 0x0a, 0x41, 0x50, 0x49, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0d, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x55, 0x49, 0x44, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x0e, 0x52, 0x03, 0x55, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x1d, 0x0a, 0x07,
	0x49, 0x73, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x11, 0x20, 0x01, 0x28, 0x05, 0x48, 0x0f, 0x52,
	0x07, 0x49, 0x73, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x53,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05, 0x48,
	0x10, 0x52, 0x0a, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x88, 0x01, 0x01,
	0x12, 0x21, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x13, 0x20,
	0x01, 0x28, 0x05, 0x48, 0x11, 0x52, 0x09, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x52, 0x65, 0x61, 0x64, 0x6f, 0x6e, 0x6c, 0x79, 0x18,
	0x14, 0x20, 0x01, 0x28, 0x05, 0x48, 0x12, 0x52, 0x08, 0x52, 0x65, 0x61, 0x64, 0x6f, 0x6e, 0x6c,
	0x79, 0x88, 0x01, 0x01, 0x1a, 0x4d, 0x0a, 0x09, 0x41, 0x72, 0x67, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x49, 0x44, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x07, 0x0a, 0x05, 0x5f, 0x50, 0x61, 0x74, 0x68, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x53, 0x69, 0x67, 0x6e, 0x42, 0x05, 0x0a, 0x03, 0x5f,
	0x49, 0x50, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x42, 0x0c,
	0x0a, 0x0a, 0x5f, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0b, 0x0a, 0x09,
	0x5f, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x41, 0x50,
	0x49, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x55, 0x49, 0x44, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x49, 0x73, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x53, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x64, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x52, 0x65, 0x61, 0x64, 0x6f, 0x6e, 0x6c,
	0x79, 0x22, 0xc3, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27,
	0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2f, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x61, 0x63, 0x6b, 0x88, 0x01, 0x01, 0x1a, 0x4d, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x32, 0x26, 0x0a, 0x07, 0x52, 0x65, 0x76, 0x65, 0x72,
	0x73, 0x65, 0x12, 0x1b, 0x0a, 0x02, 0x44, 0x6f, 0x12, 0x08, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x09, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_microservice_proto_rawDescOnce sync.Once
	file_microservice_proto_rawDescData = file_microservice_proto_rawDesc
)

func file_microservice_proto_rawDescGZIP() []byte {
	file_microservice_proto_rawDescOnce.Do(func() {
		file_microservice_proto_rawDescData = protoimpl.X.CompressGZIP(file_microservice_proto_rawDescData)
	})
	return file_microservice_proto_rawDescData
}

var file_microservice_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_microservice_proto_goTypes = []interface{}{
	(*Request)(nil),   // 0: Request
	(*Response)(nil),  // 1: Response
	nil,               // 2: Request.ArgsEntry
	nil,               // 3: Response.DataEntry
	(*anypb.Any)(nil), // 4: google.protobuf.Any
}
var file_microservice_proto_depIdxs = []int32{
	2, // 0: Request.Args:type_name -> Request.ArgsEntry
	3, // 1: Response.Data:type_name -> Response.DataEntry
	0, // 2: Response.RequestBack:type_name -> Request
	4, // 3: Request.ArgsEntry.value:type_name -> google.protobuf.Any
	4, // 4: Response.DataEntry.value:type_name -> google.protobuf.Any
	0, // 5: Reverse.Do:input_type -> Request
	1, // 6: Reverse.Do:output_type -> Response
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_microservice_proto_init() }
func file_microservice_proto_init() {
	if File_microservice_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_microservice_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_microservice_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_microservice_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_microservice_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_microservice_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_microservice_proto_goTypes,
		DependencyIndexes: file_microservice_proto_depIdxs,
		MessageInfos:      file_microservice_proto_msgTypes,
	}.Build()
	File_microservice_proto = out.File
	file_microservice_proto_rawDesc = nil
	file_microservice_proto_goTypes = nil
	file_microservice_proto_depIdxs = nil
}
