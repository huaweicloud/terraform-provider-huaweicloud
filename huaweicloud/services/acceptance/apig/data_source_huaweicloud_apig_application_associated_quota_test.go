package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApplicationAssociatedQuota_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_apig_application_associated_quota.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
		rName  = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApplicationAssociatedQuota_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "id"),
					resource.TestCheckResourceAttrSet(dcName, "app_quota_id"),
					resource.TestCheckResourceAttrSet(dcName, "name"),
					resource.TestCheckResourceAttrSet(dcName, "call_limits"),
					resource.TestCheckResourceAttrSet(dcName, "time_unit"),
					resource.TestCheckResourceAttrSet(dcName, "time_interval"),
					resource.TestCheckResourceAttrSet(dcName, "remark"),
					resource.TestCheckResourceAttrSet(dcName, "bound_app_num"),
					resource.TestMatchResourceAttr(dcName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceApplicationAssociatedQuota_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_application" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = "%[3]s"
}

resource "huaweicloud_apig_application_quota" "test" {
  instance_id   = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name          = "%[3]s"
  time_unit     = "MINUTE"
  call_limits   = 100
  time_interval = 5
  description   = "Created by terraform script for testing"
}

resource "huaweicloud_apig_application_quota_associate" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  quota_id    = huaweicloud_apig_application_quota.test.id

  applications {
    id = huaweicloud_apig_application.test.id
  }
}

data "huaweicloud_apig_application_associated_quota" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  app_id      = huaweicloud_apig_application.test.id

  depends_on = [
    huaweicloud_apig_application_quota_associate.test,
  ]
}
`, common.TestBaseNetwork(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
