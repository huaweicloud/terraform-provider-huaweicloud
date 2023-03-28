package fgs

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud FunctionGraph client: %s", err)
	}
	return function.GetMetadata(c, state.Primary.ID).Extract()
}

func TestAccFgsV2Function_basic(t *testing.T) {
	var f function.Function
	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "functiongraph_version", "v1"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.max_async_event_age_in_seconds", "3500"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.max_async_retry_attempts", "2"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.on_success.0.destination", "OBS"),
					resource.TestCheckResourceAttrSet(rName, "async_invoke.0.on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.on_failure.0.destination", "SMN"),
					resource.TestCheckResourceAttrSet(rName, "async_invoke.0.on_failure.0.param"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.enable_async_status_log", "true"),
				),
			},
			{
				Config: testAccFunction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", "fuction test update"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.max_async_event_age_in_seconds", "4000"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.max_async_retry_attempts", "3"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.on_success.0.destination", "DIS"),
					resource.TestCheckResourceAttrSet(rName, "async_invoke.0.on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.on_failure.0.destination", "FunctionGraph"),
					resource.TestCheckResourceAttrSet(rName, "async_invoke.0.on_failure.0.param"),
					resource.TestCheckResourceAttr(rName, "async_invoke.0.enable_async_status_log", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
					"xrole",
					"agency",
				},
			},
		},
	})
}

func testAccFunction_basic_step1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "function_test_trust"

  async_invoke {
    max_async_event_age_in_seconds = 3500
    max_async_retry_attempts       = 2
    enable_async_status_log        = true

    on_success {
      destination = "OBS"
      param = jsonencode({
        bucket  = huaweicloud_obs_bucket.test.bucket
        prefix  = "/success"
        expires = 5
      })
    }

    on_failure {
      destination = "SMN"
      param       = jsonencode({
        topic_urn = huaweicloud_smn_topic.test.topic_urn
      })
    }
  }
}
`, rName)
}

func testAccFunction_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dis_stream" "test" {
  stream_name     = "%[2]s"
  partition_count = 1
}

resource "huaweicloud_fgs_function" "failure_transport" {
  name        = "%[2]s-failure-transport"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  app         = "default"
  description = "fuction test update"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "function_test_trust"
  vpc_id      = huaweicloud_vpc.test.id
  network_id  = huaweicloud_vpc_subnet.test.id

  async_invoke {
    max_async_event_age_in_seconds = 4000
    max_async_retry_attempts       = 3

    on_success {
      destination = "DIS"
      param = jsonencode({
        stream_name = huaweicloud_dis_stream.test.stream_name
      })
    }

    on_failure {
      destination = "FunctionGraph"
      param       = jsonencode({
        func_urn = huaweicloud_fgs_function.failure_transport.id
      })
    }
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func TestAccFgsV2Function_withEpsId(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_withEpsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_text(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_text(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_createByImage(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	rName1 := "huaweicloud_fgs_function.create_with_vpc_access"
	rName2 := "huaweicloud_fgs_function.create_without_vpc_access"

	rc1 := acceptance.InitResourceCheck(
		rName1,
		&f,
		getResourceObj,
	)

	rc2 := acceptance.InitResourceCheck(
		rName2,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckComponentDeployment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_createByImage_step_1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", randName+"_1"),
					resource.TestCheckResourceAttr(rName1, "agency", "functiongraph_swr_trust"),
					resource.TestCheckResourceAttr(rName1, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(rName1, "handler", "-"),
					resource.TestCheckResourceAttr(rName1, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttrPair(rName1, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName1, "network_id", "huaweicloud_vpc_subnet.test", "id"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", randName+"_2"),
					resource.TestCheckResourceAttr(rName2, "agency", "functiongraph_swr_trust"),
					resource.TestCheckResourceAttr(rName2, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(rName2, "handler", "-"),
					resource.TestCheckResourceAttr(rName2, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttr(rName2, "vpc_id", ""),
					resource.TestCheckResourceAttr(rName2, "network_id", ""),
				),
			},
			{
				Config: testAccFgsV2Function_createByImage_step_2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "handler", "-"),
					resource.TestCheckResourceAttr(rName1, "vpc_id", ""),
					resource.TestCheckResourceAttr(rName1, "network_id", ""),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "handler", "-"),
					resource.TestCheckResourceAttrPair(rName2, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(rName2, "network_id", "huaweicloud_vpc_subnet.test", "id"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"xrole",
					"agency",
				},
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"xrole",
					"agency",
				},
			},
		},
	})
}

func testAccFgsV2Function_text(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}
`, rName)
}

func testAccFgsV2Function_withEpsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name                  = "%s"
  app                   = "default"
  description           = "fuction test"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  enterprise_project_id = "%s"
  func_code             = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccFgsV2Function_createByImage_step_1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "create_with_vpc_access" {
  name        = "%[2]s_1"
  app         = "default"
  handler     = "-"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"

  custom_image {
    url = "%[3]s"
  }

  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_fgs_function" "create_without_vpc_access" {
  name        = "%[2]s_2"
  app         = "default"
  handler     = "-"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"

  custom_image {
    url = "%[3]s"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_BUILD_IMAGE_URL)
}

func testAccFgsV2Function_createByImage_step_2(rName string) string {
	return fmt.Sprintf(`
%[1]s

# Closs the VPC access
resource "huaweicloud_fgs_function" "create_with_vpc_access" {
  name        = "%[2]s_1"
  app         = "default"
  handler     = "-"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"

  custom_image {
    url = "%[3]s"
  }
}

# Open the VPC access
resource "huaweicloud_fgs_function" "create_without_vpc_access" {
  name        = "%[2]s_2"
  app         = "default"
  handler     = "-"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"

  custom_image {
    url = "%[3]s"
  }

  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_BUILD_IMAGE_URL)
}
