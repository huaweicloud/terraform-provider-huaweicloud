package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v1/auth"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getDliAuthResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.DliV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Dli v1 client, err=%s", err)
	}
	obj, userName := dli.ParseAuthInfoFromId(state.Primary.ID)

	permission, pErr := dli.QueryPermission(client, obj, userName)
	if pErr != nil {
		return nil, fmt.Errorf("this resource is not exist. Id=%s", state.Primary.ID)
	}

	return permission, nil
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_basic(name, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("databases.%s", name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "SELECT"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_basic(name, "CREATE_TABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
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

func testAccDliAuthResource_basic(name string, privileges string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name = "%s"
}

resource "huaweicloud_dli_permission" "test" {
  user_name  = "%s"
  object     = "databases.${huaweicloud_dli_database.test.name}"
  privileges = ["%s"]
}
`, name, acceptance.HW_DLI_AUTHORIZED_USER_NAME, privileges)
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
			acceptance.TestAccPreCheckDliFlinkVersion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_flink(name, "GET"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestMatchResourceAttr(resourceName, "object", regexp.MustCompile(`jobs\.flink\.\d*`)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "GET"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_flink(name, "START"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
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

func testAccDliAuthResource_flink(name, privileges string) string {
	flinkJob := testAccFlinkJobResource_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = "%s"
  object     = "jobs.flink.${huaweicloud_dli_flinksql_job.test.id}"
  privileges = ["%s"]
}
`, flinkJob, acceptance.HW_DLI_AUTHORIZED_USER_NAME, privileges)
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_table(name, "SELECT"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(resourceName, "object",
						fmt.Sprintf("databases.%s.tables.%s", name, name)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "SELECT"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_table(name, "DROP_TABLE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
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

func testAccDliAuthResource_table(name, privileges string) string {
	table := testAccDliTableResource_basic(name)

	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = "%s"
  object     = "databases.${huaweicloud_dli_database.test.name}.tables.${huaweicloud_dli_table.test.name}"
  privileges = ["%s"]
}
`, table, acceptance.HW_DLI_AUTHORIZED_USER_NAME, privileges)
}

func TestAccResourceDliAuth_package(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"
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
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_package(packageGroupName, "USE_GROUP"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("groups.%s", packageGroupName)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "USE_GROUP"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_package(packageGroupName, "GET_GROUP"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
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

func testAccDliAuthResource_package(packageGroupName string, privileges string) string {
	group := testAccDliPackage_basic(packageGroupName)

	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_permission" "test" {
  user_name  = "%s"
  object     = "groups.${huaweicloud_dli_package.test.group_name}"
  privileges = ["%s"]
}

`, group, acceptance.HW_DLI_AUTHORIZED_USER_NAME, privileges)
}

func TestAccResourceDliAuth_queue(t *testing.T) {
	var obj auth.Privilege
	resourceName := "huaweicloud_dli_permission.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDliAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAuthorizedUserConfigured(t)
			acceptance.TestAccPreCheckDliSQLQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliAuthResource_queue("DROP_QUEUE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("queues.%s", acceptance.HW_DLI_SQL_QUEUE_NAME)),
					resource.TestCheckResourceAttr(resourceName, "privileges.0", "DROP_QUEUE"),
					resource.TestCheckResourceAttrSet(resourceName, "is_admin"),
				),
			},
			{
				Config: testAccDliAuthResource_queue("SUBMIT_JOB"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "user_name", acceptance.HW_DLI_AUTHORIZED_USER_NAME),
					resource.TestCheckResourceAttr(resourceName, "object", fmt.Sprintf("queues.%s", acceptance.HW_DLI_SQL_QUEUE_NAME)),
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

func testAccDliAuthResource_queue(privileges string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_permission" "test" {
  user_name  = "%s"
  object     = "queues.%s"
  privileges = ["%s"]
}
`, acceptance.HW_DLI_AUTHORIZED_USER_NAME, acceptance.HW_DLI_SQL_QUEUE_NAME, privileges)
}
