package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProtectionGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_protection_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceProtectionGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.disaster_recovery_drill_num"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.dr_type"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.priority_station"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.protected_instance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.protection_type"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.replication_num"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.server_type"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.source_availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.source_vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.target_availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.target_vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "server_groups.0.updated_at"),

					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_availability_zone_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceProtectionGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_protection_groups" "test" {
  depends_on = [huaweicloud_sdrs_protection_group.test]
}

# Filter by status
locals {
  status = data.huaweicloud_sdrs_protection_groups.test.server_groups[0].status
}

data "huaweicloud_sdrs_protection_groups" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_sdrs_protection_groups.filter_by_status.server_groups[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by name
locals {
  name = data.huaweicloud_sdrs_protection_groups.test.server_groups[0].name
}

data "huaweicloud_sdrs_protection_groups" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_sdrs_protection_groups.filter_by_name.server_groups[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by availability_zone
locals {
  availability_zone = data.huaweicloud_sdrs_protection_groups.test.server_groups[0].source_availability_zone
}

data "huaweicloud_sdrs_protection_groups" "filter_by_availability_zone" {
  availability_zone = local.availability_zone
}

locals {
  availability_zone_filter_result = [
    for v in data.huaweicloud_sdrs_protection_groups.filter_by_availability_zone.server_groups[*].source_availability_zone :
    v == local.availability_zone
  ]
}

output "is_availability_zone_filter_useful" {
  value = length(local.availability_zone_filter_result) > 0 && alltrue(local.availability_zone_filter_result)
}
`, testProtectionGroup_basic(name))
}
