package servicestage

import (
	"fmt"
	"testing"

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
	return servicestage.QueryV3Environment(client, state.Primary.ID)
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
			acceptance.TestAccPreCheckEpsID(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3EnvironmentAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.id", "huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.type", "eip"),
				),
			},
			{
				Config: testAccV3EnvironmentAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.id", acceptance.HW_CCE_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "resources.0.type", "cce"),
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

func testAccV3EnvironmentAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "0"
}
`, common.TestVpc(name), name)
}

func testAccV3EnvironmentAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = huaweicloud_vpc_eip.test.id
    type = "eip"
  }
}
`, testAccV3EnvironmentAssociate_base(name))
}

func testAccV3EnvironmentAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = "%[2]s"
    type = "cce"
  }
}
`, testAccV3EnvironmentAssociate_base(name), acceptance.HW_CCE_CLUSTER_ID)
}
