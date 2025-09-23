package dcs

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

func getDcsBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBackup: Query DCS backup
	var (
		getBackupHttpUrl = "v2/{project_id}/instances/{instance_id}/backups"
		getBackupProduct = "dcs"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<backup_id>")
	}
	instanceID := parts[0]
	backupID := parts[1]
	getBackupBasePath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupBasePath = strings.ReplaceAll(getBackupBasePath, "{project_id}", getBackupClient.ProjectID)
	getBackupBasePath = strings.ReplaceAll(getBackupBasePath, "{instance_id}", instanceID)

	getDdmSchemasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 204},
	}

	var currentTotal int
	var getBackupPath string
	for {
		getBackupPath = getBackupBasePath + buildGetDcsBackupQueryParams(currentTotal)
		getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getDdmSchemasOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving DCS backup: %s", err)
		}
		getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
		if err != nil {
			return nil, err
		}
		backups := utils.PathSearch("backup_record_response", getBackupRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("total_num", getBackupRespBody, 0)
		backup := utils.PathSearch(fmt.Sprintf("[?backup_id=='%s']|[0]", backupID), backups, nil)
		if backup != nil {
			return backup, nil
		}
		currentTotal += len(backups)
		if currentTotal == int(total.(float64)) {
			break
		}
	}
	return nil, fmt.Errorf("error get DCS backup by backup_id (%s)", backupID)
}

func buildGetDcsBackupQueryParams(offset int) string {
	return fmt.Sprintf("?limit=10&offset=%v", offset)
}

func TestAccDcsBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDcsBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_dcs_instance.instance_1", "id"),
					resource.TestCheckResourceAttr(rName, "type", "manual"),
					resource.TestCheckResourceAttr(rName, "status", "succeed"),
					resource.TestCheckResourceAttr(rName, "description", "test DCS backup remark"),
					resource.TestCheckResourceAttr(rName, "backup_format", "rdb"),
					resource.TestCheckResourceAttrSet(rName, "backup_id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "size"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "is_support_restore"),
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

func testDcsBackup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_backup" "test" {
  instance_id   = huaweicloud_dcs_instance.instance_1.id
  description   = "test DCS backup remark"
  backup_format = "rdb"
}
`, testAccDcsV1Instance_basic(name))
}
