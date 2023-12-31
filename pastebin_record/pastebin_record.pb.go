// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: pastebin_record/pastebin_record.proto

package pastebin_record

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IPAddressVersion int32

const (
	IPAddressVersion_V4 IPAddressVersion = 0
	IPAddressVersion_V6 IPAddressVersion = 1
)

// Enum value maps for IPAddressVersion.
var (
	IPAddressVersion_name = map[int32]string{
		0: "V4",
		1: "V6",
	}
	IPAddressVersion_value = map[string]int32{
		"V4": 0,
		"V6": 1,
	}
)

func (x IPAddressVersion) Enum() *IPAddressVersion {
	p := new(IPAddressVersion)
	*p = x
	return p
}

func (x IPAddressVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IPAddressVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_pastebin_record_pastebin_record_proto_enumTypes[0].Descriptor()
}

func (IPAddressVersion) Type() protoreflect.EnumType {
	return &file_pastebin_record_pastebin_record_proto_enumTypes[0]
}

func (x IPAddressVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IPAddressVersion.Descriptor instead.
func (IPAddressVersion) EnumDescriptor() ([]byte, []int) {
	return file_pastebin_record_pastebin_record_proto_rawDescGZIP(), []int{0}
}

type IPAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip      []byte           `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Version IPAddressVersion `protobuf:"varint,2,opt,name=version,proto3,enum=pastebin.IPAddressVersion" json:"version,omitempty"`
}

func (x *IPAddress) Reset() {
	*x = IPAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pastebin_record_pastebin_record_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPAddress) ProtoMessage() {}

func (x *IPAddress) ProtoReflect() protoreflect.Message {
	mi := &file_pastebin_record_pastebin_record_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPAddress.ProtoReflect.Descriptor instead.
func (*IPAddress) Descriptor() ([]byte, []int) {
	return file_pastebin_record_pastebin_record_proto_rawDescGZIP(), []int{0}
}

func (x *IPAddress) GetIp() []byte {
	if x != nil {
		return x.Ip
	}
	return nil
}

func (x *IPAddress) GetVersion() IPAddressVersion {
	if x != nil {
		return x.Version
	}
	return IPAddressVersion_V4
}

type PastebinRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body               string                 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`                                  // literal text; body of utf-8 text file
	Creator            *IPAddress             `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`                            // IP address (either v4 or v6) of client that created this paste
	TimeCreated        *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=time_created,json=timeCreated,proto3" json:"time_created,omitempty"` // when this record was inserted into database
	Expiration         *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=expiration,proto3,oneof" json:"expiration,omitempty"`                // when to purge
	Filename           *string                `protobuf:"bytes,5,opt,name=filename,proto3,oneof" json:"filename,omitempty"`
	MimeType           *string                `protobuf:"bytes,6,opt,name=mime_type,json=mimeType,proto3,oneof" json:"mime_type,omitempty"`                               // mime type string of the "file"
	SyntaxHighlighting *string                `protobuf:"bytes,7,opt,name=syntax_highlighting,json=syntaxHighlighting,proto3,oneof" json:"syntax_highlighting,omitempty"` // language for syntax highlighting
}

func (x *PastebinRecord) Reset() {
	*x = PastebinRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pastebin_record_pastebin_record_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PastebinRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PastebinRecord) ProtoMessage() {}

func (x *PastebinRecord) ProtoReflect() protoreflect.Message {
	mi := &file_pastebin_record_pastebin_record_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PastebinRecord.ProtoReflect.Descriptor instead.
func (*PastebinRecord) Descriptor() ([]byte, []int) {
	return file_pastebin_record_pastebin_record_proto_rawDescGZIP(), []int{1}
}

func (x *PastebinRecord) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *PastebinRecord) GetCreator() *IPAddress {
	if x != nil {
		return x.Creator
	}
	return nil
}

func (x *PastebinRecord) GetTimeCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.TimeCreated
	}
	return nil
}

func (x *PastebinRecord) GetExpiration() *timestamppb.Timestamp {
	if x != nil {
		return x.Expiration
	}
	return nil
}

func (x *PastebinRecord) GetFilename() string {
	if x != nil && x.Filename != nil {
		return *x.Filename
	}
	return ""
}

func (x *PastebinRecord) GetMimeType() string {
	if x != nil && x.MimeType != nil {
		return *x.MimeType
	}
	return ""
}

func (x *PastebinRecord) GetSyntaxHighlighting() string {
	if x != nil && x.SyntaxHighlighting != nil {
		return *x.SyntaxHighlighting
	}
	return ""
}

