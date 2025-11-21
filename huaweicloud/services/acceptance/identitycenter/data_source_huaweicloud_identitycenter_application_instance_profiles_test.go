package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceApplicationInstanceProfiles_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_application_instance_profiles.test"
	dc := acceptance.InitDataSourceCheck(rName)

	uuid, _ := uuid.GenerateUUID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceApplicationInstanceProfiles_basic(uuid),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "Default"),
					resource.TestCheckResourceAttr(rName, "status", "ENABLED"),
					resource.TestCheckResourceAttrSet(rName, "profile_id"),
				),
			},
		},
	})
}

func testAccDatasourceApplicationInstanceProfiles_basic(uuid string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_application_instance_profiles" "test" {
  depends_on              = [huaweicloud_identitycenter_application_instance.test]
  instance_id             = data.huaweicloud_identitycenter_instance.test.id
  application_instance_id = huaweicloud_identitycenter_application_instance.test.id
}
`, testApplicationInstance_basic(uuid))
}
