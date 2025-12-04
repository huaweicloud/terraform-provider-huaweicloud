package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDomainTemplateApply_basic(t *testing.T) {
	var (
		rName = "huaweicloud_cdn_domain_template_apply.test"
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
				Config: testAccDomainTemplateApply_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "id"),
					resource.TestCheckResourceAttrSet(rName, "template_id"),
					resource.TestCheckResourceAttrSet(rName, "resources"),
				),
			},
		},
	})
}

func testAccDomainTemplateApply_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_template" "test" {
  name        = "%[1]s"
  description = "Created by terraform for template apply test"
  configs     = jsonencode({
    "cache_rules": [
      {
        "force_cache": "on",
        "follow_origin": "off",
        "match_type": "all",
        "priority": 1,
        "stale_while_revalidate": "off",
        "ttl": 20,
        "ttl_unit": "d",
        "url_parameter_type": "full_url",
        "url_parameter_value": ""
      }
    ],
    "origin_follow302_status": "off",
    "compress": {
      "type": "gzip",
      "status": "on",
      "file_type": ".js,.html,.css"
    },
    "ip_filter": {
      "type": "white",
      "value": "1.1.1.1"
    }
  })
}

resource "huaweicloud_cdn_domain_template_apply" "test" {
  template_id = huaweicloud_cdn_domain_template.test.id
  resources   = "%[2]s"
}
`, acceptance.RandomAccResourceNameWithDash(), acceptance.HW_CDN_DOMAIN_NAME)
}
