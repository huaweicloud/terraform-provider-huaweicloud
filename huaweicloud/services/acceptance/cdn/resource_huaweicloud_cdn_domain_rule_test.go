package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getDomainRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	domainName := state.Primary.Attributes["name"]
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.QueryCdnDomainRule(client, domainName)
}

// This test case requires opening the CDN whitelist before it can be called.
func TestAccDomainRule_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain_rule.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDomainRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDN(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainRule_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.name", "test-rule"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.cache_rule.0.follow_origin", "on"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.cache_rule.0.force_cache", "on"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.cache_rule.0.ttl", "30"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.cache_rule.0.ttl_unit", "d"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.conditions.0.match.0.logic", "and"),

					resource.TestCheckResourceAttrSet(resourceName, "rules.0.rule_id"),
					resource.TestCheckResourceAttrSet(resourceName, "rules.0.conditions.0.match.0.criteria"),
				),
			},
			{
				Config: testAccDomainRule_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDomainRuleImportState(resourceName),
			},
		},
	})
}

func testAccDomainRule_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_rule" "test" {
  name = "%s"

  rules {
    name     = "test-rule"
    priority = 1
    status   = "on"

    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "on"
        ttl           = 30
        ttl_unit      = "d"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              match_pattern = [
                "HTTP",
              ]
              match_target_name = ""
              match_target_type = "scheme"
              match_type        = "contains"
              negate            = false
            },
            {
              case_sensitive = false
              match_pattern = [
                "GET",
              ]
              match_target_name = ""
              match_target_type = "method"
              match_type        = "contains"
              negate            = false
            },
          ]
        )
        logic = "and"
      }
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

func testAccDomainRule_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain_rule" "test" {
  name = "%s"

  rules {
    name     = "test-rule2"
    priority = 2
    status   = "on"

    actions {
      access_control {
        type = "block"
      }
    }
    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "off"
        ttl           = 30
        ttl_unit      = "d"
      }
    }
    actions {
      flexible_origin {
        host_name       = "test.name.com.cn"
        http_port       = 80
        https_port      = 443
        ip_or_domain    = "110.110.110.11"
        obs_bucket_type = null
        origin_protocol = "follow"
        priority        = 3
        sources_type    = "ipaddr"
        weight          = 4
      }
    }
    actions {
      http_response_header {
        action = "set"
        name   = "Content-Language"
        value  = "en-Us"
      }
    }
    actions {
      origin_request_header {
        action = "set"
        name   = "X-Token"
        value  = "abc"
      }
    }
    actions {
      origin_request_url_rewrite {
        rewrite_type = "simple"
        source_url   = null
        target_url   = "/test/*.jpg"
      }
    }
    actions {
      request_url_rewrite {
        execution_mode       = "break"
        redirect_host        = null
        redirect_status_code = 0
        redirect_url         = "/index/test.html"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              criteria = [
                {
                  case_sensitive = false
                  match_pattern = [
                    "HTTP",
                  ]
                  match_target_name = ""
                  match_target_type = "scheme"
                  match_type        = "contains"
                  negate            = false
                },
              ]
              logic  = "and"
              negate = false
            },
            {
              case_sensitive = false
              criteria = [
                {
                  case_sensitive = true
                  match_pattern = [
                    "HTTP",
                  ]
                  match_target_name = "aa"
                  match_target_type = "arg"
                  match_type        = "contains"
                  negate            = false
                },
                {
                  case_sensitive = false
                  criteria = [
                    {
                      case_sensitive = true
                      match_pattern = [
                        "HTTP",
                      ]
                      match_target_name = ""
                      match_target_type = "filename"
                      match_type        = "contains"
                      negate            = false
                    },
                    {
                      case_sensitive = true
                      match_pattern = [
                        "HTTP",
                      ]
                      match_target_name = ""
                      match_target_type = "ua"
                      match_type        = "contains"
                      negate            = true
                    },
                  ]
                  logic  = "or"
                  negate = false
                },
              ]
              logic  = "and"
              negate = false
            },
          ]
        )
        logic = "or"
      }
    }
  }
  rules {
    name     = "test-rule1"
    priority = 1
    status   = "on"

    actions {
      access_control {
        type = "trust"
      }
    }
    actions {
      cache_rule {
        follow_origin = "on"
        force_cache   = "on"
        ttl           = 30
        ttl_unit      = "d"
      }
    }

    conditions {
      match {
        criteria = jsonencode(
          [
            {
              case_sensitive = false
              match_pattern = [
                "HTTP",
              ]
              match_target_name = ""
              match_target_type = "scheme"
              match_type        = "contains"
              negate            = false
            },
            {
              case_sensitive = false
              match_pattern = [
                "GET",
              ]
              match_target_name = ""
              match_target_type = "method"
              match_type        = "contains"
              negate            = false
            },
          ]
        )
        logic = "and"
      }
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

func testDomainRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		return rs.Primary.Attributes["name"], nil
	}
}
