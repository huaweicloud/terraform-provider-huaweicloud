package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDomainStatus_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_waf_domain_status.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDomainStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "name"),
					resource.TestCheckResourceAttrSet(datasource, "status"),
					resource.TestCheckResourceAttrSet(datasource, "waf_instance_id"),
				),
			},
		},
	})
}

func testAccDataSourceDomainStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_waf_domain_status" "test" {
  host_id = "%s"
}
`, acceptance.HW_WAF_DOMAIN_ID)
}
