package xk6sftp

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSFTPClient_Connect(t *testing.T) {
	mockClient := new(MockClient)
	mockSftpClient := new(MockSFTPClient)
	mockClient.On("NewClient", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockSftpClient)

	result := mockClient.NewClient("user", "pwd", "host", 22)
	require.NotNil(t, result)
	mockClient.AssertExpectations(t)
}

func TestSFTPClient_UploadFile(t *testing.T) {
	mockClient := new(MockClient)
	mockSftpClient := new(MockSFTPClient)
	mockClient.On("NewClient", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockSftpClient)
	mockSftpClient.On("UploadFile", mock.Anything, mock.Anything).Return(&OperationResult{Success: true})
	mockSftpClient.On("DeleteFile", mock.Anything).Return(&OperationResult{Success: true})

	sftpClient := mockClient.NewClient("user", "pwd", "host", 22)
	require.NotNil(t, sftpClient)

	localPath := "testdata/upload.txt"
	remotePath := "/wwwroot/upload.txt"

	result := sftpClient.UploadFile(remotePath, localPath)
	require.Equal(t, true, result.Success)

	result = sftpClient.DeleteFile(remotePath)
	require.Equal(t, true, result.Success)

	mockClient.AssertExpectations(t)
}

func TestSFTPClient_DownloadFile(t *testing.T) {
	mockClient := new(MockClient)
	mockSftpClient := new(MockSFTPClient)
	mockClient.On("NewClient", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockSftpClient)
	mockSftpClient.On("DownloadFile", mock.Anything, mock.Anything).Return(&OperationResult{Success: true})
	mockSftpClient.On("DeleteFile", mock.Anything).Return(&OperationResult{Success: true})

	sftpClient := mockClient.NewClient("user", "pwd", "host", 22)
	require.NotNil(t, sftpClient)

	localPath := "testdata/download.txt"
	remotePath := "/wwwroot/download.txt"

	result := sftpClient.DownloadFile(remotePath, localPath)
	require.Equal(t, true, result.Success)

	result = sftpClient.DeleteFile(remotePath)
	require.Equal(t, true, result.Success)

	mockClient.AssertExpectations(t)
}
