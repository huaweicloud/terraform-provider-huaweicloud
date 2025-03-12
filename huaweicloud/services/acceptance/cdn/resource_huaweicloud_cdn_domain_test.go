package cdn

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cdn"
)

func getCdnDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		domainName = state.Primary.Attributes["name"]
		epsID      = state.Primary.Attributes["enterprise_project_id"]
	)

	client, err := cfg.CdnV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CDN v1 client: %s", err)
	}
	return cdn.ReadCdnDomainDetail(client, domainName, epsID)
}

// Try to get the domain name from the environment variable. If the environment variable does not exist, use the
// domain name in the format of `xxx.huaweicloud.com`
func generateDomainName() string {
	if acceptance.HW_CDN_DOMAIN_NAME == "" {
		return fmt.Sprintf("%s.huaweicloud.com", acceptance.RandomAccResourceNameWithDash())
	}
	return acceptance.HW_CDN_DOMAIN_NAME
}

func TestAccCdnDomain_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
		domainName   = generateDomainName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.active", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin", "100.254.53.75"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.origin_type", "ipaddr"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "443"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.weight", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "val"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccCdnDomain_update1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "download"),
					resource.TestCheckResourceAttr(resourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.weight", "100"),
				),
			},
			{
				Config: testAccCdnDomain_update2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "web"),
					resource.TestCheckResourceAttr(resourceName, "service_area", "outside_mainland_china"),
					resource.TestCheckResourceAttr(resourceName, "sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.retrieval_host", "customize.test.huaweicloud.com"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.http_port", "8001"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.https_port", "8002"),
					resource.TestCheckResourceAttr(resourceName, "sources.0.weight", "1"),
				),
			},
			{
				Config: testAccCdnDomain_update3(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "sources.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id", "configs.0.url_signing.0.key", "configs.0.https_settings.0.certificate_body",
					"configs.0.https_settings.0.private_key", "cache_settings",
				},
			},
		},
	})
}

func testAccCdnDomain_basic(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%s"
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
`, domainName)
}

func testAccCdnDomain_update1(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%s"
  type         = "download"
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
    weight      = 100
  }

  cache_settings {
    follow_origin = true
    rules {
      rule_type          = "file_extension"
      content            = ".jpg"
      ttl                = 0
      ttl_type           = "d"
      priority           = 3
      url_parameter_type = "ignore_url_params"
    }
  }
}
`, domainName)
}

func testAccCdnDomain_update2(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%s"
  type         = "web"
  service_area = "outside_mainland_china"

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
    weight         = 1
  }

  cache_settings {}
}
`, domainName)
}

func testAccCdnDomain_update3(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%s"
  type         = "web"
  service_area = "outside_mainland_china"

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
`, domainName)
}

// Prepare the HTTPS certificate before running this test case
// All configuration item modifications may trigger `CDN.0163`. This is a problem that we have no way to solve.
// When a `CDN.0163` error occurs, you can avoid this error by adjusting the test case configuration items.
func TestAccCdnDomain_configHttpSettings(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Note: The domain name, certificate and SCM ID configured in the environment variables need to match,
			// otherwise the test case will fail.
			acceptance.TestAccPreCheckCertCDN(t)
			acceptance.TestAccPreCheckCERT(t)
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_configHttpSettings,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_CERT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_name", "terraform-test"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_source", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_type", "server"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.ocsp_stapling_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.quic.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.include_subdomains", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.max_age", "0"),
					testAccCheckTLSVersion(resourceName, "TLSv1.1,TLSv1.2"),

					resource.TestCheckResourceAttrSet(resourceName, "configs.0.https_settings.0.certificate_body"),
					resource.TestCheckResourceAttrSet(resourceName, "configs.0.https_settings.0.private_key"),
				),
			},
			{
				Config: testAccCdnDomain_configHttpSettings_update1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_CERT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_source", "2"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_name", "terraform-update"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.scm_certificate_id",
						acceptance.HW_CCM_SSL_CERTIFICATE_ID),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.certificate_type", "server"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.ocsp_stapling_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.quic.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.include_subdomains", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.max_age", "63072000"),

					testAccCheckTLSVersion(resourceName, "TLSv1.1,TLSv1.2,TLSv1.3"),
				),
			},
			{
				Config: testAccCdnDomain_configHttpSettings_update2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_CDN_CERT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.https_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.https_settings.0.http2_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.hsts.0.enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id", "configs.0.url_signing.0.key", "configs.0.https_settings.0.certificate_body",
					"configs.0.https_settings.0.private_key", "cache_settings",
				},
			},
		},
	})
}

