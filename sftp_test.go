// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package xk6sftp

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	user     = "user"
	password = "pwd"
	host     = "localhost"
	port     = 3322
)

func TestSFTPClient_Connect(t *testing.T) {
	client := Client{}
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	sftpClient.Close()
}

func TestSFTPClient_Close(t *testing.T) {
	client := Client{}
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	sftpClient.Close()
}

func TestSFTPClient_ConnectErrors(t *testing.T) {
	client := Client{}
	sftpClient := client.NewClient(user, password, host, 23)
	require.Nil(t, sftpClient)
}

func TestSFTPClient_UploadFile(t *testing.T) {
	client := new(Client)
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	defer sftpClient.Close()

	localPath, err := os.Getwd()
	require.NoError(t, err)
	localPath = localPath + "/testdata/upload.txt"
	remotePath := "/wwwroot/upload.txt"

	// Create directory if not exists
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		err = os.Mkdir("testdata", 0755)
		require.NoError(t, err)
	}

	// Create a local file to upload
	file, err := os.Create(localPath)
	require.NoError(t, err)
	_, err = file.WriteString("This is a test file for upload.")
	require.NoError(t, err)
	file.Close()

	fmt.Printf("Local file created at %s\n", localPath)
	fmt.Printf("Uploading file to %s\n", remotePath)
	result := sftpClient.UploadFile(localPath, remotePath)
	require.Equal(t, true, result.Success)

	// Clean up
	os.Remove(localPath)
	result = sftpClient.DeleteFile(remotePath)
	require.Equal(t, true, result.Success)
}

func TestSFTPClient_DownloadFile(t *testing.T) {
	client := new(Client)
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	defer sftpClient.Close()

	localPath, err := os.Getwd()
	require.NoError(t, err)
	localPath = localPath + "/testdata/download.txt"
	remotePath := "/wwwroot/download.txt"

	createRemoteTestFile(sftpClient, remotePath, t)

	result := sftpClient.DownloadFile(remotePath, localPath)
	require.Equal(t, true, result.Success)

	// Verify the downloaded file
	downloadedFile, err := os.ReadFile(localPath)
	require.NoError(t, err)
	require.Equal(t, "This is a test file for download.", string(downloadedFile))

	// Clean up
	os.Remove(localPath)
	result = sftpClient.DeleteFile(remotePath)
	require.Equal(t, true, result.Success)
}

func createRemoteTestFile(sftpClient *SFTPClient, remotePath string, t *testing.T) {
	file, err := sftpClient.client.Create(remotePath)
	require.NoError(t, err)
	_, err = file.Write([]byte("This is a test file for download."))
	require.NoError(t, err)
	file.Close()
}

func TestErrorValidations(t *testing.T) {
	client := new(Client)
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	defer sftpClient.Close()

	testCases := []struct {
		client     *SFTPClient
		name       string
		localPath  string
		remotePath string
		action     func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult
		expected   bool
	}{
		{
			client:     sftpClient,
			name:       "Invalid local file path for upload",
			localPath:  "nonexistentfile.txt",
			remotePath: "/wwwroot/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.UploadFile(localPath, remotePath)
			},
			expected: false,
		},
		{
			client:     sftpClient,
			name:       "Invalid remote file path for upload",
			localPath:  func() string { p, _ := os.Getwd(); return p + "/examples/test-data/utf-8.txt" }(),
			remotePath: "/invalid|path/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.UploadFile(localPath, remotePath)
			},
			expected: false,
		},
		{
			client:     sftpClient,
			name:       "Upload a file to download tests",
			localPath:  func() string { p, _ := os.Getwd(); return p + "/examples/test-data/utf-8.txt" }(),
			remotePath: "/wwwroot/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.UploadFile(localPath, remotePath)
			},
			expected: true,
		},
		{
			client:     sftpClient,
			name:       "Invalid remote file path for download",
			localPath:  func() string { p, _ := os.Getwd(); return p + "/testdata/upload.txt" }(),
			remotePath: "/nonexistentdir/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.DownloadFile(remotePath, localPath)
			},
			expected: false,
		},
		{
			client:     sftpClient,
			name:       "Invalid local path for download",
			localPath:  "/invalid|path/upload.txt",
			remotePath: "/wwwroot/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.DownloadFile(remotePath, localPath)
			},
			expected: false,
		},
		{
			client:     sftpClient,
			name:       "Invalid remote file path for delete",
			localPath:  "",
			remotePath: "/nonexistentdir/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.DeleteFile(remotePath)
			},
			expected: false,
		},
		{
			client:     &SFTPClient{client: nil},
			name:       "Validate if client is properly connected when deleting a file",
			localPath:  "",
			remotePath: "/nonexistentdir/upload.txt",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.DeleteFile(remotePath)
			},
			expected: false,
		},
		{
			client:     &SFTPClient{client: nil},
			name:       "Validate if client is properly connected when closing the connection",
			localPath:  "",
			remotePath: "",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.Close()
			},
			expected: true,
		},
		{
			client:     &SFTPClient{client: nil},
			name:       "Validate if client is properly connected when uploading a file",
			localPath:  "",
			remotePath: "",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.UploadFile(localPath, remotePath)
			},
			expected: false,
		},
		{
			client:     &SFTPClient{client: nil},
			name:       "Validate if client is properly connected when downloading a file",
			localPath:  "",
			remotePath: "",
			action: func(sftpClient *SFTPClient, localPath, remotePath string) *OperationResult {
				return sftpClient.DownloadFile(remotePath, localPath)
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.action(tc.client, tc.localPath, tc.remotePath)
			require.Equal(t, tc.expected, result.Success)
		})
	}
}
