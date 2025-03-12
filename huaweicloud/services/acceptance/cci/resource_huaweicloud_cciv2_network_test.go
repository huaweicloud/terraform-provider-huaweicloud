package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getV2NetworkResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("cci", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI Yangtse v2 client: %s", err)
	}
	return cci.GetNetwork(client, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2Network_basic(t *testing.T) {
	var network networks.Network
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_network.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&network,
		getV2NetworkResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Network_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "namespace", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckOutput("is_warm_pool_size_pass", "true"),
					resource.TestCheckOutput("is_warm_pool_recycle_interval_pass", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "subnets.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "api_version"),
					resource.TestCheckResourceAttrSet(resourceName, "kind"),
					resource.TestCheckResourceAttrSet(resourceName, "creation_timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "finalizers.#"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_version"),
					resource.TestCheckResourceAttrSet(resourceName, "uid"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.conditions.#"),
					resource.TestCheckResourceAttrSet(resourceName, "status.0.subnet_attrs.#"),
					resource.TestCheckResourceAttrPair(resourceName, "status.0.subnet_attrs.0.network_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "status.0.subnet_attrs.0.subnet_v4_id",
						"huaweicloud_vpc_subnet.test", "subnet_id"),
				),
			},
			{
				Config: testAccV2Network_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckOutput("is_warm_pool_size_pass", "true"),
					resource.TestCheckOutput("is_warm_pool_recycle_interval_pass", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0",
						"huaweicloud_networking_secgroup.test1", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2NetworkImportStateFunc(resourceName),
			},
		},
	})
}

func testAccV2NetworkImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["namespace"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("the namespace (%s) or name(%s) or ID (%s) is nil",
				rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["namespace"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccV2Network_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cciv2_network" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[3]s"

  annotations = {
    "yangtse.io/project-id"                 = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/project-id"],
    "yangtse.io/domain-id"                  = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/domain-id"],
    "yangtse.io/warm-pool-size"             = "10",
    "yangtse.io/warm-pool-recycle-interval" = "2",
  }
  
  subnets {
    subnet_id = huaweicloud_vpc_subnet.test.subnet_id
  }

  security_group_ids = [huaweicloud_networking_secgroup.test.id]
}

output "is_warm_pool_size_pass" {
  value = huaweicloud_cciv2_network.test.annotations["yangtse.io/warm-pool-size"] == "10"
}

output "is_warm_pool_recycle_interval_pass" {
  value = huaweicloud_cciv2_network.test.annotations["yangtse.io/warm-pool-recycle-interval"] == "2"
}
`, common.TestBaseNetwork(rName), testAccV2Namespace_basic(rName), rName)
}

func testAccV2Network_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_networking_secgroup" "test1" {
  name                 = "%[3]s_update"
  delete_default_rules = true
}

resource "huaweicloud_cciv2_network" "test" {
  namespace = huaweicloud_cciv2_namespace.test.name
  name      = "%[3]s"

  annotations = {
    "yangtse.io/project-id"                 = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/project-id"],
    "yangtse.io/domain-id"                  = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/domain-id"],
    "yangtse.io/warm-pool-size"             = "8",
    "yangtse.io/warm-pool-recycle-interval" = "3",
  }
  
  subnets {
    subnet_id = huaweicloud_vpc_subnet.test.subnet_id
  }

  security_group_ids = [huaweicloud_networking_secgroup.test1.id]
}

output "is_warm_pool_size_pass" {
  value = huaweicloud_cciv2_network.test.annotations["yangtse.io/warm-pool-size"] == "8"
}

output "is_warm_pool_recycle_interval_pass" {
  value = huaweicloud_cciv2_network.test.annotations["yangtse.io/warm-pool-recycle-interval"] == "3"
}
`, common.TestBaseNetwork(rName), testAccV2Namespace_basic(rName), rName)
}
