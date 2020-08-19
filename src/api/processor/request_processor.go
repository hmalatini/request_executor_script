package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hmalatini/request_executor_script/src/api/logger"
)

const className = "RequestProcessor"

var (
	baseUrl string
	uri     string
	body    string
	headers map[string]string
	method  string
)

type RequestProcessor struct {
}

func NewRequestProcessor() *RequestProcessor {
	return &RequestProcessor{}
}

func (p *RequestProcessor) InitProcessor(baseUrlParam string, uriParam string, bodyParam string, headersParam map[string]string, methodParam string) {
	baseUrl = baseUrlParam
	uri = uriParam
	body = bodyParam
	headers = headersParam
	method = methodParam

	logger.LogDebug(className, "Request Processor Initialized")
}

func (p *RequestProcessor) Process(record []string) (bool, string) {
	parsedUri, err := p.parseString(uri, record)
	if err != nil {
		msg := fmt.Sprintf("Error parsing URI for record %s: %s", record, err.Error())
		logger.LogError(className, msg)
		return false, msg
	}

	logger.LogDebug(className, fmt.Sprintf("URI parsed successfully: %s", parsedUri))

	var requestBody *bytes.Buffer
	if body != "" {
		parsedBody, err := p.parseString(body, record)
		if err != nil {
			msg := fmt.Sprintf("Error parsing Body for record %s: %s", record, err.Error())
			logger.LogError(className, msg)
			return false, msg
		}

		logger.LogDebug(className, fmt.Sprintf("Body parsed successfully: %s", parsedBody))

		requestBody = bytes.NewBuffer([]byte(parsedBody))
	} else {
		requestBody = nil
	}

	url := fmt.Sprintf("%s%s", baseUrl, parsedUri)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, requestBody)
	if req == nil {
		return false, fmt.Sprintf("Error creating request for record %s: request is nil", record)
	}
	if err != nil {
		return false, fmt.Sprintf("Error creating request for record %s: %s", err.Error())
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	logger.LogDebug(className, "Request created successfully")

	resp, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("Error making the request for record %s: %s", record, err.Error())
		logger.LogError(className, msg)
		return false, msg
	}
	defer resp.Body.Close()

	logger.LogDebug(className, fmt.Sprintf("Request made. Response: %+v", resp))

	var bodyString string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogWarning(className, fmt.Sprintf("Cant read the body for record %s", record))
	}
	bodyString = string(bodyBytes)

	description := fmt.Sprintf("StatusCode: %d - Url: %s - Body: %s", resp.StatusCode, url, bodyString)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.LogDebug(className, fmt.Sprintf("Request get error. Response: %s", description))
		return false, description
	}

	logger.LogDebug(className, fmt.Sprintf("Request get success. Response: %s", description))
	return true, description
}

func (p *RequestProcessor) GetParsedUri(record []string) (string, error) {
	parsedUri, err := p.parseString(uri, record)
	if err != nil {
		msg := fmt.Sprintf("Error parsing URI for record %s: %s", record, err.Error())
		logger.LogError(className, msg)
		return "", err
	}

	return parsedUri, nil
}

func (p *RequestProcessor) GetParsedBody(record []string) (string, error) {
	parsedUri, err := p.parseString(body, record)
	if err != nil {
		msg := fmt.Sprintf("Error parsing BODY for record %s: %s", record, err.Error())
		logger.LogError(className, msg)
		return "", err
	}

	return parsedUri, nil
}

func (p *RequestProcessor) parseString(str string, record []string) (string, error) {
	var sb strings.Builder
	for i := 0; i < len(str); i++ {
		char := str[i]
		if char != '%' {
			sb.WriteByte(char)
		} else {
			i++
			var nb strings.Builder
			for j := i; j < len(str); j++ {
				numbChar := str[j]
				if numbChar >= '0' && numbChar <= '9' {
					nb.WriteByte(numbChar)
				} else {
					pos, err := strconv.Atoi(nb.String())
					if err != nil {
						return "", err
					}

					if pos > len(record) {
						return "", fmt.Errorf("not enough data in record")
					}

					sb.WriteString(record[pos-1])
					i = j - 1
					break
				}
			}

			if i == len(str)-1 {
				pos, err := strconv.Atoi(nb.String())
				if err != nil {
					return "", err
				}

				if pos > len(record) {
					return "", fmt.Errorf("not enough data in record")
				}

				sb.WriteString(record[pos-1])
			}
		}
	}

	return sb.String(), nil
}
