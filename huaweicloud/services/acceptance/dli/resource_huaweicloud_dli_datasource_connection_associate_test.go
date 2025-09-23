package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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

		elasticResourceName, updateElasticResourceRoolName = getElasticResourcePoolNames()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliElasticResourcePoolName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAcDatasourceConnectionAssociatec_basic_step1(name, elasticResourceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id",
						"huaweicloud_dli_datasource_connection.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.0", elasticResourceName),
				),
			},
			{
				Config: testAcDatasourceConnectionAssociatec_basic_step2(name, updateElasticResourceRoolName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "elastic_resource_pools.0", updateElasticResourceRoolName),
					// After elastic resource pool is created, it cannot be deleted within one hour.
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
%s

resource "huaweicloud_dli_datasource_connection" "test" {
  name      = "%s"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}
`, common.TestVpc(name), name)
}

func testAcDatasourceConnectionAssociatec_basic_step1(name, elasticResourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = ["%s"]
}
`, testAcDatasourceConnectionAssociatec_base(name), elasticResourceName)
}

func testAcDatasourceConnectionAssociatec_basic_step2(name, elasticResourceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_datasource_connection_associate" "test" {
  connection_id          = huaweicloud_dli_datasource_connection.test.id
  elastic_resource_pools = ["%s"]
}
`, testAcDatasourceConnectionAssociatec_base(name), elasticResourceName)
}
