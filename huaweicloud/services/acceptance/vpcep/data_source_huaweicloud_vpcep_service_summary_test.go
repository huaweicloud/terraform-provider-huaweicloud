package vpcep

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceServiceSummary_basic(t *testing.T) {
	var (
		name            = acceptance.RandomAccResourceName()
		dataSourceName1 = "data.huaweicloud_vpcep_service_summary.filter_by_id"
		dcById          = acceptance.InitDataSourceCheck(dataSourceName1)

		dataSourceName2 = "data.huaweicloud_vpcep_service_summary.filter_by_name"
		dcByName        = acceptance.InitDataSourceCheck(dataSourceName2)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceServiceSummary_basic1(name),
				Check: resource.ComposeTestCheckFunc(
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName1, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "service_name"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "is_charge"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "enable_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName1, "public_border_group"),
					resource.TestMatchResourceAttr(dataSourceName1, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDataSourceServiceSummary_basic2(name),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName2, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName2, "service_name"),
					resource.TestCheckResourceAttrSet(dataSourceName2, "service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName2, "is_charge"),
					resource.TestCheckResourceAttrSet(dataSourceName2, "enable_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName2, "public_border_group"),
					resource.TestMatchResourceAttr(dataSourceName2, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceServiceSummary_basic1(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_service_summary" "filter_by_id" {
  endpoint_service_id = huaweicloud_vpcep_service.test.id
}
`, testAccVPCEPService_Basic(name))
}

func testDataSourceServiceSummary_basic2(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_service_summary" "filter_by_name" {
  endpoint_service_name = huaweicloud_vpcep_service.test.service_name
}
`, testAccVPCEPService_Basic(name))
}
