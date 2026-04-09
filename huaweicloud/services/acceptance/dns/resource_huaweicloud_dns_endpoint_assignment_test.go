package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getEndpointAssignmentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	if state.Primary.Attributes["regional"] == "true" {
		conf.RegionClient = true
	}

	client, err := conf.NewServiceClient("dns", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating DNS client: %s", err)
	}

	return dns.GetEntpointById(client, state.Primary.ID)
}

func TestAccEndpointAssignment_basic(t *testing.T) {
	var (
		endpoint   interface{}
		rName      = "huaweicloud_dns_endpoint_assignment.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&endpoint,
		getEndpointAssignmentResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testEndpointAssignment_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "direction", "inbound"),
					resource.TestCheckResourceAttr(rName, "assignments.#", "5"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.ip_address_id"),
					resource.TestCheckResourceAttrSet(rName, "vpc_id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testEndpointAssignment_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "assignments.#", "6"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
			},
		},
	})
}

func testEndpointAssignment_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  count      = 2
  name       = "%[1]s${count.index}"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, count.index)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, count.index), 1)
}`, rName)
}

func testEndpointAssignment_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_endpoint_assignment" "test" {
  name      = "%[2]s"
  direction = "inbound"

  dynamic "assignments" {
    for_each = range(4)
    content {
      subnet_id  = huaweicloud_vpc_subnet.test[0].id
      ip_address = cidrhost(huaweicloud_vpc_subnet.test[0].cidr, assignments.key + 100)
    }
  }

  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test[1].id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test[1].cidr, 100)
  }
}`, testEndpointAssignment_base(rName), rName)
}

func testEndpointAssignment_basic_step2(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_endpoint_assignment" "test" {
  name      = "%[2]s"
  direction = "inbound"

  dynamic "assignments" {
    for_each = range(4)
    content {
      subnet_id  = huaweicloud_vpc_subnet.test[0].id
      ip_address = cidrhost(huaweicloud_vpc_subnet.test[0].cidr, assignments.key + 100)
    }
  }

  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test[1].id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test[1].cidr, 102)
  }
  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test[1].id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test[1].cidr, 103)
  }
}`, testEndpointAssignment_base(rName), updateName)
}

// Check that some regions require the region parameter to be specified, such as 'sa-brazil-1'.
func TestAccEndpointAssignment_regional(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dns_endpoint_assignment.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getEndpointAssignmentResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCustomRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointAssignment_regional_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "direction", "inbound"),
					resource.TestCheckResourceAttr(rName, "assignments.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "assignments.0.ip_address_id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr(rName, "regional", "true"),
				),
			},
			{
				Config: testAccEndpointAssignment_regional_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "assignments.#", "2"),
					resource.TestCheckResourceAttr(rName, "regional", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				// Only ignore regional in acceptance test, it will not be changed in actual use.
				ImportStateVerifyIgnore: []string{"regional"},
				ImportStateIdFunc:       testAccEndpointAssignmentImportStateFunc(rName),
			},
		},
	})
}

func testAccEndpointAssignmentImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		region := rs.Primary.Attributes["region"]
		endpointId := rs.Primary.ID
		if region == "" || endpointId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<region>/<id>', but got '%s/%s'", region, endpointId)
		}

		return fmt.Sprintf("%s/%s", region, endpointId), nil
	}
}

func testEndpointAssignment_regional_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  region = "%[1]s"
  name   = "%[2]s"
  cidr   = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  region     = "%[1]s"
  name       = "%[2]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}`, acceptance.HW_CUSTOM_REGION_NAME, rName)
}

func testAccEndpointAssignment_regional_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_endpoint_assignment" "test" {
  region    = "%[2]s"
  name      = "%[3]s"
  direction = "inbound"

  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test.id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test.cidr, 100)
  }
  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test.id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test.cidr, 101)
  }
}`, testEndpointAssignment_regional_base(rName), acceptance.HW_CUSTOM_REGION_NAME, rName)
}

func testAccEndpointAssignment_regional_step2(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_endpoint_assignment" "test" {
  region    = "%[2]s"
  name      = "%[3]s"
  direction = "inbound"

  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test.id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test.cidr, 102)
  }
  assignments {
    subnet_id  = huaweicloud_vpc_subnet.test.id
    ip_address = cidrhost(huaweicloud_vpc_subnet.test.cidr, 103)
  }
}`, testEndpointAssignment_regional_base(rName), acceptance.HW_CUSTOM_REGION_NAME, updateName)
}
