package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getWarRoomResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetWarRoom(client, state.Primary.ID)
}

func TestAccResourceWarRoom_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_war_room.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWarRoomResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckCocRoleID(t)
			acceptance.TestAccPreCheckCocSceneID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testWarRoom_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_id"),
					resource.TestCheckResourceAttrSet(resourceName, "incident.#"),
					resource.TestCheckResourceAttrSet(resourceName, "incident.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "incident.0.incident_id"),
					resource.TestCheckResourceAttrSet(resourceName, "incident.0.is_change_event"),
					resource.TestCheckResourceAttrSet(resourceName, "incident.0.failure_level"),
					resource.TestCheckResourceAttrSet(resourceName, "occur_time"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "circular_level"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_status.#"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_status.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_status.0.name_zh"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_status.0.name_en"),
					resource.TestCheckResourceAttrSet(resourceName, "war_room_status.0.type"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"schedule_group", "participant", "notification_type"},
			},
		},
	})
}

func testWarRoom_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_war_room" "test" {
  war_room_name         = "%[2]s"
  application_id_list   = [huaweicloud_coc_application.test.id]
  incident_number       = huaweicloud_coc_incident.test.id
  war_room_admin        = "%[3]s"
  enterprise_project_id = "0"

  schedule_group {
    role_id  = "%[4]s"
    scene_id = "%[5]s"
  }

  depends_on = [huaweicloud_coc_incident_handle.test]
}
`, testIncidentHandle_basic(name), name, acceptance.HW_USER_ID, acceptance.HW_COC_ROLE_ID,
		acceptance.HW_COC_SCENE_ID)
}
