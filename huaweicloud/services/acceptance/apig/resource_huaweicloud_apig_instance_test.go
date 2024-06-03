package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getInstanceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	return apig.QueryInstanceDetail(client, state.Primary.ID)
}

func TestAccInstance_basic(t *testing.T) {
	var (
		instance instances.Instance

		resourceName = "huaweicloud_apig_instance.test"
		rName        = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getInstanceFunc,
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
				Config: testAccInstance_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_provider", "elb"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_service_name", "apig"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_service_address"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				Config: testAccInstance_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "18:00:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "22:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "vpcep_service_name", "new_custom_apig"),
					resource.TestCheckResourceAttrSet(resourceName, "vpcep_service_address"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstance_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "0"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"

  tags = {}
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccInstance_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup" "new" {
  name = "%[2]s_new"
}

resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.new.id
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  vpcep_service_name    = "new_custom_apig"

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "18:00:00"

  tags = {
    foo = "bar"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccInstance_egress(t *testing.T) {
	var (
		instance instances.Instance

		resourceName = "huaweicloud_apig_instance.test"
		rName        = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getInstanceFunc,
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
				Config: testAccInstance_baseConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
				),
			},
			{
				Config: testAccInstance_egress(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccInstance_egressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "egress_address"),
				),
			},
			{
				Config: testAccInstance_baseConfig(rName), // Unbind egress nat
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "bandwidth_size", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstance_baseConfig(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  maintain_begin        = "14:00:00"
  description           = "created by acc test"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccInstance_ingress(t *testing.T) {
	var (
		instance instances.Instance

		resourceName = "huaweicloud_apig_instance.test"
		rName        = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getInstanceFunc,
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
				Config: testAccInstance_baseConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "edition", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "14:00:00"),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "ingress_address", ""),
				),
			},
			{
				Config: testAccInstance_ingress(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", "bandwidth"),
					resource.TestCheckResourceAttrSet(resourceName, "ingress_address"),
				),
			},
			{
				Config: testAccInstance_ingressUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_size", "6"),
					resource.TestCheckResourceAttr(resourceName, "ingress_bandwidth_charging_mode", "traffic"),
				),
			},
			{
				Config: testAccInstance_baseConfig(rName), // Unbind ingress eip
				Check: resource.ComposeTestCheckFunc(rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_ingress_address"),
					resource.TestCheckResourceAttr(resourceName, "ingress_address", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccInstance_egress(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  bandwidth_size        = 5 # The bandwidth value must be greater than or equal to 5
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_egressUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition               = "BASIC"
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
  bandwidth_size        = 10
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_ingress(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition                         = "BASIC"
  name                            = "%[2]s"
  enterprise_project_id           = "%[3]s"
  maintain_begin                  = "14:00:00"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccInstance_ingressUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  vpc_id             = huaweicloud_vpc.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)

  edition                         = "BASIC"
  name                            = "%[2]s"
  enterprise_project_id           = "%[3]s"
  maintain_begin                  = "14:00:00"
  ingress_bandwidth_size          = 6
  ingress_bandwidth_charging_mode = "traffic"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
