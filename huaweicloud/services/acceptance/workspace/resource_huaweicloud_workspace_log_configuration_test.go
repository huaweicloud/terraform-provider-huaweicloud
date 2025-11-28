package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getLogConfigurationResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetLogConfiguration(client)
}

func TestAccLogConfiguration_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_workspace_log_configuration.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getLogConfigurationResourceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		baseConfig = testLogConfiguration_base(name)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLogConfiguration_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test.0", "id"),
				),
			},
			{
				Config: testLogConfiguration_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id", "huaweicloud_lts_group.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_stream_id", "huaweicloud_lts_stream.test.1", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLogConfiguration_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  count = 2

  group_name  = format("%[1]s_%%d", count.index)
  ttl_in_days = 10
}

resource "huaweicloud_lts_stream" "test" {
  count = 2

  group_id    = huaweicloud_lts_group.test[count.index].id
  stream_name = format("%[1]s_%%d", count.index)
}
`, name)
}

func testLogConfiguration_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_log_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test[0].id
  log_stream_id = huaweicloud_lts_stream.test[0].id
}
`, baseConfig, name)
}

func testLogConfiguration_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_log_configuration" "test" {
  log_group_id  = huaweicloud_lts_group.test[1].id
  log_stream_id = huaweicloud_lts_stream.test[1].id
}
`, baseConfig, name)
}
