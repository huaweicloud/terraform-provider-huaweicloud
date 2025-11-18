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

func getRuleEngineRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	domainName := state.Primary.Attributes["domain_name"]
	ruleId := state.Primary.ID
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	return cdn.GetRuleEngineRuleById(client, domainName, ruleId)
}

func TestAccRuleEngineRule_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_cdn_rule_engine_rule.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getRuleEngineRuleFunc)

		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleEngineRule_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "on"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "conditions"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "11"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.ttl", "10"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.ttl_unit", "m"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.follow_origin", "min_ttl"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.force_cache", "off"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.access_control.0.type", "block"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.name", "Access-Control-Deny-Origin"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.value", "*"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.action", "delete"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.0.cache_type", "follow_origin"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.0.redirect_url", "/path/$1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.0.execution_mode", "break"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.sources_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.ip_or_domain", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.https_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.origin_protocol", "follow"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.host_name", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.weight", "10"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.1.sources_type", "third_bucket"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.1.ip_or_domain", "test.third-bucket.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.action", "delete"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.name", "test"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.value", "123"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.0.rewrite_type", "simple"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.0.target_url", "/test"),
					resource.TestCheckResourceAttr(resourceName, "actions.8.origin_range.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.8.origin_range.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.0.limit_rate_after", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.0.limit_rate_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.0.code", "403"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.0.ttl", "123"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.1.code", "404"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.1.ttl", "123"),
				),
			},
			{
				Config: testAccRuleEngineRule_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain_name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "status", "off"),
					resource.TestCheckResourceAttr(resourceName, "priority", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "conditions"),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "11"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.ttl", "360"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.ttl_unit", "s"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.follow_origin", "off"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.cache_rule.0.force_cache", "on"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.access_control.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.1.access_control.0.type", "trust"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.name", "Access-Control-Allow-Origin"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.value", "*"),
					resource.TestCheckResourceAttr(resourceName, "actions.2.http_response_header.0.action", "set"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.0.cache_type", "ttl"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.0.ttl", "86400"),
					resource.TestCheckResourceAttr(resourceName, "actions.3.browser_cache_rule.0.ttl_unit", "s"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.0.redirect_status_code", "301"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.0.redirect_url", "/new-path/$1"),
					resource.TestCheckResourceAttr(resourceName, "actions.4.request_url_rewrite.0.execution_mode", "redirect"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.sources_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.ip_or_domain", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.0.weight", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.1.sources_type", "third_bucket"),
					resource.TestCheckResourceAttr(resourceName, "actions.5.flexible_origin.1.ip_or_domain", "test.third-bucket.com"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.action", "set"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.name", "new-test"),
					resource.TestCheckResourceAttr(resourceName, "actions.6.origin_request_header.0.value", "456"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.0.rewrite_type", "wildcard"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.0.source_url", "/test"),
					resource.TestCheckResourceAttr(resourceName, "actions.7.origin_request_url_rewrite.0.target_url", "/test/new"),
					resource.TestCheckResourceAttr(resourceName, "actions.8.origin_range.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.8.origin_range.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.0.limit_rate_after", "20"),
					resource.TestCheckResourceAttr(resourceName, "actions.9.request_limit_rule.0.limit_rate_value", "30"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.0.code", "405"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.0.ttl", "456"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.1.code", "407"),
					resource.TestCheckResourceAttr(resourceName, "actions.10.error_code_cache.1.ttl", "456"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testRuleEngineRuleImportStateWithId(resourceName),
				ImportStateVerifyIgnore: []string{
					"actions.5.flexible_origin.1.bucket_secret_key",
					"conditions_origin",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testRuleEngineRuleImportStateWithName(resourceName),
				ImportStateVerifyIgnore: []string{
					"actions.5.flexible_origin.1.bucket_secret_key",
					"conditions_origin",
				},
			},
		},
	})
}

func testRuleEngineRuleImportStateWithId(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		domainName := rs.Primary.Attributes["domain_name"]
		ruleId := rs.Primary.ID
		if domainName == "" || ruleId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<domain_name>/<id>', but got '%s/%s'", domainName, ruleId)
		}
		return fmt.Sprintf("%s/%s", domainName, ruleId), nil
	}
}

func testRuleEngineRuleImportStateWithName(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		domainName := rs.Primary.Attributes["domain_name"]
		ruleName := rs.Primary.Attributes["name"]
		if domainName == "" || ruleName == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<domain_name>/<name>', but got '%s/%s'", domainName, ruleName)
		}
		return fmt.Sprintf("%s/%s", domainName, ruleName), nil
	}
}

