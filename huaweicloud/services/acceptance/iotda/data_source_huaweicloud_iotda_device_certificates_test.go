package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceCertificates_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_certificates.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.cn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.owner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.verify_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.effective_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificates.0.expiry_date"),

					resource.TestCheckOutput("is_certificate_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceDeviceCertificates_base() string {
	name := acceptance.RandomAccResourceName()
	spaceBasic := testSpace_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_device_certificate" "test" {
  content = <<EOT
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
`, spaceBasic)
}

func testAccDataSourceDeviceCertificates_basic() string {
	deviceCertificatesBase := testDataSourceDeviceCertificates_base()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_device_certificates" "test" {
  depends_on = [
    huaweicloud_iotda_device_certificate.test,
  ]
}

# Filter using certificate ID.
locals {
  certificate_id = data.huaweicloud_iotda_device_certificates.test.certificates[0].id
}

data "huaweicloud_iotda_device_certificates" "certificate_id_filter" {
  certificate_id = local.certificate_id
}

output "is_certificate_id_filter_useful" {
  value = length(data.huaweicloud_iotda_device_certificates.certificate_id_filter.certificates) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_certificates.certificate_id_filter.certificates[*].id : v == local.certificate_id]
  )
}

# Filter using status.
locals {
  status = data.huaweicloud_iotda_device_certificates.test.certificates[0].status
}

data "huaweicloud_iotda_device_certificates" "status_filter" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_iotda_device_certificates.status_filter.certificates) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_certificates.status_filter.certificates[*].status : v == local.status]
  )
}

# Filter using non existent certificate ID.
data "huaweicloud_iotda_device_certificates" "not_found" {
  certificate_id = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_device_certificates.not_found.certificates) == 0
}
`, deviceCertificatesBase)
}
