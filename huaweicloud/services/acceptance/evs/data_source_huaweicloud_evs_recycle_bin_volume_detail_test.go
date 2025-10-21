package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRecycleBinVolumeDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evs_recycle_bin_volume_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a volume ID that is in the recycle bin.
			acceptance.TestAccPreCheckEVSRecycleBinVolumeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRecycleBinVolumeDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.multiattach"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.bootable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume.0.enterprise_project_id"),
				),
			},
		},
	})
}

func testDataSourceRecycleBinVolumeDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_evs_recycle_bin_volume_detail" "test" {
  volume_id = "%s"
}
`, acceptance.HW_EVS_RECYCLE_BIN_VOLUME_ID)
}
