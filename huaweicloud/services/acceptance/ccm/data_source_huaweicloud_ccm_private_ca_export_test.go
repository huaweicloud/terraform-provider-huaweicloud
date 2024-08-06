package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateCaExport_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ccm_private_ca_export.test"
		rName      = acceptance.RandomAccResourceNameWithDash()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourcePrivateCaExport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificate"),
				),
			},
		},
	})
}

func testDataSourceDataSourcePrivateCaExport_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_ca_export" "test" {
  ca_id = huaweicloud_ccm_private_ca.test_root.id
}
`, tesPrivateCA_postpaid_root(name))
}
