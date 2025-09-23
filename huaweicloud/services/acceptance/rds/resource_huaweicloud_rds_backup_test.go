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

func TestAccBackup_mysql_basic(t *testing.T) {
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
				Config: testBackup_mysql_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
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

func TestAccBackup_sqlserver_basic(t *testing.T) {
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
				Config: testBackup_sqlserver_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
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

func TestAccBackup_pg_basic(t *testing.T) {
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
				Config: testBackup_pg_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
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

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_mysql_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "huaweicloud_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_base(), name)
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_sqlserver_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  depends_on        = [huaweicloud_networking_secgroup_rule.ingress]
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8631
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "huaweicloud_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_base(), name)
}

// disable auto_backup to prevent the instance status from changing to "BACKING UP" before manual backup creation.
func testBackup_pg_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 8
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8632
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}

resource "huaweicloud_rds_backup" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_rds_instance.test.id
}
`, testAccRdsInstance_base(), name)
}

func testAccBackupImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
