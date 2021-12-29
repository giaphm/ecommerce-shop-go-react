// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.1
// source: products.proto

package products

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type IsProductAvailableRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductUuid string `protobuf:"bytes,1,opt,name=product_uuid,json=productUuid,proto3" json:"product_uuid,omitempty"`
}

func (x *IsProductAvailableRequest) Reset() {
	*x = IsProductAvailableRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_products_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsProductAvailableRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsProductAvailableRequest) ProtoMessage() {}

func (x *IsProductAvailableRequest) ProtoReflect() protoreflect.Message {
	mi := &file_products_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsProductAvailableRequest.ProtoReflect.Descriptor instead.
func (*IsProductAvailableRequest) Descriptor() ([]byte, []int) {
	return file_products_proto_rawDescGZIP(), []int{0}
}

func (x *IsProductAvailableRequest) GetProductUuid() string {
	if x != nil {
		return x.ProductUuid
	}
	return ""
}

type IsProductAvailableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsAvailable bool `protobuf:"varint,1,opt,name=is_available,json=isAvailable,proto3" json:"is_available,omitempty"`
}

func (x *IsProductAvailableResponse) Reset() {
	*x = IsProductAvailableResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_products_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsProductAvailableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsProductAvailableResponse) ProtoMessage() {}

func (x *IsProductAvailableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_products_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsProductAvailableResponse.ProtoReflect.Descriptor instead.
func (*IsProductAvailableResponse) Descriptor() ([]byte, []int) {
	return file_products_proto_rawDescGZIP(), []int{1}
}

func (x *IsProductAvailableResponse) GetIsAvailable() bool {
	if x != nil {
		return x.IsAvailable
	}
	return false
}

type UpdateProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductUuid string `protobuf:"bytes,1,opt,name=product_uuid,json=productUuid,proto3" json:"product_uuid,omitempty"`
}

func (x *UpdateProductRequest) Reset() {
	*x = UpdateProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_products_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductRequest) ProtoMessage() {}

func (x *UpdateProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_products_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductRequest.ProtoReflect.Descriptor instead.
func (*UpdateProductRequest) Descriptor() ([]byte, []int) {
	return file_products_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateProductRequest) GetProductUuid() string {
	if x != nil {
		return x.ProductUuid
	}
	return ""
}

type EmptyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyResponse) Reset() {
	*x = EmptyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_products_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyResponse) ProtoMessage() {}

func (x *EmptyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_products_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyResponse.ProtoReflect.Descriptor instead.
func (*EmptyResponse) Descriptor() ([]byte, []int) {
	return file_products_proto_rawDescGZIP(), []int{3}
}

var File_products_proto protoreflect.FileDescriptor

var file_products_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x22, 0x3e, 0x0a, 0x19, 0x49, 0x73,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x55, 0x75, 0x69, 0x64, 0x22, 0x3f, 0x0a, 0x1a, 0x49, 0x73,
	0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x69, 0x73, 0x5f, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b,
	0x69, 0x73, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x39, 0x0a, 0x14, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x55, 0x75, 0x69, 0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xbe, 0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x61, 0x0a, 0x12, 0x49,
	0x73, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x2e, 0x49, 0x73, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x73, 0x2e, 0x49, 0x73, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x41, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48,
	0x0a, 0x0b, 0x53, 0x65, 0x6c, 0x6c, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x1e, 0x2e,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_products_proto_rawDescOnce sync.Once
	file_products_proto_rawDescData = file_products_proto_rawDesc
)

func file_products_proto_rawDescGZIP() []byte {
	file_products_proto_rawDescOnce.Do(func() {
		file_products_proto_rawDescData = protoimpl.X.CompressGZIP(file_products_proto_rawDescData)
	})
	return file_products_proto_rawDescData
}

var file_products_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_products_proto_goTypes = []interface{}{
	(*IsProductAvailableRequest)(nil),  // 0: products.IsProductAvailableRequest
	(*IsProductAvailableResponse)(nil), // 1: products.IsProductAvailableResponse
	(*UpdateProductRequest)(nil),       // 2: products.UpdateProductRequest
	(*EmptyResponse)(nil),              // 3: products.EmptyResponse
}
var file_products_proto_depIdxs = []int32{
	0, // 0: products.ProductsService.IsProductAvailable:input_type -> products.IsProductAvailableRequest
	2, // 1: products.ProductsService.SellProduct:input_type -> products.UpdateProductRequest
	1, // 2: products.ProductsService.IsProductAvailable:output_type -> products.IsProductAvailableResponse
	3, // 3: products.ProductsService.SellProduct:output_type -> products.EmptyResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_products_proto_init() }
func file_products_proto_init() {
	if File_products_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_products_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsProductAvailableRequest); i {
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
		file_products_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsProductAvailableResponse); i {
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
		file_products_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductRequest); i {
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
		file_products_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyResponse); i {
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
			RawDescriptor: file_products_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_products_proto_goTypes,
		DependencyIndexes: file_products_proto_depIdxs,
		MessageInfos:      file_products_proto_msgTypes,
	}.Build()
	File_products_proto = out.File
	file_products_proto_rawDesc = nil
	file_products_proto_goTypes = nil
	file_products_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ProductsServiceClient is the client API for ProductsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProductsServiceClient interface {
	IsProductAvailable(ctx context.Context, in *IsProductAvailableRequest, opts ...grpc.CallOption) (*IsProductAvailableResponse, error)
	SellProduct(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type productsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductsServiceClient(cc grpc.ClientConnInterface) ProductsServiceClient {
	return &productsServiceClient{cc}
}

func (c *productsServiceClient) IsProductAvailable(ctx context.Context, in *IsProductAvailableRequest, opts ...grpc.CallOption) (*IsProductAvailableResponse, error) {
	out := new(IsProductAvailableResponse)
	err := c.cc.Invoke(ctx, "/products.ProductsService/IsProductAvailable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productsServiceClient) SellProduct(ctx context.Context, in *UpdateProductRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/products.ProductsService/SellProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductsServiceServer is the server API for ProductsService service.
type ProductsServiceServer interface {
	IsProductAvailable(context.Context, *IsProductAvailableRequest) (*IsProductAvailableResponse, error)
	SellProduct(context.Context, *UpdateProductRequest) (*EmptyResponse, error)
}

// UnimplementedProductsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedProductsServiceServer struct {
}

func (*UnimplementedProductsServiceServer) IsProductAvailable(context.Context, *IsProductAvailableRequest) (*IsProductAvailableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsProductAvailable not implemented")
}
func (*UnimplementedProductsServiceServer) SellProduct(context.Context, *UpdateProductRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SellProduct not implemented")
}

func RegisterProductsServiceServer(s *grpc.Server, srv ProductsServiceServer) {
	s.RegisterService(&_ProductsService_serviceDesc, srv)
}

func _ProductsService_IsProductAvailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsProductAvailableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).IsProductAvailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/products.ProductsService/IsProductAvailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).IsProductAvailable(ctx, req.(*IsProductAvailableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductsService_SellProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductsServiceServer).SellProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/products.ProductsService/SellProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductsServiceServer).SellProduct(ctx, req.(*UpdateProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProductsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "products.ProductsService",
	HandlerType: (*ProductsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsProductAvailable",
			Handler:    _ProductsService_IsProductAvailable_Handler,
		},
		{
			MethodName: "SellProduct",
			Handler:    _ProductsService_SellProduct_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "products.proto",
}
