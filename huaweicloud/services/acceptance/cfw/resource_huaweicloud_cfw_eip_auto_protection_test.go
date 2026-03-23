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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEipAutoProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl  = "v1/{project_id}/eip/auto-protect-status/{object_id}"
		product  = "cfw"
		objectID = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{object_id}", objectID)
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

	statusResp := utils.PathSearch("data.status", respBody, float64(0)).(float64)
	if statusResp == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return utils.PathSearch("data", respBody, nil), nil
}

func TestAccEipAutoProtection_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cfw_eip_auto_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getEipAutoProtectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEipAutoProtection_basic(1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttrSet(rName, "object_id"),
					resource.TestCheckResourceAttrSet(rName, "available_eip_count"),
					resource.TestCheckResourceAttrSet(rName, "beyond_max_count"),
					resource.TestCheckResourceAttrSet(rName, "eip_protected_self"),
					resource.TestCheckResourceAttrSet(rName, "eip_total"),
					resource.TestCheckResourceAttrSet(rName, "eip_un_protected"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testEipAutoProtectionImportState(rName),
			},
		},
	})
}

func testEipAutoProtection_base() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testEipAutoProtection_basic(status int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_eip_auto_protection" "test" {
  fw_instance_id = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  status         = %[2]d
}
`, testEipAutoProtection_base(), status)
}

func testEipAutoProtectionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		fwInstanceID := rs.Primary.Attributes["fw_instance_id"]
		if fwInstanceID == "" {
			return "", fmt.Errorf("attribute (fw_instance_id) of Resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", fwInstanceID, rs.Primary.ID), nil
	}
}
