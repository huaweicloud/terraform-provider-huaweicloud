package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentitycenterSystemPolicyAttachments_basic(t *testing.T) {
	dataSource := "data.resource_huaweicloud_identitycenter_system_policy_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIdentitycenterSystemPolicyAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "attached_managed_roles.#"),
					resource.TestCheckResourceAttrSet(rName, "attached_managed_roles.0.role_id"),
					resource.TestCheckResourceAttrSet(rName, "attached_managed_roles.0.role_name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceIdentitycenterSystemPolicyAttachments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "resource_huaweicloud_identitycenter_system_policy_attachments" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id
}
`, testSystemPolicyAttachment_basic(name))
}
