package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsSnapshotGroupsTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_snapshot_groups_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS snapshot group with tags before running.
			acceptance.TestAccPreCheckEVSFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsSnapshotGroupsTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.value.#"),
				),
			},
		},
	})
}

const testDataSourceEvsSnapshotGroupsTags_basic = `data "huaweicloud_evs_snapshot_groups_tags" "test" {}`
