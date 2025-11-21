package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApplicationAssignments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_application_assignments.test"
	name := acceptance.RandomAccResourceName()
	uuid, _ := uuid.GenerateUUID()
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
				Config: testDataSourceApplicationAssignments_basic(name, uuid),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "application_assignments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "application_assignments.0.application_urn"),
					resource.TestCheckResourceAttrPair(dataSource, "application_assignments.0.principal_id",
						"huaweicloud_identitycenter_user.test", "id"),
					resource.TestCheckResourceAttr(dataSource, "application_assignments.0.principal_type", "USER"),
				),
			},
		},
	})
}

func testDataSourceApplicationAssignments_basic(name string, uuid string) string {
	return fmt.Sprintf(`
%[1]s
 
data "huaweicloud_identitycenter_application_assignments" "test" {
  depends_on     = [huaweicloud_identitycenter_application_assignment.test]
  instance_id    = data.huaweicloud_identitycenter_instance.test.id
  principal_id   = huaweicloud_identitycenter_user.test.id
  principal_type = "USER"
}
`, testApplicationAssignment_basic(name, uuid))
}
