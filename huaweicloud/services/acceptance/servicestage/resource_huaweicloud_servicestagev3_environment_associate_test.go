package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/servicestage"
)

func getV3EnvAssociatedResourcesFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("servicestage", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage client: %s", err)
	}
	return servicestage.ListV3EnvironmentAssociatedResources(client, state.Primary.ID)
}

func TestAccV3EnvironmentAssociate_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_servicestagev3_environment_associate.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3EnvAssociatedResourcesFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccV3EnvironmentAssociate_nonExistEnvironmentAndResources(),
				ExpectError: regexp.MustCompile(`error associating resources to the environment \([a-f0-9-]+\)`),
			},
			{
				Config: testAccV3EnvironmentAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.id", "huaweicloud_vpc_eip.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.name", ""),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.id", "huaweicloud_vpc_eip.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.1.type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "resources.1.name", ""),
					resource.TestCheckResourceAttrPair(resourceName, "resources.2.id", "huaweicloud_vpc_eip.test.2", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.2.type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "resources.2.name", ""),
				),
			},
			{
				Config: testAccV3EnvironmentAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.id", "huaweicloud_cce_cluster.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.type", "cce"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.name", "huaweicloud_cce_cluster.test.0", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.id", "huaweicloud_cce_cluster.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.1.type", "cce"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.name", "huaweicloud_cce_cluster.test.1", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.2.id", "huaweicloud_cce_cluster.test.2", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.2.type", "cce"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.2.name", "huaweicloud_cce_cluster.test.2", "name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources_origin",
				},
			},
		},
	})
}

func testAccV3EnvironmentAssociate_nonExistEnvironmentAndResources() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = "%[1]s"

  resources {
	id   = "%[1]s"
	type = "eip"
  }
}
`, randomUUID.String())
}

func testAccV3EnvironmentAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  count = 3

  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = format("%[2]s-%%d", count.index)
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_cce_cluster" "test" {
  count = 3

  name                   = format("%[2]s-%%d", count.index)
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
  container_network_cidr = "172.16.0.0/24"
  security_group_id      = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, common.TestBaseNetwork(name), name)
}

func testAccV3EnvironmentAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  dynamic "resources" {
    for_each = huaweicloud_vpc_eip.test

	# Associate all EIPs to the environment and without configuring the EIP name.
    content {
      id   = resources.value.id
      type = "eip"
    }
  }
}
`, testAccV3EnvironmentAssociate_base(name))
}

func testAccV3EnvironmentAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  dynamic "resources" {
    for_each = huaweicloud_cce_cluster.test

	# Associate all CCE clusters to the environment and configure the CCE cluster name.
    content {
      id   = resources.value.id
      name = resources.value.name
      type = "cce"
    }
  }
}
`, testAccV3EnvironmentAssociate_base(name))
}
