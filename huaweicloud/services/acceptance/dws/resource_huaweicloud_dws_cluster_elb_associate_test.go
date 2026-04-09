package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterAssociatedElbFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetClusterAssociatedElbById(client, state.Primary.ID, state.Primary.Attributes["elb_id"])
}

// Before testing, please ensure that the subnet of the DWS cluster is not the first IP address segment of 255.255.240.0.
func TestAccClusterElbAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dws_cluster_elb_associate.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getClusterAssociatedElbFunc)

		name          = acceptance.RandomAccResourceName()
		randomUUID, _ = uuid.GenerateUUID()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testClusterElbAssociate_nonExistentElb(randomUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`error associating ELB \(%s\) to the DWS cluster \(%s\)`,
					randomUUID, acceptance.HW_DWS_CLUSTER_ID)),
			},
			{
				Config: testClusterElbAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrPair(rName, "elb_id", "huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestMatchResourceAttr(rName, "private_ip", regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`)),
					resource.TestMatchResourceAttr(rName, "public_ip", regexp.MustCompile(`^\d{1,3}(\.\d{1,3}){3}$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testClusterElbAssociate_nonExistentElb(randomUUID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_elb_associate" "test" {
  cluster_id = "%[1]s"
  elb_id     = "%[2]s"
}
`, acceptance.HW_DWS_CLUSTER_ID, randomUUID)
}

func testClusterElbAssociate_basic(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_dws_clusters" "test" {}

locals {
  vpc_id = try([for v in data.huaweicloud_dws_clusters.test.clusters: v.vpc_id if v.id == "%[2]s"][0], "")
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpcs" "test" {
  id = local.vpc_id
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[3]s"
  vpc_id     = data.huaweicloud_vpcs.test.vpcs[0].id
  cidr       = cidrsubnet(data.huaweicloud_vpcs.test.vpcs[0].cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(data.huaweicloud_vpcs.test.vpcs[0].cidr, 4, 0), 1)
}

resource "huaweicloud_vpc_eip" "test" {
  name                  = "%[3]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[3]s"
    size        = 1
    charge_mode = "traffic"
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                  = "%[3]s"
  vpc_id                = data.huaweicloud_vpcs.test.vpcs[0].id
  ipv4_subnet_id        = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  # Bind EIP to ensure that the 'public_ip' field can be correctly obtained and verified after association.
  ipv4_eip_id           = huaweicloud_vpc_eip.test.id
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  lifecycle {
    ignore_changes = [
      l4_flavor_id, l7_flavor_id
    ]
  }
}

resource "huaweicloud_dws_cluster_elb_associate" "test" {
  cluster_id = "%[2]s"
  elb_id     = huaweicloud_elb_loadbalancer.test.id
}
`, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST, acceptance.HW_DWS_CLUSTER_ID, name)
}
