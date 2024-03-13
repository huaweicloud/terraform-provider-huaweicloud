package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cdn/v1/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceExtensionOpts(epsId string) *domains.ExtensionOpts {
	if epsId != "" {
		return &domains.ExtensionOpts{
			EnterpriseProjectId: epsId,
		}
	}
	return nil
}

func getCdnDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CdnV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDN v1 client: %s", err)
	}

	opts := getResourceExtensionOpts(state.Primary.Attributes["enterprise_project_id"])
	return domains.Get(client, state.Primary.ID, opts).Extract()
}

func TestAccCdnDomain_basic(t *testing.T) {
	var (
		domain       domains.CdnDomain
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckCDN(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.active", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin", "100.254.53.75"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCdnDomain_cache,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.rule_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.ttl", "180"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.ttl_type", "d"),
					resource.TestCheckResourceAttr(resourceName, "cache_settings.0.rules.0.priority", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
				),
			},
			{
				Config: testAccCdnDomain_retrievalHost,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "sources.0.retrieval_host", "customize.test.huaweicloud.com"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "8002"),
				),
			},
			{
				Config: testAccCdnDomain_standby,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "sources.0.active", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin", "14.215.177.39"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.active", "0"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.origin", "220.181.28.52"),
					resource.TestCheckResourceAttr(resourceName, "sources.1.origin_type", "ipaddr"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
				},
			},
		},
	})
}

var testAccCdnDomain_basic = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

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

  tags = {
    key = "val"
    foo = "bar"
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_cache = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

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
      rule_type = 0
      ttl       = 180
      ttl_type  = 4
      priority  = 2
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_retrievalHost = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  configs {
    origin_protocol = "http"
  }

  sources {
    active         = 1
    origin         = "100.254.53.75"
    origin_type    = "ipaddr"
    retrieval_host = "customize.test.huaweicloud.com"
    http_port      = 8001
    https_port     = 8002
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

var testAccCdnDomain_standby = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  sources {
    active      = 1
    origin      = "14.215.177.39"
    origin_type = "ipaddr"
  }
  sources {
    active      = 0
    origin      = "220.181.28.52"
    origin_type = "ipaddr"
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME)

func TestAccCdnDomain_configs(t *testing.T) {
	var (
		domain       domains.CdnDomain
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckCDN(t)
			acceptance.TestAccPreCheckCERT(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_configs,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_name", "terraform-test"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.cache_url_parameter_filter.0.type", "ignore_url_params"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),
				),
			},
		},
	})
}

var testAccCdnDomain_configs = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "web"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol               = "http"
    ipv6_enable                   = true
    range_based_retrieval_enabled = "true"

    https_settings {
      certificate_name = "terraform-test"
      certificate_body = file("%s")
      http2_enabled    = true
      https_enabled    = true
      private_key      = file("%s")
    }

    cache_url_parameter_filter {
      type = "ignore_url_params"
    }

    retrieval_request_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    http_response_header {
      name   = "test-name"
      value  = "test-val"
      action = "set"
    }

    url_signing {
      enabled = false
    }

    compress {
      enabled = false
    }

    force_redirect {
      enabled = true
      type   = "http"
    }
  }
}
`, acceptance.HW_CDN_DOMAIN_NAME, acceptance.HW_CDN_CERT_PATH, acceptance.HW_CDN_PRIVATE_KEY_PATH)
