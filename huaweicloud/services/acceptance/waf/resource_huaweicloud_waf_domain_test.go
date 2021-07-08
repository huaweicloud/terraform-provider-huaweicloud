package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/golangsdk/openstack/waf_hw/v1/domains"
)

func TestAccWafDomainV1_basic(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acctest.RandString(8)
	certificateName := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_basic(certificateName, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
				),
			},
			{
				Config: testAccWafDomainV1_update(certificateName, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
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
	randName := acctest.RandString(8)
	certificateName := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckWafDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainV1_policy(certificateName, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
		},
	})
}

func testAccCheckWafDomainV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_domain" {
			continue
		}

		_, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf domain still exists")
		}
	}

	return nil
}

func testAccCheckWafDomainV1Exists(n string, domain *domains.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
		}

		found, err := domains.Get(wafClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Waf domain not found")
		}

		*domain = *found

		return nil
	}
}

func testAccWafDomainV1_basic(certificateName string, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "www.%s.com"
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
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDomainV1_update(certificateName string, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "www.%s.com"
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
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDomainV1_policy(certificateName string, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "policy_%s"
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "www.%s.com"
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
`, testAccWafCertificateV1_conf(certificateName), name, name)
}
