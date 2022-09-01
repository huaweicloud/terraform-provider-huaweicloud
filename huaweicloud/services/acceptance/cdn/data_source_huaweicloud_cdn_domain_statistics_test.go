package cdn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceStatistics_basic(t *testing.T) {
	rName := "data.huaweicloud_cdn_domain_statistics.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "result", "{}"),
				),
			},
		},
	})
}

func testAccDatasourceStatistics_basic() string {
	return `
data "huaweicloud_cdn_domain_statistics" "test" {
  action      = "location_detail"
  start_time  = 1662019200000
  end_time    = 1662021000000
  domain_name = "terraform.test.huaweicloud.com"
  stat_type   = "req_num"
}
`
}
