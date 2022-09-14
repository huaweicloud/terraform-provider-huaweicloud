package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func getRdsDatabasePrivilegeFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcRdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<database_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]
	return rds.QueryDatabaseUsers(client, instanceId, dbName)
}

func TestAccRdsDatabasePrivilege_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_database_privilege.test"
	var users []model.UserWithPrivilege
	rc := acceptance.InitResourceCheck(
		resourceName,
		&users,
		getRdsDatabasePrivilegeFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDatabasePrivilege_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "users.0.name",
						"huaweicloud_rds_account.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "users.0.readonly", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsDatabasePrivilege_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%s"
  password    = "Test@12345678"
}

resource "huaweicloud_rds_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_database.test.name

  users {
    name = huaweicloud_rds_account.test.name
  }
}
`, testRdsDatabase_basic(rName, rName), rName)
}
