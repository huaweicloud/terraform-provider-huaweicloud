package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCdnDomain_basic(t *testing.T) {
	var domain domains.CdnDomain

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCDN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCdnDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomainV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCdnDomainV1Exists("huaweicloud_cdn_domain.domain_1", &domain),
					resource.TestCheckResourceAttr(
						"huaweicloud_cdn_domain.domain_1", "name", HW_CDN_DOMAIN_NAME),
				),
			},
		},
	})
}

func testAccCheckCdnDomainV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cdn_domain" {
			continue
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err == nil && found.DomainStatus != "deleting" {
			return fmtp.Errorf("Destroying CDN domain failed or domain still exists")
		}
	}

	return nil
}

func testAccCheckCdnDomainV1Exists(n string, domain *domains.CdnDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("CDN Domain Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		cdnClient, err := config.CdnV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud CDN Domain client: %s", err)
		}

		found, err := domains.Get(cdnClient, rs.Primary.ID, nil).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("CDN Domain not found")
		}

		*domain = *found
		return nil
	}
}

var testAccCdnDomainV1_basic = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "domain_1" {
  name   = "%s"
  type   = "wholeSite"
  enterprise_project_id = 0
  sources {
      active = 1
      origin = "100.254.53.75"
      origin_type  = "ipaddr"
  }
}
`, HW_CDN_DOMAIN_NAME)
