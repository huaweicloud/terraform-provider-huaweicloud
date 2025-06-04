package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// MAXFieldLength is the maximum string length of single field when logging
const MAXFieldLength int = 1024

var logAtomicId int64
var maxTimeout = 10 * time.Minute

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default http client RoundTripper to allow for logging.
type LogRoundTripper struct {
	Rt         http.RoundTripper
	MaxRetries int
}

func retryTimeout(count int) time.Duration {
	seconds := math.Pow(2, float64(count))
	timeout := time.Duration(seconds) * time.Second
	if timeout > maxTimeout { // won't wait more than maxTimeout
		timeout = maxTimeout
	}
	return timeout
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	var response *http.Response
	var bs bytes.Buffer

	atomicId := atomic.AddInt64(&logAtomicId, 1)
	logId := fmt.Sprintf("%d-%d", time.Now().UnixMilli(), atomicId)

	defer func() {
		// logging the API request and response
		var logErr error
		if request != nil {
			log.Printf("[DEBUG] [%s] API Request URL: %s %s\nAPI Request Headers:\n%s",
				logId, request.Method, request.URL, FormatHeaders(request.Header, "\n"))

			if request.Body != nil {
				logErr = lrt.logRequest(&bs, request.Header.Get("Content-Type"), logId)
				if logErr != nil {
					log.Printf("[WARN] [%s] failed to log API Request Body: %s", logId, logErr)
				}

				request.Body.Close()
			}
		}

		if response != nil {
			log.Printf("[DEBUG] [%s] API Response Code: %d\nAPI Response Headers:\n%s",
				logId, response.StatusCode, FormatHeaders(response.Header, "\n"))

			if response.Body != nil {
				response.Body, logErr = lrt.logResponse(response.Body, response.Header.Get("Content-Type"), logId)
				if logErr != nil {
					log.Printf("[WARN] [%s] failed to log API Response Body: %s", logId, logErr)
				}
			}
		}
	}()

	if request.Body != nil {
		request.Body, err = lrt.dumpRequest(request.Body, &bs)
		if err != nil {
			return nil, err
		}
	}

	// executes a single HTTP transaction
	response, err = lrt.Rt.RoundTrip(request)
	if response == nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, "no such host") {
			return nil, err
		}
	}

	// Retrying connection
	retry := 1
	for response == nil {
		if retry > lrt.MaxRetries {
			log.Printf("[DEBUG] [%s] connection error, retries exhausted. Aborting", logId)
			err = fmt.Errorf("connection error, retries exhausted. Aborting. Last error was: %s", err)
			return nil, err
		}

		log.Printf("[DEBUG] [%s] connection error, retry number %d: %s", logId, retry, err)

		//lintignore:R018
		time.Sleep(retryTimeout(retry))
		response, err = lrt.Rt.RoundTrip(request)
		retry++
	}

	// retry connection reset by peer error
	retry = 1
	for err != nil && strings.Contains(err.Error(), "connection reset by peer") {
		if retry > lrt.MaxRetries {
			log.Printf("[DEBUG] [%s] connection error, retries exhausted. Aborting", logId)
			err = fmt.Errorf("connection error, retries exhausted. Aborting. Last error was: %s", err)
			return nil, err
		}

		log.Printf("[DEBUG] [%s] connection error, retry number %d: %s", logId, retry, err)

		// lintignore:R018
		time.Sleep(retryTimeout(retry))
		response, err = lrt.Rt.RoundTrip(request)
		retry++
	}

	return response, err
}

// dumpRequest will copy the HTTP Request details to buffer, then close the original.
func (*LogRoundTripper) dumpRequest(original io.ReadCloser, bs *bytes.Buffer) (io.ReadCloser, error) {
	defer original.Close()

	_, err := io.Copy(bs, original)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(strings.NewReader(bs.String())), nil
}

// logRequest will log the HTTP Request details.
// If the body is JSON, it will attempt to be pretty-formatted.
func (*LogRoundTripper) logRequest(bs *bytes.Buffer, contentType, logId string) error {
	isJSONFormat := regexp.MustCompile(`^application/(merge-patch\+)?json`).Match([]byte(contentType))
	isXMLFormat := strings.HasPrefix(bs.String(), "<") && !strings.HasPrefix(bs.String(), "<html>")
	// Handle request contentType
	switch {
	case isJSONFormat:
		debugInfo := formatJSON(bs.Bytes(), logId, true)
		log.Printf("[DEBUG] [%s] API Request Body: %s", logId, debugInfo)
	case isXMLFormat:
		log.Printf("[DEBUG] [%s] API Request Body: %s", logId, bs.String())
	default:
		log.Printf("[DEBUG] [%s] Not logging because the request body isn't JSON or XML format", logId)
	}

	return nil
}

