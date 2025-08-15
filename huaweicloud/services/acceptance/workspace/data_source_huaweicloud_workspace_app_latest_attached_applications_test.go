package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLatestAttachedApplications_basic(t *testing.T) {
	var (
		name        = acceptance.RandomAccResourceName()
		dataSource  = "data.huaweicloud_workspace_app_latest_attached_applications.test"
		dc          = acceptance.InitDataSourceCheck(dataSource)
		randUUID, _ = uuid.GenerateUUID()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceLatestAttachedApplications_notFound(randUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf("'%s' is a non-existing cloud application server", randUUID)),
			},
			{
				Config: testAccDataSourceLatestAttachedApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.app_id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.record_id"),
				),
			},
		},
	})
}

func testAccDataSourceLatestAttachedApplications_notFound(randUUID string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_latest_attached_applications" "test" {
  server_id = "%s"
}
`, randUUID)
}

func testAccDataSourceLatestAttachedApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_latest_attached_applications" "test" {
  depends_on = [huaweicloud_workspace_app_application_batch_attach.test]

  server_id = huaweicloud_workspace_app_image_server.test.id
}
`, testAccResourceAppApplicationBatchAttach_basic(name))
}
