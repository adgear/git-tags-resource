// Code generated by MockGen. DO NOT EDIT.
// Source: services/gitTagsService.go

// Package mocks is a generated GoMock package.
package mocks

import (
	semver "github.com/Masterminds/semver"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGitTagsService is a mock of GitTagsService interface
type MockGitTagsService struct {
	ctrl     *gomock.Controller
	recorder *MockGitTagsServiceMockRecorder
}

// MockGitTagsServiceMockRecorder is the mock recorder for MockGitTagsService
type MockGitTagsServiceMockRecorder struct {
	mock *MockGitTagsService
}

// NewMockGitTagsService creates a new mock instance
func NewMockGitTagsService(ctrl *gomock.Controller) *MockGitTagsService {
	mock := &MockGitTagsService{ctrl: ctrl}
	mock.recorder = &MockGitTagsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitTagsService) EXPECT() *MockGitTagsServiceMockRecorder {
	return m.recorder
}

// FetchTags mocks base method
func (m *MockGitTagsService) FetchTags(uri, repositoryName, tagFilter string) ([]string, error) {
	ret := m.ctrl.Call(m, "FetchTags", uri, repositoryName, tagFilter)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTags indicates an expected call of FetchTags
func (mr *MockGitTagsServiceMockRecorder) FetchTags(uri, repositoryName, tagFilter interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTags", reflect.TypeOf((*MockGitTagsService)(nil).FetchTags), uri, repositoryName, tagFilter)
}

// ExtractTags mocks base method
func (m *MockGitTagsService) ExtractTags(latestOnly bool, tagList []semver.Version) ([]map[string]string, error) {
	ret := m.ctrl.Call(m, "ExtractTags", latestOnly, tagList)
	ret0, _ := ret[0].([]map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExtractTags indicates an expected call of ExtractTags
func (mr *MockGitTagsServiceMockRecorder) ExtractTags(latestOnly, tagList interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExtractTags", reflect.TypeOf((*MockGitTagsService)(nil).ExtractTags), latestOnly, tagList)
}

// CloneRef mocks base method
func (m *MockGitTagsService) CloneRef(uri, destination, version string) error {
	ret := m.ctrl.Call(m, "CloneRef", uri, destination, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloneRef indicates an expected call of CloneRef
func (mr *MockGitTagsServiceMockRecorder) CloneRef(uri, destination, version interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneRef", reflect.TypeOf((*MockGitTagsService)(nil).CloneRef), uri, destination, version)
}

// CheckoutATag mocks base method
func (m *MockGitTagsService) CheckoutATag(destination, version string) (map[string]string, error) {
	ret := m.ctrl.Call(m, "CheckoutATag", destination, version)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckoutATag indicates an expected call of CheckoutATag
func (mr *MockGitTagsServiceMockRecorder) CheckoutATag(destination, version interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutATag", reflect.TypeOf((*MockGitTagsService)(nil).CheckoutATag), destination, version)
}

// CheckoutLWTag mocks base method
func (m *MockGitTagsService) CheckoutLWTag(destination, version string) (map[string]string, error) {
	ret := m.ctrl.Call(m, "CheckoutLWTag", destination, version)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckoutLWTag indicates an expected call of CheckoutLWTag
func (mr *MockGitTagsServiceMockRecorder) CheckoutLWTag(destination, version interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutLWTag", reflect.TypeOf((*MockGitTagsService)(nil).CheckoutLWTag), destination, version)
}
