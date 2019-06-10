// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	sync "sync"

	backend "github.com/IBM/ibmcloud-storage-volume-lib/volume-providers/softlayer/backend"
	datatypes "github.com/softlayer/softlayer-go/datatypes"
)

type NetworkStorageService struct {
	RemoveAccessFromIPAddressListStub       func([]datatypes.Network_Subnet_IpAddress) (bool, error)
	RemoveAccessFromSubnetListStub          func([]datatypes.Network_Subnet) (bool, error)
	AllowAccessFromIPAddressListStub        func([]datatypes.Network_Subnet_IpAddress) (bool, error)
	allowAccessFromIPAddressListMutex       sync.RWMutex
	allowAccessFromIPAddressListArgsForCall []struct {
		arg1 []datatypes.Network_Subnet_IpAddress
	}
	allowAccessFromIPAddressListReturns struct {
		result1 bool
		result2 error
	}
	allowAccessFromIPAddressListReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	AllowAccessFromSubnetListStub        func([]datatypes.Network_Subnet) (bool, error)
	allowAccessFromSubnetListMutex       sync.RWMutex
	allowAccessFromSubnetListArgsForCall []struct {
		arg1 []datatypes.Network_Subnet
	}
	allowAccessFromSubnetListReturns struct {
		result1 bool
		result2 error
	}
	allowAccessFromSubnetListReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	CreateSnapshotStub        func(*string) (datatypes.Network_Storage, error)
	createSnapshotMutex       sync.RWMutex
	createSnapshotArgsForCall []struct {
		arg1 *string
	}
	createSnapshotReturns struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	createSnapshotReturnsOnCall map[int]struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	DeleteObjectStub        func() (bool, error)
	deleteObjectMutex       sync.RWMutex
	deleteObjectArgsForCall []struct {
	}
	deleteObjectReturns struct {
		result1 bool
		result2 error
	}
	deleteObjectReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	EditObjectStub        func(*datatypes.Network_Storage) (bool, error)
	editObjectMutex       sync.RWMutex
	editObjectArgsForCall []struct {
		arg1 *datatypes.Network_Storage
	}
	editObjectReturns struct {
		result1 bool
		result2 error
	}
	editObjectReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	FilterStub        func(string) backend.NetworkStorageService
	filterMutex       sync.RWMutex
	filterArgsForCall []struct {
		arg1 string
	}
	filterReturns struct {
		result1 backend.NetworkStorageService
	}
	filterReturnsOnCall map[int]struct {
		result1 backend.NetworkStorageService
	}
	GetObjectStub        func() (datatypes.Network_Storage, error)
	getObjectMutex       sync.RWMutex
	getObjectArgsForCall []struct {
	}
	getObjectReturns struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	getObjectReturnsOnCall map[int]struct {
		result1 datatypes.Network_Storage
		result2 error
	}
	GetSnapshotsStub        func() ([]datatypes.Network_Storage, error)
	getSnapshotsMutex       sync.RWMutex
	getSnapshotsArgsForCall []struct {
	}
	getSnapshotsReturns struct {
		result1 []datatypes.Network_Storage
		result2 error
	}
	getSnapshotsReturnsOnCall map[int]struct {
		result1 []datatypes.Network_Storage
		result2 error
	}
	IDStub        func(int) backend.NetworkStorageService
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
		arg1 int
	}
	iDReturns struct {
		result1 backend.NetworkStorageService
	}
	iDReturnsOnCall map[int]struct {
		result1 backend.NetworkStorageService
	}
	MaskStub        func(string) backend.NetworkStorageService
	maskMutex       sync.RWMutex
	maskArgsForCall []struct {
		arg1 string
	}
	maskReturns struct {
		result1 backend.NetworkStorageService
	}
	maskReturnsOnCall map[int]struct {
		result1 backend.NetworkStorageService
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *NetworkStorageService) RemoveAccessFromSubnetList(arg1 []datatypes.Network_Subnet_IpAddress) (bool, error) {
	return true, nil
}

