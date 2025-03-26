// SPDX-FileCopyrightText: 2025 Industria de Diseño Textil S.A. INDITEX
//
// SPDX-License-Identifier: Apache-2.0

package xk6sftp

type OperationResult struct {
	Success bool
	Message string
}

type SFTPClientInterface interface {
	UploadFile(localPath, remotePath string) *OperationResult
	DownloadFile(remotePath, localPath string) *OperationResult
	DeleteFile(remotePath string) *OperationResult
	Close() *OperationResult
}
