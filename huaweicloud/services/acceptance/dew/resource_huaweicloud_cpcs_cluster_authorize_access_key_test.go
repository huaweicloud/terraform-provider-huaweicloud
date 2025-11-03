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

func getCpcsClusterAuthorizeAccessKeyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("kms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DEW client: %s", err)
	}

	clusterId := state.Primary.Attributes["cluster_id"]
	accessKeyId := state.Primary.Attributes["access_key_id"]
	return dew.QueryCpcsClusterAuthorizeAccessKey(client, clusterId, accessKeyId)
}

// Currently, this resource is valid only in cn-north-9 region.
func TestAccCpcsClusterAuthorizeAccessKey_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_cpcs_cluster_authorize_access_key.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCpcsClusterAuthorizeAccessKeyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid cluster ID, app ID, and access key ID in the environment variables.
			acceptance.TestAccPrecheckCpcsClusterId(t)
			acceptance.TestAccPrecheckCpcsAppId(t)
			acceptance.TestAccPrecheckCpcsAccessKeyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCpcsClusterAuthorizeAccessKey_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_CPCS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "app_id", acceptance.HW_CPCS_APP_ID),
					resource.TestCheckResourceAttr(rName, "access_key_id", acceptance.HW_CPCS_ACCESS_KEY_ID),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "app_name"),
					resource.TestCheckResourceAttrSet(rName, "key_name"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: testCpcsClusterAuthorizeAccessKeyImportState(rName),
				ImportStateVerifyIgnore: []string{
					"app_id",
				},
			},
		},
	})
}

func testCpcsClusterAuthorizeAccessKey_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cpcs_cluster_authorize_access_key" "test" {
  cluster_id    = "%s"
  app_id        = "%s"
  access_key_id = "%s"
}
`, acceptance.HW_CPCS_CLUSTER_ID, acceptance.HW_CPCS_APP_ID, acceptance.HW_CPCS_ACCESS_KEY_ID)
}

func testCpcsClusterAuthorizeAccessKeyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		if rs.Primary.Attributes["cluster_id"] == "" || rs.Primary.Attributes["access_key_id"] == "" {
			return "", fmt.Errorf("required attributes (cluster_id, access_key_id) of resource (%s) not found", name)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["access_key_id"]), nil
	}
}
