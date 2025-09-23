package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
)

func getQuotaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		epsId   = hss.QueryAllEpsValue
		quotaId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	quotaResp, err := hss.GetQuotaById(client, epsId, quotaId)
	if err != nil {
		return nil, err
	}

	if quotaResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return quotaResp, nil
}

func TestAccQuota_basic(t *testing.T) {
	var (
		obj          interface{}
		rName        = "huaweicloud_hss_quota.test"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getQuotaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a non default enterprise project ID.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.premium"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "used_status"),
					resource.TestCheckResourceAttrSet(rName, "charging_mode"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "shared_quota"),
					resource.TestCheckResourceAttrSet(rName, "is_trial_quota"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_name"),
				),
			},
			{
				Config: testAccQuota_basic_update(migrateEpsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", migrateEpsId),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"period_unit", "period", "auto_renew",
				},
			},
		},
	})
}

const testAccQuota_basic string = `
resource "huaweicloud_hss_quota" "test" {
  version               = "hss.version.premium"
  period_unit           = "month"
  period                = 1
  auto_renew            = "true"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`

func testAccQuota_basic_update(migrateEpsId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_quota" "test" {
  version               = "hss.version.premium"
  period_unit           = "month"
  period                = 1
  auto_renew            = "false"
  enterprise_project_id = "%s"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, migrateEpsId)
}

func TestAccQuota_periodUnitIsYear(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_hss_quota.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getQuotaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuota_periodUnitIsYear_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.premium"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "used_status"),
					resource.TestCheckResourceAttrSet(rName, "charging_mode"),
					resource.TestCheckResourceAttrSet(rName, "expire_time"),
					resource.TestCheckResourceAttrSet(rName, "shared_quota"),
					resource.TestCheckResourceAttrSet(rName, "is_trial_quota"),
					resource.TestCheckResourceAttrSet(rName, "enterprise_project_name"),
				),
			},
		},
	})
}

func testAccQuota_periodUnitIsYear_basic() string {
	return `
resource "huaweicloud_hss_quota" "test" {
  version               = "hss.version.premium"
  period_unit           = "year"
  period                = 1
  auto_renew            = "false"
  enterprise_project_id = "0"
}
`
}
