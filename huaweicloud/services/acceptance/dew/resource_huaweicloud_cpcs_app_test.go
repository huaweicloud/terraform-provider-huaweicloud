package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
)

func getCpcsAppResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("kms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DEW Client: %s", err)
	}

	return dew.QueryCpcsAppByAppName(client, state.Primary.Attributes["app_name"])
}

// Currently, this resource is valid only in cn-north-9 region.
func TestAccCpcsApp_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_cpcs_app.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCpcsAppResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCpcsApp_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", name),
					resource.TestCheckResourceAttr(rName, "description", "test application"),
					resource.TestCheckResourceAttrSet(rName, "vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "vpc_name"),
					resource.TestCheckResourceAttrSet(rName, "subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "subnet_name"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCpcsAppImportState(rName),
			},
		},
	})
}

func testCpcsApp_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpcs" "test" {}
data "huaweicloud_vpc_subnets" "test" {
  vpc_id = data.huaweicloud_vpcs.test.vpcs[0].id
}

resource "huaweicloud_cpcs_app" "test" {
  app_name    = "%[2]s"
  vpc_id      = data.huaweicloud_vpcs.test.vpcs[0].id
  vpc_name    = data.huaweicloud_vpcs.test.vpcs[0].name
  subnet_id   = data.huaweicloud_vpc_subnets.test.subnets[0].id
  subnet_name = data.huaweicloud_vpc_subnets.test.subnets[0].name
  description = "test application"
}
`, common.TestVpc(name), name)
}

func testCpcsAppImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.Attributes["app_name"] == "" {
			return "", fmt.Errorf("attribute (app_name) of resource (%s) not found", name)
		}

		return rs.Primary.Attributes["app_name"], nil
	}
}
