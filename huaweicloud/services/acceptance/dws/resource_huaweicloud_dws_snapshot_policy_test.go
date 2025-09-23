package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dws/v1/cluster"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDwsSnapshotPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsSnapshotPolicy: Query the DWS snapshot policy.
	var (
		getDwsSnapshotPolicyHttpUrl = "v2/{project_id}/clusters/{cluster_id}/snapshot-policies"
		getDwsSnapshotPolicyProduct = "dws"
	)
	getDwsSnapshotPolicyClient, err := cfg.NewServiceClient(getDwsSnapshotPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsSnapshotPolicyPath := getDwsSnapshotPolicyClient.Endpoint + getDwsSnapshotPolicyHttpUrl
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{project_id}", getDwsSnapshotPolicyClient.ProjectID)
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{cluster_id}",
		fmt.Sprintf("%v", state.Primary.Attributes["cluster_id"]))

	getDwsSnapshotPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		OkCodes: []int{
			200,
		},
	}
	getDwsSnapshotPolicyResp, err := getDwsSnapshotPolicyClient.Request("GET", getDwsSnapshotPolicyPath, &getDwsSnapshotPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsSnapshotPolicy: %s", err)
	}

	getDwsSnapshotPolicyRespBody, err := utils.FlattenResponse(getDwsSnapshotPolicyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DwsSnapshotPolicy: %s", err)
	}

	jsonPath := fmt.Sprintf("backup_strategies[?policy_id=='%s']|[0]", state.Primary.ID)
	rawData := utils.PathSearch(jsonPath, getDwsSnapshotPolicyRespBody, nil)
	if rawData == nil {
		return nil, fmt.Errorf("error retrieving DwsSnapshotPolicy: %s", err)
	}

	return rawData, nil
}

func TestAccDwsSnapshotPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_snapshot_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsSnapshotPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsSnapshotPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "full"),
					resource.TestCheckResourceAttr(rName, "strategy", "0 8 6 4 * ?"),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_dws_cluster.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDwsSnapshotPolicyImportState(rName),
			},
		},
	})
}

func testDwsSnapshotPolicy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_snapshot_policy" "test" {
  name       = "%s"
  cluster_id = huaweicloud_dws_cluster.test.id
  type       = "full"
  strategy   = "0 8 6 4 * ?"
}
`, testAccDwsCluster_basic_step1(name, 3, cluster.PublicBindTypeNotUse, "cluster123@!"), name)
}

func testDwsSnapshotPolicyImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["cluster_id"] == "" {
			return "", fmt.Errorf("Attribute (cluster_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["cluster_id"] + "/" +
			rs.Primary.ID, nil
	}
}
