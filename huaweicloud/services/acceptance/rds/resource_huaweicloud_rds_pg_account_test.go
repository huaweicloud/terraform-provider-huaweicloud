package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPgAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPgAccount: query RDS PostgreSQL account
	var (
		getPgAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getPgAccountProduct = "rds"
	)
	getPgAccountClient, err := cfg.NewServiceClient(getPgAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getPgAccountPath := getPgAccountClient.Endpoint + getPgAccountHttpUrl
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{project_id}", getPgAccountClient.ProjectID)
	getPgAccountPath = strings.ReplaceAll(getPgAccountPath, "{instance_id}", instanceId)

	getPgAccountResp, err := pagination.ListAllItems(
		getPgAccountClient,
		"page",
		getPgAccountPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	getPgAccountRespJson, err := json.Marshal(getPgAccountResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	var getPgAccountRespBody interface{}
	err = json.Unmarshal(getPgAccountRespJson, &getPgAccountRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL account: %s", err)
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getPgAccountRespBody, nil)

	if account != nil {
		return account, nil
	}

	return nil, fmt.Errorf("error retrieving RDS PostgreSQL account by instanceID %s and account %s", instanceId,
		accountName)
}

func TestAccPgAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
					resource.TestCheckResourceAttrPair(rName, "memberof.0",
						"huaweicloud_rds_pg_account.member", "name"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_super"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_inherit"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_create_role"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_create_db"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_can_login"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_conn_limit"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_replication"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.rol_bypass_rls"),
				),
			},
			{
				Config: testPgAccount_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
					resource.TestCheckResourceAttrPair(rName, "memberof.0",
						"huaweicloud_rds_pg_account.member_update", "name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testPgAccount_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "16"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testPgAccount_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account" "member" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_member"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_pg_account" "test" {
  depends_on  = [huaweicloud_rds_pg_account.member]
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Terraform145@"
  description = "test_description"
  memberof    = [huaweicloud_rds_pg_account.member.name]
}
`, testPgAccount_base(name), name)
}

func testPgAccount_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account" "member" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_member"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_pg_account" "member_update" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_member_update"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_pg_account" "test" {
  depends_on  = [huaweicloud_rds_pg_account.member_update]
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Terraform145@"
  description = "test_description_update"
  memberof    = [huaweicloud_rds_pg_account.member_update.name]
}
`, testPgAccount_base(name), name)
}
