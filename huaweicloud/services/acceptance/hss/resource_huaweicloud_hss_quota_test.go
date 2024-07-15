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

func getQuotaFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "hss"
		epsId   = state.Primary.Attributes["enterprise_project_id"]
		id      = state.Primary.ID
	)

	// If the enterprise project ID is not set during query, set to query all enterprise projects.
	if epsId == "" {
		epsId = hss.QueryAllEpsValue
	}

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS client: %s", err)
	}

	quotas, err := hss.GetQuotaById(client, id, epsId)
	if err != nil {
		return nil, err
	}

	if len(quotas) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return quotas[0], nil
}

func TestAccQuota_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_hss_quota.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getQuotaFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccQuota_basic(),
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
				Config: testAccQuota_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func testAccQuota_basic() string {
	return `
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
}

func testAccQuota_basic_update() string {
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
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccQuota_periodUnitIsYear(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_hss_quota.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getQuotaFunc,
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
