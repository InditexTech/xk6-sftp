package xk6sftp

import "github.com/stretchr/testify/mock"

type MockClient struct {
	mock.Mock
}

func (m *MockClient) NewClient(user, password, host string, port int) SFTPClientInterface {
	args := m.Called(user, password, host, port)
	return args.Get(0).(SFTPClientInterface)
}

var _ SFTPClientInterface = &MockSFTPClient{}

type MockSFTPClient struct {
	mock.Mock
}

func (m *MockSFTPClient) UploadFile(localPath, remotePath string) *OperationResult {
	args := m.Called(localPath, remotePath)
	return args.Get(0).(*OperationResult)
}

func (m *MockSFTPClient) DownloadFile(remotePath, localPath string) *OperationResult {
	args := m.Called(remotePath, localPath)
	return args.Get(0).(*OperationResult)
}

func (m *MockSFTPClient) DeleteFile(remotePath string) *OperationResult {
	args := m.Called(remotePath)
	return args.Get(0).(*OperationResult)
}

func (m *MockSFTPClient) Close() *OperationResult {
	args := m.Called()
	return args.Get(0).(*OperationResult)
}
