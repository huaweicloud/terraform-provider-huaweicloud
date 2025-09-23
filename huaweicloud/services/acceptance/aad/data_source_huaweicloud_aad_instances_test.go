package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAadInstancesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_instances.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare AAD instances before running this test cases.
			acceptance.TestAccPrecheckAadEnable(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAadInstancesDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.0.ip_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.0.basic_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.0.elastic_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.ips.0.ip_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.expire_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.service_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.overseas_type"),
				),
			},
		},
	})
}

const testAadInstancesDataSource_basic = `data "huaweicloud_aad_instances" "test" {}`
