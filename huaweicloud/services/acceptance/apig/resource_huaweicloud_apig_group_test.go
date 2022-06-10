package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apigroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigGroupV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				Config: testAccApigGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func TestAccApigGroupV2_variables(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_group.test"
		group        apigroups.Group
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				// Bind two environment to group, and create some variables.
				Config: testAccApigGroup_variables(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "environment.#", "2"),
				),
			},
			{
				// Update the variables for two environments.
				Config: testAccApigGroup_variablesUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigGroupExists(resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "environment.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigInstanceSubResourceImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigGroupDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_group" {
			continue
		}
		_, err := apigroups.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API group (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigGroupExists(groupName string, app *apigroups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[groupName]
		if !ok {
			return fmt.Errorf("Resource %s not found", groupName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No APIG V2 API group Id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := apigroups.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("APIG v2 API group not exist: %s", err)
		}
		*app = *found
		return nil
	}
}

func testAccApigGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigGroup_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s_update"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Updated by script"
}
`, testAccApigApplication_base(rName), rName)
}

func testAccApigGroup_variablesBase(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_environment" "test1" {
  name        = "%s_1"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}

resource "huaweicloud_apig_environment" "test2" {
  name        = "%s_2"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}
`, rName, rName)
}

// Create two environments for the group, and add a total of three variables to the two environments.
// Each of the two environments has a variable with the same name and different value.
func testAccApigGroup_variables(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  environment {
    environment_id = huaweicloud_apig_environment.test1.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
  }
  environment {
    environment_id = huaweicloud_apig_environment.test2.id

    variable {
      name  = "TERRAFORM"
      value = "/res/terraform"
    }
    variable {
      name  = "DEMO"
      value = "/stage/demo"
    }
  }
}
`, testAccApigApplication_base(rName), testAccApigGroup_variablesBase(rName), rName)
}

func testAccApigGroup_variablesUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"

  environment {
    environment_id = huaweicloud_apig_environment.test1.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
    variable {
      name  = "TEST"
      value = "/stage/test"
    }
  }
  environment {
    environment_id = huaweicloud_apig_environment.test2.id

    variable {
      name  = "TERRAFORM"
      value = "/stage/terraform"
    }
  }
}
`, testAccApigApplication_base(rName), testAccApigGroup_variablesBase(rName), rName)
}
