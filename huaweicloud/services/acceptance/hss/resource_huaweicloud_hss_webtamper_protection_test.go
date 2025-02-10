package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getWebTamperProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		epsId   = hss.QueryAllEpsValue
		hostId  = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	return hss.GetWebTamperProtectionHost(client, region, epsId, hostId)
}

func TestAccWebTamperProtection_basic(t *testing.T) {
	var (
		wtpProtectHost interface{}
		rName          = "huaweicloud_hss_webtamper_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&wtpProtectHost,
		getWebTamperProtectionResourceFunc,
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

const testAccWebTamperProtection_base string = `
resource "huaweicloud_hss_quota" "test" {
  version     = "hss.version.wtp"
  period_unit = "month"
  period      = 1
}
`

func testAccWebTamperProtection_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  host_id               = "%[2]s"
  quota_id              = huaweicloud_hss_quota.test.id
  is_dynamics_protect   = true
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testAccWebTamperProtection_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  host_id               = "%[2]s"
  quota_id              = huaweicloud_hss_quota.test.id
  is_dynamics_protect   = false
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func TestAccWebTamperProtection_noFillQuotaId(t *testing.T) {
	var (
		wtpProtectHost interface{}
		rName          = "huaweicloud_hss_webtamper_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&wtpProtectHost,
		getWebTamperProtectionResourceFunc,
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
`, testAccWebTamperProtection_base, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}

func testAccWebTamperProtection_noFillQuotaId_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_hss_webtamper_protection" "test" {
  depends_on = [huaweicloud_hss_quota.test]

  host_id               = "%[2]s"
  is_dynamics_protect   = false
  enterprise_project_id = "0"
}
`, testAccWebTamperProtection_base, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