// The response value order of field `tls_version` will be modified.
// For example `TLSv1.1,TLSv1.2` will be modified to `TLSv1.2,TLSv1.1`.
func testAccCheckTLSVersion(n string, tlsVersion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource (%s) not found: %s", n, rs)
		}

		tlsVersionAttr := rs.Primary.Attributes["configs.0.https_settings.0.tls_version"]
		if tlsVersionAttr == "" {
			return fmt.Errorf("attribute `configs.0.https_settings.0.tls_version` is not found from state")
		}

		tlsVersionArray := strings.Split(tlsVersion, ",")
		tlsVersionAttrArray := strings.Split(tlsVersionAttr, ",")
		sort.Strings(tlsVersionArray)
		sort.Strings(tlsVersionAttrArray)
		if reflect.DeepEqual(tlsVersionArray, tlsVersionAttrArray) {
			return nil
		}
		return fmt.Errorf("attribute 'configs.0.https_settings.0.tls_version' expected (%s), got (%s)",
			tlsVersion, tlsVersionAttr)
	}
}

var testAccCdnDomain_configHttpSettings = fmt.Sprintf(`
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
      certificate_name     = "terraform-test"
      certificate_body     = file("%s")
      http2_enabled        = true
      https_enabled        = true
      private_key          = file("%s")
      tls_version          = "TLSv1.1,TLSv1.2"
      certificate_source   = 0
      certificate_type     = "server"
      ocsp_stapling_status = "on"
    }

    quic {
      enabled = true
    }

    hsts {
      enabled            = true
      include_subdomains = "off"
      max_age            = 0
    }
  }
}
`, acceptance.HW_CDN_CERT_DOMAIN_NAME, acceptance.HW_CDN_CERT_PATH, acceptance.HW_CDN_PRIVATE_KEY_PATH)

