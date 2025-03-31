// SPDX-FileCopyrightText: 2025 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

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

type SFTPClient struct {
	client *sftp.Client
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

func (s *SFTPClient) UploadFile(localPath, remotePath string) *OperationResult {
	if !s.existsConnection() {
		return &OperationResult{Success: false, Message: "sftp client is not connected"}
	}

	absLocalPath, _ := filepath.Abs(localPath)
	srcFile, err := os.Open(absLocalPath)
	if err != nil {
		logger.Errorf("failed to open local file (%s): %v", absLocalPath, err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to open local file (%s): %v", absLocalPath, err)}
	}
	defer srcFile.Close()

	err = s.client.MkdirAll(filepath.Dir(remotePath))
	if err != nil {
		logger.Errorf("failed to create directories for remote file: %v", err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to create directories for remote file: %v", err)}
	}

	dstFile, err := s.client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to create remote file (%s): %v", remotePath, err)}
	}
	defer dstFile.Close()

	bytes, err := io.ReadAll(srcFile)
	if err != nil {
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to read local file (%s): %v", localPath, err)}
	}

	if _, err := dstFile.Write(bytes); err != nil {
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to write to remote file (%s): %v", remotePath, err)}
	}

	return &OperationResult{Success: true, Message: "file uploaded successfully"}
}

func (s *SFTPClient) DownloadFile(remotePath, localPath string) *OperationResult {
	if !s.existsConnection() {
		return &OperationResult{Success: false, Message: "sftp client is not connected"}
	}

	srcFile, err := s.client.Open(remotePath)
	if err != nil {
		logger.Errorf("failed to open remote file (%s): %v", remotePath, err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to open remote file (%s): %v", remotePath, err)}
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		logger.Errorf("failed to create directories for local file: %v", err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to create directories for local file: %v", err)}
	}

	dstFile, err := os.Create(localPath)
	if err != nil {
		logger.Errorf("failed to create local file (%s): %v", localPath, err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to create local file (%s): %v", localPath, err)}
	}
	defer dstFile.Close()

	bytes, err := io.ReadAll(srcFile)
	if err != nil {
		logger.Errorf("failed to read remote file (%s): %v", remotePath, err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to read remote file (%s): %v", remotePath, err)}
	}

	if _, err := dstFile.Write(bytes); err != nil {
		logger.Errorf("failed to write to local file: %v", err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to write to local file: %v", err)}
	}

	return &OperationResult{Success: true, Message: "file downloaded successfully"}
}

func (s *SFTPClient) DeleteFile(remotePath string) *OperationResult {
	if !s.existsConnection() {
		return &OperationResult{Success: false, Message: "sftp client is not connected"}
	}

	err := s.client.Remove(remotePath)
	if err != nil {
		logger.Errorf("failed to delete remote file (%s): %v", remotePath, err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("failed to delete remote file (%s): %v", remotePath, err)}
	}
	return &OperationResult{Success: true, Message: "file deleted successfully"}
}

func (s *SFTPClient) Close() *OperationResult {
	if !s.existsConnection() {
		return &OperationResult{Success: true, Message: "sftp client is already closed"}
	}
	err := s.client.Close()
	if err != nil {
		logger.Errorf("error closing sftp client: %v", err)
		return &OperationResult{Success: false, Message: fmt.Sprintf("error closing sftp client: %v", err)}
	}
	return &OperationResult{Success: true, Message: "sftp client closed successfully"}
}

func (s *SFTPClient) existsConnection() bool {
	if s.client == nil {
		logger.Errorf("sftp client is not connected")
		return false
	}
	return true
}
