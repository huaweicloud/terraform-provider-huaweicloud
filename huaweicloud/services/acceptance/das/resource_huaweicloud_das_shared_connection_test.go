package das

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/das"
)

func getSharedConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid shared connection ID format: %s", state.Primary.ID)
	}

	result, err := das.GetSharedConnectionById(client, parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func TestAccSharedConnection_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName             = "huaweicloud_das_shared_connection.test"
		rcCreateSharedConnection = acceptance.InitResourceCheck(resourceName, &obj, getSharedConnectionResourceFunc)

		name        = acceptance.RandomAccResourceName()
		password    = acceptance.RandomPassword()
		currentTime = time.Now().Local().Format(time.RFC3339)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcCreateSharedConnection.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccSharedConnection_nonExistentSharedConnection(name, password),
				ExpectError: regexp.MustCompile("error creating DAS shared connection"),
			},
			{
				Config: testAccSharedConnection_basic(name, password, currentTime),
				Check: resource.ComposeTestCheckFunc(
					rcCreateSharedConnection.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id"},
			},
		},
	})
}

func testAccSharedConnection_basic_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_das_database_instance_connection" "test" {
  depends_on = [
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id      = "%[1]s"
  engine_type      = "mysql"
  network_type     = "rds"
  username         = "%[2]s"
  password         = "%[3]s"
  is_save_password = true
  node_ids         = ["%[1]s"]
  description      = "Created by terraform script"
  database_name    = "%[2]s"
  sql_record_flag  = true
}

data "huaweicloud_identity_users" "test" {
  depends_on = [
    huaweicloud_das_database_instance_connection.test,
  ]
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}

func testAccSharedConnection_nonExistentSharedConnection(name, password string) string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_shared_connection" "non_existent_shared_connection" {
  connection_id = "%[2]s"
  user_id       = "%[3]s"
  user_name     = "%[4]s"
}
`, testAccSharedConnection_basic_base(name, password), randUUID.String(), name, password)
}

func testAccSharedConnection_basic(name, password, currentTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_shared_connection" "test" {
  depends_on = [
    data.huaweicloud_identity_users.test,
  ]

  connection_id = huaweicloud_das_database_instance_connection.test.id
  user_id       = data.huaweicloud_identity_users.test.users[0].id
  user_name     = data.huaweicloud_identity_users.test.users[0].name
  expired_at    = format("%%s+08:00", split("+", timeadd("%[2]s", "1h"))[0])
}
`, testAccSharedConnection_basic_base(name, password), currentTime)
}
