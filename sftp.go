package xk6sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct{}

type OperationResult struct {
	Success bool
	Message string
}

type SFTPClient struct {
	client *sftp.Client
}

func (*Client) NewClientWithRSAKey(user, host string, port int, privateKeyPath string) *SFTPClient {
	absPrivateKeyPath, err := filepath.Abs(privateKeyPath)
	if err != nil {
		logger.Errorf("failed to get absolute path for private key (%s): %v", privateKeyPath, err)
		return nil
	}

	key, err := os.ReadFile(absPrivateKeyPath)
	if err != nil {
		logger.Errorf("unable to read private key: %v", err)
		return nil
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Errorf("unable to parse private key: %v", err)
		return nil
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Errorf("failed to dial to %s: %v", addr, err)
		return nil
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		logger.Errorf("failed to create sftp client: %v", err)
		return nil
	}

	return &SFTPClient{client: client}
}

func (*Client) NewClient(user, password, host string, port int) *SFTPClient {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Errorf("failed to dial to %s: %v", addr, err)
		return nil
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		logger.Errorf("failed to create sftp client: %v", err)
		return nil
	}

	return &SFTPClient{client: client}
}

func (s *SFTPClient) UploadFile(localPath, remotePath string) OperationResult {
	absLocalPath, err := filepath.Abs(localPath)
	if err != nil {
		logger.Errorf("failed to get absolute path for local file (%s): %v", localPath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to get absolute path for local file (%s): %v", localPath, err)}
	}

	srcFile, err := os.Open(absLocalPath)
	if err != nil {
		logger.Errorf("failed to open local file (%s): %v", absLocalPath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to open local file (%s): %v", absLocalPath, err)}
	}
	defer srcFile.Close()

	err = s.client.MkdirAll(filepath.Dir(remotePath))
	if err != nil {
		logger.Errorf("failed to create directories for remote file: %v", err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to create directories for remote file: %v", err)}
	}

	dstFile, err := s.client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to create remote file (%s): %v", remotePath, err)}
	}
	defer dstFile.Close()

	bytes, err := io.ReadAll(srcFile)
	if err != nil {
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to read local file (%s): %v", localPath, err)}
	}

	if _, err := dstFile.Write(bytes); err != nil {
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to write to remote file (%s): %v", remotePath, err)}
	}

	return OperationResult{Success: true, Message: "file uploaded successfully"}
}

func (s *SFTPClient) DownloadFile(remotePath, localPath string) OperationResult {
	srcFile, err := s.client.Open(remotePath)
	if err != nil {
		logger.Errorf("failed to open remote file (%s): %v", remotePath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to open remote file (%s): %v", remotePath, err)}
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		logger.Errorf("failed to create directories for local file: %v", err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to create directories for local file: %v", err)}
	}

	dstFile, err := os.Create(localPath)
	if err != nil {
		logger.Errorf("failed to create local file (%s): %v", localPath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to create local file (%s): %v", localPath, err)}
	}
	defer dstFile.Close()

	bytes, err := io.ReadAll(srcFile)
	if err != nil {
		logger.Errorf("failed to read remote file (%s): %v", remotePath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to read remote file (%s): %v", remotePath, err)}
	}

	if _, err := dstFile.Write(bytes); err != nil {
		logger.Errorf("failed to write to local file: %v", err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to write to local file: %v", err)}
	}

	return OperationResult{Success: true, Message: "file downloaded successfully"}
}

func (s *SFTPClient) DeleteFile(remotePath string) OperationResult {
	err := s.client.Remove(remotePath)
	if err != nil {
		logger.Errorf("failed to delete remote file (%s): %v", remotePath, err)
		return OperationResult{Success: false, Message: fmt.Sprintf("failed to delete remote file (%s): %v", remotePath, err)}
	}
	return OperationResult{Success: true, Message: "file deleted successfully"}
}

func (s *SFTPClient) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}
