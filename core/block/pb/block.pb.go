// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: github.com/yahuizhan/dappley-metrics-go-api/core/block/pb/block.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	pb "github.com/yahuizhan/dappley-metrics-go-api/core/transaction/pb"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header       *BlockHeader      `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Transactions []*pb.Transaction `protobuf:"bytes,2,rep,name=transactions,proto3" json:"transactions,omitempty"`
	ParentHash   []byte            `protobuf:"bytes,3,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetHeader() *BlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Block) GetTransactions() []*pb.Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Block) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

type BlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash         []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	PreviousHash []byte `protobuf:"bytes,2,opt,name=previous_hash,json=previousHash,proto3" json:"previous_hash,omitempty"`
	Nonce        int64  `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Timestamp    int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature    []byte `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
	Height       uint64 `protobuf:"varint,6,opt,name=height,proto3" json:"height,omitempty"`
	Producer     string `protobuf:"bytes,7,opt,name=producer,proto3" json:"producer,omitempty"`
}

func (x *BlockHeader) Reset() {
	*x = BlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockHeader) ProtoMessage() {}

func (x *BlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockHeader.ProtoReflect.Descriptor instead.
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescGZIP(), []int{1}
}

func (x *BlockHeader) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *BlockHeader) GetPreviousHash() []byte {
	if x != nil {
		return x.PreviousHash
	}
	return nil
}

func (x *BlockHeader) GetNonce() int64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *BlockHeader) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BlockHeader) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *BlockHeader) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *BlockHeader) GetProducer() string {
	if x != nil {
		return x.Producer
	}
	return ""
}

var File_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto protoreflect.FileDescriptor

var file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDesc = []byte{
	0x0a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x61, 0x68,
	0x75, 0x69, 0x7a, 0x68, 0x61, 0x6e, 0x2f, 0x64, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x79, 0x2d, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2d, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x62, 0x2f, 0x62, 0x6c, 0x6f, 0x63,
	0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x70, 0x62,
	0x1a, 0x51, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x61, 0x68,
	0x75, 0x69, 0x7a, 0x68, 0x61, 0x6e, 0x2f, 0x64, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x79, 0x2d, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2d, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70,
	0x62, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x96, 0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x2c, 0x0a,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x70, 0x62, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70,
	0x62, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x70,
	0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0a, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x22, 0xcc, 0x01, 0x0a,
	0x0b, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75,
	0x73, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x42, 0x3b, 0x5a, 0x39, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x61, 0x68, 0x75, 0x69, 0x7a,
	0x68, 0x61, 0x6e, 0x2f, 0x64, 0x61, 0x70, 0x70, 0x6c, 0x65, 0x79, 0x2d, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x2d, 0x67, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescOnce sync.Once
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescData = file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDesc
)

func file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescGZIP() []byte {
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescOnce.Do(func() {
		file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescData)
	})
	return file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDescData
}

var file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_goTypes = []interface{}{
	(*Block)(nil),          // 0: blockpb.Block
	(*BlockHeader)(nil),    // 1: blockpb.BlockHeader
	(*pb.Transaction)(nil), // 2: transactionpb.Transaction
}
var file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_depIdxs = []int32{
	1, // 0: blockpb.Block.header:type_name -> blockpb.BlockHeader
	2, // 1: blockpb.Block.transactions:type_name -> transactionpb.Transaction
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_init() }
func file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_init() {
	if File_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockHeader); i {
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
			RawDescriptor: file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_goTypes,
		DependencyIndexes: file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_depIdxs,
		MessageInfos:      file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_msgTypes,
	}.Build()
	File_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto = out.File
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_rawDesc = nil
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_goTypes = nil
	file_github_com_yahuizhan_dappley_metrics_go_api_core_block_pb_block_proto_depIdxs = nil
}
