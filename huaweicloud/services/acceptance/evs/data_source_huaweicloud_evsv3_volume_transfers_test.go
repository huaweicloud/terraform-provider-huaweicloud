package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsV3VolumeTransfers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_evsv3_volume_transfers.test"
	name := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEvsV3VolumeTransfers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSource, "transfers.0.links.#"),
				),
			},
		},
	})
}

func testDataSourceEvsV3VolumeTransfers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_evsv3_volume_transfers" "test" {
  depends_on = [huaweicloud_evsv3_volume_transfer.test]
}
`, testAccV3VolumeTransfer_basic(name))
}
