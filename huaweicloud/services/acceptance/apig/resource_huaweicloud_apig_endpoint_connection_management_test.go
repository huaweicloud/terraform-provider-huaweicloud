package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getEndpointConnectionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetEndpointConntionByEndpointId(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccEndpointConnectionManagement_basic(t *testing.T) {
	var (
		connection      interface{}
		rName           = "huaweicloud_apig_endpoint_connection_management.test"
		name            = acceptance.RandomAccResourceName()
		nameWithNetwork = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&connection,
		getEndpointConnectionFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointConnectionManagement_basic_step1(name, nameWithNetwork),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "endpoint_id", "huaweicloud_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(rName, "action", "receive"),
					resource.TestCheckResourceAttr(rName, "status", "accepted"),
				),
			},
			{
				Config: testAccEndpointConnectionManagement_basic_step2(name, nameWithNetwork),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "action", "reject"),
					resource.TestCheckResourceAttr(rName, "status", "rejected"),
				),
			},
		},
	})
}

func testAccEndpointConnectionManagement_base(name, nameWithNetwork string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

// Only resource returns parameter 'vpcep_service_address'.
resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_vpc_subnet" "test2" {
  name       = "%[3]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.100.0/24"
  gateway_ip = "192.168.100.1"
}

# The vpcep_service_address format is "{region}.{vpcep_service_name}.{service_id}"
# The subnet of the instance and endpoint service cannot be the same.
resource "huaweicloud_vpcep_endpoint" "test" {
  service_id = element(split(".", huaweicloud_apig_instance.test.vpcep_service_address), 2)
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test2.id
}
`, common.TestBaseNetwork(name), name, nameWithNetwork)
}

func testAccEndpointConnectionManagement_basic_step1(name, nameWithNetwork string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_endpoint_connection_management" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "receive"
  endpoint_id = huaweicloud_vpcep_endpoint.test.id
}
`, testAccEndpointConnectionManagement_base(name, nameWithNetwork))
}

func testAccEndpointConnectionManagement_basic_step2(name, nameWithNetwork string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_endpoint_connection_management" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "reject"
  endpoint_id = huaweicloud_vpcep_endpoint.test.id
}
`, testAccEndpointConnectionManagement_base(name, nameWithNetwork))
}
