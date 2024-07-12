package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getWebTamperProtectionFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcHssV5Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS v5 client: %s", err)
	}

	var (
		region = acceptance.HW_REGION_NAME
		epsId  = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		hostId = state.Primary.ID
	)

	// If the enterprise project ID is not set during query, query all enterprise projects.
	if epsId == "" {
		epsId = hss.QueryAllEpsValue
	}
	listOpts := hssv5model.ListWtpProtectHostRequest{
		Region:              region,
		EnterpriseProjectId: utils.String(epsId),
		HostId:              utils.String(hostId),
		ProtectStatus:       utils.String(string(hss.ProtectStatusOpened)),
	}

	resp, err := client.ListWtpProtectHost(&listOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying HSS web tamper protection hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return nil, fmt.Errorf("the host (%s) for HSS web tamper protection does not exist", hostId)
	}

	wtpProtectHostList := *resp.DataList
	if len(wtpProtectHostList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return wtpProtectHostList[0], nil
}

func TestAccWebTamperProtection_basic(t *testing.T) {
	var (
		wtpProtectHost *hssv5model.WtpProtectHostResponseInfo
		rName          = "huaweicloud_hss_webtamper_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&wtpProtectHost,
		getWebTamperProtectionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWebTamperProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "os_bit"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttr(rName, "protect_status", string(hss.ProtectStatusOpened)),
					resource.TestCheckResourceAttr(rName, "rasp_protect_status", string(hss.ProtectStatusOpened)),
				),
			},
			{
				Config: testAccWebTamperProtection_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "os_bit"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttr(rName, "protect_status", string(hss.ProtectStatusOpened)),
					resource.TestCheckResourceAttr(rName, "rasp_protect_status", string(hss.ProtectStatusClosed)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"quota_id", "is_dynamics_protect", "enterprise_project_id",
				},
			},
		},
	})
}

func testAccWebTamperProtection_base() string {
	return `
resource "huaweicloud_hss_quota" "test" {
  version     = "hss.version.wtp"
  period_unit = "month"
  period      = 1
}`
}

func testAccWebTamperProtection_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  host_id               = "%[2]s"
  quota_id              = huaweicloud_hss_quota.test.id
  is_dynamics_protect   = true
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base(), acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testAccWebTamperProtection_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  host_id               = "%[2]s"
  quota_id              = huaweicloud_hss_quota.test.id
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base(), acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func TestAccWebTamperProtection_noFillQuotaId(t *testing.T) {
	var (
		wtpProtectHost *hssv5model.WtpProtectHostResponseInfo
		rName          = "huaweicloud_hss_webtamper_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&wtpProtectHost,
		getWebTamperProtectionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWebTamperProtection_noFillQuotaId(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "os_bit"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttr(rName, "protect_status", string(hss.ProtectStatusOpened)),
					resource.TestCheckResourceAttr(rName, "rasp_protect_status", string(hss.ProtectStatusOpened)),
				),
			},
			{
				Config: testAccWebTamperProtection_noFillQuotaId_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "os_bit"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttr(rName, "protect_status", string(hss.ProtectStatusOpened)),
					resource.TestCheckResourceAttr(rName, "rasp_protect_status", string(hss.ProtectStatusClosed)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"quota_id", "is_dynamics_protect", "enterprise_project_id",
				},
			},
		},
	})
}

func testAccWebTamperProtection_noFillQuotaId() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  depends_on = [huaweicloud_hss_quota.test]

  host_id               = "%[2]s"
  is_dynamics_protect   = true
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base(), acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testAccWebTamperProtection_noFillQuotaId_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  depends_on = [huaweicloud_hss_quota.test]

  host_id               = "%[2]s"
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base(), acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
