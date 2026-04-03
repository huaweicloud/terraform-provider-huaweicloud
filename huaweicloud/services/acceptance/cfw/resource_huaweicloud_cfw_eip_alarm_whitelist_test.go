package cfw

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEipAlarmWhitelistResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl  = "v1/{project_id}/eip/alarm-whitelist/{fw_instance_id}"
		product  = "cfw"
		publicIP = state.Primary.Attributes["public_ip"]
		epsId    = state.Primary.Attributes["enterprise_project_id"]
	)

	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", state.Primary.ID)
	requestPath += cfw.BuildEipAlarmWhitelistQueryParams(epsId, publicIP)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	listResp := utils.PathSearch("data.list", respBody, make([]interface{}, 0)).([]interface{})
	if len(listResp) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	typeResp := utils.PathSearch("type", listResp[0], float64(0)).(float64)
	if typeResp == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return listResp[0], nil
}

func TestAccEipAlarmWhitelist_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cfw_eip_alarm_whitelist.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEipAlarmWhitelistResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPrecheckEipIDAndIP(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// Destroy does not remove the whitelist entry from the cloud (no delete API).
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testEipAlarmWhitelist_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "fw_instance_id", acceptance.HW_CFW_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "eip_id", acceptance.HW_EIP_ID),
					resource.TestCheckResourceAttr(rName, "public_ip", acceptance.HW_EIP_ADDRESS),
					resource.TestCheckResourceAttrSet(rName, "type"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"object_id", "enterprise_project_id"},
				ImportStateIdFunc:       testEipAlarmWhitelistImportStateId(rName),
			},
		},
	})
}

func testEipAlarmWhitelist_base() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testEipAlarmWhitelist_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_eip_alarm_whitelist" "test" {
  fw_instance_id = "%[2]s"
  eip_id         = "%[3]s"
  public_ip      = "%[4]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}
`, testEipAlarmWhitelist_base(), acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_EIP_ID, acceptance.HW_EIP_ADDRESS)
}

func testEipAlarmWhitelistImportStateId(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		resourceId := rs.Primary.ID
		if resourceId == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found", name)
		}

		publicIP := rs.Primary.Attributes["public_ip"]
		if publicIP == "" {
			return "", fmt.Errorf("attribute (public_ip) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", resourceId, publicIP), nil
	}
}
