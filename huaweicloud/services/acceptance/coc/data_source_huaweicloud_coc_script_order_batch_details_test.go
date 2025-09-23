package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScriptOrderBatchDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_script_order_batch_details.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScriptOrderBatchDetails_basic(rName, acceptance.HW_COC_INSTANCE_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.gmt_finished"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.execute_costs"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.message"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.provider"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.agent_sn"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.properties.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.properties.0.zone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.properties.0.fixed_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_instances.0.target_instance.0.properties.0.host_name"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScriptOrderBatchDetails_basic(name, instanceId string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_script_order_batch_details" "test" {
  execute_uuid = huaweicloud_coc_script_execute.test.id
  batch_index  = 1
}

data "huaweicloud_coc_script_order_batch_details" "status_filter" {
  execute_uuid = huaweicloud_coc_script_execute.test.id
  batch_index  = 1
  status       = data.huaweicloud_coc_script_order_batch_details.test.execute_instances[0].status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_script_order_batch_details.status_filter.execute_instances) > 0 && alltrue(
    [for v in data.huaweicloud_coc_script_order_batch_details.status_filter.execute_instances[*].status :
      v == data.huaweicloud_coc_script_order_batch_details.status_filter.status]
  )
}
`, tesScriptExecute_basic(name, instanceId))
}
