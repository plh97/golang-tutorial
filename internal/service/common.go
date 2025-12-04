package service

import "go-nunu/pkg/aws" // 导入 R2 客户端

type CommonService interface {
	UploadPresignedUrl(fileExt string, scene int) (string, string, error)
}

type commonService struct {
	// R2 客户端应该是直接依赖，而不是通过 *Service 间接依赖
	R2Client *aws.CloudflareR2
	// *Service // 如果不需要 *Service 的其他字段，就移除它
}

// NewCommonService 的构造函数，直接接收 R2 客户端
func NewCommonService(r2Client *aws.CloudflareR2) CommonService {
	return &commonService{
		R2Client: r2Client,
		// Service: service, // 如果保留 *Service，需要确保它也作为参数传入
	}
}

// UploadPresignedUrl generates a presigned URL for uploading a file to R2
func (c *commonService) UploadPresignedUrl(fileExt string, scene int) (string, string, error) {
	// TODO: Implement the presigned URL generation logic
	return c.R2Client.UploadPresignedUrl(fileExt, scene)
}
