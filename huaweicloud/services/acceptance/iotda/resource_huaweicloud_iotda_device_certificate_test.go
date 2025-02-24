package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
)

func getDeviceCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	return iotda.QueryDeviceCertificate(client, state.Primary.ID, "")
}

func TestAccDeviceCertificate_basic(t *testing.T) {
	var (
		deviceCaObj interface{}
		rName       = "huaweicloud_iotda_device_certificate.test"
		name        = acceptance.RandomAccResourceName()
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&deviceCaObj,
		getDeviceCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceCertificate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cn", "huaweiIoT"),
					resource.TestCheckResourceAttr(rName, "status", "Unverified"),
					resource.TestCheckResourceAttrSet(rName, "verify_code"),
					resource.TestCheckResourceAttrSet(rName, "effective_date"),
					resource.TestCheckResourceAttrSet(rName, "expiry_date"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"content", "space_id"},
			},
		},
	})
}

func testDeviceCertificate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_certificate" "test" {
  space_id = huaweicloud_iotda_space.test.id
  content  = <<EOT
-----BEGIN CERTIFICATE-----
MIIDlTCCAn0CFDksqsC4D2sWd5aIJ/3kveD9bi1VMA0GCSqGSIb3DQEBCwUAMIGF
MQswCQYDVQQGEwJDTjEQMA4GA1UECAwHYmVpamluZzEQMA4GA1UEBwwHYmVpamlu
ZzEPMA0GA1UECgwGaHVhd2VpMQ4wDAYDVQQLDAVjbG91ZDESMBAGA1UEAwwJaHVh
d2VpSW9UMR0wGwYJKoZIhvcNAQkBFg54eHhAaHVhd2VpLmNvbTAgFw0yMjA2MjIw
MjAyNTVaGA8yMTU5MDUxNTAyMDI1NVowgYUxCzAJBgNVBAYTAkNOMRAwDgYDVQQI
DAdiZWlqaW5nMRAwDgYDVQQHDAdiZWlqaW5nMQ8wDQYDVQQKDAZodWF3ZWkxDjAM
BgNVBAsMBWNsb3VkMRIwEAYDVQQDDAlodWF3ZWlJb1QxHTAbBgkqhkiG9w0BCQEW
Dnh4eEBodWF3ZWkuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
wKVJ6SpnhcBj3A5AEBrykKgZHSnU6+EqmTHoKGHFLjIMQuVr8cv8a8j7gDINc0gQ
L+UE3D3JdnIxarUH7/WWLjqKqqzTkcNOFi5Mue6JVtGKLJ/3mEsTFAlEz9ArRz+0
SCYgMJDQ9Z7EYY82NSjpVf8/0YyqOZeQMOGn58i/yxAxqzK4ux1l410JOH55kz7k
bwW9oXmnyTaZqfvT4JsQ4U9JwkdkRcJlzxR8w24t/s+qE4/KM7avlCw1oMvTnkTO
uJy9UEWCTCOVyC7+CB9QmysjkK3DZgH/ER2Q6sjSmzIEkpcqJpt4cYTcFL3nDLC7
TfIr5SmDZ4Q3Z2CYqI6/oQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCffi7j5Zsz
LTCJsvwplZqjXPnTHSzqvMoD96mf7qSObqsdVh9KaM9sjg49wYtKANqIa4l/oO+t
qtu+C7evY9M4X4EBOdDTG6OH7U8YOoIPAUgVZ81zhvFmXbCUb+ekOLwBWlccIwtN
l4223PfNnA8HZhQjpAtfT0DS5pAEmQZ8X2ByqAXXMWgWGgQRh4vyvbd7MMTCr6fp
zTZSnp58+Rqlf3Sqjps3BJydjxURvLpoB+WvByWhafXoxeOu3yEZQrZF+rzu6yaB
BWlaAkI9PmPD+h2ozFiXH7nzABrAXFOk1AaRFM+cAJLnuVyWjX/hv5rz9Z4FfWJu
RQocSUkUw0EW
-----END CERTIFICATE-----
EOT
}
`, testSpace_basic(name))
}
