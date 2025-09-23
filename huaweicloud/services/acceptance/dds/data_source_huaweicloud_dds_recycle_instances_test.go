package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsRecycleInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_recycle_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSRecycleInstancesEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdsRecycleInstances_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.backup_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.deleted_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.retained_until"),
				),
			},
		},
	})
}

const testDataSourceDdsRecycleInstances_basic = `data "huaweicloud_dds_recycle_instances" "test" {}`
