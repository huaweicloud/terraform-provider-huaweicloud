package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dew"
)

func getCpcsAppClusterAssociationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("kms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DEW client: %s", err)
	}

	appId := state.Primary.Attributes["app_id"]
	clusterId := state.Primary.Attributes["cluster_id"]
	return dew.QueryCpcsAppClusterAssociation(client, appId, clusterId)
}

// Currently, this resource is valid only in cn-north-9 region.
func TestAccCpcsAppClusterAssociation_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cpcs_app_cluster_association.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCpcsAppClusterAssociationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an unbound App and Cluster, and configure the ID in the environment variables.
			acceptance.TestAccPrecheckCpcsAppClusterAssociation(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCpcsAppClusterAssociation_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_id", acceptance.HW_CPCS_APP_ID),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CPCS_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(rName, "app_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_name"),
					resource.TestCheckResourceAttrSet(rName, "vpc_name"),
					resource.TestCheckResourceAttrSet(rName, "subnet_name"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testCpcsAppClusterAssociationImportState(rName),
			},
		},
	})
}

func testCpcsAppClusterAssociation_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cpcs_app_cluster_association" "test" {
  app_id     = "%s"
  cluster_id = "%s"
}
`, acceptance.HW_CPCS_APP_ID, acceptance.HW_CPCS_CLUSTER_ID)
}

func testCpcsAppClusterAssociationImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.Attributes["app_id"] == "" || rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("required attributes (app_id, cluster_id) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["app_id"], rs.Primary.Attributes["cluster_id"]), nil
	}
}
