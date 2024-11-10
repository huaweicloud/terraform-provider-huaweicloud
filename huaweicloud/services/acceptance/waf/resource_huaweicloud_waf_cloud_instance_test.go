package waf

import (
	"fmt"
	"regexp"
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
	epsId := state.Primary.Attributes["enterprise_project_id"]
	instance, _, err := waf.QueryCloudInstance(client, state.Primary.ID, epsId)
	return instance, err
}

// Only one WAF cloud mode instance can be created in each region.
func TestAccCloudInstance_prepaid_basic(t *testing.T) {
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
			// The conditions for creating a cloud mode WAF are a bit strict, with a separate environment variable
			// configured to control whether to skip test cases
			acceptance.TestAccPreCheckWafCloudInstance(t)
			// Configure two enterprise projects to test enterprise project migration.
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstance_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeIntroduction)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.#", "0"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
				),
			},
			{
				Config: testAccCloudInstance_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "resource_spec_code", string(waf.SpecCodeStandard)),
					resource.TestCheckResourceAttr(rName, "bandwidth_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "domain_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "rule_expack_product.0.resource_size", "1"),
					resource.TestCheckResourceAttr(rName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
					"period_unit",
					"period",
					"auto_renew",
				},
				ImportStateIdFunc: testWAFResourceImportState(rName),
			},
		},
	})
}

func testAccCloudInstance_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "detection"
  enterprise_project_id = "%s"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCloudInstance_update() string {
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
`, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST)
}

// Currently, this test case only supported on HuaweiCloud International Station.
// Only one WAF cloud mode instance can be created in each region.
func TestAccCloudInstance_postpaid_basic(t *testing.T) {
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
			// The conditions for creating a cloud mode WAF are a bit strict, with a separate environment variable
			// configured to control whether to skip test cases
			acceptance.TestAccPreCheckWafCloudInstance(t)
			// This test case can only be tested at international state, so a separate switch is configured.
			acceptance.TestAccPreCheckWafInternationalInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudInstance_postPaid_validation(),
				ExpectError: regexp.MustCompile("`website` must be specified in postpaid charging mode"),
			},
			{
				Config: testAccCloudInstance_postpaid_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enterprise_project_id",
					"website",
				},
				ImportStateIdFunc: testWAFResourceImportState(rName),
			},
		},
	})
}

func testAccCloudInstance_postPaid_validation() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  charging_mode         = "postPaid"
  enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccCloudInstance_postpaid_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  charging_mode         = "postPaid"
  website               = "hec-hk"
  enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
