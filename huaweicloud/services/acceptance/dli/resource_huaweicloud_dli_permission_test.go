package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dli/v1/auth"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDliAuthResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating Dli v1 client, err=%s", err)
	}
	obj, userName := dli.ParseAuthInfoFromId(state.Primary.ID)

	permission, pErr := dli.QueryPermission(client, obj, userName)
	if pErr == nil {
		return permission, nil
	}

	return nil, fmtp.Errorf("This resource is not exist. Id=%s", state.Primary.ID)
}

// test database permissions
func TestAccResourceDliAuth_basic(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_basic(name, acceptance.HW_PROJECT_ID, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("databases.%s", name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "SELECT"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_basic(name, acceptance.HW_PROJECT_ID, "CREATE_TABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("databases.%s", name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "CREATE_TABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
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

func testAccDliAuthResource_basic(name string, projectId string, privileges string) string {
	database := testAccDliDatabase_basic(name)
	userConfig := testAccDliAuthUserConfig(name, projectId)

	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = huaweicloud_identity_user.test.name
  object     = "databases.${huaweicloud_dli_database.test.name}"
  privileges = ["%s"]
}
`, database, userConfig, privileges)
}

func TestAccResourceDliAuth_flink(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_flink(name, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME, "GET"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestMatchResourceAttr(resourceName, "object", regexp.MustCompile(`jobs\.flink\.\d*`)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "GET"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_flink(name, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME, "START"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestMatchResourceAttr(resourceName, "object", regexp.MustCompile(`jobs\.flink\.\d*`)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "START"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
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

func testAccDliAuthResource_flink(name, projectId, region string, privileges string) string {
	flinkJob := testAccFlinkJobResource_basic(name, region)
	userConfig := testAccDliAuthUserConfig(name, projectId)

	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = huaweicloud_identity_user.test.name
  object     = "jobs.flink.${huaweicloud_dli_flinksql_job.test.id}"
  privileges = ["%s"]
}
`, flinkJob, userConfig, privileges)
}

func TestAccResourceDliAuth_table(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_table(name, acceptance.HW_PROJECT_ID, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object",
						fmt.Sprintf("databases.%s.tables.%s", name, name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "SELECT"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_table(name, acceptance.HW_PROJECT_ID, "DROP_TABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object",
						fmt.Sprintf("databases.%s.tables.%s", name, name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "DROP_TABLE"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
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

func testAccDliAuthResource_table(name string, projectId string, privileges string) string {
	table := testAccDliTableResource_basic(name)
	userConfig := testAccDliAuthUserConfig(name, projectId)

	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = huaweicloud_identity_user.test.name
  object     = "databases.${huaweicloud_dli_database.test.name}.tables.${huaweicloud_dli_table.test.name}"
  privileges = ["%s"]
}
`, table, userConfig, privileges)
}

func TestAccResourceDliAuth_package(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
	name := acceptance.RandomAccResourceName()
	packageGroupName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_package(name, acceptance.HW_PROJECT_ID, packageGroupName, "USE_GROUP"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("groups.%s", packageGroupName)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "USE_GROUP"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_package(name, acceptance.HW_PROJECT_ID, packageGroupName, "GET_GROUP"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("groups.%s", packageGroupName)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "GET_GROUP"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
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

func testAccDliAuthResource_package(name, projectId, packageGroupName string, privileges string) string {
	group := testAccDliPackage_basic(packageGroupName)
	userConfig := testAccDliAuthUserConfig(name, projectId)

	return fmt.Sprintf(`
%s
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = huaweicloud_identity_user.test.name
  object     = "groups.${huaweicloud_dli_package.test.group_name}"
  privileges = ["%s"]
}

`, group, userConfig, privileges)
}

func TestAccResourceDliAuth_queue(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_queue(name, acceptance.HW_PROJECT_ID, "DROP_QUEUE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("queues.%s", name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "DROP_QUEUE"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_queue(name, acceptance.HW_PROJECT_ID, "SUBMIT_JOB"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", name),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("queues.%s", name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "SUBMIT_JOB"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
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

func testAccDliAuthResource_queue(name string, projectId string, privileges string) string {
	queue := testAccDliQueue_basic(name, dli.CU_16)
	userConfig := testAccDliAuthUserConfig(name, projectId)

	return fmt.Sprintf(`
%s
%s
resource "huaweicloud_dli_permission" "test" {
  user_name  = huaweicloud_identity_user.test.name
  object     = "queues.${huaweicloud_dli_queue.test.name}"
  privileges = ["%s"]
}
`, queue, userConfig, privileges)
}

func testAccDliAuthUserConfig(name string, projectId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
}

resource "huaweicloud_identity_role" "role_1" {
  name        = "%s"
  description = "DLI readOnly role created by terraform"
  type        = "AX"
  policy      = <<EOF
{
    "Version": "1.1",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "dli:jobs:get",
                "dli:queue:show_privileges",
                "dli:database:show_functions",
                "dli:table:show_segments",
                "dli:table:describe_table",
                "dli:table:show_privileges",
                "dli:table:select",
                "dli:database:displayAllDatabases",
                "dli:database:show_roles",
                "dli:database:show_users",
                "dli:datasourceauth:show_privileges",
                "dli:database:show_privileges",
                "dli:database:show_all_roles",
                "dli:column:select",
                "dli:datasourceauth:use_auth",
                "dli:table:show_partitions",
                "dli:database:displayAllTables",
                "dli:table:show_table_properties",
                "dli:table:show_create_table"
            ]
        }
    ]
}
EOF
}

resource "huaweicloud_identity_role_assignment" "role_assignment_1" {
  role_id    = huaweicloud_identity_role.role_1.id
  group_id   = huaweicloud_identity_group.group_1.id
  project_id = "%s"
}

resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.test.id]
}

`, name, name, name, projectId)
}
