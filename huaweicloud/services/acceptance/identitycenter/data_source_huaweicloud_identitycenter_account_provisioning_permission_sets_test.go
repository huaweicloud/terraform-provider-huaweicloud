package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentitycenterAccountProvisioningPermissionSets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_account_provisioning_permission_sets.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIdentitycenterAccountProvisioningPermissionSets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceIdentitycenterAccountProvisioningPermissionSets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_account_provisioning_permission_sets" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.system.id
  account_id  = huaweicloud_identitycenter_user.test.id
}
`, testIdentityCenterUser_basic(name))
}
