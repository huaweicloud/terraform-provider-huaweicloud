package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getExecutionPlanV2ResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("rfs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate UUID: %s", err)
	}

	stackName := state.Primary.Attributes["stack_name"]
	executionPlanName := state.Primary.Attributes["execution_plan_name"]
	return rfs.ReadExecutionPlanV2Detail(client, uuid, stackName, executionPlanName)
}

func TestAccExecutionPlanV2_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_execution_plan_v2.test"
		name  = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getExecutionPlanV2ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccExecutionPlanV2_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "stack_name", name),
					resource.TestCheckResourceAttr(rName, "execution_plan_name", name),
					resource.TestCheckResourceAttr(rName, "description", "test execution plan"),
					resource.TestCheckResourceAttrPair(rName, "stack_id", "huaweicloud_rfs_stack.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "template_body"),
					resource.TestCheckResourceAttrSet(rName, "vars_body"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "execution_plan_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "summary.0.resource_add"),
					resource.TestCheckResourceAttrSet(rName, "summary.0.resource_delete"),
					resource.TestCheckResourceAttrSet(rName, "summary.0.resource_import"),
					resource.TestCheckResourceAttrSet(rName, "summary.0.resource_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testExecutionPlanV2ImportState(rName),
				ImportStateVerifyIgnore: []string{"template_body"},
			},
		},
	})
}

func testAccExecutionPlanV2_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rfs_execution_plan_v2" "test" {
  stack_name          = huaweicloud_rfs_stack.test.name
  execution_plan_name = "%[1]s"
  description         = "test execution plan"
  stack_id            = huaweicloud_rfs_stack.test.id
  template_body       = %[2]s
  vars_body           = %[3]s

  vars_structure {
    var_key   = "vpc_name"
    var_value = "%[1]s-vpc"
  }

  vars_structure {
    var_key   = "subnet_name"
    var_value = "%[1]s-subnet"
  }

  vars_structure {
    var_key   = "instance_count"
    var_value = "2"
  }
}
`, name, updateTemplateInJsonFormat(), basicVariablesInVarsFormat(name))
}

func testExecutionPlanV2ImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		stackName := rs.Primary.Attributes["stack_name"]
		executionPlanName := rs.Primary.Attributes["execution_plan_name"]
		if stackName == "" || executionPlanName == "" {
			return "", fmt.Errorf("the stack_name (%s) or execution_plan_name (%s) is nil",
				stackName, executionPlanName)
		}

		return stackName + "/" + executionPlanName, nil
	}
}
