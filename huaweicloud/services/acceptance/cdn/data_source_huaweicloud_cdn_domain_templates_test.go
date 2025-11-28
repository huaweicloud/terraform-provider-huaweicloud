package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomainTemplates_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_domain_templates.test"
		dc    = acceptance.InitDataSourceCheck(rName)
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainTemplates_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "templates.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.type"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.account_id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.configs"),
					resource.TestMatchResourceAttr(rName, "templates.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "templates.0.modify_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_domain_template_found", "true"),
				),
			},
		},
	})
}

func testAccDataDomainTemplates_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_template" "test" {
  name        = "%[1]s"
  description = "Created by terraform"
  configs = jsonencode({
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
    "http_response_header": [
      {
        "action": "set",
        "name": "Content-Disposition",
        "value": "1235"
      }
    ],
    "origin_follow302_status": "off",
    "compress": {
      "type": "gzip,br",
      "status": "on",
      "file_type": ".js,.html,.css,.xml,.json,.shtml,.htmx"
    },
    "origin_range_status": "on",
    "referer": {
      "type": "black",
      "value": "1.2.1.1",
      "include_empty": false
    },
    "ip_filter": {
      "type": "white",
      "value": "1.1.2.2"
    },
    "user_agent_filter": {
      "type": "white",
      "ua_list": [
        "1.1.3.3"
      ],
      "include_empty": false
    },
    "flow_limit_strategy": [
      {
        "strategy_type": "instant",
        "item_type": "bandwidth",
        "limit_value": 1000001,
        "alarm_percent_threshold": null,
        "ban_time": 60
      }
    ]
  })
}

data "huaweicloud_cdn_domain_templates" "test" {
  depends_on = [
    huaweicloud_cdn_domain_template.test
  ]
}

locals {
  domain_template_id           = huaweicloud_cdn_domain_template.test.id
  domain_template_query_result = try([
    for v in data.huaweicloud_cdn_domain_templates.test.templates : v if v.id == local.domain_template_id
  ][0], null)
}

output "is_domain_template_found" {
  value = local.domain_template_query_result != null && local.domain_template_query_result.id == local.domain_template_id
}
`, name)
}
