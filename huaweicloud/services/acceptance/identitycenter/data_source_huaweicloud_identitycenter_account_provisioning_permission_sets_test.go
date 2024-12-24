package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccountProvisioningPermissionSets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_account_provisioning_permission_sets.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccountProvisioningPermissionSets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.#"),
					resource.TestCheckOutput("is_provisioning_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAccountProvisioningPermissionSets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_account_provisioning_permission_sets" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.test.id
  account_id  = "%[2]s"
}

data "huaweicloud_identitycenter_account_provisioning_permission_sets" "provision" {
  instance_id         = data.huaweicloud_identitycenter_instance.test.id
  account_id          = "%[2]s"
  provisioning_status = "LATEST_PERMISSION_SET_PROVISIONED"
}

data "huaweicloud_identitycenter_account_provisioning_permission_sets" "not_provision" {
  instance_id         = data.huaweicloud_identitycenter_instance.test.id
  account_id          = "%[2]s"
  provisioning_status = "LATEST_PERMISSION_SET_NOT_PROVISIONED"
}

locals {
  list_by_provision     = data.huaweicloud_identitycenter_account_provisioning_permission_sets.provision.permission_sets
  list_by_not_provision = data.huaweicloud_identitycenter_account_provisioning_permission_sets.not_provision.permission_sets
}

output "is_provisioning_status_filter_useful" {
  value = length(local.list_by_provision) == 1 && length(local.list_by_not_provision) == 0
}
`, testProvisionPermissionSet_basic(name), acceptance.HW_IDENTITY_CENTER_ACCOUNT_ID)
}
