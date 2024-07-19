package dds

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

func getDdsBackupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBackup: Query DDS backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "dds"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)

	backupId := state.Primary.ID
	getBackupQueryParams := buildGetBackupQueryParams(backupId)
	getBackupPath += getBackupQueryParams

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DDS backup: %s", err)
	}
	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return nil, err
	}
	backups := utils.PathSearch("backups", getBackupRespBody, make([]interface{}, 0)).([]interface{})
	if len(backups) == 0 {
		return nil, fmt.Errorf("error retrieving DDS backup by backup ID: %s", backupId)
	}

	return backups[0], nil
}

func buildGetBackupQueryParams(backupId string) string {
	return fmt.Sprintf("?backup_id=%s", backupId)
}

func TestAccDdsBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdsBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdsBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "this is a test dds instance"),
					resource.TestCheckResourceAttr(rName, "type", "Manual"),
					resource.TestCheckResourceAttr(rName, "status", "COMPLETED"),
					resource.TestCheckResourceAttr(rName, "datastore.0.type", "DDS-Community"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "instance_name",
						"${huaweicloud_dds_instance.instance.name}"),
					acceptance.TestCheckResourceAttrWithVariable(rName, "datastore.0.version",
						"${huaweicloud_dds_instance.instance.datastore.0.version}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccBackupImportStateFunc(rName),
			},
		},
	})
}

func testDdsBackup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_backup" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
  name        = "%s"
  description = "this is a test dds instance"
}
`, testAccDDSInstanceV3Config_basic(name, 8800), name)
}

func testAccBackupImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
