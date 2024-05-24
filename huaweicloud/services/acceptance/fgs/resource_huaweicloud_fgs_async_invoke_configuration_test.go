package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getAsyncInvokeConfigFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph v2 client: %s", err)
	}
	return function.GetAsyncInvokeConfig(c, state.Primary.ID)
}

func TestAccAsyncInvokeConfig_basic(t *testing.T) {
	var cfg function.AsyncInvokeConfig
	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_fgs_async_invoke_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&cfg,
		getAsyncInvokeConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please read the instructions carefully before use to ensure sufficient permissions.
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAsyncInvokeConfig_basic_step1(name, acceptance.HW_FGS_AGENCY_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "max_async_event_age_in_seconds", "3500"),
					resource.TestCheckResourceAttr(rName, "max_async_retry_attempts", "2"),
					resource.TestCheckResourceAttr(rName, "on_success.0.destination", "OBS"),
					resource.TestCheckResourceAttrSet(rName, "on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "on_failure.0.destination", "SMN"),
					resource.TestCheckResourceAttrSet(rName, "on_failure.0.param"),
					resource.TestCheckResourceAttr(rName, "enable_async_status_log", "true"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccAsyncInvokeConfig_basic_step2(name, acceptance.HW_FGS_AGENCY_NAME),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "function_urn",
						"huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "max_async_event_age_in_seconds", "4000"),
					resource.TestCheckResourceAttr(rName, "max_async_retry_attempts", "3"),
					resource.TestCheckResourceAttr(rName, "on_success.0.destination", "DIS"),
					resource.TestCheckResourceAttrSet(rName, "on_success.0.param"),
					resource.TestCheckResourceAttr(rName, "on_failure.0.destination", "FunctionGraph"),
					resource.TestCheckResourceAttrSet(rName, "on_failure.0.param"),
					resource.TestCheckResourceAttr(rName, "enable_async_status_log", "false"),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
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

func testAccAsyncInvokeConfig_basic_step1(name, agency string) string {
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
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "%[2]s"
}

resource "huaweicloud_fgs_async_invoke_configuration" "test" {
  function_urn                   = huaweicloud_fgs_function.test.urn
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
`, name, agency)
}

func testAccAsyncInvokeConfig_basic_step2(name, agency string) string {
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
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "e42a37a22f4988ba7a681e3042e5c7d13c04e6c1"
  agency      = "%[3]s"
}

resource "huaweicloud_fgs_async_invoke_configuration" "test" {
  function_urn                   = huaweicloud_fgs_function.test.urn
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
`, common.TestBaseNetwork(name), name, agency)
}
