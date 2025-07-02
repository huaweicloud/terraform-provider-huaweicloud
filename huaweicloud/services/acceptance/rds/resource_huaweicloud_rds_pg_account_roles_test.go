package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPgAccountRolesResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getResp, err := pagination.ListAllItems(
		client,
		"page",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account roles: %s", err)
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account roles: %s", err)
	}

	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account roles: %s", err)
	}

	roles := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0].memberof", accountName), getRespBody, nil)

	if roles == nil || len(roles.([]interface{})) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccPgAccountRoles_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_account_roles.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgAccountRolesResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgAccountRoles_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "user", "root"),
					resource.TestCheckResourceAttrPair(rName, "roles.0",
						"huaweicloud_rds_pg_account.test.0", "name"),
				),
			},
			{
				Config: testPgAccountRoles_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "user", "root"),
					resource.TestCheckResourceAttrPair(rName, "roles.0",
						"huaweicloud_rds_pg_account.test.1", "name"),
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

func testPgAccountRoles_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "ha"
  group_type    = "dedicated"
  vcpus         = 2
}

resource "huaweicloud_rds_instance" "test" {
  name                = "%[2]s"
  flavor              = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  ha_replication_mode = "sync"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1]
  ]

  db {
    type    = "PostgreSQL"
    version = "12"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_pg_account" "test" {
  count = 2

  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_${count.index}"
  password    = "Terraform145@"
}
`, common.TestBaseNetwork(name), name)
}

func testPgAccountRoles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account_roles" "test" {
  depends_on = [
    huaweicloud_rds_pg_account.test[0],
    huaweicloud_rds_pg_account.test[1]
  ]

  instance_id = huaweicloud_rds_instance.test.id
  user        = "root"
  roles       = [huaweicloud_rds_pg_account.test[0].name]
}
`, testPgAccountRoles_base(name), name)
}

func testPgAccountRoles_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account_roles" "test" {
  depends_on = [
    huaweicloud_rds_pg_account.test[0],
    huaweicloud_rds_pg_account.test[1]
  ]

  instance_id = huaweicloud_rds_instance.test.id
  user        = "root"
  roles       = [huaweicloud_rds_pg_account.test[1].name]
}
`, testPgAccountRoles_base(name), name)
}