var File_pastebin_record_pastebin_record_proto protoreflect.FileDescriptor

var file_pastebin_record_pastebin_record_proto_rawDesc = []byte{
	0x0a, 0x25, 0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69,
	0x6e, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x51, 0x0a, 0x09, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x70, 0x12,
	0x34, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1a, 0x2e, 0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x2e, 0x49, 0x50, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x8e, 0x03, 0x0a, 0x0e, 0x50, 0x61, 0x73, 0x74, 0x65, 0x62,
	0x69, 0x6e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x2d, 0x0a, 0x07,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x2e, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x3d, 0x0a, 0x0c, 0x74,
	0x69, 0x6d, 0x65, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x74,
	0x69, 0x6d, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x3f, 0x0a, 0x0a, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x0a, 0x65, 0x78,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52,
	0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09,
	0x6d, 0x69, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x02, 0x52, 0x08, 0x6d, 0x69, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x12, 0x34,
	0x0a, 0x13, 0x73, 0x79, 0x6e, 0x74, 0x61, 0x78, 0x5f, 0x68, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x03, 0x52, 0x12, 0x73,
	0x79, 0x6e, 0x74, 0x61, 0x78, 0x48, 0x69, 0x67, 0x68, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x69, 0x6e,
	0x67, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x6d, 0x69, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x16,
	0x0a, 0x14, 0x5f, 0x73, 0x79, 0x6e, 0x74, 0x61, 0x78, 0x5f, 0x68, 0x69, 0x67, 0x68, 0x6c, 0x69,
	0x67, 0x68, 0x74, 0x69, 0x6e, 0x67, 0x2a, 0x22, 0x0a, 0x10, 0x49, 0x50, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x06, 0x0a, 0x02, 0x56, 0x34,
	0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x56, 0x36, 0x10, 0x01, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x70, 0x72, 0x69, 0x65,
	0x74, 0x61, 0x72, 0x79, 0x2f, 0x70, 0x61, 0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x2f, 0x70, 0x61,
	0x73, 0x74, 0x65, 0x62, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pastebin_record_pastebin_record_proto_rawDescOnce sync.Once
	file_pastebin_record_pastebin_record_proto_rawDescData = file_pastebin_record_pastebin_record_proto_rawDesc
)

func file_pastebin_record_pastebin_record_proto_rawDescGZIP() []byte {
	file_pastebin_record_pastebin_record_proto_rawDescOnce.Do(func() {
		file_pastebin_record_pastebin_record_proto_rawDescData = protoimpl.X.CompressGZIP(file_pastebin_record_pastebin_record_proto_rawDescData)
	})
	return file_pastebin_record_pastebin_record_proto_rawDescData
}

var file_pastebin_record_pastebin_record_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pastebin_record_pastebin_record_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pastebin_record_pastebin_record_proto_goTypes = []interface{}{
	(IPAddressVersion)(0),         // 0: pastebin.IPAddressVersion
	(*IPAddress)(nil),             // 1: pastebin.IPAddress
	(*PastebinRecord)(nil),        // 2: pastebin.PastebinRecord
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_pastebin_record_pastebin_record_proto_depIdxs = []int32{
	0, // 0: pastebin.IPAddress.version:type_name -> pastebin.IPAddressVersion
	1, // 1: pastebin.PastebinRecord.creator:type_name -> pastebin.IPAddress
	3, // 2: pastebin.PastebinRecord.time_created:type_name -> google.protobuf.Timestamp
	3, // 3: pastebin.PastebinRecord.expiration:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pastebin_record_pastebin_record_proto_init() }
func file_pastebin_record_pastebin_record_proto_init() {
	if File_pastebin_record_pastebin_record_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pastebin_record_pastebin_record_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPAddress); i {
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
		file_pastebin_record_pastebin_record_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PastebinRecord); i {
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
	file_pastebin_record_pastebin_record_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pastebin_record_pastebin_record_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pastebin_record_pastebin_record_proto_goTypes,
		DependencyIndexes: file_pastebin_record_pastebin_record_proto_depIdxs,
		EnumInfos:         file_pastebin_record_pastebin_record_proto_enumTypes,
		MessageInfos:      file_pastebin_record_pastebin_record_proto_msgTypes,
	}.Build()
	File_pastebin_record_pastebin_record_proto = out.File
	file_pastebin_record_pastebin_record_proto_rawDesc = nil
	file_pastebin_record_pastebin_record_proto_goTypes = nil
	file_pastebin_record_pastebin_record_proto_depIdxs = nil
}
