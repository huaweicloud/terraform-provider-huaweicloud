package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getDataServiceInstanceLogDumpResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return dataarts.GetDataServiceInstanceLogDump(client, state.Primary.Attributes["workspace_id"],
		state.Primary.Attributes["instance_id"])
}

func TestAccDataServiceInstanceLogDump_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dataarts_dataservice_instance_log_dump.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDataServiceInstanceLogDumpResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataServiceInstanceLogDump_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "obs"),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DATAARTS_INSTANCE_ID),
				),
			},
			{
				Config: testAccDataServiceInstanceLogDump_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "lts"),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DATAARTS_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_group_name", "huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_name", "huaweicloud_lts_stream.test", "stream_name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccDataServiceInstanceLogDumpImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{"enable_force_new"},
			},
		},
	})
}

func testAccDataServiceInstanceLogDumpImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found in state", rName)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		instanceId := rs.Primary.ID
		if workspaceId == "" || instanceId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<instance_id>', but got '%s/%s'",
				workspaceId, instanceId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, instanceId), nil
	}
}

func testAccDataServiceInstanceLogDump_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_dataservice_instance_log_dump" "test" {
  workspace_id = "%[1]s"
  instance_id  = "%[2]s"
  type         = "obs"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_INSTANCE_ID)
}

func testAccDataServiceInstanceLogDump_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 7
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}

# The API requires a minimum interval of 30 seconds between API calls.
resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    huaweicloud_lts_group.test,
    huaweicloud_lts_stream.test,
  ]

  create_duration = "30s"
}

resource "huaweicloud_dataarts_dataservice_instance_log_dump" "test" {
  workspace_id     = "%[2]s"
  instance_id      = "%[3]s"
  type             = "lts"
  log_group_id     = huaweicloud_lts_group.test.id
  log_group_name   = huaweicloud_lts_group.test.group_name
  log_stream_id    = huaweicloud_lts_stream.test.id
  log_stream_name  = huaweicloud_lts_stream.test.stream_name
  enable_force_new = "true"

  depends_on = [time_sleep.wait_30_seconds]
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_INSTANCE_ID)
}
