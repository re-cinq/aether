// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: proto/plugin.proto

package proto

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

type ResourceEmissions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=value,proto3" json:"value,omitempty"`
	Unit  string  `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *ResourceEmissions) Reset() {
	*x = ResourceEmissions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceEmissions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceEmissions) ProtoMessage() {}

func (x *ResourceEmissions) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceEmissions.ProtoReflect.Descriptor instead.
func (*ResourceEmissions) Descriptor() ([]byte, []int) {
	return file_proto_plugin_proto_rawDescGZIP(), []int{0}
}

func (x *ResourceEmissions) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *ResourceEmissions) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Usage        float64            `protobuf:"fixed64,2,opt,name=usage,proto3" json:"usage,omitempty"`
	UnitAmount   float64            `protobuf:"fixed64,3,opt,name=unit_amount,json=unitAmount,proto3" json:"unit_amount,omitempty"`
	Emissions    *ResourceEmissions `protobuf:"bytes,4,opt,name=emissions,proto3" json:"emissions,omitempty"`
	Labels       map[string]string  `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Unit         string             `protobuf:"bytes,6,opt,name=unit,proto3" json:"unit,omitempty"`
	UpdatedAt    int64              `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	ResourceType string             `protobuf:"bytes,8,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	Energy       float64            `protobuf:"fixed64,9,opt,name=energy,proto3" json:"energy,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_plugin_proto_rawDescGZIP(), []int{1}
}

func (x *Metric) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metric) GetUsage() float64 {
	if x != nil {
		return x.Usage
	}
	return 0
}

func (x *Metric) GetUnitAmount() float64 {
	if x != nil {
		return x.UnitAmount
	}
	return 0
}

func (x *Metric) GetEmissions() *ResourceEmissions {
	if x != nil {
		return x.Emissions
	}
	return nil
}

func (x *Metric) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *Metric) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *Metric) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Metric) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *Metric) GetEnergy() float64 {
	if x != nil {
		return x.Energy
	}
	return 0
}

type InstanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                string             `protobuf:"bytes,11,opt,name=id,proto3" json:"id,omitempty"`
	Provider          string             `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	Service           string             `protobuf:"bytes,2,opt,name=service,proto3" json:"service,omitempty"`
	Name              string             `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Region            string             `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
	Zone              string             `protobuf:"bytes,5,opt,name=zone,proto3" json:"zone,omitempty"`
	Kind              string             `protobuf:"bytes,6,opt,name=kind,proto3" json:"kind,omitempty"`
	EmbodiedEmissions *ResourceEmissions `protobuf:"bytes,8,opt,name=EmbodiedEmissions,proto3" json:"EmbodiedEmissions,omitempty"`
	Metrics           map[string]*Metric `protobuf:"bytes,9,rep,name=metrics,proto3" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Labels            map[string]string  `protobuf:"bytes,10,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Status            string             `protobuf:"bytes,12,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *InstanceRequest) Reset() {
	*x = InstanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceRequest) ProtoMessage() {}

func (x *InstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceRequest.ProtoReflect.Descriptor instead.
func (*InstanceRequest) Descriptor() ([]byte, []int) {
	return file_proto_plugin_proto_rawDescGZIP(), []int{2}
}

func (x *InstanceRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *InstanceRequest) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *InstanceRequest) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *InstanceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InstanceRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *InstanceRequest) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

func (x *InstanceRequest) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *InstanceRequest) GetEmbodiedEmissions() *ResourceEmissions {
	if x != nil {
		return x.EmbodiedEmissions
	}
	return nil
}

