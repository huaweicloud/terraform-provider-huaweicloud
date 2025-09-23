package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/certificates/{certificate_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{certificate_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccElbV3Certificate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_certificate.server"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3CertificateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "server"),
					resource.TestCheckResourceAttr(resourceName, "domain", "www.elb.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "common_name"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
					resource.TestCheckResourceAttrSet(resourceName, "subject_alternative_names.#"),
				),
			},
			{
				Config: testAccElbV3CertificateConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "domain", "www.elbupdate.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test certificate update"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func TestAccElbV3Certificate_client(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_certificate.client"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3CertificateConfig_client(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "client"),
				),
			},
		},
	})
}

func TestAccElbV3Certificate_serverSm(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_certificate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3CertificateConfig_serverSm(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "server_sm"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccElbV3CertificateConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_certificate" "server" {
  name                  = "%s"
  description           = "terraform test certificate"
  domain                = "www.elb.com"
  enterprise_project_id = "0"

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

func testAccElbV3CertificateConfig_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_certificate" "server" {
  name                  = "%s"
  description           = "terraform test certificate update"
  domain                = "www.elbupdate.com"
  enterprise_project_id = "0"

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
i1YhgnQbn5E0hz55OLu5jvOkKQjPCW+9Aa==
-----END CERTIFICATE-----
EOT
}
`, name)
}

func testAccElbV3CertificateConfig_client(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_certificate" "client" {
  name                  = "%s"
  description           = "terraform CA certificate"
  type                  = "client"
  enterprise_project_id = "0"

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

func testAccElbV3CertificateConfig_serverSm(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_certificate" "test" {
  name                  = "%s"
  type                  = "server_sm"
  description           = "terraform test"
  domain                = "www.aaa.com"
  scm_certificate_id    = "scs1702951276162"
  enterprise_project_id = "0"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIICWzCCAgGgAwIBAgIJAOWoGwJCndPoMAoGCCqBHM9VAYN1MGcxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdIYWlEaWFuMRMwEQYDVQQK
DApHTUNlcnQub3JnMR8wHQYDVQQDDBZHTUNlcnQgR00gUm9vdCBDQSAtIDAxMB4X
DTIzMTIxOTAxNTkyNFoXDTI2MDkxMzAxNTkyNFowcDELMAkGA1UEBhMCQ04xEDAO
BgNVBAgMB2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmh1YXdl
aTEOMAwGA1UECwwFY2xvdWQxHDAaBgNVBAMME3d3dy5odWF3ZWljbG91ZC5jb20w
WTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAARL01NwuN3VhdL0j1/6Rdq8TgDw2fNa
EEpbeGbJBOkpEFGsX8cmcD2d3PkcLDerBWDDKm8Gnrnoz9lFP3Sz1Ed2o4GMMIGJ
MAwGA1UdEwEB/wQCMAAwCwYDVR0PBAQDAgeAMCwGCWCGSAGG+EIBDQQfFh1HTUNl
cnQub3JnIFNpZ25lZCBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQU1t/uXw29dzsZo4Wz
2xbOmv8rF4owHwYDVR0jBBgwFoAUf1peOwCEWSoPmL6hDm85lUMQTQcwCgYIKoEc
z1UBg3UDSAAwRQIhAJpg0O1rK2zYaZatY70uxt7i25b58DV9yEv+DFs8YxvwAiBW
G9nkFL9VXQBEl/NOIOI0ZLN2FY+9n/giQ+UTRJXI+A==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICIzCCAcigAwIBAgIJAKun/ZLoSXfeMAoGCCqBHM9VAYN1MGcxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdIYWlEaWFuMRMwEQYDVQQK
DApHTUNlcnQub3JnMR8wHQYDVQQDDBZHTUNlcnQgR00gUm9vdCBDQSAtIDAxMB4X
DTE5MTAyNDEyMzEzM1oXDTM5MDcxMTEyMzEzM1owZzELMAkGA1UEBhMCQ04xEDAO
BgNVBAgMB0JlaWppbmcxEDAOBgNVBAcMB0hhaURpYW4xEzARBgNVBAoMCkdNQ2Vy
dC5vcmcxHzAdBgNVBAMMFkdNQ2VydCBHTSBSb290IENBIC0gMDEwWTATBgcqhkjO
PQIBBggqgRzPVQGCLQNCAASXWWtv+ifV7dJHqPNXwcmioh/48Wg3IuI+o11nLEOD
zljxL2yMxoQM6xfNJHuqadXXNZv3D2rml5Pk0W/tmfHEo10wWzAdBgNVHQ4EFgQU
f1peOwCEWSoPmL6hDm85lUMQTQcwHwYDVR0jBBgwFoAUf1peOwCEWSoPmL6hDm85
lUMQTQcwDAYDVR0TBAUwAwEB/zALBgNVHQ8EBAMCAQYwCgYIKoEcz1UBg3UDSQAw
RgIhAJ7AZAC0i+4OyfxDuvPIg0I7ZtqL2kII2f1syaIW4C6iAiEAlHuUu0TMrOAr
sU47scL1B9BhyEh5tbEjsKLHia3K0YU=
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEOySLz8iJ8e4X/TbVRMImutZ3XCaAvfIXlTL3A4lEddoAoGCCqBHM9V
AYItoUQDQgAES9NTcLjd1YXS9I9f+kXavE4A8NnzWhBKW3hmyQTpKRBRrF/HJnA9
ndz5HCw3qwVgwypvBp656M/ZRT90s9RHdg==
-----END EC PRIVATE KEY-----
EOT

  enc_certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIICXzCCAgSgAwIBAgIJAOWoGwJCndPpMAoGCCqBHM9VAYN1MGcxCzAJBgNVBAYT
AkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAwDgYDVQQHDAdIYWlEaWFuMRMwEQYDVQQK
DApHTUNlcnQub3JnMR8wHQYDVQQDDBZHTUNlcnQgR00gUm9vdCBDQSAtIDAxMB4X
DTIzMTIxOTAxNTk1OFoXDTI2MDkxMzAxNTk1OFowczELMAkGA1UEBhMCQ04xEDAO
BgNVBAgMB2JlaWppbmcxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmh1YXdl
aTEOMAwGA1UECwwFY2xvdWQxHzAdBgNVBAMMFnd3dy5odWF3ZWljbG91ZGVuYy5j
b20wWTATBgcqhkjOPQIBBggqgRzPVQGCLQNCAARL2cd+bGglePnmI6/pOskNJMdH
VzWrhEf+jRZajcygzfOvT8B3825NoQc2nWS1oHcGXmRROPPxKjgCnSMfq6x7o4GM
MIGJMAwGA1UdEwEB/wQCMAAwCwYDVR0PBAQDAgeAMCwGCWCGSAGG+EIBDQQfFh1H
TUNlcnQub3JnIFNpZ25lZCBDZXJ0aWZpY2F0ZTAdBgNVHQ4EFgQU1mQx5L2JnD1e
+N73tW0E45RDUlwwHwYDVR0jBBgwFoAUf1peOwCEWSoPmL6hDm85lUMQTQcwCgYI
KoEcz1UBg3UDSQAwRgIhAJ7OMzuyD1xNGmlEOl5I8KKnKlXo/51YryDbRw+nQLoM
AiEAuEO8RweqtvcykyE64RqXOXQTnxcgn/klcI/6i+4ceU4=
-----END CERTIFICATE-----
EOT

  enc_private_key = <<EOT
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIKXkn17PSXAeG41avIEFIJ8diAhT+E3h+901uzwCp0ldoAoGCCqBHM9V
AYItoUQDQgAES9nHfmxoJXj55iOv6TrJDSTHR1c1q4RH/o0WWo3MoM3zr0/Ad/Nu
TaEHNp1ktaB3Bl5kUTjz8So4Ap0jH6usew==
-----END EC PRIVATE KEY-----
EOT
}
`, name)
}
