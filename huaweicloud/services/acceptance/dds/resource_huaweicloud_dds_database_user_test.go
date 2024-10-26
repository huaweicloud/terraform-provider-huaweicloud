package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dds/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDatabaseUserFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s ", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	opts := users.ListOpts{
		Name:   state.Primary.Attributes["name"],
		DbName: state.Primary.Attributes["db_name"],
	}
	resp, err := users.List(client, instanceId, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting user (%s) from DDS instance (%s): %v", name, instanceId, err)
	}
	if len(resp) < 1 {
		return nil, fmt.Errorf("unable to find user (%s) from DDS instance (%s)", name, instanceId)
	}
	user := resp[0]
	return &user, nil
}

func TestAccDatabaseUser_basic(t *testing.T) {
	var user users.UserResp
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_database_user.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getDatabaseUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseUser_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.name",
						"huaweicloud_dds_database_role.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "inherited_privileges",
						"huaweicloud_dds_database_role.test", "inherited_privileges"),
				),
			},
			{
				Config: testAccDatabaseUser_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseUserImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccDatabaseUserImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, dbName, userName string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_dds_database_user" {
				instanceId = rs.Primary.Attributes["instance_id"]
				dbName = rs.Primary.Attributes["db_name"]
				userName = rs.Primary.Attributes["name"]
			}
		}
		if instanceId == "" || dbName == "" || userName == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", instanceId, dbName, userName)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, dbName, userName), nil
	}
}

func testAccDatabaseUser_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_database_role" "test" {
  instance_id = huaweicloud_dds_instance.instance.id

  name    = "%[2]s"
  db_name = "admin"
}

resource "huaweicloud_dds_database_user" "test" {
  instance_id = huaweicloud_dds_instance.instance.id

  name     = "%[2]s"
  password = "HuaweiTest@12345678"
  db_name  = "admin"

  roles {
    name    = huaweicloud_dds_database_role.test.name
    db_name = "admin"
  }
}
`, testAccDDSInstanceV3Config_basic(rName, 8800), rName)
}

func testAccDatabaseUser_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_database_role" "test" {
  instance_id = huaweicloud_dds_instance.instance.id

  name    = "%[2]s"
  db_name = "admin"
}

resource "huaweicloud_dds_database_user" "test" {
  instance_id = huaweicloud_dds_instance.instance.id

  name     = "%[2]s"
  password = "HuaweiTest@123"
  db_name  = "admin"

  roles {
    name    = huaweicloud_dds_database_role.test.name
    db_name = "admin"
  }
}
`, testAccDDSInstanceV3Config_basic(rName, 8800), rName)
}
