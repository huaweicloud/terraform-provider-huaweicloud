package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getDatabasePrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dli", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI client: %s", err)
	}

	_, privilege, err := dli.GetObjectPrivilegesForSpecifiedUser(client, state.Primary.Attributes["object"],
		state.Primary.Attributes["user_name"])
	return privilege, err
}

func TestAccDatabasePrivilege_basic(t *testing.T) {
	var (
		obj interface{}

		name        = acceptance.RandomAccResourceName()
		rNameToUser = "huaweicloud_dli_database_privilege.test"

		rcToUser = acceptance.InitResourceCheck(rNameToUser, &obj, getDatabasePrivilegeResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcToUser.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatabasePrivilege_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Resource that authorized IAM user.
					rcToUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameToUser, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.#", "1"),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.0", "SELECT"),
				),
			},
			{
				Config: testDatabasePrivilege_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcToUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.#", "3"),
				),
			},
			{
				ResourceName:      rNameToUser,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDatabasePrivilege_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[1]s"
  data_location = "DLI"

  columns {
    name = "user_name"
    type = "string"
  }
  columns {
    name = "age"
    type = "int"
  }
}
`, name)
}

func testDatabasePrivilege_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_database_privilege" "test" {
  object    = format("databases.%%s.tables.%%s", huaweicloud_dli_database.test.name, huaweicloud_dli_table.test.name)
  user_name = "%[2]s" # Make sure the user has correct permissions and has logged in to the DLI console.
}
`, testDatabasePrivilege_base(name), acceptance.HW_DLI_AUTHORIZED_USER_NAME)
}

func testDatabasePrivilege_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_database_privilege" "test" {
  object    = format("databases.%%s.tables.%%s", huaweicloud_dli_database.test.name, huaweicloud_dli_table.test.name)
  user_name = "%[2]s" # Make sure the user has correct permissions and has logged in to the DLI console.

  privileges = [
    "SELECT", "DESCRIBE_TABLE", "DISPLAY_TABLE"
  ]
}
`, testDatabasePrivilege_base(name), acceptance.HW_DLI_AUTHORIZED_USER_NAME)
}

func TestAccDatabasePrivilege_database(t *testing.T) {
	var (
		obj interface{}

		name        = acceptance.RandomAccResourceName()
		rNameToUser = "huaweicloud_dli_database_privilege.test"

		rcToUser = acceptance.InitResourceCheck(rNameToUser, &obj, getDatabasePrivilegeResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcToUser.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDatabasePrivilege_database_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// Resource that authorized IAM user.
					rcToUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameToUser, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.#", "1"),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.0", "SELECT"),
				),
			},
			{
				Config: testDatabasePrivilege_database_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcToUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameToUser, "privileges.#", "3"),
				),
			},
			{
				ResourceName:      rNameToUser,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testDatabasePrivilege_database_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_database_privilege" "test" {
  object    = format("databases.%%s", huaweicloud_dli_database.test.name)
  user_name = "%[2]s" # Make sure the user has correct permissions and has logged in to the DLI console.
}
`, testDatabasePrivilege_base(name), acceptance.HW_DLI_AUTHORIZED_USER_NAME)
}

func testDatabasePrivilege_database_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_database_privilege" "test" {
  object    = format("databases.%%s", huaweicloud_dli_database.test.name)
  user_name = "%[2]s" # Make sure the user has correct permissions and has logged in to the DLI console.

  privileges = [
    "SELECT", "CREATE_TABLE", "EXPLAIN"
  ]
}
`, testDatabasePrivilege_base(name), acceptance.HW_DLI_AUTHORIZED_USER_NAME)
}
