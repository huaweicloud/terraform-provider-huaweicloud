package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf/v1/clouds"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getCloudInstanceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF v1 client: %s", err)
	}
	return waf.QueryCloudInstance(client, state.Primary.ID)
}

func TestAccCloudInstance_basic(t *testing.T) {
	var instance clouds.Instance
	rName := "huaweicloud_waf_cloud_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&instance,
		getCloudInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeIntroduction)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCloudInstance_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeStandard)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
					"charging_mode",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func TestAccCloudInstance_withEpsID(t *testing.T) {
	var instance clouds.Instance
	rName := "huaweicloud_waf_cloud_instance.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&instance,
		getCloudInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstance_basic_withEpsID(acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeIntroduction)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccCloudInstance_update_withEpsID(acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeStandard)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
				),
			},
		},
	})
}

const testAccCloudInstance_basic = `
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "detection"
  enterprise_project_id = "0"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`

const testAccCloudInstance_update = `
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "professional"
  enterprise_project_id = "0"

  bandwidth_expack_product {
    resource_size = 1
  }
  domain_expack_product {
    resource_size = 1
  }
  rule_expack_product {
    resource_size = 1
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`

func testAccCloudInstance_basic_withEpsID(epsID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "detection"
  enterprise_project_id = "%s"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, epsID)
}

func testAccCloudInstance_update_withEpsID(epsID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "professional"
  enterprise_project_id = "%s"

  bandwidth_expack_product {
    resource_size = 1
  }
  domain_expack_product {
    resource_size = 1
  }
  rule_expack_product {
    resource_size = 1
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"
}
`, epsID)
}
