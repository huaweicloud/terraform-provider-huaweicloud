package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbCertificates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_certificates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataaSourceElbCertificates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.subject_alternative_names.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.fingerprint"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.domain"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.private_key"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.common_name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enc_certificate"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.scm_certificate_id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enc_private_key"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.expire_time"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.updated_at"),
					resource.TestCheckOutput("certificate_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("domain_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("common_name_filter_is_useful", "true"),
					resource.TestCheckOutput("fingerprint_filter_is_useful", "true"),
					resource.TestCheckOutput("scm_certificate_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataaSourceElbCertificates_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_certificate" "test" {
  name               = "%s"
  description        = "terraform test certificate"
  domain             = "www.elb.com"
  scm_certificate_id = "scs1702951276162"

  private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwZ5UJULAjWr7p6FVwGRQRjFN2s8tZ/6LC3X82fajpVsYqF1x
qEuUDndDXVD09E4u83MS6HO6a3bIVQDp6/klnYldiE6Vp8HH5BSKaCWKVg8lGWg1
UM9wZFnlryi14KgmpIFmcu9nA8yV/6MZAe6RSDmb3iyNBmiZ8aZhGw2pI1YwR+15
MVqFFGB+7ExkziROi7L8CFCyCezK2/oOOvQsH1dzQ8z1JXWdg8/9Zx7Ktvgwu5PQ
M3cJtSHX6iBPOkMU8Z8TugLlTqQXKZOEgwajwvQ5mf2DPkVgM08XAgaLJcLigwD5
13koAdtJd5v+9irw+5LAuO3JclqwTvwy7u/YwwIDAQABAoIBACU9S5fjD9/jTMXA
DRs08A+gGgZUxLn0xk+NAPX3LyB1tfdkCaFB8BccLzO6h3KZuwQOBPv6jkdvEDbx
Nwyw3eA/9GJsIvKiHc0rejdvyPymaw9I8MA7NbXHaJrY7KpqDQyk6sx+aUTcy5jg
iMXLWdwXYHhJ/1HVOo603oZyiS6HZeYU089NDUcX+1SJi3e5Ke0gPVXEqCq1O11/
rh24bMxnwZo4PKBWdcMBN5Zf/4ij9vrZE+fFzW7vGBO48A5lvZxWU2U5t/OZQRtN
1uLOHmMFa0FIF2aWbTVfwdUWAFsvAOkHj9VV8BXOUwKOUuEktdkfAlvrxmsFrO/H
yDeYYPkCgYEA/S55CBbR0sMXpSZ56uRn8JHApZJhgkgvYr+FqDlJq/e92nAzf01P
RoEBUajwrnf1ycevN/SDfbtWzq2XJGqhWdJmtpO16b7KBsC6BdRcH6dnOYh31jgA
vABMIP3wzI4zSVTyxRE8LDuboytF1mSCeV5tHYPQTZNwrplDnLQhywcCgYEAw8Yc
Uk/eiFr3hfH/ZohMfV5p82Qp7DNIGRzw8YtVG/3+vNXrAXW1VhugNhQY6L+zLtJC
aKn84ooup0m3YCg0hvINqJuvzfsuzQgtjTXyaE0cEwsjUusOmiuj09vVx/3U7siK
Hdjd2ICPCvQ6Q8tdi8jV320gMs05AtaBkZdsiWUCgYEAtLw4Kk4f+xTKDFsrLUNf
75wcqhWVBiwBp7yQ7UX4EYsJPKZcHMRTk0EEcAbpyaJZE3I44vjp5ReXIHNLMfPs
uvI34J4Rfot0LN3n7cFrAi2+wpNo+MOBwrNzpRmijGP2uKKrq4JiMjFbKV/6utGF
Up7VxfwS904JYpqGaZctiIECgYA1A6nZtF0riY6ry/uAdXpZHL8ONNqRZtWoT0kD
79otSVu5ISiRbaGcXsDExC52oKrSDAgFtbqQUiEOFg09UcXfoR6HwRkba2CiDwve
yHQLQI5Qrdxz8Mk0gIrNrSM4FAmcW9vi9z4kCbQyoC5C+4gqeUlJRpDIkQBWP2Y4
2ct/bQKBgHv8qCsQTZphOxc31BJPa2xVhuv18cEU3XLUrVfUZ/1f43JhLp7gynS2
ep++LKUi9D0VGXY8bqvfJjbECoCeu85vl8NpCXwe/LoVoIn+7KaVIZMwqoGMfgNl
nEqm7HWkNxHhf8A6En/IjleuddS1sf9e/x+TJN1Xhnt9W6pe7Fk1
-----END RSA PRIVATE KEY-----
EOT

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIJAKdmmOBYnFvoMA0GCSqGSIb3DQEBCwUAMGkxCzAJBgNV
BAYTAnh4MQswCQYDVQQIDAJ4eDELMAkGA1UEBwwCeHgxCzAJBgNVBAoMAnh4MQsw
CQYDVQQLDAJ4eDELMAkGA1UEAwwCeHgxGTAXBgkqhkiG9w0BCQEWCnh4QDE2My5j
b20wHhcNMTcxMjA0MDM0MjQ5WhcNMjAxMjAzMDM0MjQ5WjBpMQswCQYDVQQGEwJ4
eDELMAkGA1UECAwCeHgxCzAJBgNVBAcMAnh4MQswCQYDVQQKDAJ4eDELMAkGA1UE
CwwCeHgxCzAJBgNVBAMMAnh4MRkwFwYJKoZIhvcNAQkBFgp4eEAxNjMuY29tMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwZ5UJULAjWr7p6FVwGRQRjFN
2s8tZ/6LC3X82fajpVsYqF1xqEuUDndDXVD09E4u83MS6HO6a3bIVQDp6/klnYld
iE6Vp8HH5BSKaCWKVg8lGWg1UM9wZFnlryi14KgmpIFmcu9nA8yV/6MZAe6RSDmb
3iyNBmiZ8aZhGw2pI1YwR+15MVqFFGB+7ExkziROi7L8CFCyCezK2/oOOvQsH1dz
Q8z1JXWdg8/9Zx7Ktvgwu5PQM3cJtSHX6iBPOkMU8Z8TugLlTqQXKZOEgwajwvQ5
mf2DPkVgM08XAgaLJcLigwD513koAdtJd5v+9irw+5LAuO3JclqwTvwy7u/YwwID
AQABo1AwTjAdBgNVHQ4EFgQUo5A2tIu+bcUfvGTD7wmEkhXKFjcwHwYDVR0jBBgw
FoAUo5A2tIu+bcUfvGTD7wmEkhXKFjcwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAQEAWJ2rS6Mvlqk3GfEpboezx2J3X7l1z8Sxoqg6ntwB+rezvK3mc9H0
83qcVeUcoH+0A0lSHyFN4FvRQL6X1hEheHarYwJK4agb231vb5erasuGO463eYEG
r4SfTuOm7SyiV2xxbaBKrXJtpBp4WLL/s+LF+nklKjaOxkmxUX0sM4CTA7uFJypY
c8Tdr8lDDNqoUtMD8BrUCJi+7lmMXRcC3Qi3oZJW76ja+kZA5mKVFPd1ATih8TbA
i34R7EQDtFeiSvBdeKRsPp8c0KT8H1B4lXNkkCQs2WX5p4lm99+ZtLD4glw8x6Ic
i1YhgnQbn5E0hz55OLu5jvOkKQjPCW+8Kg==
-----END CERTIFICATE-----
EOT
}
`, name)
}

func testDataaSourceElbCertificates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_certificates" "test" {}

locals {
  certificate_id = huaweicloud_elb_certificate.test.id
}
data "huaweicloud_elb_certificates" "certificate_id_filter" {
  certificate_id = [huaweicloud_elb_certificate.test.id]
}
output "certificate_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.certificate_id_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.certificate_id_filter.certificates[*].id : v == local.certificate_id]
  )  
}

locals {
  name = huaweicloud_elb_certificate.test.name
}
data "huaweicloud_elb_certificates" "name_filter" {
  depends_on = [huaweicloud_elb_certificate.test]

  name = [huaweicloud_elb_certificate.test.name]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.name_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.name_filter.certificates[*].name : v == local.name]
  )  
}

locals {
  description = huaweicloud_elb_certificate.test.description
}
data "huaweicloud_elb_certificates" "description_filter" {
  depends_on = [huaweicloud_elb_certificate.test]

  description = [huaweicloud_elb_certificate.test.description]
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.description_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.description_filter.certificates[*].description : v == local.description]
  )  
}

locals {
  domain = huaweicloud_elb_certificate.test.domain
}
data "huaweicloud_elb_certificates" "domain_filter" {
  depends_on = [huaweicloud_elb_certificate.test]

  domain = [huaweicloud_elb_certificate.test.domain]
}
output "domain_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.domain_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.domain_filter.certificates[*].domain : v == local.domain]
  )  
}

locals {
  type = huaweicloud_elb_certificate.test.type
}
data "huaweicloud_elb_certificates" "type_filter" {
  type = [huaweicloud_elb_certificate.test.type]
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.type_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.type_filter.certificates[*].type : v == local.type]
  )  
}

locals {
  common_name = huaweicloud_elb_certificate.test.common_name
}
data "huaweicloud_elb_certificates" "common_name_filter" {
  common_name = [huaweicloud_elb_certificate.test.common_name]
}
output "common_name_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.common_name_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.common_name_filter.certificates[*].common_name : v == local.common_name]
  )  
}

locals {
  fingerprint = huaweicloud_elb_certificate.test.fingerprint
}
data "huaweicloud_elb_certificates" "fingerprint_filter" {
  fingerprint = [huaweicloud_elb_certificate.test.fingerprint]
}
output "fingerprint_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.fingerprint_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.fingerprint_filter.certificates[*].fingerprint : v == local.fingerprint]
  )  
}

locals {
  scm_certificate_id = huaweicloud_elb_certificate.test.scm_certificate_id
}
data "huaweicloud_elb_certificates" "scm_certificate_id_filter" {
  scm_certificate_id = [huaweicloud_elb_certificate.test.scm_certificate_id]
}
output "scm_certificate_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_certificates.scm_certificate_id_filter.certificates) > 0 && alltrue(
  [for v in data.huaweicloud_elb_certificates.scm_certificate_id_filter.certificates[*].scm_certificate_id :
  v == local.scm_certificate_id]
  )  
}
`, testDataaSourceElbCertificates_base(name))
}
