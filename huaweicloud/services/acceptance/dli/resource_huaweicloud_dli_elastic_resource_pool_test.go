package dli

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getElasticResourcePoolFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dli", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI client: %s", err)
	}

	return dli.GetElasticResourcePoolByName(client, state.Primary.Attributes["name"])
}

func TestAccElasticResourcePool_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj            interface{}
		withStandard   = "huaweicloud_dli_elastic_resource_pool.test"
		rcWithStandard = acceptance.InitResourceCheck(withStandard, &obj, getElasticResourcePoolFunc)

		withBasic   = "huaweicloud_dli_elastic_resource_pool.basic"
		rcWithBasic = acceptance.InitResourceCheck(withBasic, &obj, getElasticResourcePoolFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithStandard.CheckResourceDestroy(),
			rcWithBasic.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccElasticResourcePool_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithStandard.CheckResourceExists(),
					resource.TestCheckResourceAttr(withStandard, "name", name),
					resource.TestCheckResourceAttr(withStandard, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withStandard, "min_cu", "80"),
					resource.TestCheckResourceAttr(withStandard, "max_cu", "96"),
					resource.TestCheckResourceAttr(withStandard, "cidr", "172.16.0.0/12"),
					resource.TestCheckResourceAttr(withStandard, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(withStandard, "status"),
					resource.TestCheckResourceAttrSet(withStandard, "current_cu"),
					resource.TestCheckResourceAttrSet(withStandard, "created_at"),
					rcWithBasic.CheckResourceExists(),
					resource.TestCheckResourceAttr(withBasic, "min_cu", "32"),
					resource.TestCheckResourceAttr(withBasic, "max_cu", "48"),
					resource.TestCheckResourceAttr(withBasic, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(withBasic, "tags.%", "1"),
					resource.TestCheckResourceAttr(withBasic, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withBasic, "label.spec", "basic"),
					// The number of parameter 'actual_cu' needs to be returned after the queues are created, and will
					// not be tested for the time being.
				),
			},
			{
				Config: testAccElasticResourcePool_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithStandard.CheckResourceExists(),
					resource.TestCheckResourceAttr(withStandard, "name", name),
					resource.TestCheckResourceAttr(withStandard, "description", ""),
					resource.TestCheckResourceAttr(withStandard, "min_cu", "64"),
					resource.TestCheckResourceAttr(withStandard, "max_cu", "112"),
					resource.TestCheckResourceAttr(withStandard, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					rcWithBasic.CheckResourceExists(),
					resource.TestCheckResourceAttr(withBasic, "min_cu", "16"),
					resource.TestCheckResourceAttr(withBasic, "max_cu", "64"),
					resource.TestCheckResourceAttr(withBasic, "description", "Updated basic resource pool by terraform script"),
					resource.TestCheckResourceAttr(withBasic, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					waitForDeletionCooldownComplete(),
					waitForCUModificationComplete(withStandard),
					waitForCUModificationComplete(withBasic),
				),
			},
			{
				ResourceName:      withStandard,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccElasticResourcePoolImportStateFunc(withStandard),
			},
			{
				ResourceName:      withBasic,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccElasticResourcePoolImportStateFunc(withBasic),
				ImportStateVerifyIgnore: []string{
					"tags",
				},
			},
		},
	})
}

func waitForDeletionCooldownComplete() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// After elastic resource pool is created, it cannot be deleted within 15 minete.
		// lintignore:R018
		time.Sleep(15 * time.Minute)
		return nil
	}
}

func waitForCUModificationComplete(resName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("resource (%s) not found", resName)
		}
		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := cfg.NewServiceClient("dli", acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DLI client: %s", err)
		}

		startTime := time.Now() // Record the start time of the loop.
		for {
			respBody, err := dli.GetElasticResourcePoolByName(client, rs.Primary.Attributes["name"])
			if err != nil {
				return fmt.Errorf("error querying elastic resource pool (%s): %s", rs.Primary.Attributes["name"], err)
			}
			status := utils.PathSearch("status", respBody, "").(string)
			if status == "FAILED" {
				failReason := utils.PathSearch("fail_reason", respBody, "").(string)
				return fmt.Errorf("the CU number modification failed: %s", failReason)
			}
			if status == "AVAILABLE" {
				return nil
			}
			if time.Since(startTime) > 15*time.Minute {
				break
			}
			// lintignore:R018
			time.Sleep(30 * time.Second)
		}
		return fmt.Errorf("modification timeout for the CU number")
	}
}

func testAccElasticResourcePoolImportStateFunc(resName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resName, rs)
		}
		poolName := rs.Primary.Attributes["name"]
		if poolName == "" {
			return "", fmt.Errorf("the resource pool name is missing")
		}
		return poolName, nil
	}
}

func testAccElasticResourcePool_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = upper("%[1]s")
  description           = "Created by terraform script"
  max_cu                = 96
  min_cu                = 80
  enterprise_project_id = "0"
}

resource "huaweicloud_dli_elastic_resource_pool" "basic" {
  name        = upper("%[1]s_basic")
  description = "Created basic resource pool by terraform script"
  min_cu      = 32
  max_cu      = 48

  tags = {
    foo = "bar"
  }
 
  label = {
    spec = "basic"
  }
}
`, name)
}

func testAccElasticResourcePool_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[1]s"
  max_cu                = 112
  min_cu                = 64
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_dli_elastic_resource_pool" "basic" {
  name                  = upper("%[1]s_basic")
  description           = "Updated basic resource pool by terraform script"
  min_cu                = 16
  max_cu                = 64
  enterprise_project_id = "%[2]s"

  tags = {
    foo = "bar"
  }

  label = {
    spec = "basic"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccElasticResourcePool_prePaid(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj          interface{}
		resourceName = "huaweicloud_dli_elastic_resource_pool.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getElasticResourcePoolFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccElasticResourcePool_prePaid_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "min_cu", "16"),
					resource.TestCheckResourceAttr(resourceName, "max_cu", "32"),
					resource.TestCheckResourceAttr(resourceName, "cidr", "172.16.0.0/12"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "current_cu"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "prepay_cu"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "label.spec", "basic"),
					// The number of parameter 'actual_cu' needs to be returned after the queues are created, and will
					// not be tested for the time being.
				),
			},
			{
				Config: testAccElasticResourcePool_prePaid_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "min_cu", "16"),
					resource.TestCheckResourceAttr(resourceName, "max_cu", "48"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					// waitForDeletionCooldownComplete(),
					waitForCUModificationComplete(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccElasticResourcePoolImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"tags",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}

func testAccElasticResourcePool_prePaid_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name          = "%[1]s"
  description   = "Created by terraform script"
  min_cu        = 16
  max_cu        = 32
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  tags = {
    foo = "bar"
  }

  label = {
    spec = "basic"
  }
}
`, name)
}

func testAccElasticResourcePool_prePaid_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_elastic_resource_pool" "test" {
  name                  = "%[1]s"
  min_cu                = 16
  max_cu                = 48
  enterprise_project_id = "%[2]s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  auto_renew            = "false"

  tags = {
    foo = "bar"
  }

  label = {
    spec = "basic"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
