package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReplicationPairs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_replication_pairs.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReplicationPairs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.volume_ids"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.replication_model"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.record_metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.record_metadata.0.multiattach"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.record_metadata.0.bootable"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.record_metadata.0.volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.record_metadata.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.server_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_pairs.0.priority_station"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceReplicationPairs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_replication_pairs" "test" {
  depends_on = [huaweicloud_sdrs_replication_pair.test]
}

# Filter by name
locals {
  name = data.huaweicloud_sdrs_replication_pairs.test.replication_pairs[0].name
}

data "huaweicloud_sdrs_replication_pairs" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_sdrs_replication_pairs.filter_by_name.replication_pairs[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by server_group_id
locals {
  server_group_id = data.huaweicloud_sdrs_replication_pairs.test.replication_pairs[0].server_group_id
}

data "huaweicloud_sdrs_replication_pairs" "filter_by_server_group_id" {
  server_group_id = local.server_group_id
}

locals {
  server_group_id_filter_result = [
    for v in data.huaweicloud_sdrs_replication_pairs.filter_by_server_group_id.replication_pairs[*].server_group_id :
    v == local.server_group_id
  ]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_filter_result) > 0 && alltrue(local.server_group_id_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_sdrs_replication_pairs.test.replication_pairs[0].status
}

data "huaweicloud_sdrs_replication_pairs" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_sdrs_replication_pairs.filter_by_status.replication_pairs[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testReplicationPair_basic(name))
}
