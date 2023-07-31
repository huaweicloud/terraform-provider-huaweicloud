package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenter_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_instance.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenter_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "id"),
					resource.TestCheckResourceAttrSet(rName, "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenter_basic() string {
	return `
data "huaweicloud_identitycenter_instance" "test" {}
`
}
