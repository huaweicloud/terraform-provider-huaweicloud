package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplicationInstanceProfileDelete_basic(t *testing.T) {
	applicationInstanceId := acceptance.HW_IDENTITY_CENTER_APPLICATION_INSTANCE_ID
	profileId := acceptance.HW_IDENTITY_CENTER_APPLICATION_INSTANCE_PROFILE_ID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterApplicationInstanceProfileId(t)
			acceptance.TestAccPreCheckIdentityCenterApplicationInstanceId(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testApplicationInstanceProfileDelete_basic(applicationInstanceId, profileId),
			},
		},
	})
}

func testApplicationInstanceProfileDelete_basic(applicationInstanceId string, profileId string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_application_instance_profile_delete" "test" {
  instance_id             = data.huaweicloud_identitycenter_instance.test.id
  application_instance_id = "%[1]s"
  profile_id              = "%[2]s"
}
`, applicationInstanceId, profileId)
}
