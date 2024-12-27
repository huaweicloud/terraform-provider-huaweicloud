package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePermissionSetProvisioningAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_permission_set_provisioning_accounts.test"
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
				Config: testDataSourcePermissionSetProvisioningAccounts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "account_ids.#"),
					resource.TestCheckOutput("is_provisioning_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePermissionSetProvisioningAccounts_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_permission_set_provisioning_accounts" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.test.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
}

data "huaweicloud_identitycenter_permission_set_provisioning_accounts" "provision" {
  instance_id         = data.huaweicloud_identitycenter_instance.test.id
  permission_set_id   = huaweicloud_identitycenter_permission_set.test.id
  provisioning_status = "LATEST_PERMISSION_SET_PROVISIONED"
}

data "huaweicloud_identitycenter_permission_set_provisioning_accounts" "not_provision" {
  instance_id         = data.huaweicloud_identitycenter_instance.test.id
  permission_set_id   = huaweicloud_identitycenter_permission_set.test.id
  provisioning_status = "LATEST_PERMISSION_SET_NOT_PROVISIONED"
}

locals {
  list_by_provision     = data.huaweicloud_identitycenter_permission_set_provisioning_accounts.provision.account_ids
  list_by_not_provision = data.huaweicloud_identitycenter_permission_set_provisioning_accounts.not_provision.account_ids
}

output "is_provisioning_status_filter_useful" {
  value = length(local.list_by_provision) == 1 && length(local.list_by_not_provision) == 0
}
`, testProvisionPermissionSet_basic(name))
}
