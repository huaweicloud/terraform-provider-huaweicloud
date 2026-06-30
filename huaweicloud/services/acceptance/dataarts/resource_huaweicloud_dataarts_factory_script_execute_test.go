package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getFactoryScriptExecuteResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts-dlf", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}
	return dataarts.GetScriptInstanceById(client, state.Primary.Attributes["workspace_id"],
		state.Primary.Attributes["script_name"], state.Primary.Attributes["id"])
}

func TestAccFactoryScriptExecute_basic(t *testing.T) {
	var (
		obj interface{}

		successExecuteName = "huaweicloud_dataarts_factory_script_execute.success"
		rcSuccessExecute   = acceptance.InitResourceCheck(successExecuteName, &obj, getFactoryScriptExecuteResourceFunc)
		failedExecuteName  = "huaweicloud_dataarts_factory_script_execute.failed"

		rcFailedExecute = acceptance.InitResourceCheck(failedExecuteName, &obj, getFactoryScriptExecuteResourceFunc)
		rName           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and the instance record will not delete even if resource is deleted.
		// So we need to ignore the check destroy.
		// lintignore:AT001
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcSuccessExecute.CheckResourceDestroy(),
			rcFailedExecute.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config:      testAccFactoryScriptExecute_nonExistentWorkspaceAndScript(),
				ExpectError: regexp.MustCompile(`error executing DataArts script`),
			},
			{
				Config: testAccFactoryScriptExecute_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rcSuccessExecute.CheckResourceExists(),
					resource.TestCheckResourceAttr(successExecuteName, "status", "FINISHED"),
					resource.TestCheckResourceAttr(successExecuteName, "message", ""),
					rcFailedExecute.CheckResourceExists(),
					resource.TestCheckResourceAttr(failedExecuteName, "status", "FAILED"),
					resource.TestMatchResourceAttr(failedExecuteName, "message", regexp.MustCompile(`Database 'NULL' does not exist`)),
				),
			},
		},
	})
}

func testAccFactoryScriptExecute_nonExistentWorkspaceAndScript() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_script_execute" "non_existent_script" {
  workspace_id = "%[1]s"
  script_name  = "non_existent_script"
}
`, randUUID.String())
}

func testAccFactoryScriptExecute_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name        = "%[1]s"
  description = "Created by Terraform Script"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[1]s"
  data_location = "DLI"
  description   = "Created by Terraform Script"

  columns {
    name = "name"
    type = "string"
  }

  columns {
    name = "age"
    type = "int"
  }
}

resource "huaweicloud_dataarts_factory_script" "bind_database" {
  depends_on = [
    huaweicloud_dli_database.test,
    huaweicloud_dli_table.test,
  ]

  workspace_id    = "%[2]s"
  name            = "%[1]s_bind_database"
  type            = "DLISQL"
  content         = "SELECT * FROM ${huaweicloud_dli_database.test.name}.${huaweicloud_dli_table.test.name}"
  connection_name = "%[3]s"
  directory       = "/basic"
  database        = huaweicloud_dli_database.test.name
  queue_name      = "default"
  description     = "Created by Terraform Script"
}

resource "huaweicloud_dataarts_factory_script_execute" "success" {
  depends_on = [
    huaweicloud_dataarts_factory_script.bind_database,
  ]

  workspace_id = "%[2]s"
  script_name  = huaweicloud_dataarts_factory_script.bind_database.name
  params       = jsonencode({
    "spark.sql.adaptive.enabled"                                   = "true"
    "spark.sql.adaptive.join.enabled"                              = "true"
    "spark.sql.adaptive.join.skewedJoin.enabled"                   = "true"
    "spark.sql.forcePartitionPredicatesOnPartitionedTable.enabled" = "true"
    "spark.sql.mergeSmallFiles.enabled"                            = "true"
  })
}

resource "huaweicloud_dataarts_factory_script" "unbind_database" {
  depends_on = [
    huaweicloud_dli_database.test,
    huaweicloud_dli_table.test,
  ]

  workspace_id    = "%[2]s"
  name            = "%[1]s_unbind_database"
  type            = "DLISQL"
  content         = "SELECT * FROM ${huaweicloud_dli_database.test.name}.${huaweicloud_dli_table.test.name}"
  connection_name = "%[3]s"
  directory       = "/basic"
  queue_name      = "default"
  description     = "Created by Terraform Script"
}

# The error message is: Database 'NULL' does not exist.
resource "huaweicloud_dataarts_factory_script_execute" "failed" {
  depends_on = [
    huaweicloud_dataarts_factory_script.unbind_database,
  ]

  workspace_id = "%[2]s"
  script_name  = huaweicloud_dataarts_factory_script.unbind_database.name
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CONNECTION_NAME)
}
