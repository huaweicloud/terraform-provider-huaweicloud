package gaussdb

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceStatusStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_instance_status_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck:          func() { acceptance.TestAccPreCheck(t) },
			ProviderFactories: acceptance.TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testDataSourceInstanceStatusStatistics_basic(),
					Check: resource.ComposeTestCheckFunc(
						dc.CheckResourceExists(),
						resource.TestMatchResourceAttr(dataSource, "instances_statistics.#", regexp.MustCompile(`^[0-9]+$`)),
						resource.TestCheckResourceAttrSet(dataSource, "instance_status_statistics.0.status"),
						resource.TestCheckResourceAttrSet(dataSource, "instance_status_statistics.0.count"),
					),
				},
			},
		},
	)
}

func testDataSourceInstanceStatusStatistics_basic() string {
	return `
data "huaweicloud_gaussdb_instance_status_statistics" "test" {}
`
}
