package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataServiceAuthorizedApps_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_dataservice_authorized_apps.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataServiceAuthorizedApps_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "apps.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				// Cancel the authorization, prepare to destroy the environment.
				Config: testAccDataSourceDataServiceAuthorizedApps_basic_step2(name),
			},
		},
	})
}

func testAccDataSourceDataServiceAuthorizedApps_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_ids     = slice(huaweicloud_dataarts_dataservice_app.test[*].id, 0, 2)
}

data "huaweicloud_dataarts_dataservice_authorized_apps" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_auth.test]

  workspace_id = "%[2]s"
  dlm_type     = "EXCLUSIVE"

  api_id = huaweicloud_dataarts_dataservice_api.test.id
}
`, testAccDataServiceApiAuth_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataSourceDataServiceAuthorizedApps_basic_step2(name string) string {
	return testAccDataServiceApiAuth_basic_step2(name)
}