func (x *InstanceRequest) GetMetrics() map[string]*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *InstanceRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *InstanceRequest) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type ListInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instances []*InstanceRequest `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
}

func (x *ListInstanceResponse) Reset() {
	*x = ListInstanceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInstanceResponse) ProtoMessage() {}

func (x *ListInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInstanceResponse.ProtoReflect.Descriptor instead.
func (*ListInstanceResponse) Descriptor() ([]byte, []int) {
	return file_proto_plugin_proto_rawDescGZIP(), []int{3}
}

func (x *ListInstanceResponse) GetInstances() []*InstanceRequest {
	if x != nil {
		return x.Instances
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_plugin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_plugin_proto_msgTypes[4]
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
	return file_proto_plugin_proto_rawDescGZIP(), []int{4}
}

var File_proto_plugin_proto protoreflect.FileDescriptor

var file_proto_plugin_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x11, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x22, 0xe9, 0x02, 0x0a, 0x06, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x75, 0x6e, 0x69, 0x74, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x75, 0x6e, 0x69, 0x74, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x36, 0x0a, 0x09, 0x65, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x09, 0x65,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x31, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x23,
	0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x06, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x1a, 0x39, 0x0a, 0x0b, 0x4c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8c, 0x04, 0x0a, 0x0f, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x7a, 0x6f, 0x6e, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x46, 0x0a, 0x11, 0x45, 0x6d, 0x62, 0x6f, 0x64, 0x69, 0x65, 0x64,
	0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x11, 0x45, 0x6d, 0x62, 0x6f, 0x64,
	0x69, 0x65, 0x64, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3d, 0x0a, 0x07,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x3a, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a,
	0x49, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x23, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4c, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a,
	0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x73, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x38, 0x0a, 0x08,
	0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64,
	0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0x60, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x32, 0x0a, 0x05, 0x46, 0x65, 0x74, 0x63, 0x68, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x0c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_plugin_proto_rawDescOnce sync.Once
	file_proto_plugin_proto_rawDescData = file_proto_plugin_proto_rawDesc
)

func file_proto_plugin_proto_rawDescGZIP() []byte {
	file_proto_plugin_proto_rawDescOnce.Do(func() {
		file_proto_plugin_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_plugin_proto_rawDescData)
	})
	return file_proto_plugin_proto_rawDescData
}

var file_proto_plugin_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_plugin_proto_goTypes = []interface{}{
	(*ResourceEmissions)(nil),    // 0: proto.ResourceEmissions
	(*Metric)(nil),               // 1: proto.Metric
	(*InstanceRequest)(nil),      // 2: proto.InstanceRequest
	(*ListInstanceResponse)(nil), // 3: proto.ListInstanceResponse
	(*Empty)(nil),                // 4: proto.Empty
	nil,                          // 5: proto.Metric.LabelsEntry
	nil,                          // 6: proto.InstanceRequest.MetricsEntry
	nil,                          // 7: proto.InstanceRequest.LabelsEntry
}
var file_proto_plugin_proto_depIdxs = []int32{
	0,  // 0: proto.Metric.emissions:type_name -> proto.ResourceEmissions
	5,  // 1: proto.Metric.labels:type_name -> proto.Metric.LabelsEntry
	0,  // 2: proto.InstanceRequest.EmbodiedEmissions:type_name -> proto.ResourceEmissions
	6,  // 3: proto.InstanceRequest.metrics:type_name -> proto.InstanceRequest.MetricsEntry
	7,  // 4: proto.InstanceRequest.labels:type_name -> proto.InstanceRequest.LabelsEntry
	2,  // 5: proto.ListInstanceResponse.instances:type_name -> proto.InstanceRequest
	1,  // 6: proto.InstanceRequest.MetricsEntry.value:type_name -> proto.Metric
	2,  // 7: proto.Exporter.Send:input_type -> proto.InstanceRequest
	4,  // 8: proto.Source.Fetch:input_type -> proto.Empty
	4,  // 9: proto.Source.Stop:input_type -> proto.Empty
	4,  // 10: proto.Exporter.Send:output_type -> proto.Empty
	3,  // 11: proto.Source.Fetch:output_type -> proto.ListInstanceResponse
	4,  // 12: proto.Source.Stop:output_type -> proto.Empty
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_proto_plugin_proto_init() }
func file_proto_plugin_proto_init() {
	if File_proto_plugin_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_plugin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceEmissions); i {
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
		file_proto_plugin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_proto_plugin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceRequest); i {
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
		file_proto_plugin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListInstanceResponse); i {
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
		file_proto_plugin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_proto_plugin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_proto_plugin_proto_goTypes,
		DependencyIndexes: file_proto_plugin_proto_depIdxs,
		MessageInfos:      file_proto_plugin_proto_msgTypes,
	}.Build()
	File_proto_plugin_proto = out.File
	file_proto_plugin_proto_rawDesc = nil
	file_proto_plugin_proto_goTypes = nil
	file_proto_plugin_proto_depIdxs = nil
}
