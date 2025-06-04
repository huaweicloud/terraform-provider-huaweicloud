package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsDrInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_dr_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsDrInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.slave_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.slave_instances.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.slave_instances.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.slave_instances.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.slave_instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_dr_relations.0.master_instance.#"),
				),
			},
		},
	})
}

func testDataSourceRdsDrInstances_basic() string {
	return `
data "huaweicloud_rds_dr_instances" "test" {}
`
}
