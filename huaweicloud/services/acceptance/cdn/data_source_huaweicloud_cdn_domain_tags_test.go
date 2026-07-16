package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomainTags_basic(t *testing.T) {
	var (
		domainName = generateDomainName()

		dcName = "data.huaweicloud_cdn_domain_tags.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainTags_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "tags.#", regexp.MustCompile("^[1-9]([0-9]+)?$")),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.value"),
					resource.TestCheckResourceAttrSet(dcName, "tags.1.key"),
					resource.TestCheckResourceAttrSet(dcName, "tags.1.value"),
				),
			},
		},
	})
}

func testAccDataDomainTags_base(domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_domain" "test" {
  name         = "%[1]s"
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

func testAccDataDomainTags_basic(domainName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cdn_domain_tags" "test" {
  resource_id = huaweicloud_cdn_domain.test.id
}`, testAccDataDomainTags_base(domainName))
}
