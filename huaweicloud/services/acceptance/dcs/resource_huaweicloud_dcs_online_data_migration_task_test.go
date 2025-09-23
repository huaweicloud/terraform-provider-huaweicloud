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

func getOnlineDataMigrationTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccDcsOnlineDataMigrationTask_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_online_data_migration_task.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOnlineDataMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDcsOnlineDataMigrationTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "migration_method", "full_amount_migration"),
					resource.TestCheckResourceAttr(rName, "resume_mode", "manual"),
					resource.TestCheckResourceAttrPair(rName, "source_instance.0.id",
						"huaweicloud_dcs_instance.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "source_instance.0.name",
						"huaweicloud_dcs_instance.test.0", "name"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.id",
						"huaweicloud_dcs_instance.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.name",
						"huaweicloud_dcs_instance.test.1", "name"),
					resource.TestCheckResourceAttrSet(rName, "ecs_tenant_private_ip"),
					resource.TestCheckResourceAttrSet(rName, "network_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "supported_features.#"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testDcsOnlineDataMigrationTask_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"data.huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "subnet_id",
						"data.huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "migration_method", "incremental_migration"),
					resource.TestCheckResourceAttr(rName, "resume_mode", "auto"),
					resource.TestCheckResourceAttrPair(rName, "source_instance.0.id",
						"huaweicloud_dcs_instance.test.2", "id"),
					resource.TestCheckResourceAttrPair(rName, "source_instance.0.name",
						"huaweicloud_dcs_instance.test.2", "name"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.id",
						"huaweicloud_dcs_instance.test.3", "id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance.0.name",
						"huaweicloud_dcs_instance.test.3", "name"),
					resource.TestCheckResourceAttrSet(rName, "ecs_tenant_private_ip"),
					resource.TestCheckResourceAttrSet(rName, "network_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "supported_features.#"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source_instance.0.password", "target_instance.0.password"},
			},
		},
	})
}

func testDcsOnlineDataMigrationTask_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "6.0"
  capacity       = 1
  name           = "redis.ha.au1.large.r4.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dcs_instance" "test" {
  count = 4

  name               = "%[1]s_${count.index}"
  engine_version     = "6.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
`, name)
}

func testDcsOnlineDataMigrationTask_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dcs_online_data_migration_task" "test" {
  task_name         = "%[2]s"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  description       = "terraform test"
  migration_method  = "full_amount_migration"
  resume_mode       = "manual"

  source_instance {
    id       = huaweicloud_dcs_instance.test[0].id
    password = "Huawei_test"
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test[1].id
    password = "Huawei_test"
  }

  lifecycle {
    ignore_changes = [
      source_instance.0.addrs, target_instance.0.addrs,
    ]
  }
}
`, testDcsOnlineDataMigrationTask_base(name), name)
}

func testDcsOnlineDataMigrationTask_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dcs_online_data_migration_task" "test" {
  task_name          = "%[2]s"
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  description        = "terraform test"
  migration_method   = "incremental_migration"
  resume_mode        = "auto"
  bandwidth_limit_mb = "100"

  source_instance {
    id       = huaweicloud_dcs_instance.test[2].id
    password = "Huawei_test"
  }

  target_instance {
    id       = huaweicloud_dcs_instance.test[3].id
    password = "Huawei_test"
  }

  lifecycle {
    ignore_changes = [
      source_instance.0.addrs, target_instance.0.addrs,
    ]
  }
}
`, testDcsOnlineDataMigrationTask_base(name), name)
}
