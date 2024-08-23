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

func getDwsSnapshotResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDwsSnapshot: Query the DWS snapshot.
	var (
		getDwsSnapshotHttpUrl = "v1.0/{project_id}/snapshots/{snapshot_id}"
		getDwsSnapshotProduct = "dws"
	)
	getDwsSnapshotClient, err := cfg.NewServiceClient(getDwsSnapshotProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	getDwsSnapshotPath := getDwsSnapshotClient.Endpoint + getDwsSnapshotHttpUrl
	getDwsSnapshotPath = strings.ReplaceAll(getDwsSnapshotPath, "{project_id}", getDwsSnapshotClient.ProjectID)
	getDwsSnapshotPath = strings.ReplaceAll(getDwsSnapshotPath, "{snapshot_id}", state.Primary.ID)

	getDwsSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		OkCodes: []int{
			200,
		},
	}
	getDwsSnapshotResp, err := getDwsSnapshotClient.Request("GET", getDwsSnapshotPath, &getDwsSnapshotOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS snapshot: %s", err)
	}
	return utils.FlattenResponse(getDwsSnapshotResp)
}

func TestAccDwsSnapshot_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dws_snapshot.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDwsSnapshotResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDwsSnapshot_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "AVAILABLE"),
					resource.TestCheckResourceAttr(rName, "type", "MANUAL"),
					resource.TestCheckResourceAttrPair(rName, "cluster_id", "huaweicloud_dws_cluster.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "started_at"),
					resource.TestCheckResourceAttrSet(rName, "finished_at"),
					resource.TestCheckResourceAttrSet(rName, "size"),
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

func testDwsSnapshot_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dws_snapshot" "test" {
  name       = "%s"
  cluster_id = huaweicloud_dws_cluster.test.id
}
`, testAccDwsCluster_basic_step1(name, 3, cluster.PublicBindTypeNotUse, "cluster123@!"), name)
}
