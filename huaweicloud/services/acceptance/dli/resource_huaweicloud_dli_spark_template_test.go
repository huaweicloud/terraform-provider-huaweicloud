package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSparkTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSparkTemplate: Query the Spark template.
	var (
		getSparkTemplateHttpUrl = "v3/{project_id}/templates/{id}"
		getSparkTemplateProduct = "dli"
	)
	getSparkTemplateClient, err := cfg.NewServiceClient(getSparkTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getSparkTemplatePath := getSparkTemplateClient.Endpoint + getSparkTemplateHttpUrl
	getSparkTemplatePath = strings.ReplaceAll(getSparkTemplatePath, "{project_id}", getSparkTemplateClient.ProjectID)
	getSparkTemplatePath = strings.ReplaceAll(getSparkTemplatePath, "{id}", state.Primary.ID)

	getSparkTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSparkTemplateResp, err := getSparkTemplateClient.Request("GET", getSparkTemplatePath, &getSparkTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving SparkTemplate: %s", err)
	}
	return utils.FlattenResponse(getSparkTemplateResp)
}

func TestAccSparkTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dli_spark_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSparkTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSparkTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "group", "demo"),
					resource.TestCheckResourceAttr(rName, "body.0.queue_name", "queue_demo"),
					resource.TestCheckResourceAttr(rName, "body.0.name", "demo"),
					resource.TestCheckResourceAttr(rName, "body.0.app_name", "jar_package/demo.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.main_class", "com.demo.main"),
					resource.TestCheckResourceAttr(rName, "body.0.app_parameters.0", "abc"),
					resource.TestCheckResourceAttr(rName, "body.0.specification", "A"),
					resource.TestCheckResourceAttr(rName, "body.0.jars.0", "jar_package/demo.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.python_files.0", "python_package/demo.py"),
					resource.TestCheckResourceAttr(rName, "body.0.files.0", "file_package/demo.txt"),
					resource.TestCheckResourceAttr(rName, "body.0.modules.#", "2"),
					resource.TestCheckResourceAttr(rName, "body.0.modules.0", "sys.res.dli"),
					resource.TestCheckResourceAttr(rName, "body.0.modules.1", "sys.datasource.dws"),
					resource.TestCheckResourceAttr(rName, "body.0.resources.0.name", "jar_package/demo_dep.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.resources.0.type", "jar"),
					resource.TestCheckResourceAttr(rName, "body.0.dependent_packages.0.name", "driver_package"),
					resource.TestCheckResourceAttr(rName, "body.0.configurations.a", "b"),
					resource.TestCheckResourceAttr(rName, "body.0.driver_memory", "2G"),
					resource.TestCheckResourceAttr(rName, "body.0.driver_cores", "3"),
					resource.TestCheckResourceAttr(rName, "body.0.executor_cores", "4"),
					resource.TestCheckResourceAttr(rName, "body.0.executor_memory", "5G"),
					resource.TestCheckResourceAttr(rName, "body.0.num_executors", "6"),
					resource.TestCheckResourceAttr(rName, "body.0.obs_bucket", "demo"),
				),
			},
			{
				Config: testSparkTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo2"),
					resource.TestCheckResourceAttr(rName, "group", "demo2"),
					resource.TestCheckResourceAttr(rName, "body.0.queue_name", "queue_demo2"),
					resource.TestCheckResourceAttr(rName, "body.0.name", "demo2"),
					resource.TestCheckResourceAttr(rName, "body.0.app_name", "jar_package/demo.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.main_class", "com.demo.main2"),
					resource.TestCheckResourceAttr(rName, "body.0.app_parameters.0", "abcd"),
					resource.TestCheckResourceAttr(rName, "body.0.specification", "B"),
					resource.TestCheckResourceAttr(rName, "body.0.jars.0", "jar_package/demo2.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.python_files.0", "python_package/demo2.py"),
					resource.TestCheckResourceAttr(rName, "body.0.files.0", "file_package/demo2.txt"),
					resource.TestCheckResourceAttr(rName, "body.0.modules.#", "1"),
					resource.TestCheckResourceAttr(rName, "body.0.modules.0", "sys.datasource.dws"),
					resource.TestCheckResourceAttr(rName, "body.0.resources.0.name", "jar_package/demo_dep2.jar"),
					resource.TestCheckResourceAttr(rName, "body.0.resources.0.type", "jar"),
					resource.TestCheckResourceAttr(rName, "body.0.dependent_packages.0.name", "driver_package2"),
					resource.TestCheckResourceAttr(rName, "body.0.configurations.c", "d"),
					resource.TestCheckResourceAttr(rName, "body.0.driver_memory", "3G"),
					resource.TestCheckResourceAttr(rName, "body.0.driver_cores", "4"),
					resource.TestCheckResourceAttr(rName, "body.0.executor_cores", "5"),
					resource.TestCheckResourceAttr(rName, "body.0.executor_memory", "6G"),
					resource.TestCheckResourceAttr(rName, "body.0.num_executors", "7"),
					resource.TestCheckResourceAttr(rName, "body.0.obs_bucket", "demo2"),
					resource.TestCheckResourceAttr(rName, "body.0.auto_recovery", "true"),
					resource.TestCheckResourceAttr(rName, "body.0.max_retry_times", "5"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testSparkTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_spark_template" "test" {
  name        = "%s"
  description = "This is a demo"
  group       = "demo"

  body {
    queue_name     = "queue_demo"
    name           = "demo"
    app_name       = "jar_package/demo.jar"
    main_class     = "com.demo.main"
    app_parameters = ["abc"]
    specification  = "A"

    jars = [
      "jar_package/demo.jar"
    ]

    python_files = [
      "python_package/demo.py"
    ]

    files = [
      "file_package/demo.txt"
    ]

    modules = [
      "sys.res.dli",
      "sys.datasource.dws"
    ]

    resources {
      name = "jar_package/demo_dep.jar"
      type = "jar"
    }

    dependent_packages {
      name = "driver_package"
    }

    configurations = {
      "a" : "b"
    }

    driver_memory   = "2G"
    driver_cores    = 3
    executor_cores  = 4
    executor_memory = "5G"
    num_executors   = 6
    obs_bucket      = "demo"
  }
}
`, name)
}

func testSparkTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_spark_template" "test" {
  name        = "%s"
  description = "This is a demo2"
  group       = "demo2"

  body {
    queue_name     = "queue_demo2"
    name           = "demo2"
    app_name       = "jar_package/demo.jar"
    main_class     = "com.demo.main2"
    app_parameters = ["abcd"]
    specification  = "B"

    jars = [
      "jar_package/demo2.jar"
    ]

    python_files = [
      "python_package/demo2.py"
    ]

    files = [
      "file_package/demo2.txt"
    ]

    modules = [
      "sys.datasource.dws"
    ]

    resources {
      name = "jar_package/demo_dep2.jar"
      type = "jar"
    }

    dependent_packages {
      name = "driver_package2"
    }

    configurations = {
      "c" : "d"
    }

    driver_memory   = "3G"
    driver_cores    = 4
    executor_cores  = 5
    executor_memory = "6G"
    num_executors   = 7
    obs_bucket      = "demo2"
    auto_recovery   = true
    max_retry_times = 5
  }
}
`, name)
}