var testAccCdnDomain_configHttpSettings_update1 = fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%[1]s"
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
      certificate_source   = 2
      certificate_name     = "terraform-update"
      scm_certificate_id   = "%[2]s"
      certificate_type     = "server"
      http2_enabled        = true
      https_enabled        = true
      tls_version          = "TLSv1.1,TLSv1.2,TLSv1.3"
      ocsp_stapling_status = "off"
    }

    quic {
      enabled = false
    }

    hsts {
      enabled            = true
      include_subdomains = "on"
      max_age            = 63072000
    }
  }
}
`, acceptance.HW_CDN_CERT_DOMAIN_NAME, acceptance.HW_CCM_SSL_CERTIFICATE_ID)

var testAccCdnDomain_configHttpSettings_update2 = fmt.Sprintf(`
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
      https_enabled = false
    }

    hsts {
      enabled = false
    }
  }
}
`, acceptance.HW_CDN_CERT_DOMAIN_NAME)

// All configuration item modifications may trigger `CDN.0163`. This is a problem that we have no way to solve.
// When a `CDN.0163` error occurs, you can avoid this error by adjusting the test case configuration items.
func TestAccCdnDomain_configs(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
		domainName   = generateDomainName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_configs(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.slice_etag_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_receive_timeout", "60"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_follow302_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.type", "type_a"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.sign_method", "md5"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.match_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.sign_arg", "Psd_123"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.key", "A27jtfSTy13q7A0UnTA9vpxYXEb"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.backup_key", "S36klgTFa60q3V8DmSK2hwfBOYp"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.time_format", "dec"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.expire_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.inherit_type", "m3u8,mpd"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.inherit_time_type", "sys_time"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.type", "gzip"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.file_type", ".js,.html,.css,.xml,.json,.shtml,.htm"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.redirect_code", "301"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.referer.0.type", "white"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.referer.0.include_empty", "false"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.enable_video_seek", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.enable_flv_by_time_seek", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.start_parameter", "test-start"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.end_parameter", "test-end"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_cache.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_filter.0.type", "black"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_filter.0.value", "5.12.3.65,35.2.65.21"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.#", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.server_name", "back.allium.cn.com"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.status", "on"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.#", "3"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.#", "3"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.type", "white"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.include_empty", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.ua_list.#", "3"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_failed_status", "403"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_server", "https://testdomain.com"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_success_status", "200"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.file_type_setting", "specific_file"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_args", "k1|k2|key33"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_args_setting", "reserve_specific_args"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_headers_setting", "reserve_specific_headers"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_headers", "key1|key2"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.response_status", "403"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.specified_file_type", "jpg|mp4"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.timeout", "50"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.timeout_action", "pass"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_args_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_headers_rules.#", "2"),
				),
			},
			{
				Config: testAccCdnDomain_configsUpdate1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_protocol", "follow"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.range_based_retrieval_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.description", "update description"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.slice_etag_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_receive_timeout", "30"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_follow302_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.name", "test-name-update"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.value", "test-val-update"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.retrieval_request_header.0.action", "set"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.type", "type_c2"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.sign_method", "sha256"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.match_type", "all"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.sign_arg", "Dma_001"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.key", "3P7k9s4r0aey9CB1mvvDHG2"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.backup_key", "5F8a6c3r1xgp7DL0jkeBYZ4"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.time_format", "hex"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.expire_time", "31536000"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.inherit_type", "mpd"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.inherit_time_type", "parent_url_time"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.type", "br"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.file_type", ".js,.html"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.type", "http"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.redirect_code", "302"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "on"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.referer.0.type", "black"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.referer.0.include_empty", "true"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.enable_video_seek", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.enable_flv_by_time_seek", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.start_parameter", "test-startUpdate"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.end_parameter", ""),

					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.match_type", "file_path"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.match_pattern", "/test/folder01;/test/folder02"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.priority", "3"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.back_sources.0.http_port", "83"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.back_sources.0.https_port", "470"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.back_sources.0.ip_or_domain", "www.hshs.cdd"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.0.back_sources.0.sources_type", "domain"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.limit_rate_after", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.limit_rate_value", "104857600"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.match_type", "catalog"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.match_value", "/test/ff"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.priority", "4"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.0.type", "size"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_cache.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_cache.0.code", "403"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_cache.0.ttl", "70"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_filter.0.type", "white"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_filter.0.value", "5.12.3.66"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.0.match_type", "file_path"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.0.priority", "10"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.0.source_url", "/tt/abc.txt"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.0.target_url", "/new/$1/$2.html"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.type", "black"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.include_empty", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.ua_list.0", "t1*"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.0.error_code", "416"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.0.target_code", "301"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.0.target_link", "http://example.com"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.0.area", "HK,TW,AE,LB"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.0.content_type", "file_directory"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.0.content_value", "/sdf/wer/ww/qq"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.0.exception_ip", "3.5.6.8,32.4.3.12,11.23.44.32"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.0.type", "white"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.server_name", "backupdate.allium.cn.com"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.status", "on"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.execution_mode", "break"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.redirect_host", ""),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.redirect_status_code", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.redirect_url", "/test/index.html"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.condition.0.match_type", "catalog"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.condition.0.match_value", "/test/folder/1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.0.condition.0.priority", "10"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.cache_type", "ttl"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.ttl", "30"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.ttl_unit", "m"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.condition.0.match_type", "file_extension"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.condition.0.match_value", ".jpg,.zip,.gz"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.0.condition.0.priority", "2"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_failed_status", "503"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_server", "https://testdomain-update.com"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.auth_success_status", "302"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.file_type_setting", "all"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.request_method", "POST"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_args_setting", "reserve_all_args"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.reserve_headers_setting", "reserve_all_headers"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.response_status", "206"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.timeout", "3000"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.timeout_action", "forbid"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_args_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_headers_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_args_rules.0.key", "http_user_agent"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_args_rules.0.type", "nginx_preset_var"),
					resource.TestCheckResourceAttr(resourceName,
						"configs.0.remote_auth.0.remote_auth_rules.0.add_custom_args_rules.0.value", "$server_protocol"),
				),
			},
			{
				Config: testAccCdnDomain_configsUpdate2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.description", ""),
					resource.TestCheckResourceAttr(resourceName, "configs.0.slice_etag_status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_receive_timeout", "5"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_follow302_status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.request_limit_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_cache.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.ip_filter.0.type", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.origin_request_url_rewrite.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.user_agent_filter.0.type", "off"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.error_code_redirect_rules.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.access_area_filter.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.sni.0.status", "off"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.request_url_rewrite.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.browser_cache_rules.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.remote_auth_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "on"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.inherit_config.0.status", "off"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.compress.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.video_seek.0.enable_video_seek", "false"),
				),
			},
			{
				Config: testAccCdnDomain_configsUpdate3(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.flexible_origin.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.remote_auth.0.remote_auth_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.url_signing.0.status", "off"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.force_redirect.0.status", "off"),

					resource.TestCheckResourceAttr(resourceName, "configs.0.referer.0.type", "off"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id", "configs.0.url_signing.0.key", "configs.0.url_signing.0.backup_key",
					"configs.0.https_settings.0.certificate_body", "configs.0.https_settings.0.private_key", "cache_settings",
				},
			},
		},
	})
}

func testAccCdnDomain_configs(domainName string) string {
	return fmt.Sprintf(`
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
    range_based_retrieval_enabled = true
    description                   = "test description"
    origin_receive_timeout        = "60"

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
      enabled     = true
      type        = "type_a"
      sign_method = "md5"
      match_type  = "all"
      sign_arg    = "Psd_123"
      key         = "A27jtfSTy13q7A0UnTA9vpxYXEb"
      backup_key  = "S36klgTFa60q3V8DmSK2hwfBOYp"
      time_format = "dec"
      expire_time = 0

      inherit_config {
        enabled           = true
        inherit_type      = "m3u8,mpd"
        inherit_time_type = "sys_time"
      }
    }

    compress {
      enabled = true
      type    = "gzip"
    }

    force_redirect {
      enabled       = true
      type          = "http"
      redirect_code = 301
    }

    referer {
      type          = "white"
      value         = "*.common.com,192.187.2.43,www.test.top:4990"
      include_empty = false
    }

    video_seek {
      enable_video_seek       = true
      enable_flv_by_time_seek = true
      start_parameter         = "test-start"
      end_parameter           = "test-end"
    }

    flexible_origin {
      match_type = "all"
      priority   = 1

      back_sources {
        http_port    = 1
        https_port   = 65535
        ip_or_domain = "165.132.12.2"
        sources_type = "ipaddr"
      }
    }

    flexible_origin {
      match_type    = "file_extension"
      match_pattern = ".jpg;.zip;.exe"
      priority      = 2

      back_sources {
        http_port    = 65535
        https_port   = 1
        ip_or_domain = "165.5.1.4"
        sources_type = "ipaddr"
      }
    }

    request_limit_rules {
      limit_rate_after = 0
      limit_rate_value = 0
      match_type       = "all"
      priority         = 2
      type             = "size"
    }
    
    request_limit_rules {
      limit_rate_after = 1073741824
      limit_rate_value = 104857600
      match_type       = "catalog"
      match_value      = "/test/ff"
      priority         = 5
      type             = "size"
    }

    error_code_cache {
      code = 301
      ttl  = 0
    }

    error_code_cache {
      code = 500
      ttl  = 31536000
    }

    ip_filter {
      type  = "black"
      value = "5.12.3.65,35.2.65.21"
    }

    origin_request_url_rewrite {
      match_type = "all"
      priority   = 2
      target_url = "/nn.tx"
    }

    origin_request_url_rewrite {
      match_type = "file_path"
      priority   = 5
      source_url = "/tt/ab.txt"
      target_url = "/new/$1/$2.html"
    }

    user_agent_filter {
      type    = "white"
      ua_list = [
        "t1",
        "t2",
        "t3*",
      ]
    }

    error_code_redirect_rules {
      error_code  = 416
      target_code = 301
      target_link = "http://example.com"
    }

    error_code_redirect_rules {
      error_code  = 502
      target_code = 302
      target_link = "https://xxx.cn/"
    }

    access_area_filter {
      area          = "HK,TW,AE,LB,TJ,MY"
      content_type  = "file_directory"
      content_value = "/sdf/wer/ww"
      exception_ip  = "3.5.6.8,32.4.3.12,11.23.44.32,102.34.4.12,192.68.1.2"
      type          = "white"
    }
    access_area_filter {
      area         = "MO,HK,TW,AE,MY,BN,TR"
      content_type = "all"
      type         = "black"
    }

    sni {
      enabled     = true
      server_name = "back.allium.cn.com"
    }

    request_url_rewrite {
      execution_mode = "break"
      redirect_url   = "/test/index.html"

      condition {
        match_type  = "catalog"
        match_value = "/test/folder/1"
        priority    = 10
      }
    }
    request_url_rewrite {
      execution_mode       = "redirect"
      redirect_host        = "https://www.example.com"
      redirect_status_code = 303
      redirect_url         = "/test/index.jsp"

      condition {
        match_type  = "full_path"
        match_value = "/test.jpg"
        priority    = 5
      }
    }
    request_url_rewrite {
      execution_mode = "break"
      redirect_url   = "/test/demo.html"

      condition {
        match_type = "home_page"
        priority   = 6
      }
    }

    browser_cache_rules {
      cache_type = "follow_origin"

      condition {
        match_type  = "catalog"
        match_value = "/test/folder/1"
        priority    = 5
      }
    }
    browser_cache_rules {
      cache_type = "never"

      condition {
        match_type = "all"
        priority   = 7
      }
    }
    browser_cache_rules {
      cache_type = "ttl"
      ttl        = 30
      ttl_unit   = "m"

      condition {
        match_type  = "file_extension"
        match_value = ".jpg,.zip"
        priority    = 2
      }
    }

    remote_auth {
      enabled = true

      remote_auth_rules {
        auth_failed_status      = "403"
        auth_server             = "https://testdomain.com"
        auth_success_status     = "200"
        file_type_setting       = "specific_file"
        request_method          = "GET"
        reserve_args            = "k1|k2|key33"
        reserve_args_setting    = "reserve_specific_args"
        reserve_headers_setting = "reserve_specific_headers"
        reserve_headers         = "key1|key2"
        response_status         = "403"
        specified_file_type     = "jpg|mp4"
        timeout                 = 50
        timeout_action          = "pass"

        add_custom_args_rules {
          key   = "http_user_agent"
          type  = "nginx_preset_var"
          value = "$http_host"
        }

        add_custom_args_rules {
          key   = "args_custom_key"
          type  = "custom_var"
          value = "args_custom_value"
        }

        add_custom_headers_rules {
          key   = "http_user_agent"
          type  = "nginx_preset_var"
          value = "$remote_addr"
        }

        add_custom_headers_rules {
          key   = "headers_custom_key"
          type  = "custom_var"
          value = "headers_custom_value"
        }
      }
    }
  }
}
`, domainName)
}

func testAccCdnDomain_configsUpdate1(domainName string) string {
	return fmt.Sprintf(`
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
    origin_protocol               = "follow"
    ipv6_enable                   = false
    range_based_retrieval_enabled = false
    description                   = "update description"
    slice_etag_status             = "off"
    origin_receive_timeout        = "30"
    origin_follow302_status       = "on"

    retrieval_request_header {
      name   = "test-name-update"
      value  = "test-val-update"
      action = "set"
    }

    http_response_header {
      name   = "Content-Disposition"
      value  = "test-val-update"
      action = "set"
    }

    url_signing {
      enabled     = true
      type        = "type_c2"
      sign_method = "sha256"
      match_type  = "all"
      sign_arg    = "Dma_001"
      key         = "3P7k9s4r0aey9CB1mvvDHG2"
      backup_key  = "5F8a6c3r1xgp7DL0jkeBYZ4"
      time_format = "hex"
      expire_time = 31536000

      inherit_config {
        enabled           = true
        inherit_type      = "mpd"
        inherit_time_type = "parent_url_time"
      }
    }

    compress {
      enabled   = true
      type      = "br"
      file_type = ".js,.html"
    }

    force_redirect {
      enabled       = true
      type          = "http"
      redirect_code = 302
    }

    referer {
      type          = "black"
      value         = "*.common.com,192.187.2.43"
      include_empty = true
    }

    video_seek {
      enable_video_seek       = true
      enable_flv_by_time_seek = false
      start_parameter         = "test-startUpdate"
    }

    flexible_origin {
      match_type    = "file_path"
      match_pattern = "/test/folder01;/test/folder02"
      priority      = 3

      back_sources {
        http_port    = 83
        https_port   = 470
        ip_or_domain = "www.hshs.cdd"
        sources_type = "domain"
      }
    }

    request_limit_rules {
      limit_rate_after = 0
      limit_rate_value = 104857600
      match_type       = "catalog"
      match_value      = "/test/ff"
      priority         = 4
      type             = "size"
    }

    error_code_cache {
      code = 403
      ttl  = 70
    }

    ip_filter {
      type  = "white"
      value = "5.12.3.66"
    }

    origin_request_url_rewrite {
      match_type = "file_path"
      priority   = 10
      source_url = "/tt/abc.txt"
      target_url = "/new/$1/$2.html"
    }

    user_agent_filter {
      type          = "black"
      include_empty = "false"
      ua_list = [
        "t1*",
      ]
    }

    error_code_redirect_rules {
      error_code  = 416
      target_code = 301
      target_link = "http://example.com"
    }

    access_area_filter {
      area          = "HK,TW,AE,LB"
      content_type  = "file_directory"
      content_value = "/sdf/wer/ww/qq"
      exception_ip  = "3.5.6.8,32.4.3.12,11.23.44.32"
      type          = "white"
    }

    sni {
      enabled     = true
      server_name = "backupdate.allium.cn.com"
    }

    request_url_rewrite {
      execution_mode = "break"
      redirect_url   = "/test/index.html"

      condition {
        match_type  = "catalog"
        match_value = "/test/folder/1"
        priority    = 10
      }
    }

    browser_cache_rules {
      cache_type = "ttl"
      ttl        = 30
      ttl_unit   = "m"

      condition {
        match_type  = "file_extension"
        match_value = ".jpg,.zip,.gz"
        priority    = 2
      }
    }

    remote_auth {
      enabled = true

      remote_auth_rules {
        auth_failed_status      = "503"
        auth_server             = "https://testdomain-update.com"
        auth_success_status     = "302"
        file_type_setting       = "all"
        request_method          = "POST"
        reserve_args_setting    = "reserve_all_args"
        reserve_headers_setting = "reserve_all_headers"
        response_status         = "206"
        timeout                 = 3000
        timeout_action          = "forbid"

        add_custom_args_rules {
          key   = "http_user_agent"
          type  = "nginx_preset_var"
          value = "$server_protocol"
        }
      }
    }
  }
}
`, domainName)
}

func testAccCdnDomain_configsUpdate2(domainName string) string {
	return fmt.Sprintf(`
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
    origin_protocol               = "follow"
    ipv6_enable                   = false
    range_based_retrieval_enabled = false
    slice_etag_status             = "on"
    origin_receive_timeout        = "5"
    origin_follow302_status       = "off"

    remote_auth {
      enabled = false
    }

    url_signing {
      enabled     = true
      type        = "type_c2"
      sign_method = "sha256"
      match_type  = "all"
      sign_arg    = "web506"
      key         = "3P7k9s4r0aey9CB1mvvDHG2"
      backup_key  = "5F8a6c3r1xgp7DL0jkeBYZ4"
      time_format = "hex"
      expire_time = 31536000

      inherit_config {
        enabled = false
      }
    }

    compress {
      enabled = false
    }

    video_seek {
      enable_video_seek = false
    }

    sni {
      enabled = false
    }

    ip_filter {
      type = "off"
    }

    user_agent_filter {
      type = "off"
    }
  }
}
`, domainName)
}

func testAccCdnDomain_configsUpdate3(domainName string) string {
	return fmt.Sprintf(`
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
    origin_protocol               = "follow"
    ipv6_enable                   = false
    range_based_retrieval_enabled = false

    remote_auth {
      enabled = false
    }

    url_signing {
      enabled = false
    }

    force_redirect {
      enabled = false
    }

    referer {
      type = "off"
    }
  }
}
`, domainName)
}

// All configuration item modifications may trigger `CDN.0163`. This is a problem that we have no way to solve.
// When a `CDN.0163` error occurs, you can avoid this error by adjusting the test case configuration items.
func TestAccCdnDomain_configTypeWholeSite(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
		domainName   = generateDomainName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_wholeSite(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.timeout", "1"),
				),
			},
			{
				Config: testAccCdnDomain_wholeSiteUpdate1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.timeout", "300"),
				),
			},
			{
				Config: testAccCdnDomain_wholeSiteUpdate2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "type", "wholeSite"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.websocket.0.enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id", "configs.0.url_signing.0.key", "configs.0.https_settings.0.certificate_body",
					"configs.0.https_settings.0.private_key", "cache_settings",
				},
			},
		},
	})
}

func testAccCdnDomain_wholeSite(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = true
      timeout = 1
    }
  }
}
`, domainName)
}

