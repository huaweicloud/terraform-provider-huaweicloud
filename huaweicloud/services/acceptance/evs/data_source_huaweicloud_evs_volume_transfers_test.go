package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsVolumeTransfers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evs_volume_transfers.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceEvsVolumeTransfers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.volume_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceEvsVolumeTransfers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evs_volume_transfers" "test" {
  depends_on = [huaweicloud_evs_volume_transfer.test]
}
`, testAccVolumeTransfer_basic(name))
}
