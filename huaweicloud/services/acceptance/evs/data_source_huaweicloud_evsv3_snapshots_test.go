package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV3Snapshots_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_evsv3_snapshots.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceV3Snapshots_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.status"),

					resource.TestCheckOutput("is_volume_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceV3Snapshots_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "Created by acc test"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 12
}

resource "huaweicloud_evsv3_snapshot" "test" {
  volume_id   = huaweicloud_evs_volume.test.id
  name        = "%[1]s"
  description = "Daily backup"

  metadata = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccDataSourceV3Snapshots_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evsv3_snapshots" "test" {
  depends_on = [huaweicloud_evsv3_snapshot.test]
}

# Filter using volume ID.
data "huaweicloud_evsv3_snapshots" "volume_id_filter" {
  volume_id = data.huaweicloud_evsv3_snapshots.test.snapshots.0.volume_id
}

locals {
  volume_id = data.huaweicloud_evsv3_snapshots.test.snapshots.0.volume_id
}

output "is_volume_id_filter_useful" {
  value = length(data.huaweicloud_evsv3_snapshots.volume_id_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evsv3_snapshots.volume_id_filter.snapshots[*].volume_id : v == local.volume_id]
  )
}

# Filter using name.
data "huaweicloud_evsv3_snapshots" "name_filter" {
  name = data.huaweicloud_evsv3_snapshots.test.snapshots.0.name
}

# The name parameter is a fuzzy query, so only the number of query results is verified.
output "is_name_filter_useful" {
  value = length(data.huaweicloud_evsv3_snapshots.name_filter.snapshots) > 0
}

# Filter using status.
data "huaweicloud_evsv3_snapshots" "status_filter" {
  status = data.huaweicloud_evsv3_snapshots.test.snapshots.0.status
}

locals {
  status = data.huaweicloud_evsv3_snapshots.test.snapshots.0.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_evsv3_snapshots.volume_id_filter.snapshots) > 0 && alltrue(
    [for v in data.huaweicloud_evsv3_snapshots.volume_id_filter.snapshots[*].status : v == local.status]
  )
}

# Filter using non existent name.
data "huaweicloud_evsv3_snapshots" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_evsv3_snapshots.not_found.snapshots) == 0
}
`, testAccDataSourceV3Snapshots_base(name))
}
