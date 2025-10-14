package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDomainBatchCopy_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNURL(t)
			acceptance.TestAccPreCheckCDNTargetDomainUrls(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainBatchCopy_basic_step1(),
			},
			{
				Config:      testAccDomainBatchCopy_basic_step2(),
				ExpectError: regexp.MustCompile("Parameter invalidProperty is invalid. Reason: the config don`t support copy."),
			},
		},
	})
}

func testAccDomainBatchCopy_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "source" {
  name         = "%[1]s"
  type         = "web"
  service_area = "outside_mainland_china"

  configs {
    origin_protocol = "http"

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }
  }

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
    http_port   = 80
    https_port  = 443
  }

  cache_settings {
    rules {
      rule_type           = "all"
      ttl                 = 365
      ttl_type            = "d"
      priority            = 2
      url_parameter_type  = "del_params"
      url_parameter_value = "test_value"
    }
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}

resource "huaweicloud_cdn_domain" "target" {
  name         = "%[2]s"
  type         = "web"
  service_area = "outside_mainland_china"

  configs {
    origin_protocol = "http"
  }

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
    http_port   = 80
    https_port  = 443
  }

  cache_settings {
    rules {
      rule_type           = "all"
      ttl                 = 365
      ttl_type            = "d"
      priority            = 2
      url_parameter_type  = "del_params"
      url_parameter_value = "test_value"
    }
  }

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, acceptance.HW_CDN_DOMAIN_URL, acceptance.HW_CDN_TARGET_DOMAIN_URLS)
}

func testAccDomainBatchCopy_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cdn_domain_batch_copy" "test" {
  source_domain  = "%[2]s"
  target_domains = "%[3]s"
  config_list    = ["http_response_header"]

  depends_on = [
    huaweicloud_cdn_domain.source,
    huaweicloud_cdn_domain.target,
  ]
}
`, testAccDomainBatchCopy_base(),
		acceptance.HW_CDN_DOMAIN_URL, acceptance.HW_CDN_TARGET_DOMAIN_URLS)
}

func testAccDomainBatchCopy_basic_step2() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cdn_domain_batch_copy" "test_with_invalid_property" {
  source_domain  = "%[2]s"
  target_domains = "%[3]s"
  config_list    = ["invalid_property"]

  depends_on = [
    huaweicloud_cdn_domain.source,
    huaweicloud_cdn_domain.target,
  ]
}
`, testAccDomainBatchCopy_base(),
		acceptance.HW_CDN_DOMAIN_URL, acceptance.HW_CDN_TARGET_DOMAIN_URLS)
}
