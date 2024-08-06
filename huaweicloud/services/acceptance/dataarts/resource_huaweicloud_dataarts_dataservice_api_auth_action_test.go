package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataServiceApiAuthAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Just test whether the authorize request was executed successful.
				Config: testAccDataServiceApiAuthAction_basic_step1(name),
			},
			{
				// Just test whether the authorize status was cancelled successful.
				Config: testAccDataServiceApiAuthAction_basic_step2(name),
			},
		},
	})
}

// Authorize some APPs to access the API.
func testAccDataServiceApiAuthAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]
  count      = 2

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id      = element(huaweicloud_dataarts_dataservice_app.test[*].id, count.index)
  type        = "APPLY_TYPE_AUTHORIZE"
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

// Cancel the API access permissions for all APPs.
func testAccDataServiceApiAuthAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]
  count      = 2

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id      = element(huaweicloud_dataarts_dataservice_app.test[*].id, count.index)
  type        = "APPLY_TYPE_APP_CANCEL_AUTHORIZE"
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}
