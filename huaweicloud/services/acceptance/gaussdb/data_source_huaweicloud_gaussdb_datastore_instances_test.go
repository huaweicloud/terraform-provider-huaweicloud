package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatastoreInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_datastore_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDatastoreInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.solution"),
					resource.TestCheckResourceAttrSet(dataSource, "engine_instance_details.0.instances.0.hotfix_versions"),
				),
			},
		},
	})
}

const testDataSourceDatastoreInstances_basic = `data "huaweicloud_gaussdb_datastore_instances" "test" {}`
