package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSystemPolicyAttachments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_system_policy_attachments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSystemPolicyAttachments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
				),
			},
		},
	})
}

func testDataSourceSystemPolicyAttachments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_system_policy_attachments" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = huaweicloud_identitycenter_permission_set.test.id

  depends_on = [huaweicloud_identitycenter_system_policy_attachment.test]
}
`, testSystemPolicyAttachment_basic(name))
}
