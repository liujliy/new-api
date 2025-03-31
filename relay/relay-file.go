package relay

import (
	"fmt"
	"net/http"
	"one-api/dto"
	"one-api/model"
	relaycommon "one-api/relay/common"
	"one-api/service"

	"github.com/gin-gonic/gin"
)

func FileHelper(c *gin.Context) (openaiErr *dto.OpenAIErrorWithStatusCode) {
	relayInfo := relaycommon.GenRelayInfo(c)
	// 设置file_id
	if c.Request.Method == http.MethodGet {
		relayInfo.FileID = c.Param("id")
	}
	adaptor := GetAdaptor(relayInfo.ApiType)
	if adaptor == nil {
		return service.OpenAIErrorWrapperLocal(fmt.Errorf("invalid api type: %d", relayInfo.ApiType), "invalid_api_type", http.StatusBadRequest)
	}

	adaptor.Init(relayInfo)

	ioReader, err := adaptor.ConvertFileRequest(c, relayInfo)
	if err != nil {
		return service.OpenAIErrorWrapperLocal(err, "convert_request_failed", http.StatusInternalServerError)
	}

	resp, err := adaptor.DoRequest(c, relayInfo, ioReader)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "do_request_failed", http.StatusInternalServerError)
	}
	statusCodeMappingStr := c.GetString("status_code_mapping")

	var httpResp *http.Response
	if resp != nil {
		httpResp = resp.(*http.Response)
		if httpResp.StatusCode != http.StatusOK {
			openaiErr = service.RelayErrorHandler(httpResp, false)
			// reset status code 重置状态码
			service.ResetStatusCode(openaiErr, statusCodeMappingStr)
			return openaiErr
		}
	}
	fileInfo, openaiErr := adaptor.DoResponse(c, httpResp, relayInfo)
	if openaiErr != nil {
		// reset status code 重置状态码
		service.ResetStatusCode(openaiErr, statusCodeMappingStr)
		return openaiErr
	}
	if fileInfo == nil {
		return nil
	}
	// 保存用户上传的文件
	// 渠道对应的文件ID存在则更新状态，不存在则新增
	file := &model.File{
		UserID:      relayInfo.UserId,
		Username:    c.GetString("username"),
		ChannelId:   relayInfo.ChannelId,
		ChannelName: c.GetString("channel_name"),
		FileID:      fileInfo.(*dto.FileResponse).Id,
		FileName:    fileInfo.(*dto.FileResponse).Filename,
		Status:      fileInfo.(*dto.FileResponse).Status,
	}
	file.Insert()

	return nil
}