// logResponse will log the HTTP Response details, then close the original and build a new ReadCloser.
// If the body is JSON, it will attempt to be pretty-formatted.
func (*LogRoundTripper) logResponse(original io.ReadCloser, contentType, logId string) (io.ReadCloser, error) {
	defer original.Close()

	var bs bytes.Buffer
	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	isJSONFormat := strings.HasPrefix(contentType, "application/json")
	isXMLFormat := strings.HasPrefix(contentType, "application/xml")
	switch {
	case isJSONFormat:
		debugInfo := formatJSON(bs.Bytes(), logId, true)
		log.Printf("[DEBUG] [%s] API Response Body: %s", logId, debugInfo)
	case isXMLFormat:
		log.Printf("[DEBUG] [%s] API Response Body: %s", logId, bs.String())
	default:
		log.Printf("[DEBUG] [%s] Not logging because the response body isn't JSON or XML format", logId)
	}

	return io.NopCloser(strings.NewReader(bs.String())), nil
}

// formatJSON will try to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func formatJSON(raw []byte, logId string, maskBody bool) string {
	var data map[string]interface{}

	if len(raw) == 0 {
		return ""
	}

	err := json.Unmarshal(raw, &data)
	if err != nil {
		log.Printf("[DEBUG] [%s] Unable to parse JSON: %s", logId, err)
		return string(raw)
	}

	// Mask known password fields
	if maskBody {
		maskSecurityFields(data)
	}

	// Ignore the catalog
	if _, ok := data["catalog"]; ok {
		return "{ **skipped** }"
	}
	if v, ok := data["token"].(map[string]interface{}); ok {
		if _, ok := v["catalog"]; ok {
			return ""
		}
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("[DEBUG] [%s] Unable to re-marshal JSON: %s", logId, err)
		return string(raw)
	}

	return string(pretty)
}

// RedactHeaders processes a headers object, returning a redacted list.
func RedactHeaders(headers http.Header) (processedHeaders []string) {
	// sensitiveWords is a list of headers that need to be redacted.
	var sensitiveWords = []string{"token", "authorization"}

	for name, header := range headers {
		for _, v := range header {
			if utils.IsStrContainsSliceElement(name, sensitiveWords, true, false) {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, "***"))
			} else {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, v))
			}
		}
	}
	return
}

// FormatHeaders processes a headers object plus a deliminator, returning a string
func FormatHeaders(headers http.Header, seperator string) string {
	redactedHeaders := RedactHeaders(headers)
	sort.Strings(redactedHeaders)

	return strings.Join(redactedHeaders, seperator)
}

func maskSecurityFields(data map[string]interface{}) {
	for k, val := range data {
		switch val := val.(type) {
		case string:
			if isSecurityFields(k) {
				data[k] = "***"
			} else if len(val) > MAXFieldLength {
				data[k] = "** large string **"
			}
		case map[string]interface{}:
			if isSecurityFields(k) {
				data[k] = map[string]string{"***": "***"}
			} else {
				maskSecurityFields(val)
			}
		}
	}
}

func isSecurityFields(field string) bool {
	checkField := strings.ToLower(field)
	// 'password' is apply to the most request JSON body.
	// 'secret' is apply to the AK/SK response JSON body.
	// 'pwd' and 'token' is the high frequency sensitive keywords in the request and response bodies.
	if strings.Contains(checkField, "password") || strings.Contains(checkField, "secret") ||
		strings.HasSuffix(checkField, "pwd") || strings.HasSuffix(checkField, "token") {
		return true
	}

	// 'adminpass' is apply to the ecs/bms instance request JSON body
	// 'encrypted_user_data' is apply to the function request JSON body of FunctionGraph
	// 'nonce' is apply to the random string for authorization methods.
	// 'email', 'phone', 'phone_number', 'phone_num' and 'sip_number' can uniquely identify a person.
	// 'signature' are used for encryption.
	// 'user_passwd' is apply to the dms/kafka user request JSON body
	// 'auth' is apply to kms keypairs associate or disassociate request JSON body
	// 'cert_content', 'private_key' and 'trusted_root_ca' are both sensitive parameters of the SSL certificate for APIG
	// 'sk', 'src_sk' and 'dst_sk' are used in oms_task and oms_task_group
	// request JSON body
	securityFields := []string{"adminpass", "encrypted_user_data", "nonce", "email", "phone", "phone_number", "phone_num",
		"sip_number", "signature", "user_passwd", "auth", "cert_content", "private_key", "trusted_root_ca", "sk", "src_sk",
		"dst_sk", "pwd", "plain_text", "plain_text_base64", "cipher_text"}
	return utils.StrSliceContains(securityFields, checkField)
}
