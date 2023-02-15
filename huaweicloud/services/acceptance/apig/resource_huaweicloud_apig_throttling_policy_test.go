package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getThrottlingPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return throttles.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccThrottlingPolicy_basic(t *testing.T) {
	var (
		policy throttles.ThrottlingPolicy

		rName      = "huaweicloud_apig_throttling_policy.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		appCode    = acctest.RandString(64)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&policy,
		getThrottlingPolicyFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigThrottlingPolicy_basic(name, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "type", "API-based"),
					resource.TestCheckResourceAttr(rName, "period", "15000"),
					resource.TestCheckResourceAttr(rName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(rName, "max_api_requests", "100"),
					resource.TestCheckResourceAttr(rName, "max_user_requests", "60"),
					resource.TestCheckResourceAttr(rName, "max_app_requests", "60"),
					resource.TestCheckResourceAttr(rName, "max_ip_requests", "60"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_update(updateName, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(rName, "type", "API-shared"),
					resource.TestCheckResourceAttr(rName, "period", "10"),
					resource.TestCheckResourceAttr(rName, "period_unit", "MINUTE"),
					resource.TestCheckResourceAttr(rName, "max_api_requests", "70"),
					resource.TestCheckResourceAttr(rName, "max_user_requests", "45"),
					resource.TestCheckResourceAttr(rName, "max_app_requests", "45"),
					resource.TestCheckResourceAttr(rName, "max_ip_requests", "45"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccThrottlingPolicyImportStateFunc(),
			},
		},
	})
}

func TestAccThrottlingPolicy_spec(t *testing.T) {
	var (
		policy throttles.ThrottlingPolicy

		rName   = "huaweicloud_apig_throttling_policy.test"
		name    = acceptance.RandomAccResourceName()
		appCode = acctest.RandString(64)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&policy,
		getThrottlingPolicyFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigThrottlingPolicy_basic(name, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "API-based"),
					resource.TestCheckResourceAttr(rName, "period", "15000"),
					resource.TestCheckResourceAttr(rName, "period_unit", "SECOND"),
					resource.TestCheckResourceAttr(rName, "max_api_requests", "100"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_spec(name, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_throttles.#", "1"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_specUpdate(name, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_throttles.#", "1"),
				),
			},
			{
				Config: testAccApigThrottlingPolicy_basic(name, appCode),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_throttles.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccThrottlingPolicyImportStateFunc(),
			},
		},
	})
}

func testAccThrottlingPolicyImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_throttling_policy.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApigThrottlingPolicy_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[1]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
  ]
}
`, name)
}

func testAccApigThrottlingPolicy_basic(name, appCode string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application" "test" {
  name        = "%[3]s"
  instance_id = huaweicloud_apig_instance.test.id

  app_codes = [
    base64encode("%[2]s"),
  ]
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%[3]s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100
  max_user_requests = 60
  max_app_requests  = 60
  max_ip_requests   = 60
  description       = "Created by script"
}
`, testAccApigThrottlingPolicy_base(name), appCode, name)
}

func testAccApigThrottlingPolicy_update(name, appCode string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_application" "test" {
  name        = "%[3]s"
  instance_id = huaweicloud_apig_instance.test.id

  app_codes = [
    base64encode("%[2]s"),
  ]
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%[3]s"
  type              = "API-shared"
  period            = 10
  period_unit       = "MINUTE"
  max_api_requests  = 70
  max_user_requests = 45
  max_app_requests  = 45
  max_ip_requests   = 45
  description       = "Updated by script"
}
`, testAccApigThrottlingPolicy_base(name), appCode, name)
}

func testAccApigThrottlingPolicy_spec(name, appCode string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  randCodes = [
    base64encode("%[2]s"), base64encode(strrev("%[2]s")),
  ]
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  name        = "%[3]s_${count.index}"
  instance_id = huaweicloud_apig_instance.test.id
  app_codes   = slice(local.randCodes, count.index, count.index+1)
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%[3]s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100

  dynamic "app_throttles" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 0, 1)

    content {
      max_api_requests     = 30
      throttling_object_id = app_throttles.value
	}
  }
}
`, testAccApigThrottlingPolicy_base(name), appCode, name)
}

func testAccApigThrottlingPolicy_specUpdate(name, appCode string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  randCodes = [
    "%[2]s", strrev("%[2]s"),
  ]
}

resource "huaweicloud_apig_application" "test" {
  count = 2

  name        = "%[3]s_${count.index}"
  instance_id = huaweicloud_apig_instance.test.id
  app_codes   = slice(local.randCodes, count.index, count.index+1)
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id       = huaweicloud_apig_instance.test.id
  name              = "%[3]s"
  type              = "API-based"
  period            = 15000
  period_unit       = "SECOND"
  max_api_requests  = 100

  dynamic "app_throttles" {
    for_each = slice(huaweicloud_apig_application.test[*].id, 0, 1)

    content {
      max_api_requests     = 30
      throttling_object_id = app_throttles.value
	}
  }
}
`, testAccApigThrottlingPolicy_base(name), appCode, name)
}
