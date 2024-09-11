package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCEndpoint_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_vpcep_endpoint.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpoint_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
				),
			},
			{
				Config: testAccVPCEndpoint_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "false"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "0"),
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

func TestAccVPCEndpoint_Public(t *testing.T) {
	var endpoint endpoints.Endpoint
	resourceName := "huaweicloud_vpcep_endpoint.myendpoint"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointPublic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "routetables.#"),
				),
			},
		},
	})
}

func getVpcepEndpointResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	vpcepClient, err := conf.VPCEPClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPCEP client: %s", err)
	}

	return endpoints.Get(vpcepClient, state.Primary.ID).Extract()
}

func testAccVPCEndpoint_Precondition(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

resource "huaweicloud_compute_instance" "ecs" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccVPCEndpoint_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id       = huaweicloud_vpcep_service.test.id
  vpc_id           = data.huaweicloud_vpc.myvpc.id
  network_id       = data.huaweicloud_vpc_subnet.test.id
  enable_dns       = true
  description      = "test description"
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24"]

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

func testAccVPCEndpoint_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id       = huaweicloud_vpcep_service.test.id
  vpc_id           = data.huaweicloud_vpc.myvpc.id
  network_id       = data.huaweicloud_vpc_subnet.test.id
  enable_dns       = true
  description      = "test description"
  enable_whitelist = false

  tags = {
    owner = "tf-acc-update"
    foo   = "bar"
  }
}
`, testAccVPCEndpoint_Precondition(rName))
}

var testAccVPCEndpointPublic = `
data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "mynet" {
  vpc_id = data.huaweicloud_vpc.myvpc.id
  name   = "subnet-default"
}

data "huaweicloud_vpcep_public_services" "cloud_service" {
  service_name = "dis"
}

resource "huaweicloud_vpcep_endpoint" "myendpoint" {
  service_id       = data.huaweicloud_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = data.huaweicloud_vpc.myvpc.id
  network_id       = data.huaweicloud_vpc_subnet.mynet.id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.10"]
}
`

func TestAccVPCEndpoint_gatewayEndpoint(t *testing.T) {
	var (
		endpoint     endpoints.Endpoint
		rName        = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_vpcep_endpoint.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare a gateway VPC endpoint service ID in advance.
			acceptance.TestAccPreCheckVPCEPServiceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpoint_gatewayEndpoint(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "routetables.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "routetables.0", "data.huaweicloud_vpc_route_table.custom", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "service_type"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_statement"),
				),
			},
			{
				Config: testAccVPCEndpoint_gatewayEndpointUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "routetables.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "routetables.0", "data.huaweicloud_vpc_route_table.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_id"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "service_type"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_statement"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_dns"},
			},
		},
	})
}

func testAccVPCEndpoint_RouteTables_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  name              = "%[1]s"
  cidr              = "192.168.1.0/24"
  gateway_ip        = "192.168.1.1"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_vpc_subnet" "retest" {
  depends_on = [
    huaweicloud_vpc_subnet.test
  ]

  vpc_id            = huaweicloud_vpc.test.id
  name              = "%[1]s-rt"
  cidr              = "192.168.5.0/24"
  gateway_ip        = "192.168.5.1"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

data "huaweicloud_vpc_subnet_ids" "test" {
  depends_on = [
    huaweicloud_vpc_subnet.retest
  ]

  vpc_id = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_route_table" "test" {
  name    = "%[1]s-rtb"
  vpc_id  = huaweicloud_vpc.test.id
  subnets = data.huaweicloud_vpc_subnet_ids.test.ids
}

data "huaweicloud_vpc_route_table" "test" {
  vpc_id = huaweicloud_vpc.test.id
}

data "huaweicloud_vpc_route_table" "custom" {
  vpc_id = huaweicloud_vpc.test.id
  name   = huaweicloud_vpc_route_table.test.name
}
`, rName)
}

func testAccVPCEndpoint_gatewayEndpoint(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_endpoint" "test" {
  depends_on = [
    huaweicloud_vpc_route_table.test
  ]

  service_id  = "%[2]s"
  vpc_id      = huaweicloud_vpc.test.id
  description = "created by terraform"

  routetables = [
    data.huaweicloud_vpc_route_table.custom.id,
  ]

  policy_statement = <<EOF
  [
    {
      "Effect": "Allow",
      "Action": [
        "obs:bucket:ListBucket"
      ],
      "Resource": [
        "obs:*:*:*:*/*",
        "obs:*:*:*:*"
      ]
    },
    {
      "Effect": "Deny",
      "Action": [
        "obs:object:DeleteObject"
      ],
      "Resource": [
        "obs:*:*:*:*/*",
        "obs:*:*:*:*"
      ]
    }
  ]
EOF

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpoint_RouteTables_base(rName), acceptance.HW_VPCEP_SERVICE_ID)
}

func testAccVPCEndpoint_gatewayEndpointUpdate(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_endpoint" "test" {
  depends_on = [
    huaweicloud_vpc_route_table.test
  ]

  service_id  = "%[2]s"
  vpc_id      = huaweicloud_vpc.test.id
  description = "created by terraform"

  routetables = [
    data.huaweicloud_vpc_route_table.test.id
  ]

  policy_statement = <<EOF
  [
    {
      "Effect": "Deny",
      "Action": [
        "obs:bucket:ListBucket"
      ],
      "Resource": [
        "obs:*:*:*:*/*",
        "obs:*:*:*:*"
      ]
    }
  ]
EOF

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpoint_RouteTables_base(rName), acceptance.HW_VPCEP_SERVICE_ID)
}
