package hss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEventSystemUserWhiteListResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		epsId   = hss.QueryAllEpsValue
		hostID  = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	queryPath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath += fmt.Sprintf("?host_id=%s", hostID)
	queryPath += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS event system user white list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(dataList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataList[0], nil
}

func TestAccEventSystemUserWhiteList_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_hss_event_system_user_white_list.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getEventSystemUserWhiteListResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEventSystemUserWhiteList_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.0", "test_user1"),
					resource.TestCheckResourceAttr(resourceName, "remarks", "remarks_test"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_name"),
					resource.TestCheckResourceAttrSet(resourceName, "host_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
				),
			},
			{
				Config: testEventSystemUserWhiteList_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.0", "test_user1"),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.1", "test_user2"),
					resource.TestCheckResourceAttr(resourceName, "system_user_name_list.2", "test_user3"),
					resource.TestCheckResourceAttr(resourceName, "remarks", "remarks_update"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enterprise_project_id"},
			},
		},
	})
}

func testEventSystemUserWhiteList_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_event_system_user_white_list" "test" {
  host_id               = "%[1]s"
  system_user_name_list = ["test_user1"]
  remarks               = "remarks_test"
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testEventSystemUserWhiteList_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_event_system_user_white_list" "test" {
  host_id               = "%[1]s"
  system_user_name_list = ["test_user1", "test_user2", "test_user3"]
  remarks               = "remarks_update"
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
