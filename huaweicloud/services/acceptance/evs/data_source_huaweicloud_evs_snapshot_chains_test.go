package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsSnapshotChains_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_snapshot_chains.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS standard snapshot with tags before running.
			acceptance.TestAccPreCheckEVSFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsSnapshotChains_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.#"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.snapshot_count"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.capacity"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "snapshot_chains.0.updated_at"),
				),
			},
		},
	})
}

const testDataSourceEvsSnapshotChains_basic = `data "huaweicloud_evs_snapshot_chains" "test" {}`
