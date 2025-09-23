package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDDSV3InstanceModifyIP_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance_internal_ip_modify.test"

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3ModifyIP_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "new_ip", "192.168.0.2"),
				),
			},
			{
				Config: testAccDDSInstanceV3ModifyIP_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "new_ip", "192.168.0.3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceDDSInstanceNodeImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccDDSInstanceV3ModifyIP_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_internal_ip_modify" "test" {
  depends_on = [huaweicloud_dds_instance.instance]
  
  instance_id = huaweicloud_dds_instance.instance.id
  node_id     = huaweicloud_dds_instance.instance.nodes.0.id
  new_ip      = "192.168.0.2"
}`, testAccDDSInstanceReplicaSetBasic(rName))
}

func testAccDDSInstanceV3ModifyIP_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance_internal_ip_modify" "test" {
  depends_on = [huaweicloud_dds_instance.instance]
  
  instance_id = huaweicloud_dds_instance.instance.id
  node_id     = huaweicloud_dds_instance.instance.nodes.0.id
  new_ip      = "192.168.0.3"
}`, testAccDDSInstanceReplicaSetBasic(rName))
}

func testAccDDSInstanceReplicaSetBasic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "ReplicaSet"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 1
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.repset"
  }

}`, common.TestBaseNetwork(rName), rName)
}

func testAccResourceDDSInstanceNodeImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		nodeID := rs.Primary.Attributes["node_id"]
		if instanceID == "" || nodeID == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<node_id>, but got '%s/%s'",
				instanceID, nodeID)
		}
		return fmt.Sprintf("%s/%s", instanceID, nodeID), nil
	}
}
