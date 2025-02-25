package config

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	versionURL     = "https://provider.obs.cn-north-4.myhuaweicloud.com/terraform-provider-huaweicloud/version.json"
	defaultSummary = "Newer Version Available"
	defaultDetail  = "The provider version is out of date, please upgrade to the latest version."
)

type VersionConfig struct {
	Latest  string                   `json:"latest"`
	Version map[string]VersionConfig `json:"version"`
	Warning string                   `json:"warning"`
	Message string                   `json:"message"`
}

func CheckUpgrade(d *schema.ResourceData, version string) diag.Diagnostics {
	if d.Get("skip_check_upgrade").(bool) {
		log.Printf("[WARN] check upgrade skipped")
		return nil
	}

	log.Printf("[DEBUG] current version: %s", version)
	if version == "" {
		return nil
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] CheckUpgrade recover: %#v", r)
		}
	}()

	body, _, err := httpRequest("GET", versionURL, nil, nil)
	if err != nil {
		log.Printf("[WARN] failed to query version config : %s", err)
		return nil
	}

	verCfg := &VersionConfig{}
	err = json.Unmarshal(body, verCfg)
	if err != nil {
		log.Printf("[WARN] failed to unmarshal version config : %s", err)
		return nil
	}

	if verCfg.Latest == "" {
		log.Printf("[WARN] failed to get latest version")
		return nil
	}

	if compareVersion(version, verCfg.Latest) != -1 {
		return nil
	}

	return buildMsg(version, verCfg)
}

func buildMsg(version string, verCfg *VersionConfig) diag.Diagnostics {
	detail := defaultDetail
	summary := defaultSummary
	latest := verCfg.Latest

	if v, ok := verCfg.Version[version]; ok {
		verCfg = &v
	}

	if verCfg.Warning != "" {
		summary = verCfg.Warning
	}
	if verCfg.Message != "" {
		detail = verCfg.Message
	}

	summary = strings.ReplaceAll(summary, "${latest}", latest)
	summary = strings.ReplaceAll(summary, "${version}", version)
	detail = strings.ReplaceAll(detail, "${latest}", latest)
	detail = strings.ReplaceAll(detail, "${version}", version)
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  summary,
			Detail:   detail,
		},
	}
}

// compareVersion compares two version strings and returns:
// ver1 < ver2 -> -1
// ver1 > ver2 -> 1
// ver1 == ver2 -> 0
func compareVersion(ver1, ver2 string) int {
	parts1 := strings.Split(ver1, ".")
	parts2 := strings.Split(ver2, ".")

	if len(parts1) != len(parts2) {
		if len(parts1) > len(parts2) {
			return 1
		}
		return -1
	}

	for i := range parts1 {
		num1, _ := strconv.Atoi(parts1[i])
		num2, _ := strconv.Atoi(parts2[i])

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	return 0
}

func httpRequest(method, requestUrl string, jsonBody any, headers map[string]string) ([]byte, int, error) {
	var body io.Reader
	var contentType string
	if jsonBody != nil {
		rendered, err := json.Marshal(jsonBody)
		if err != nil {
			return nil, 0, err
		}
		body = bytes.NewReader(rendered)
		contentType = "application/json"
	}

	client := http.Client{
		Transport: &LogRoundTripper{
			Rt: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: true, //nolint:gosec
				},
			},
			MaxRetries: 1,
		},
	}
	defer client.CloseIdleConnections()

	request, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		return nil, 0, err
	}

	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("[WARN] failed to close response: %s", err)
		}
	}()

	bodyByte, err := io.ReadAll(response.Body)

	return bodyByte, response.StatusCode, err
}
