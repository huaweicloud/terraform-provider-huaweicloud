package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getBackupResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBackup: Query the RDS manual backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "rds"
	)
	getBackupClient, err := config.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Backup Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.Replace(getBackupPath, "{project_id}", getBackupClient.ProjectID, -1)

	getBackupqueryParams := fmt.Sprintf("?instance_id=%s&backup_id=%s",
		state.Primary.Attributes["instance_id"], state.Primary.ID)
	getBackupPath = getBackupPath + getBackupqueryParams
	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	count := utils.PathSearch("total_count", getBackupRespBody, 0)
	if fmt.Sprintf("%v", count) == "0" {
		return nil, fmt.Errorf("error retrieving Backup: %s", err)
	}

	return getBackupRespBody, nil
}

func TestAccBackup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_backup.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBackup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "begin_time"),
					resource.TestCheckResourceAttrSet(rName, "end_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "size"),
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

func testBackup_basic(name string) string {
	rdsInstanceConfig := testAccBackupRdsInstance(name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_backup" "test" {
  name        = "%s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, rdsInstanceConfig, name)
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testAccBackupRdsInstance(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 0
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccBackupImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
