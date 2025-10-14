package cdn

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCdnQuotas_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_cdn_quotas.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCdnQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "quotas.#", regexp.MustCompile("^[1-9]([0-9]*)?$")),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.limit"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.user_domain_id"),
				),
			},
		},
	})
}

const testAccDataSourceCdnQuotas_basic = `data "huaweicloud_cdn_quotas" "test" {}`
