package ali

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"one-api/dto"
	"one-api/relay/channel"
	"one-api/relay/channel/openai"
	relaycommon "one-api/relay/common"
	"one-api/relay/constant"

	"github.com/gin-gonic/gin"
)

type Adaptor struct {
}

func (a *Adaptor) ConvertClaudeRequest(*gin.Context, *relaycommon.RelayInfo, *dto.ClaudeRequest) (any, error) {
	//TODO implement me
	panic("implement me")
	return nil, nil
}

func (a *Adaptor) Init(info *relaycommon.RelayInfo) {
}

func (a *Adaptor) GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	var fullRequestURL string
	switch info.RelayMode {
	case constant.RelayModeEmbeddings:
		fullRequestURL = fmt.Sprintf("%s/api/v1/services/embeddings/text-embedding/text-embedding", info.BaseUrl)
	case constant.RelayModeImagesGenerations:
		fullRequestURL = fmt.Sprintf("%s/api/v1/services/aigc/text2image/image-synthesis", info.BaseUrl)
	case constant.RelayModeFile:
		fullRequestURL = fmt.Sprintf("%s/compatible-mode/v1/files", info.BaseUrl)
	default:
		fullRequestURL = fmt.Sprintf("%s/compatible-mode/v1/chat/completions", info.BaseUrl)
	}
	return fullRequestURL, nil
}

func (a *Adaptor) SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error {
	channel.SetupApiRequestHeader(info, c, req)
	req.Set("Authorization", "Bearer "+info.ApiKey)
	if info.IsStream {
		req.Set("X-DashScope-SSE", "enable")
	}
	if c.GetString("plugin") != "" {
		req.Set("X-DashScope-Plugin", c.GetString("plugin"))
	}
	if info.RelayMode == constant.RelayModeImagesGenerations {
		req.Set("X-DashScope-Async", "enable")
	}
	return nil
}

func (a *Adaptor) ConvertOpenAIRequest(c *gin.Context, info *relaycommon.RelayInfo, request *dto.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	// 处理文本对话的请求
	fileMsg2Ai(request)

	switch info.RelayMode {
	default:
		aliReq := requestOpenAI2Ali(*request)
		return aliReq, nil
	}
}

func (a *Adaptor) ConvertImageRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.ImageRequest) (any, error) {
	aliRequest := oaiImage2Ali(request)
	return aliRequest, nil
}

func (a *Adaptor) ConvertRerankRequest(c *gin.Context, relayMode int, request dto.RerankRequest) (any, error) {
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertEmbeddingRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.EmbeddingRequest) (any, error) {
	return embeddingRequestOpenAI2Ali(request), nil
}

func (a *Adaptor) ConvertAudioRequest(c *gin.Context, info *relaycommon.RelayInfo, request dto.AudioRequest) (io.Reader, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (a *Adaptor) ConvertFileRequest(c *gin.Context, info *relaycommon.RelayInfo) (io.Reader, error) {
	if c.Request.Method == http.MethodGet {
		return nil, nil
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 获取所有表单字段
	formData := c.Request.PostForm

	// 遍历表单字段并打印输出
	for key, values := range formData {
		if key == "model" {
			continue
		}
		for _, value := range values {
			writer.WriteField(key, value)
		}
	}

	// 添加文件字段
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return nil, errors.New("file is required")
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return nil, errors.New("create form file failed")
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, errors.New("copy file failed")
	}

	// 关闭 multipart 编写器以设置分界线
	writer.Close()
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	return &requestBody, nil
}

func (a *Adaptor) DoRequest(c *gin.Context, info *relaycommon.RelayInfo, requestBody io.Reader) (any, error) {
	if info.RelayMode == constant.RelayModeFile {
		return channel.DoFormRequest(a, c, info, requestBody)
	}
	return channel.DoApiRequest(a, c, info, requestBody)
}

func (a *Adaptor) DoResponse(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (usage any, err *dto.OpenAIErrorWithStatusCode) {
	switch info.RelayMode {
	case constant.RelayModeImagesGenerations:
		err, usage = aliImageHandler(c, resp, info)
	case constant.RelayModeEmbeddings:
		err, usage = aliEmbeddingHandler(c, resp)
	case constant.RelayModeFile:
		err, usage = openai.OpenaiFileHandler(c, resp, info)
	default:
		if info.IsStream {
			err, usage = openai.OaiStreamHandler(c, resp, info)
		} else {
			err, usage = openai.OpenaiHandler(c, resp, info)
		}
	}
	return
}

func (a *Adaptor) GetModelList() []string {
	return ModelList
}

func (a *Adaptor) GetChannelName() string {
	return ChannelName
}
