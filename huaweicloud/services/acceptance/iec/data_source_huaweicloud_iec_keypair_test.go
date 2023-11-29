package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKeypairDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("KeyPair-%s", acctest.RandString(4))
	dataSourceName := "data.huaweicloud_iec_keypair.by_name"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKeypair_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "public_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "fingerprint"),
				),
			},
		},
	})
}

func testAccDataSourceKeypair_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_keypair" "kp_1" {
  name = "%s"
}

data "huaweicloud_iec_keypair" "by_name" {
  name = huaweicloud_iec_keypair.kp_1.name
}
`, rName)
}
