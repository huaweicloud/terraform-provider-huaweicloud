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

func getPgAccountPrivilegesResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getPgAccountResp, err := pagination.ListAllItems(
		client,
		"page",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account privileges: %s", err)
	}

	respJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return nil, err
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return nil, err
	}

	username := state.Primary.Attributes["user_name"]
	attributes := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0].attributes", username), respBody, nil)

	if attributes == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccPgAccountPrivileges_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_account_privileges.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgAccountPrivilegesResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPgAccountPrivileges_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_name", name),
					resource.TestCheckResourceAttr(rName, "role_privileges.#", "2"),
				),
			},
			{
				Config: testPgAccountPrivileges_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_name", name),
					resource.TestCheckResourceAttr(rName, "role_privileges.#", "2"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"system_role_privileges"},
				ImportStateIdFunc:       testPgAccountPrivilegesImportState(rName),
			},
		},
	})
}

func testPgAccountPrivileges_base(name string) string {
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
    version = "14"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Terraform145@"
}
`, common.TestBaseNetwork(name), name)
}

func testPgAccountPrivileges_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_account_privileges" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  user_name              = huaweicloud_rds_pg_account.test.name
  role_privileges        = ["CREATEDB","LOGIN"]
  system_role_privileges = ["pg_signal_backend"]
}
`, testPgAccountPrivileges_base(name))
}

func testPgAccountPrivileges_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_account_privileges" "test" {
  instance_id            = huaweicloud_rds_instance.test.id
  user_name              = huaweicloud_rds_pg_account.test.name
  role_privileges        = ["CREATEROLE","LOGIN"]
  system_role_privileges = ["pg_monitor"]
}
`, testPgAccountPrivileges_base(name))
}

func testPgAccountPrivilegesImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found: %s", name, rs)
		}
		instanceId := rs.Primary.Attributes["instance_id"]
		username := rs.Primary.Attributes["user_name"]
		return fmt.Sprintf("%s/%s", instanceId, username), nil
	}
}