func testAccCdnDomain_wholeSiteUpdate1(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = true
      timeout = 300
    }
  }
}
`, domainName)
}

func testAccCdnDomain_wholeSiteUpdate2(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = 0

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    origin_protocol = "http"
    ipv6_enable     = true

    websocket {
      enabled = false
    }
  }
}
`, domainName)
}

func TestAccCdnDomain_epsID_migrate(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
		domainName   = generateDomainName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_epsID_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccCdnDomain_epsID_update1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccCdnDomain_epsID_update2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id", "configs.0.url_signing.0.key", "configs.0.https_settings.0.certificate_body",
					"configs.0.https_settings.0.private_key", "cache_settings",
				},
			},
		},
	})
}

func testAccCdnDomain_epsID_basic(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }
}
`, domainName)
}

func testAccCdnDomain_epsID_update1(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "%s"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }
}
`, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCdnDomain_epsID_update2(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name                  = "%s"
  type                  = "wholeSite"
  service_area          = "outside_mainland_china"
  enterprise_project_id = "0"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }
}
`, domainName)
}

func TestAccCdnDomain_client_cert(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_cdn_domain.test"
		domainName   = generateDomainName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCdnDomainFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCaCertificate(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCdnDomain_client_cert_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.hosts", "demo1.com.cn|demo2.com.cn|demo3.com.cn"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.status", "on"),
					resource.TestCheckResourceAttrSet(resourceName, "configs.0.client_cert.0.trusted_cert"),
				),
			},
			{
				Config: testAccCdnDomain_client_cert_update1(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.hosts", "demo1.com.cn"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.status", "on"),
					resource.TestCheckResourceAttrSet(resourceName, "configs.0.client_cert.0.trusted_cert"),
				),
			},
			{
				Config: testAccCdnDomain_client_cert_update2(domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.hosts", ""),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.status", "off"),
					resource.TestCheckResourceAttr(resourceName, "configs.0.client_cert.0.trusted_cert", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCDNDomainImportState(resourceName),
			},
		},
	})
}

func testAccCdnDomain_client_cert_basic(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%[1]s"
  type         = "wholeSite"
  service_area = "outside_mainland_china"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    client_cert {
      enabled      = true
      hosts        = "demo1.com.cn|demo2.com.cn|demo3.com.cn"
      trusted_cert = file("%[2]s")
    }
  }
}
`, domainName, acceptance.HW_CCM_CA_CERTIFICATE_PATH)
}

func testAccCdnDomain_client_cert_update1(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%[1]s"
  type         = "wholeSite"
  service_area = "outside_mainland_china"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    client_cert {
      enabled      = true
      hosts        = "demo1.com.cn"
      trusted_cert = file("%[2]s")
    }
  }
}
`, domainName, acceptance.HW_CCM_CA_CERTIFICATE_PATH)
}

func testAccCdnDomain_client_cert_update2(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%[1]s"
  type         = "wholeSite"
  service_area = "outside_mainland_china"

  sources {
    active      = 1
    origin      = "100.254.53.75"
    origin_type = "ipaddr"
  }

  configs {
    client_cert {
      enabled      = false
    }
  }
}
`, domainName)
}

// testCDNDomainImportState use to return an ID using `name`
func testCDNDomainImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		return rs.Primary.Attributes["name"], nil
	}
}
