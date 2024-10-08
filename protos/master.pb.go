// Copyright 2015 gRPC authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.0
// source: master.proto

package protos

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

type MapperInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *MapperInfo) Reset() {
	*x = MapperInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperInfo) ProtoMessage() {}

func (x *MapperInfo) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperInfo.ProtoReflect.Descriptor instead.
func (*MapperInfo) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{0}
}

func (x *MapperInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type MapperRegisterReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *MapperRegisterReply) Reset() {
	*x = MapperRegisterReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MapperRegisterReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MapperRegisterReply) ProtoMessage() {}

func (x *MapperRegisterReply) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MapperRegisterReply.ProtoReflect.Descriptor instead.
func (*MapperRegisterReply) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{1}
}

func (x *MapperRegisterReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type ReplicationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seq              int64               `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	UserRequestQueue []*ContainerRequest `protobuf:"bytes,2,rep,name=userRequestQueue,proto3" json:"userRequestQueue,omitempty"`
	WorkerPool       *WorkerPool         `protobuf:"bytes,3,opt,name=workerPool,proto3" json:"workerPool,omitempty"`
	TaskPool         *TaskPool           `protobuf:"bytes,4,opt,name=taskPool,proto3" json:"taskPool,omitempty"`
}

func (x *ReplicationRequest) Reset() {
	*x = ReplicationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicationRequest) ProtoMessage() {}

func (x *ReplicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicationRequest.ProtoReflect.Descriptor instead.
func (*ReplicationRequest) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{2}
}

func (x *ReplicationRequest) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *ReplicationRequest) GetUserRequestQueue() []*ContainerRequest {
	if x != nil {
		return x.UserRequestQueue
	}
	return nil
}

func (x *ReplicationRequest) GetWorkerPool() *WorkerPool {
	if x != nil {
		return x.WorkerPool
	}
	return nil
}

func (x *ReplicationRequest) GetTaskPool() *TaskPool {
	if x != nil {
		return x.TaskPool
	}
	return nil
}

type ReplicationReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *ReplicationReply) Reset() {
	*x = ReplicationReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicationReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicationReply) ProtoMessage() {}

func (x *ReplicationReply) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicationReply.ProtoReflect.Descriptor instead.
func (*ReplicationReply) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{3}
}

func (x *ReplicationReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type ContainerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerName string `protobuf:"bytes,1,opt,name=containerName,proto3" json:"containerName,omitempty"`
	M             int32  `protobuf:"varint,2,opt,name=m,proto3" json:"m,omitempty"`
	R             int32  `protobuf:"varint,3,opt,name=r,proto3" json:"r,omitempty"`
}

func (x *ContainerRequest) Reset() {
	*x = ContainerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerRequest) ProtoMessage() {}

func (x *ContainerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerRequest.ProtoReflect.Descriptor instead.
func (*ContainerRequest) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{4}
}

func (x *ContainerRequest) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *ContainerRequest) GetM() int32 {
	if x != nil {
		return x.M
	}
	return 0
}

func (x *ContainerRequest) GetR() int32 {
	if x != nil {
		return x.R
	}
	return 0
}

type WorkerPool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queue           []string               `protobuf:"bytes,1,rep,name=queue,proto3" json:"queue,omitempty"`
	State           map[string]*WorkerInfo `protobuf:"bytes,2,rep,name=state,proto3" json:"state,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DeadWorkerQueue []string               `protobuf:"bytes,3,rep,name=deadWorkerQueue,proto3" json:"deadWorkerQueue,omitempty"`
}

func (x *WorkerPool) Reset() {
	*x = WorkerPool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerPool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerPool) ProtoMessage() {}

func (x *WorkerPool) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerPool.ProtoReflect.Descriptor instead.
func (*WorkerPool) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{5}
}

func (x *WorkerPool) GetQueue() []string {
	if x != nil {
		return x.Queue
	}
	return nil
}

func (x *WorkerPool) GetState() map[string]*WorkerInfo {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *WorkerPool) GetDeadWorkerQueue() []string {
	if x != nil {
		return x.DeadWorkerQueue
	}
	return nil
}

type WorkerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address   string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Heartbeat bool   `protobuf:"varint,2,opt,name=heartbeat,proto3" json:"heartbeat,omitempty"`
	Status    int32  `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *WorkerInfo) Reset() {
	*x = WorkerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerInfo) ProtoMessage() {}

func (x *WorkerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerInfo.ProtoReflect.Descriptor instead.
func (*WorkerInfo) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{6}
}

func (x *WorkerInfo) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *WorkerInfo) GetHeartbeat() bool {
	if x != nil {
		return x.Heartbeat
	}
	return false
}

func (x *WorkerInfo) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

type TaskPool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MapTasks    map[int32]*Task  `protobuf:"bytes,1,rep,name=mapTasks,proto3" json:"mapTasks,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ReduceTasks map[int32]*Task  `protobuf:"bytes,2,rep,name=reduceTasks,proto3" json:"reduceTasks,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ActiveTasks map[string]*Task `protobuf:"bytes,3,rep,name=activeTasks,proto3" json:"activeTasks,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TaskPool) Reset() {
	*x = TaskPool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskPool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskPool) ProtoMessage() {}

func (x *TaskPool) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskPool.ProtoReflect.Descriptor instead.
func (*TaskPool) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{7}
}

func (x *TaskPool) GetMapTasks() map[int32]*Task {
	if x != nil {
		return x.MapTasks
	}
	return nil
}

func (x *TaskPool) GetReduceTasks() map[int32]*Task {
	if x != nil {
		return x.ReduceTasks
	}
	return nil
}

func (x *TaskPool) GetActiveTasks() map[string]*Task {
	if x != nil {
		return x.ActiveTasks
	}
	return nil
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int32       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Input  []*BlobInfo `protobuf:"bytes,2,rep,name=input,proto3" json:"input,omitempty"`
	Output []*BlobInfo `protobuf:"bytes,3,rep,name=output,proto3" json:"output,omitempty"`
	Status int32       `protobuf:"varint,4,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_master_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_master_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_master_proto_rawDescGZIP(), []int{8}
}

func (x *Task) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Task) GetInput() []*BlobInfo {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Task) GetOutput() []*BlobInfo {
	if x != nil {
		return x.Output
	}
	return nil
}

func (x *Task) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_master_proto protoreflect.FileDescriptor

var file_master_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x0b, 0x61, 0x7a, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a, 0x0a, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2f, 0x0a, 0x13, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xce, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x65, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x12,
	0x44, 0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x10, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x32, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50,
	0x6f, 0x6f, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x0a, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x2c, 0x0a, 0x08, 0x74, 0x61, 0x73,
	0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x61,
	0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x52, 0x08, 0x74,
	0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x22, 0x2c, 0x0a, 0x10, 0x52, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x54, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x0c, 0x0a, 0x01, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x6d, 0x12, 0x0c, 0x0a,
	0x01, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x72, 0x22, 0xcf, 0x01, 0x0a, 0x0a,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x12, 0x33, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1d, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x50,
	0x6f, 0x6f, 0x6c, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x64, 0x65, 0x61, 0x64, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x51, 0x75, 0x65, 0x75, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0f,
	0x64, 0x65, 0x61, 0x64, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x51, 0x75, 0x65, 0x75, 0x65, 0x1a,
	0x4c, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x5c, 0x0a,
	0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65,
	0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62,
	0x65, 0x61, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xb7, 0x03, 0x0a, 0x08,
	0x54, 0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x3a, 0x0a, 0x08, 0x6d, 0x61, 0x70, 0x54,
	0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6d, 0x61, 0x73,
	0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x2e, 0x4d, 0x61, 0x70,
	0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x61, 0x70, 0x54,
	0x61, 0x73, 0x6b, 0x73, 0x12, 0x43, 0x0a, 0x0b, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x54, 0x61,
	0x73, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c, 0x2e, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0b, 0x72, 0x65,
	0x64, 0x75, 0x63, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x43, 0x0a, 0x0b, 0x61, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x50, 0x6f, 0x6f, 0x6c,
	0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x1a, 0x49,
	0x0a, 0x0d, 0x4d, 0x61, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x22, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4c, 0x0a, 0x10, 0x52, 0x65, 0x64,
	0x75, 0x63, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x22, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4c, 0x0a, 0x10, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x22, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x7e, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61,
	0x7a, 0x75, 0x72, 0x65, 0x2e, 0x42, 0x6c, 0x6f, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x7a, 0x75, 0x72, 0x65, 0x2e, 0x42, 0x6c, 0x6f,
	0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x92, 0x01, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x43, 0x0a, 0x0e, 0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x12, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x4d, 0x61, 0x70, 0x70,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x1b, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e,
	0x4d, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x12, 0x1a, 0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18,
	0x2e, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_master_proto_rawDescOnce sync.Once
	file_master_proto_rawDescData = file_master_proto_rawDesc
)

