package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getDomainTemplateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	tmlId := state.Primary.ID
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.GetDomainTemplateById(client, tmlId)
}

func TestAccDomainTemplate_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cdn_domain_template.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDomainTemplateFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainTemplate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "configs"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestMatchResourceAttr(resourceName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccDomainTemplate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "configs"),
					resource.TestMatchResourceAttr(resourceName, "modify_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDomainTemplateImportStateWithName(resourceName),
			},
		},
	})
}

func testDomainTemplateImportStateWithName(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		templateName := rs.Primary.Attributes["name"]
		if templateName == "" {
			return "", fmt.Errorf("template name is missing, want '<name>', but got '%s'", templateName)
		}
		return templateName, nil
	}
}

func testAccDomainTemplate_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_template" "test" {
  name        = "%[1]s"
  description = "Created by terraform"
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
`, name)
}

func testAccDomainTemplate_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_template" "test" {
  name        = "%[1]s"
  description = "Updated by terraform"
  configs     = jsonencode({
    "cache_rules": [
      {
        "force_cache": "off",
        "follow_origin": "on",
        "match_type": "all",
        "priority": 1,
        "stale_while_revalidate": "on",
        "ttl": 30,
        "ttl_unit": "d",
        "url_parameter_type": "full_url",
        "url_parameter_value": ""
      }
    ],
    "http_response_header": [
      {
        "action": "set",
        "name": "Content-Disposition",
        "value": "updated_value"
      }
    ],
    "origin_follow302_status": "on",
    "compress": {
      "type": "gzip",
      "status": "off",
      "file_type": ".js,.html"
    },
    "origin_range_status": "off",
    "referer": {
      "type": "white",
      "value": "2.2.2.2",
      "include_empty": true
    },
    "ip_filter": {
      "type": "black",
      "value": "2.2.3.3"
    },
    "user_agent_filter": {
      "type": "black",
      "ua_list": [
        "2.2.4.4"
      ],
      "include_empty": true
    },
    "flow_limit_strategy": [
      {
        "strategy_type": "instant",
        "item_type": "bandwidth",
        "limit_value": 2000002,
        "alarm_percent_threshold": null,
        "ban_time": 720
      }
    ]
  })
}
`, name)
}
