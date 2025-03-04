package controller

import (
	"net/http"
	"os"
	"time"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"

	"github.com/gin-gonic/gin"
)

func GetCosSTS(c *gin.Context) {
	// 获取临时密钥
	client := sts.NewClient(
		os.Getenv("COS_SECRET_ID"),
		os.Getenv("COS_SECRET_KEY"),
		nil,
	)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          os.Getenv("COS_REGION"),
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						"*", // 代表所有资源
					},
				},
			},
		},
	}

	resp, err := client.GetCredential(opt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "",
		"success": true,
		"data":    resp,
	})
}