func testAccRuleEngineRule_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_rule_engine_rule" "test" {
  domain_name = "%[1]s"
  name        = "%[2]s"
  status      = "on"
  priority    = 1

  conditions = jsonencode({
    "match": {
      "logic": "and",
      "criteria": [
        {
          "match_target_type": "extension",
          "match_type": "contains",
          "match_pattern": [".txt", ".png"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "scheme",
          "match_type": "contains",
          "match_pattern": ["HTTPS"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "method",
          "match_type": "contains",
          "match_pattern": ["PATCH"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "path",
          "match_type": "contains",
          "match_pattern": ["/test"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "arg",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "filename",
          "match_type": "contains",
          "match_pattern": ["test", "123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "header",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "clientip",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["1.1.1.1"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "clientip_version",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["IPv4"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "ua",
          "match_type": "contains",
          "match_pattern": ["test"],
          "negate": false,
          "case_sensitive": true
        }
      ]
    }
  })

  actions {
    cache_rule {
      ttl           = 10
      ttl_unit      = "m"
      follow_origin = "min_ttl"
      force_cache   = "off"
    }
  }

  actions {
    access_control {
      type = "block"
    }
  }

  actions {
    http_response_header {
      name   = "Access-Control-Deny-Origin"
      value  = "*"
      action = "delete"
    }
  }

  actions {
    browser_cache_rule {
      cache_type = "follow_origin"
    }
  }

  actions {
    request_url_rewrite {
      execution_mode = "break"
      redirect_url   = "/path/$1"
    }
  }

  actions {
    flexible_origin {
      sources_type      = "ipaddr"
      ip_or_domain      = "1.1.1.1"
      http_port         = 80
      https_port        = 443
      origin_protocol  = "follow"
      host_name         = "1.1.1.1"
      priority          = 2
      weight            = 10
    }
    flexible_origin {
      sources_type      = "third_bucket"
      ip_or_domain      = "test.third-bucket.com"
      bucket_access_key = "test-ak"
      bucket_secret_key = "test-sk"
      bucket_region     = "cn-north-4"
      bucket_name       = "test-third-bucket-name"
      http_port         = 80
      https_port        = 443
      origin_protocol   = "follow"
      host_name         = "1.1.1.1"
      priority          = 1
      weight            = 2
    }
  }

  actions {
    origin_request_header {
      action = "delete"
      name   = "test"
      value  = "123"
    }
  }

  actions {
    origin_request_url_rewrite {
      rewrite_type = "simple"
      target_url   = "/test"
    }
  }

  actions {
    origin_range {
      status = "on"
    }
  }

  actions {
    request_limit_rule {
      limit_rate_after = 2
      limit_rate_value = 3
    }
  }

  actions {
    error_code_cache {
      code = 403
      ttl  = 123
    }

    error_code_cache {
      code = 404
      ttl  = 123
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME, name)
}

func testAccRuleEngineRule_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_rule_engine_rule" "test" {
  domain_name = "%[1]s"
  name        = "%[2]s"
  status      = "off"
  priority    = 2

  conditions = jsonencode({
    "match": {
      "logic": "and",
      "criteria": [
        {
          "match_target_type": "extension",
          "match_type": "contains",
          "match_pattern": [".txt", ".png", ".gif"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "scheme",
          "match_type": "contains",
          "match_pattern": ["HTTP"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "method",
          "match_type": "contains",
          "match_pattern": ["GET"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "path",
          "match_type": "contains",
          "match_pattern": ["/test/new"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "arg",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "filename",
          "match_type": "contains",
          "match_pattern": ["test", "123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "header",
          "match_target_name": "test",
          "match_type": "contains",
          "match_pattern": ["123"],
          "negate": false,
          "case_sensitive": true
        },
        {
          "match_target_type": "clientip",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["1.1.1.1"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "clientip_version",
          "match_target_name": "connect",
          "match_type": "contains",
          "match_pattern": ["IPv4"],
          "negate": false,
          "case_sensitive": false
        },
        {
          "match_target_type": "ua",
          "match_type": "contains",
          "match_pattern": ["test"],
          "negate": false,
          "case_sensitive": true
        }
      ]
    }
  })

  actions {
    cache_rule {
      ttl          = 360
      ttl_unit     = "s"
      follow_origin = "off"
      force_cache   = "on"
    }
  }

  actions {
    access_control {
      type = "trust"
    }
  }

  actions {
    http_response_header {
      name   = "Access-Control-Allow-Origin"
      value  = "*"
      action = "set"
    }
  }

  actions {
    browser_cache_rule {
      cache_type = "ttl"
      ttl        = 86400
      ttl_unit   = "s"
    }
  }

  actions {
    request_url_rewrite {
      redirect_status_code = 301
      redirect_url         = "/new-path/$1"
      execution_mode       = "redirect"
    }
  }

  actions {
    flexible_origin {
      sources_type      = "ipaddr"
      ip_or_domain      = "1.1.1.1"
      bucket_access_key = ""
      bucket_secret_key = ""
      bucket_region     = ""
      bucket_name       = ""
      http_port         = 80
      https_port        = 443
      origin_protocol  = "follow"
      host_name         = "1.1.1.1"
      priority          = 1
      weight            = 2
    }
    flexible_origin {
      sources_type      = "third_bucket"
      ip_or_domain      = "test.third-bucket.com"
      bucket_access_key = "test-ak"
      bucket_secret_key = "test-sk"
      bucket_region     = "cn-north-4"
      bucket_name       = "test-third-bucket-name"
      http_port         = 80
      https_port        = 443
      origin_protocol   = "follow"
      host_name         = "1.1.1.1"
      priority          = 1
      weight            = 2
    }
  }

  actions {
    origin_request_header {
      action = "set"
      name   = "new-test"
      value  = "456"
    }
  }

  actions {
    origin_request_url_rewrite {
      rewrite_type = "wildcard"
	  source_url   = "/test"
      target_url   = "/test/new"
    }
  }

  actions {
    origin_range {
      status = "off"
    }
  }

  actions {
    request_limit_rule {
      limit_rate_after = 20
      limit_rate_value = 30
    }
  }

  actions {
    error_code_cache {
      code = 405
      ttl  = 456
    }

    error_code_cache {
      code = 407
      ttl  = 456
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME, name)
}
