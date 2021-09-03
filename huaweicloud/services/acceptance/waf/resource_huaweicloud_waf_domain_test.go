package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"
)

func getResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}
	return domains.Get(c, state.Primary.ID).Extract()
}

func TestAccWafDomainV1_basic(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_basic(randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
				),
			},
			{
				Config: testAccWafDomainV1_update(randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
			},
		},
	})
}

func TestAccWafDomainV1_policy(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_policy(randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),

					acceptance.TestCheckResourceAttrWithVariable(resourceName, "policy_id",
						"${huaweicloud_waf_policy.policy_1.id}"),
				),
			},
		},
	})
}

func testAccWafDomainV1_basic(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  keep_policy      = false
  proxy            = false

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
  }
}
`, testAccWafCertificateV1_conf(randName), domainName)
}

func testAccWafDomainV1_update(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  keep_policy      = false
  proxy            = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
  }

}
`, testAccWafCertificateV1_conf(randName), domainName)
}

func testAccWafDomainV1_policy(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "%s"
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  policy_id        = huaweicloud_waf_policy.policy_1.id
  proxy            = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
  }
}
`, testAccWafCertificateV1_conf(randName), randName, domainName)
}
