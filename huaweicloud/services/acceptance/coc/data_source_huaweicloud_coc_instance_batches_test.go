package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocInstanceBatches_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_instance_batches.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocInstanceBatches_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.batch_index"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.cloud_service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.properties.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.properties.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.properties.0.fixed_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.properties.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.target_instances.0.properties.0.zone_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocInstanceBatches_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_instance_batches" "test" {
  batch_strategy = "AUTO_BATCH"
  target_instances {
    resource_id        = huaweicloud_compute_instance.test.id
    cloud_service_name = "ecs"
    region_id          = huaweicloud_compute_instance.test.region
    type               = "CLOUDSERVER"
    custom_attributes {
      key   = "key"
      value = "value"
    }
    properties {
      host_name = huaweicloud_compute_instance.test.hostname
      fixed_ip  = huaweicloud_compute_instance.test.access_ip_v4
      region_id = huaweicloud_compute_instance.test.region
      zone_id   = huaweicloud_compute_instance.test.availability_zone
    }
  }
}
`, testAccComputeInstance_basic(name))
}
