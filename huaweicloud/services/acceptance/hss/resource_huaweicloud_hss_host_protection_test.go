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

func getHostProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region = acceptance.HW_REGION_NAME
		epsId  = hss.QueryAllEpsValue
		id     = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + "v5/{project_id}/host-management/hosts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?enterprise_project_id=%v&host_id=%v", epsId, id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS host, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	hostResp := utils.PathSearch("data_list[0]", getRespBody, nil)
	protectStatus := utils.PathSearch("protect_status", hostResp, "").(string)
	if hostResp == nil || protectStatus == string(hss.ProtectStatusClosed) {
		return nil, golangsdk.ErrDefault404{}
	}

	return hostResp, nil
}

func TestAccHostProtection_basic(t *testing.T) {
	var (
		host  interface{}
		rName = "huaweicloud_hss_host_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&host,
		getHostProtectionResourceFunc,
	)

	// Because after closing the protection, the ECS instance will automatically switch to free basic protection,
	// so avoid CheckDestroy here.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHostProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.basic"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "host_status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "agent_id"),
					resource.TestCheckResourceAttrSet(rName, "agent_status"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "detect_result"),
					resource.TestCheckResourceAttrSet(rName, "asset_value"),
					resource.TestCheckResourceAttrSet(rName, "open_time"),
					resource.TestCheckResourceAttrPair(rName, "quota_id", "huaweicloud_hss_quota.test", "id"),
				),
			},
			{
				Config: testAccHostProtection_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.enterprise"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "host_status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "agent_id"),
					resource.TestCheckResourceAttrSet(rName, "agent_status"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "detect_result"),
					resource.TestCheckResourceAttrSet(rName, "asset_value"),
					resource.TestCheckResourceAttrSet(rName, "open_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"quota_id", "is_wait_host_available",
				},
			},
		},
	})
}

func testAccHostProtection_base() string {
	return `
resource "huaweicloud_hss_quota" "test" {
  version     = "hss.version.basic"
  period_unit = "month"
  period      = 1
}`
}

func testAccHostProtection_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_host_protection" "test" {
  host_id                = "%[2]s"
  version                = "hss.version.basic"
  charging_mode          = "prePaid"
  quota_id               = huaweicloud_hss_quota.test.id
  is_wait_host_available = true
}
`, testAccHostProtection_base(), acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testAccHostProtection_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_host_protection" "test" {
  host_id       = "%[1]s"
  version       = "hss.version.enterprise"
  charging_mode = "postPaid"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
