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
	port     = 22
)

func TestSFTPClient_Connect(t *testing.T) {
	client := Client{}
	sftpClient := client.NewClient(user, password, host, port)
	require.NotNil(t, sftpClient)
	defer sftpClient.Close()
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
