// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: server/state.proto

package generated

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PodId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PodId) Reset() {
	*x = PodId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodId) ProtoMessage() {}

func (x *PodId) ProtoReflect() protoreflect.Message {
	mi := &file_server_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodId.ProtoReflect.Descriptor instead.
func (*PodId) Descriptor() ([]byte, []int) {
	return file_server_state_proto_rawDescGZIP(), []int{0}
}

func (x *PodId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	State       []byte `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	ContentType string `protobuf:"bytes,3,opt,name=contentType,proto3" json:"contentType,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_state_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_server_state_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_server_state_proto_rawDescGZIP(), []int{1}
}

func (x *State) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *State) GetState() []byte {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *State) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

type RestoreSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromId string `protobuf:"bytes,1,opt,name=fromId,proto3" json:"fromId,omitempty"`
	Ip     string `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Mode   string `protobuf:"bytes,3,opt,name=mode,proto3" json:"mode,omitempty"`
	Path   string `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Port   int32  `protobuf:"varint,5,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *RestoreSpec) Reset() {
	*x = RestoreSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_state_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreSpec) ProtoMessage() {}

func (x *RestoreSpec) ProtoReflect() protoreflect.Message {
	mi := &file_server_state_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreSpec.ProtoReflect.Descriptor instead.
func (*RestoreSpec) Descriptor() ([]byte, []int) {
	return file_server_state_proto_rawDescGZIP(), []int{2}
}

func (x *RestoreSpec) GetFromId() string {
	if x != nil {
		return x.FromId
	}
	return ""
}

func (x *RestoreSpec) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *RestoreSpec) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *RestoreSpec) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *RestoreSpec) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_state_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_server_state_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_server_state_proto_rawDescGZIP(), []int{3}
}

var File_server_state_proto protoreflect.FileDescriptor

var file_server_state_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x17, 0x0a, 0x05, 0x50, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4f, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x22, 0x71, 0x0a, 0x0b, 0x52, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x53, 0x70, 0x65, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x72, 0x6f, 0x6d,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x72, 0x6f, 0x6d, 0x49, 0x64,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70,
	0x12, 0x12, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6d, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x07, 0x0a, 0x05,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xa4, 0x01, 0x0a, 0x06, 0x48, 0x65, 0x6c, 0x70, 0x65, 0x72,
	0x12, 0x30, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x70,
	0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x1a, 0x10,
	0x2e, 0x70, 0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x35, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x16, 0x2e,
	0x70, 0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x53, 0x70, 0x65, 0x63, 0x1a, 0x10, 0x2e, 0x70, 0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x09, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x12, 0x10, 0x2e, 0x70, 0x75, 0x6c, 0x6c, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x2e, 0x50, 0x6f, 0x64, 0x49, 0x64, 0x1a, 0x10, 0x2e, 0x70, 0x75, 0x6c, 0x6c, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b,
	0x2e, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_server_state_proto_rawDescOnce sync.Once
	file_server_state_proto_rawDescData = file_server_state_proto_rawDesc
)

func file_server_state_proto_rawDescGZIP() []byte {
	file_server_state_proto_rawDescOnce.Do(func() {
		file_server_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_state_proto_rawDescData)
	})
	return file_server_state_proto_rawDescData
}

var file_server_state_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_server_state_proto_goTypes = []interface{}{
	(*PodId)(nil),       // 0: pullstate.PodId
	(*State)(nil),       // 1: pullstate.State
	(*RestoreSpec)(nil), // 2: pullstate.RestoreSpec
	(*Empty)(nil),       // 3: pullstate.Empty
}
var file_server_state_proto_depIdxs = []int32{
	1, // 0: pullstate.Helper.NewState:input_type -> pullstate.State
	2, // 1: pullstate.Helper.Restore:input_type -> pullstate.RestoreSpec
	0, // 2: pullstate.Helper.DeletePod:input_type -> pullstate.PodId
	3, // 3: pullstate.Helper.NewState:output_type -> pullstate.Empty
	3, // 4: pullstate.Helper.Restore:output_type -> pullstate.Empty
	3, // 5: pullstate.Helper.DeletePod:output_type -> pullstate.Empty
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_server_state_proto_init() }
func file_server_state_proto_init() {
	if File_server_state_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodId); i {
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
		file_server_state_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_server_state_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreSpec); i {
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
		file_server_state_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_server_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_server_state_proto_goTypes,
		DependencyIndexes: file_server_state_proto_depIdxs,
		MessageInfos:      file_server_state_proto_msgTypes,
	}.Build()
	File_server_state_proto = out.File
	file_server_state_proto_rawDesc = nil
	file_server_state_proto_goTypes = nil
	file_server_state_proto_depIdxs = nil
}