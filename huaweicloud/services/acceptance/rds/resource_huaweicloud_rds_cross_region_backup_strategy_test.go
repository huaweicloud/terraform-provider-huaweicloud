package rds

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

func getBackupStrategyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBackupStrategy: Query the RDS cross region backup strategy
	var (
		getBackupStrategyHttpUrl = "v3/{project_id}/instances/{instance_id}/backups/offsite-policy"
		getBackupStrategyProduct = "rds"
	)
	getBackupStrategyClient, err := cfg.NewServiceClient(getBackupStrategyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getBackupStrategyPath := getBackupStrategyClient.Endpoint + getBackupStrategyHttpUrl
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{project_id}", getBackupStrategyClient.ProjectID)
	getBackupStrategyPath = strings.ReplaceAll(getBackupStrategyPath, "{instance_id}", state.Primary.ID)

	getBackupStrategyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBackupStrategyResp, err := getBackupStrategyClient.Request("GET", getBackupStrategyPath, &getBackupStrategyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	getBackupStrategyRespBody, err := utils.FlattenResponse(getBackupStrategyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	policyPara := utils.PathSearch("policy_para", getBackupStrategyRespBody, nil)
	if policyPara == nil {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	backupStrategies := policyPara.([]interface{})
	if len(backupStrategies) == 0 || utils.PathSearch("keep_days", backupStrategies[0], float64(0)).(float64) == 0 {
		return nil, fmt.Errorf("error retrieving RDS cross region backup strategy: %s", err)
	}

	return getBackupStrategyRespBody, nil
}

func TestAccBackupStrategy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_cross_region_backup_strategy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupStrategyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackupStrategy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "auto"),
					resource.TestCheckResourceAttr(rName, "keep_days", "5"),
					resource.TestCheckResourceAttr(rName, "destination_region", acceptance.HW_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", acceptance.HW_DEST_PROJECT_ID),
				),
			},
			{
				Config: testBackupStrategy_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "all"),
					resource.TestCheckResourceAttr(rName, "keep_days", "8"),
					resource.TestCheckResourceAttr(rName, "destination_region", acceptance.HW_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", acceptance.HW_DEST_PROJECT_ID),
				),
			},
			{
				Config: testBackupStrategy_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "backup_type", "auto"),
					resource.TestCheckResourceAttr(rName, "keep_days", "10"),
					resource.TestCheckResourceAttr(rName, "destination_region", acceptance.HW_DEST_REGION),
					resource.TestCheckResourceAttr(rName, "destination_project_id", acceptance.HW_DEST_PROJECT_ID),
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

func testBackupStrategy_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_cross_region_backup_strategy" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  backup_type            = "auto"
  keep_days              = "5"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsInstance_mysql_step1(name), acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}

func testBackupStrategy_basic_update1(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_cross_region_backup_strategy" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  backup_type            = "all"
  keep_days              = "8"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsInstance_mysql_step1(name), acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}

func testBackupStrategy_basic_update2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_cross_region_backup_strategy" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  backup_type            = "auto"
  keep_days              = "10"
  destination_region     = "%s"
  destination_project_id = "%s"
}
`, testAccRdsInstance_mysql_step1(name), acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}
