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

func getScalingPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetAppServerGroupScalingPolicy(client, state.Primary.ID)
}

func TestAccAppServerGroupScalingPolicy_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_workspace_app_server_group_scaling_policy.test"

		rc = acceptance.InitResourceCheck(rName, &obj, getScalingPolicyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerGroupScalingPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "max_scaling_amount", "10"),
					resource.TestCheckResourceAttr(rName, "single_expansion_count", "2"),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "scaling_policy_by_session.0.session_usage_threshold", "80"),
					resource.TestCheckResourceAttr(rName, "scaling_policy_by_session.0.shrink_after_session_idle_minutes", "30"),
				),
			},
			{
				Config: testAccAppServerGroupScalingPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "max_scaling_amount", "20"),
					resource.TestCheckResourceAttr(rName, "single_expansion_count", "3"),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttr(rName, "scaling_policy_by_session.0.session_usage_threshold", "90"),
					resource.TestCheckResourceAttr(rName, "scaling_policy_by_session.0.shrink_after_session_idle_minutes", "60"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppServerGroupScalingPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_server_group_scaling_policy" "test" {
  server_group_id        = "%[1]s"
  max_scaling_amount     = 10
  single_expansion_count = 2

  scaling_policy_by_session {
    session_usage_threshold           = 80
    shrink_after_session_idle_minutes = 30
  }
}`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}

func testAccAppServerGroupScalingPolicy_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_server_group_scaling_policy" "test" {
  server_group_id        = "%[1]s"
  max_scaling_amount     = 20
  single_expansion_count = 3

  scaling_policy_by_session {
    session_usage_threshold           = 90
    shrink_after_session_idle_minutes = 60
  }
}`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}
