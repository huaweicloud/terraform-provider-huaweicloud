package ges

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGesBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBackup: Query the GES backup.
	var (
		getBackupHttpUrl = "v2/{project_id}/graphs/{graph_id}/backups"
		getBackupProduct = "ges"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GES Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)
	getBackupPath = strings.ReplaceAll(getBackupPath, "{graph_id}", fmt.Sprintf("%v", state.Primary.Attributes["graph_id"]))
	getBackupPath += "?limit=120"

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GesBackup: %s", err)
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GesBackup: %s", err)
	}

	jsonPath := fmt.Sprintf("backup_list[?id =='%s']|[0]", state.Primary.ID)
	getBackupRespBody = utils.PathSearch(jsonPath, getBackupRespBody, nil)
	if getBackupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getBackupRespBody, nil
}

func TestAccGesBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ges_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGesBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGesBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "graph_id", "huaweicloud_ges_graph.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "backup_method"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "start_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "size"),
					resource.TestCheckResourceAttrSet(rName, "duration"),
					resource.TestCheckResourceAttrSet(rName, "encrypted"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testGesBackupImportState(rName),
			},
		},
	})
}

func testGesBackup_basic(name string) string {
	graphConf := testGesGraph_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_ges_backup" "test" {
  graph_id = huaweicloud_ges_graph.test.id
}
`, graphConf)
}

func testGesBackupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["graph_id"] == "" {
			return "", fmt.Errorf("attribute (graph_id) of resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["graph_id"] + "/" +
			rs.Primary.ID, nil
	}
}
