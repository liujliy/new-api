package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/volcengine/volc-sdk-golang/base"
)

func ProxyAIGC(c *gin.Context) {
	action := c.Query("Action")
	version := c.Query("Version")

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://rtc.volcengineapi.com?Action="+action+"&Version="+version, strings.NewReader(string(bodyBytes)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	credential := base.Credentials{
		Region:          "cn-north-1",
		Service:         "rtc",
		AccessKeyID:     os.Getenv("VOLC_ACCESSKEY"),
		SecretAccessKey: os.Getenv("VOLC_SECRETKEY"),
	}
	signedReq := credential.Sign(req)

	resp, err := client.Do(signedReq)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"message": resp.Status,
			"success": false,
		})
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	err = resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
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
		"data":    result,
	})
}
