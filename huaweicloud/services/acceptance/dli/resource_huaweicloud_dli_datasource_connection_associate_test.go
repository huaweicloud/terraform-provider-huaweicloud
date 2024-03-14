package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getAssociatedElasticResourcePoolsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dli", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI client: %s", err)
	}
	return dli.GetDatasourceConnectionAssociatedPoolNames(client, state.Primary.ID)
}

func TestAccDatasourceConnectionAssociate_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dli_datasource_connection_associate.test"
		name         = acceptance.RandomAccResourceName()
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getAssociatedElasticResourcePoolsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAcDatasourceConnectionAssociatec_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id",
						"huaweicloud_dli_datasource_connection.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.#", "2"),
				),
			},
			{
				Config: testAcDatasourceConnectionAssociatec_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.#", "2"),
					// After elastic resource pool is created, it cannot be deleted within one hour.
					waitForDeletionCooldownComplete(),
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

func testAcDatasourceConnectionAssociatec_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 3, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 3, 1), 1)
}

resource "huaweicloud_dli_elastic_resource_pool" "test" {
  count = 3

  name                  = format("%[1]s_%%d", count.index)
  min_cu                = 64
  max_cu                = 64
  enterprise_project_id = "0"

  # The minimum CIDR range allowed for elastic resource pools is /19.
  cidr = cidrsubnet(huaweicloud_vpc.test.cidr, 3, count.index+2)  # 192.168.x.0/19
}

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = "%[1]s"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}
`, name)
}

func testAcDatasourceConnectionAssociatec_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = slice(huaweicloud_dli_elastic_resource_pool.test[*].name, 0, 2)
}
`, testAcDatasourceConnectionAssociatec_base(name), name)
}

func testAcDatasourceConnectionAssociatec_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = slice(huaweicloud_dli_elastic_resource_pool.test[*].name, 1, 3)
}
`, testAcDatasourceConnectionAssociatec_base(name), name)
}
