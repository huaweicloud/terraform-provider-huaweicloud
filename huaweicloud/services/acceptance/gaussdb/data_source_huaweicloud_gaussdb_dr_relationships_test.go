package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbDrRelationships_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_dr_relationships.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDbDrRelationships_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "relations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.disaster_type"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.disaster_role"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.synchronization_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.actions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.0.region_code"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.slave_region_instance_info.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.0.region_code"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "relations.0.master_region_instance_info.0.ip_address"),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("dr_role_filter_is_useful", "true"),
					resource.TestCheckOutput("dr_type_filter_is_useful", "true"),
					resource.TestCheckOutput("dr_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDbDrRelationships_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_dr_relationships" "test" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]
}

data "huaweicloud_gaussdb_dr_relationships" "instance_name_filter" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  instance_name = huaweicloud_gaussdb_instance.test[0].name
}
locals {
  instance_name = huaweicloud_gaussdb_instance.test[0].name
}
output "instance_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_dr_relationships.instance_name_filter.relations) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_dr_relationships.instance_name_filter.relations[*].instance_name :
  v == local.instance_name]
  )
}

data "huaweicloud_gaussdb_dr_relationships" "instance_id_filter" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  instance_id = huaweicloud_gaussdb_instance.test[0].id
}
locals {
  instance_id = huaweicloud_gaussdb_instance.test[0].id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_dr_relationships.instance_id_filter.relations) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_dr_relationships.instance_id_filter.relations[*].instance_id : v == local.instance_id]
  )
}

data "huaweicloud_gaussdb_dr_relationships" "dr_role_filter" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  dr_role = "master"
}
output "dr_role_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_dr_relationships.dr_role_filter.relations) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_dr_relationships.dr_role_filter.relations[*].disaster_role : v == "master"]
  )
}

data "huaweicloud_gaussdb_dr_relationships" "dr_type_filter" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  dr_type = "stream"
}
output "dr_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_dr_relationships.dr_type_filter.relations) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_dr_relationships.dr_type_filter.relations[*].disaster_type : v == "stream"]
  )
}

data "huaweicloud_gaussdb_dr_relationships" "dr_status_filter" {
  depends_on = [huaweicloud_gaussdb_dr_relationship.test]

  dr_status = huaweicloud_gaussdb_dr_relationship.test.status
}
locals {
  dr_status = huaweicloud_gaussdb_dr_relationship.test.status
}
output "dr_status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_dr_relationships.dr_status_filter.relations) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_dr_relationships.dr_status_filter.relations[*].status : v == local.dr_status]
  )
}
`, testGaussDbDrRelationship_basic(name))
}
