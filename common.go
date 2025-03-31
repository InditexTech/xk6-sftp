// SPDX-FileCopyrightText: 2025 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

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
