package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStatisticConfiguration_basic(t *testing.T) {
	var (
		rName = "huaweicloud_cdn_statistic_configuration.test"
	)

	// Avoid CheckDestroy, because there is nothing in the resource destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStatisticConfiguration_invalid(),
				ExpectError: regexp.MustCompile(`config type is invalid`),
			},
			{
				Config: testAccStatisticConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "resource_type", "domain"),
					resource.TestCheckResourceAttr(rName, "resource_name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "config_type", "0"),
					resource.TestCheckResourceAttr(rName, "config_info.0.url.0.enable", "true"),
					resource.TestCheckResourceAttr(rName, "config_info.0.ua.0.enable", "true"),
				),
			},
		},
	})
}

func testAccStatisticConfiguration_invalid() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_statistic_configuration" "test" {
  resource_type = "domain"
  resource_name = "%[1]s"
  config_type   = 100      # Invalid config type

  config_info {
    url {
      enable = true
    }
    ua {
      enable = true
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

func testAccStatisticConfiguration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_statistic_configuration" "test" {
  resource_type = "domain"
  resource_name = "%[1]s"
  config_type   = 0

  config_info {
    url {
      enable = true
    }
    ua {
      enable = true
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}