func file_master_proto_rawDescGZIP() []byte {
	file_master_proto_rawDescOnce.Do(func() {
		file_master_proto_rawDescData = protoimpl.X.CompressGZIP(file_master_proto_rawDescData)
	})
	return file_master_proto_rawDescData
}

var file_master_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_master_proto_goTypes = []interface{}{
	(*MapperInfo)(nil),          // 0: master.MapperInfo
	(*MapperRegisterReply)(nil), // 1: master.MapperRegisterReply
	(*ReplicationRequest)(nil),  // 2: master.ReplicationRequest
	(*ReplicationReply)(nil),    // 3: master.ReplicationReply
	(*ContainerRequest)(nil),    // 4: master.ContainerRequest
	(*WorkerPool)(nil),          // 5: master.WorkerPool
	(*WorkerInfo)(nil),          // 6: master.WorkerInfo
	(*TaskPool)(nil),            // 7: master.TaskPool
	(*Task)(nil),                // 8: master.Task
	nil,                         // 9: master.WorkerPool.StateEntry
	nil,                         // 10: master.TaskPool.MapTasksEntry
	nil,                         // 11: master.TaskPool.ReduceTasksEntry
	nil,                         // 12: master.TaskPool.ActiveTasksEntry
	(*BlobInfo)(nil),            // 13: azure.BlobInfo
}
var file_master_proto_depIdxs = []int32{
	4,  // 0: master.ReplicationRequest.userRequestQueue:type_name -> master.ContainerRequest
	5,  // 1: master.ReplicationRequest.workerPool:type_name -> master.WorkerPool
	7,  // 2: master.ReplicationRequest.taskPool:type_name -> master.TaskPool
	9,  // 3: master.WorkerPool.state:type_name -> master.WorkerPool.StateEntry
	10, // 4: master.TaskPool.mapTasks:type_name -> master.TaskPool.MapTasksEntry
	11, // 5: master.TaskPool.reduceTasks:type_name -> master.TaskPool.ReduceTasksEntry
	12, // 6: master.TaskPool.activeTasks:type_name -> master.TaskPool.ActiveTasksEntry
	13, // 7: master.Task.input:type_name -> azure.BlobInfo
	13, // 8: master.Task.output:type_name -> azure.BlobInfo
	6,  // 9: master.WorkerPool.StateEntry.value:type_name -> master.WorkerInfo
	8,  // 10: master.TaskPool.MapTasksEntry.value:type_name -> master.Task
	8,  // 11: master.TaskPool.ReduceTasksEntry.value:type_name -> master.Task
	8,  // 12: master.TaskPool.ActiveTasksEntry.value:type_name -> master.Task
	0,  // 13: master.Master.MapperRegister:input_type -> master.MapperInfo
	2,  // 14: master.Master.Replicate:input_type -> master.ReplicationRequest
	1,  // 15: master.Master.MapperRegister:output_type -> master.MapperRegisterReply
	3,  // 16: master.Master.Replicate:output_type -> master.ReplicationReply
	15, // [15:17] is the sub-list for method output_type
	13, // [13:15] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_master_proto_init() }
func file_master_proto_init() {
	if File_master_proto != nil {
		return
	}
	file_azure_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_master_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperInfo); i {
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
		file_master_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MapperRegisterReply); i {
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
		file_master_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicationRequest); i {
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
		file_master_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicationReply); i {
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
		file_master_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerRequest); i {
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
		file_master_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerPool); i {
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
		file_master_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerInfo); i {
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
		file_master_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskPool); i {
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
		file_master_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
			RawDescriptor: file_master_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_master_proto_goTypes,
		DependencyIndexes: file_master_proto_depIdxs,
		MessageInfos:      file_master_proto_msgTypes,
	}.Build()
	File_master_proto = out.File
	file_master_proto_rawDesc = nil
	file_master_proto_goTypes = nil
	file_master_proto_depIdxs = nil
}
