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

func getBackupImportTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/migration-task/{task_id}"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", getRespBody, "").(string)
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccDcsBackupImportTask_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_backup_import_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupImportTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcsObsBucketName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsBackupImportTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "migration_type", "backupfile_import"),
					resource.TestCheckResourceAttr(rName, "migration_method", "full_amount_migration"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "backup_files.0.file_source", "self_build_obs"),
					resource.TestCheckResourceAttr(rName, "backup_files.0.bucket_name", acceptance.HW_DCS_OBS_BUCKET_NAME),
					resource.TestCheckResourceAttr(rName, "backup_files.0.files.0.file_name", "appendonly.aof"),
					resource.TestCheckResourceAttr(rName, "backup_files.0.files.0.size", "0.09KB"),
					resource.TestCheckResourceAttr(rName, "backup_files.0.files.0.update_at", "2025-07-22T11:57:24.401Z"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.id",
						"huaweicloud_dcs_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.name",
						"huaweicloud_dcs_instance.test", "name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "released_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_instance.0.password"},
			},
		},
	})
}

func TestAccDcsBackupImportTask_backupId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_backup_import_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getBackupImportTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsBackupImportTask_backupId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "migration_type", "backupfile_import"),
					resource.TestCheckResourceAttr(rName, "migration_method", "full_amount_migration"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "backup_files.0.file_source", "backup_record"),
					resource.TestCheckResourceAttrPair(rName, "backup_files.0.backup_id",
						"huaweicloud_dcs_backup.test", "backup_id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.id",
						"huaweicloud_dcs_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.name",
						"huaweicloud_dcs_instance.test", "name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "released_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"target_instance.0.password"},
			},
		},
	})
}

func testDcsBackupImportTask_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "6.0"
  capacity       = 1
  name           = "redis.ha.au1.large.r4.1"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "6.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  action            = "allow"
  ethertype         = "IPv4"
  ports             = "6379"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = data.huaweicloud_networking_secgroup.test.id
}
`, name)
}

func testDcsBackupImportTask_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dcs_backup_import_task" "test" {
  depends_on = [huaweicloud_networking_secgroup_rule.egress]

  task_name        = "%[2]s"
  migration_type   = "backupfile_import"
  migration_method = "full_amount_migration"
  description      = "terraform test"

  backup_files {
    file_source = "self_build_obs"
    bucket_name = "%[3]s"

    files {
      file_name = "appendonly.aof"
      size      = "0.09KB"
      update_at = "2025-07-22T11:57:24.401Z"
    }
  }

  target_instance{
    id       = huaweicloud_dcs_instance.test.id
    password = "Huawei_test"
  }

  lifecycle {
    ignore_changes = [
      target_instance.0.password
    ]
  }
}
`, testDcsBackupImportTask_base(name), name, acceptance.HW_DCS_OBS_BUCKET_NAME)
}

func testDcsBackupImportTask_backupId(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dcs_backup" "test" {
  instance_id   = huaweicloud_dcs_instance.test.id
  backup_format = "rdb"
}

resource "huaweicloud_dcs_backup_import_task" "test" {
  depends_on = [huaweicloud_networking_secgroup_rule.egress]

  task_name        = "%[2]s"
  migration_type   = "backupfile_import"
  migration_method = "full_amount_migration"
  description      = "terraform test"

  backup_files {
    file_source = "backup_record"
    backup_id   = huaweicloud_dcs_backup.test.backup_id
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test.id
    password = "Huawei_test"
  }

  lifecycle {
    ignore_changes = [
      target_instance.0.password
    ]
  }
}
`, testDcsBackupImportTask_base(name), name)
}
