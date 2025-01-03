package identitycenter

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePermissionSetProvisionings_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_permission_set_provisionings.test"
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
				Config: testDataSourcePermissionSetProvisionings_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "provisionings.#"),
					resource.TestCheckResourceAttrSet(dataSource, "provisionings.0.request_id"),
					resource.TestCheckResourceAttrSet(dataSource, "provisionings.0.status"),
					resource.TestMatchResourceAttr(dataSource,
						"provisionings.0.created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePermissionSetProvisionings_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identitycenter_permission_set_provisionings" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}

locals {
  status = data.huaweicloud_identitycenter_permission_set_provisionings.test.provisionings[0].status
}

data "huaweicloud_identitycenter_permission_set_provisionings" "filter_by_status" {
  instance_id = data.huaweicloud_identitycenter_instance.test.id
  status      = local.status
}

locals {
  list_by_status = data.huaweicloud_identitycenter_permission_set_provisionings.filter_by_status.provisionings
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}
`, testProvisionPermissionSet_basic(name))
}
