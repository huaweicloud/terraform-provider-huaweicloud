package dew

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getKmsDedicatedKeystoreResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1.0/{project_id}/keystores/{keystore_id}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{keystore_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS dedicated keystore: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccKmsDedicatedKeystore_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_kms_dedicated_keystore.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsDedicatedKeystoreResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsHsmClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKmsDedicatedKeystore_basic(name),
				ExpectError: regexp.MustCompile(` error creating KMS dedicated keystore: Internal Server Error`),
			},
		},
	})
}

func testKmsDedicatedKeystore_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_dedicated_keystore" "test" {
  alias          = "%s"
  hsm_cluster_id = "%s"
  hsm_ca_cert    = <<EOT
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUehx07qc7un7IB7/X9lHCLkt/jPowDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTA1MzEwOTI1NTJaFw0yMjA1
MzEwOTI1NTJaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCvmuH5ViGtGOlevJ8vOoN3Ak4pp3SescdAfQa/r4cO
z/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc8OrcaIsjns92XITVDpFW0ThGyjhT
ZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK73OUYQY2E6l44U9G8Id763Bnws9N
Rn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu9zu8c8i2+8qLjEsonx5PrwzNlYP3
JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd5F4Rdxl+SAIY+6mr7qu1dAlcVMLS
QcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iKaedPo2GfAgMBAAGjUzBRMB0GA1Ud
DgQWBBR5yzB/GujpSlLrn0l2p+BslakGzjAfBgNVHSMEGDAWgBR5yzB/GujpSlLr
n0l2p+BslakGzjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCj
TqvcIk0Au/yOxOfIGUZzVkTiORxbwATAfRN6n/+mrgWnIbHG4XFqqjmFr7gGvHeH
+BuyU06VXJgKYaPUqbYl7eBd4Spm5v3Wq7C7i96dOHmG8fVcjQnTWleyEmUsEarv
A6/lhTqXV1+AuNUaH+9EbBUBsrCHGLkECBMKl0+cJN8lo5XncAtp7z1+O/Mn0Zi6
XyNOyvqcmmn8HUkSIS4RlJ2ohuZN6oFC3sYX9g9Vo++IkjGl3dRbf/7JutqBGHNE
RVKoPyaivymDDIIL/qSy/Pi2s0hzUhwc1M8td0K/AMxyeigwNG7mTH0RzX32bUkf
ZoURg5WiRskhtHEvBsLF
-----END CERTIFICATE-----
EOT
}
`, name, acceptance.HW_KMS_HSM_CLUSTER_ID)
}