func (fake *NetworkStorageService) RemoveAccessFromIPAddressList(arg1 []datatypes.Network_Subnet_IpAddress) (bool, error) {
	return true, nil
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressList(arg1 []datatypes.Network_Subnet_IpAddress) (bool, error) {
	var arg1Copy []datatypes.Network_Subnet_IpAddress
	if arg1 != nil {
		arg1Copy = make([]datatypes.Network_Subnet_IpAddress, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.allowAccessFromIPAddressListMutex.Lock()
	ret, specificReturn := fake.allowAccessFromIPAddressListReturnsOnCall[len(fake.allowAccessFromIPAddressListArgsForCall)]
	fake.allowAccessFromIPAddressListArgsForCall = append(fake.allowAccessFromIPAddressListArgsForCall, struct {
		arg1 []datatypes.Network_Subnet_IpAddress
	}{arg1Copy})
	fake.recordInvocation("AllowAccessFromIPAddressList", []interface{}{arg1Copy})
	fake.allowAccessFromIPAddressListMutex.Unlock()
	if fake.AllowAccessFromIPAddressListStub != nil {
		return fake.AllowAccessFromIPAddressListStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.allowAccessFromIPAddressListReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressListCallCount() int {
	fake.allowAccessFromIPAddressListMutex.RLock()
	defer fake.allowAccessFromIPAddressListMutex.RUnlock()
	return len(fake.allowAccessFromIPAddressListArgsForCall)
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressListCalls(stub func([]datatypes.Network_Subnet_IpAddress) (bool, error)) {
	fake.allowAccessFromIPAddressListMutex.Lock()
	defer fake.allowAccessFromIPAddressListMutex.Unlock()
	fake.AllowAccessFromIPAddressListStub = stub
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressListArgsForCall(i int) []datatypes.Network_Subnet_IpAddress {
	fake.allowAccessFromIPAddressListMutex.RLock()
	defer fake.allowAccessFromIPAddressListMutex.RUnlock()
	argsForCall := fake.allowAccessFromIPAddressListArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressListReturns(result1 bool, result2 error) {
	fake.allowAccessFromIPAddressListMutex.Lock()
	defer fake.allowAccessFromIPAddressListMutex.Unlock()
	fake.AllowAccessFromIPAddressListStub = nil
	fake.allowAccessFromIPAddressListReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) AllowAccessFromIPAddressListReturnsOnCall(i int, result1 bool, result2 error) {
	fake.allowAccessFromIPAddressListMutex.Lock()
	defer fake.allowAccessFromIPAddressListMutex.Unlock()
	fake.AllowAccessFromIPAddressListStub = nil
	if fake.allowAccessFromIPAddressListReturnsOnCall == nil {
		fake.allowAccessFromIPAddressListReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.allowAccessFromIPAddressListReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) AllowAccessFromSubnetList(arg1 []datatypes.Network_Subnet) (bool, error) {
	var arg1Copy []datatypes.Network_Subnet
	if arg1 != nil {
		arg1Copy = make([]datatypes.Network_Subnet, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.allowAccessFromSubnetListMutex.Lock()
	ret, specificReturn := fake.allowAccessFromSubnetListReturnsOnCall[len(fake.allowAccessFromSubnetListArgsForCall)]
	fake.allowAccessFromSubnetListArgsForCall = append(fake.allowAccessFromSubnetListArgsForCall, struct {
		arg1 []datatypes.Network_Subnet
	}{arg1Copy})
	fake.recordInvocation("AllowAccessFromSubnetList", []interface{}{arg1Copy})
	fake.allowAccessFromSubnetListMutex.Unlock()
	if fake.AllowAccessFromSubnetListStub != nil {
		return fake.AllowAccessFromSubnetListStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.allowAccessFromSubnetListReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) AllowAccessFromSubnetListCallCount() int {
	fake.allowAccessFromSubnetListMutex.RLock()
	defer fake.allowAccessFromSubnetListMutex.RUnlock()
	return len(fake.allowAccessFromSubnetListArgsForCall)
}

func (fake *NetworkStorageService) AllowAccessFromSubnetListCalls(stub func([]datatypes.Network_Subnet) (bool, error)) {
	fake.allowAccessFromSubnetListMutex.Lock()
	defer fake.allowAccessFromSubnetListMutex.Unlock()
	fake.AllowAccessFromSubnetListStub = stub
}

func (fake *NetworkStorageService) AllowAccessFromSubnetListArgsForCall(i int) []datatypes.Network_Subnet {
	fake.allowAccessFromSubnetListMutex.RLock()
	defer fake.allowAccessFromSubnetListMutex.RUnlock()
	argsForCall := fake.allowAccessFromSubnetListArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) AllowAccessFromSubnetListReturns(result1 bool, result2 error) {
	fake.allowAccessFromSubnetListMutex.Lock()
	defer fake.allowAccessFromSubnetListMutex.Unlock()
	fake.AllowAccessFromSubnetListStub = nil
	fake.allowAccessFromSubnetListReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) AllowAccessFromSubnetListReturnsOnCall(i int, result1 bool, result2 error) {
	fake.allowAccessFromSubnetListMutex.Lock()
	defer fake.allowAccessFromSubnetListMutex.Unlock()
	fake.AllowAccessFromSubnetListStub = nil
	if fake.allowAccessFromSubnetListReturnsOnCall == nil {
		fake.allowAccessFromSubnetListReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.allowAccessFromSubnetListReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) CreateSnapshot(arg1 *string) (datatypes.Network_Storage, error) {
	fake.createSnapshotMutex.Lock()
	ret, specificReturn := fake.createSnapshotReturnsOnCall[len(fake.createSnapshotArgsForCall)]
	fake.createSnapshotArgsForCall = append(fake.createSnapshotArgsForCall, struct {
		arg1 *string
	}{arg1})
	fake.recordInvocation("CreateSnapshot", []interface{}{arg1})
	fake.createSnapshotMutex.Unlock()
	if fake.CreateSnapshotStub != nil {
		return fake.CreateSnapshotStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createSnapshotReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) CreateSnapshotCallCount() int {
	fake.createSnapshotMutex.RLock()
	defer fake.createSnapshotMutex.RUnlock()
	return len(fake.createSnapshotArgsForCall)
}

func (fake *NetworkStorageService) CreateSnapshotCalls(stub func(*string) (datatypes.Network_Storage, error)) {
	fake.createSnapshotMutex.Lock()
	defer fake.createSnapshotMutex.Unlock()
	fake.CreateSnapshotStub = stub
}

func (fake *NetworkStorageService) CreateSnapshotArgsForCall(i int) *string {
	fake.createSnapshotMutex.RLock()
	defer fake.createSnapshotMutex.RUnlock()
	argsForCall := fake.createSnapshotArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) CreateSnapshotReturns(result1 datatypes.Network_Storage, result2 error) {
	fake.createSnapshotMutex.Lock()
	defer fake.createSnapshotMutex.Unlock()
	fake.CreateSnapshotStub = nil
	fake.createSnapshotReturns = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) CreateSnapshotReturnsOnCall(i int, result1 datatypes.Network_Storage, result2 error) {
	fake.createSnapshotMutex.Lock()
	defer fake.createSnapshotMutex.Unlock()
	fake.CreateSnapshotStub = nil
	if fake.createSnapshotReturnsOnCall == nil {
		fake.createSnapshotReturnsOnCall = make(map[int]struct {
			result1 datatypes.Network_Storage
			result2 error
		})
	}
	fake.createSnapshotReturnsOnCall[i] = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) DeleteObject() (bool, error) {
	fake.deleteObjectMutex.Lock()
	ret, specificReturn := fake.deleteObjectReturnsOnCall[len(fake.deleteObjectArgsForCall)]
	fake.deleteObjectArgsForCall = append(fake.deleteObjectArgsForCall, struct {
	}{})
	fake.recordInvocation("DeleteObject", []interface{}{})
	fake.deleteObjectMutex.Unlock()
	if fake.DeleteObjectStub != nil {
		return fake.DeleteObjectStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteObjectReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) DeleteObjectCallCount() int {
	fake.deleteObjectMutex.RLock()
	defer fake.deleteObjectMutex.RUnlock()
	return len(fake.deleteObjectArgsForCall)
}

func (fake *NetworkStorageService) DeleteObjectCalls(stub func() (bool, error)) {
	fake.deleteObjectMutex.Lock()
	defer fake.deleteObjectMutex.Unlock()
	fake.DeleteObjectStub = stub
}

func (fake *NetworkStorageService) DeleteObjectReturns(result1 bool, result2 error) {
	fake.deleteObjectMutex.Lock()
	defer fake.deleteObjectMutex.Unlock()
	fake.DeleteObjectStub = nil
	fake.deleteObjectReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) DeleteObjectReturnsOnCall(i int, result1 bool, result2 error) {
	fake.deleteObjectMutex.Lock()
	defer fake.deleteObjectMutex.Unlock()
	fake.DeleteObjectStub = nil
	if fake.deleteObjectReturnsOnCall == nil {
		fake.deleteObjectReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.deleteObjectReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) EditObject(arg1 *datatypes.Network_Storage) (bool, error) {
	fake.editObjectMutex.Lock()
	ret, specificReturn := fake.editObjectReturnsOnCall[len(fake.editObjectArgsForCall)]
	fake.editObjectArgsForCall = append(fake.editObjectArgsForCall, struct {
		arg1 *datatypes.Network_Storage
	}{arg1})
	fake.recordInvocation("EditObject", []interface{}{arg1})
	fake.editObjectMutex.Unlock()
	if fake.EditObjectStub != nil {
		return fake.EditObjectStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.editObjectReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) EditObjectCallCount() int {
	fake.editObjectMutex.RLock()
	defer fake.editObjectMutex.RUnlock()
	return len(fake.editObjectArgsForCall)
}

func (fake *NetworkStorageService) EditObjectCalls(stub func(*datatypes.Network_Storage) (bool, error)) {
	fake.editObjectMutex.Lock()
	defer fake.editObjectMutex.Unlock()
	fake.EditObjectStub = stub
}

func (fake *NetworkStorageService) EditObjectArgsForCall(i int) *datatypes.Network_Storage {
	fake.editObjectMutex.RLock()
	defer fake.editObjectMutex.RUnlock()
	argsForCall := fake.editObjectArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) EditObjectReturns(result1 bool, result2 error) {
	fake.editObjectMutex.Lock()
	defer fake.editObjectMutex.Unlock()
	fake.EditObjectStub = nil
	fake.editObjectReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) EditObjectReturnsOnCall(i int, result1 bool, result2 error) {
	fake.editObjectMutex.Lock()
	defer fake.editObjectMutex.Unlock()
	fake.EditObjectStub = nil
	if fake.editObjectReturnsOnCall == nil {
		fake.editObjectReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.editObjectReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) Filter(arg1 string) backend.NetworkStorageService {
	fake.filterMutex.Lock()
	ret, specificReturn := fake.filterReturnsOnCall[len(fake.filterArgsForCall)]
	fake.filterArgsForCall = append(fake.filterArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Filter", []interface{}{arg1})
	fake.filterMutex.Unlock()
	if fake.FilterStub != nil {
		return fake.FilterStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.filterReturns
	return fakeReturns.result1
}

func (fake *NetworkStorageService) FilterCallCount() int {
	fake.filterMutex.RLock()
	defer fake.filterMutex.RUnlock()
	return len(fake.filterArgsForCall)
}

func (fake *NetworkStorageService) FilterCalls(stub func(string) backend.NetworkStorageService) {
	fake.filterMutex.Lock()
	defer fake.filterMutex.Unlock()
	fake.FilterStub = stub
}

func (fake *NetworkStorageService) FilterArgsForCall(i int) string {
	fake.filterMutex.RLock()
	defer fake.filterMutex.RUnlock()
	argsForCall := fake.filterArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) FilterReturns(result1 backend.NetworkStorageService) {
	fake.filterMutex.Lock()
	defer fake.filterMutex.Unlock()
	fake.FilterStub = nil
	fake.filterReturns = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) FilterReturnsOnCall(i int, result1 backend.NetworkStorageService) {
	fake.filterMutex.Lock()
	defer fake.filterMutex.Unlock()
	fake.FilterStub = nil
	if fake.filterReturnsOnCall == nil {
		fake.filterReturnsOnCall = make(map[int]struct {
			result1 backend.NetworkStorageService
		})
	}
	fake.filterReturnsOnCall[i] = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) GetObject() (datatypes.Network_Storage, error) {
	fake.getObjectMutex.Lock()
	ret, specificReturn := fake.getObjectReturnsOnCall[len(fake.getObjectArgsForCall)]
	fake.getObjectArgsForCall = append(fake.getObjectArgsForCall, struct {
	}{})
	fake.recordInvocation("GetObject", []interface{}{})
	fake.getObjectMutex.Unlock()
	if fake.GetObjectStub != nil {
		return fake.GetObjectStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getObjectReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) GetObjectCallCount() int {
	fake.getObjectMutex.RLock()
	defer fake.getObjectMutex.RUnlock()
	return len(fake.getObjectArgsForCall)
}

func (fake *NetworkStorageService) GetObjectCalls(stub func() (datatypes.Network_Storage, error)) {
	fake.getObjectMutex.Lock()
	defer fake.getObjectMutex.Unlock()
	fake.GetObjectStub = stub
}

func (fake *NetworkStorageService) GetObjectReturns(result1 datatypes.Network_Storage, result2 error) {
	fake.getObjectMutex.Lock()
	defer fake.getObjectMutex.Unlock()
	fake.GetObjectStub = nil
	fake.getObjectReturns = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) GetObjectReturnsOnCall(i int, result1 datatypes.Network_Storage, result2 error) {
	fake.getObjectMutex.Lock()
	defer fake.getObjectMutex.Unlock()
	fake.GetObjectStub = nil
	if fake.getObjectReturnsOnCall == nil {
		fake.getObjectReturnsOnCall = make(map[int]struct {
			result1 datatypes.Network_Storage
			result2 error
		})
	}
	fake.getObjectReturnsOnCall[i] = struct {
		result1 datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) GetSnapshots() ([]datatypes.Network_Storage, error) {
	fake.getSnapshotsMutex.Lock()
	ret, specificReturn := fake.getSnapshotsReturnsOnCall[len(fake.getSnapshotsArgsForCall)]
	fake.getSnapshotsArgsForCall = append(fake.getSnapshotsArgsForCall, struct {
	}{})
	fake.recordInvocation("GetSnapshots", []interface{}{})
	fake.getSnapshotsMutex.Unlock()
	if fake.GetSnapshotsStub != nil {
		return fake.GetSnapshotsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getSnapshotsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *NetworkStorageService) GetSnapshotsCallCount() int {
	fake.getSnapshotsMutex.RLock()
	defer fake.getSnapshotsMutex.RUnlock()
	return len(fake.getSnapshotsArgsForCall)
}

func (fake *NetworkStorageService) GetSnapshotsCalls(stub func() ([]datatypes.Network_Storage, error)) {
	fake.getSnapshotsMutex.Lock()
	defer fake.getSnapshotsMutex.Unlock()
	fake.GetSnapshotsStub = stub
}

func (fake *NetworkStorageService) GetSnapshotsReturns(result1 []datatypes.Network_Storage, result2 error) {
	fake.getSnapshotsMutex.Lock()
	defer fake.getSnapshotsMutex.Unlock()
	fake.GetSnapshotsStub = nil
	fake.getSnapshotsReturns = struct {
		result1 []datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) GetSnapshotsReturnsOnCall(i int, result1 []datatypes.Network_Storage, result2 error) {
	fake.getSnapshotsMutex.Lock()
	defer fake.getSnapshotsMutex.Unlock()
	fake.GetSnapshotsStub = nil
	if fake.getSnapshotsReturnsOnCall == nil {
		fake.getSnapshotsReturnsOnCall = make(map[int]struct {
			result1 []datatypes.Network_Storage
			result2 error
		})
	}
	fake.getSnapshotsReturnsOnCall[i] = struct {
		result1 []datatypes.Network_Storage
		result2 error
	}{result1, result2}
}

func (fake *NetworkStorageService) ID(arg1 int) backend.NetworkStorageService {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("ID", []interface{}{arg1})
	fake.iDMutex.Unlock()
	if fake.IDStub != nil {
		return fake.IDStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.iDReturns
	return fakeReturns.result1
}

func (fake *NetworkStorageService) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *NetworkStorageService) IDCalls(stub func(int) backend.NetworkStorageService) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *NetworkStorageService) IDArgsForCall(i int) int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	argsForCall := fake.iDArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) IDReturns(result1 backend.NetworkStorageService) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) IDReturnsOnCall(i int, result1 backend.NetworkStorageService) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 backend.NetworkStorageService
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) Mask(arg1 string) backend.NetworkStorageService {
	fake.maskMutex.Lock()
	ret, specificReturn := fake.maskReturnsOnCall[len(fake.maskArgsForCall)]
	fake.maskArgsForCall = append(fake.maskArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Mask", []interface{}{arg1})
	fake.maskMutex.Unlock()
	if fake.MaskStub != nil {
		return fake.MaskStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.maskReturns
	return fakeReturns.result1
}

func (fake *NetworkStorageService) MaskCallCount() int {
	fake.maskMutex.RLock()
	defer fake.maskMutex.RUnlock()
	return len(fake.maskArgsForCall)
}

func (fake *NetworkStorageService) MaskCalls(stub func(string) backend.NetworkStorageService) {
	fake.maskMutex.Lock()
	defer fake.maskMutex.Unlock()
	fake.MaskStub = stub
}

func (fake *NetworkStorageService) MaskArgsForCall(i int) string {
	fake.maskMutex.RLock()
	defer fake.maskMutex.RUnlock()
	argsForCall := fake.maskArgsForCall[i]
	return argsForCall.arg1
}

func (fake *NetworkStorageService) MaskReturns(result1 backend.NetworkStorageService) {
	fake.maskMutex.Lock()
	defer fake.maskMutex.Unlock()
	fake.MaskStub = nil
	fake.maskReturns = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) MaskReturnsOnCall(i int, result1 backend.NetworkStorageService) {
	fake.maskMutex.Lock()
	defer fake.maskMutex.Unlock()
	fake.MaskStub = nil
	if fake.maskReturnsOnCall == nil {
		fake.maskReturnsOnCall = make(map[int]struct {
			result1 backend.NetworkStorageService
		})
	}
	fake.maskReturnsOnCall[i] = struct {
		result1 backend.NetworkStorageService
	}{result1}
}

func (fake *NetworkStorageService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.allowAccessFromIPAddressListMutex.RLock()
	defer fake.allowAccessFromIPAddressListMutex.RUnlock()
	fake.allowAccessFromSubnetListMutex.RLock()
	defer fake.allowAccessFromSubnetListMutex.RUnlock()
	fake.createSnapshotMutex.RLock()
	defer fake.createSnapshotMutex.RUnlock()
	fake.deleteObjectMutex.RLock()
	defer fake.deleteObjectMutex.RUnlock()
	fake.editObjectMutex.RLock()
	defer fake.editObjectMutex.RUnlock()
	fake.filterMutex.RLock()
	defer fake.filterMutex.RUnlock()
	fake.getObjectMutex.RLock()
	defer fake.getObjectMutex.RUnlock()
	fake.getSnapshotsMutex.RLock()
	defer fake.getSnapshotsMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.maskMutex.RLock()
	defer fake.maskMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *NetworkStorageService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ backend.NetworkStorageService = new(NetworkStorageService)
