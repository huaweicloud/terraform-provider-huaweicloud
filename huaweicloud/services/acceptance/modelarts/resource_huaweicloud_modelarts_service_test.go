package modelarts

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

func getModelartsServiceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getService: Query the ModelArts service.
	var (
		getServiceHttpUrl = "v1/{project_id}/services/{id}"
		getServiceProduct = "modelarts"
	)
	getServiceClient, err := cfg.NewServiceClient(getServiceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts Client: %s", err)
	}

	getServicePath := getServiceClient.Endpoint + getServiceHttpUrl
	getServicePath = strings.ReplaceAll(getServicePath, "{project_id}", getServiceClient.ProjectID)
	getServicePath = strings.ReplaceAll(getServicePath, "{id}", state.Primary.ID)

	getServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getServiceResp, err := getServiceClient.Request("GET", getServicePath, &getServiceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ModelartsService: %s", err)
	}

	getServiceRespBody, err := utils.FlattenResponse(getServiceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ModelartsService: %s", err)
	}

	return getServiceRespBody, nil
}

func TestAccModelartsService_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_modelarts_service.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsService_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "infer_type", "real-time"),
					resource.TestCheckResourceAttr(rName, "description", "This is a acc test"),
					resource.TestCheckResourceAttr(rName, "status", "running"),
					resource.TestCheckResourceAttr(rName, "config.0.specification", "modelarts.vm.gpu.p4u8.container"),
					resource.TestCheckResourceAttr(rName, "config.0.instance_count", "1"),
					resource.TestCheckResourceAttr(rName, "config.0.weight", "100"),
					resource.TestCheckResourceAttr(rName, "config.0.envs.a", "1"),
					resource.TestCheckResourceAttr(rName, "config.0.envs.b", "2"),
					resource.TestCheckResourceAttrPair(rName, "additional_properties.0.smn_notification.0.topic_urn",
						"huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttr(rName, "additional_properties.0.smn_notification.0.events.0", "3"),
					resource.TestCheckResourceAttr(rName, "additional_properties.0.log_report_channels.0.type", "LTS"),
					resource.TestCheckResourceAttr(rName, "schedule.0.type", "stop"),
					resource.TestCheckResourceAttr(rName, "schedule.0.duration", "1"),
					resource.TestCheckResourceAttr(rName, "schedule.0.time_unit", "HOURS"),
					resource.TestCheckResourceAttrSet(rName, "owner"),
					resource.TestCheckResourceAttrSet(rName, "access_address"),
					resource.TestCheckResourceAttrSet(rName, "invocation_times"),
					resource.TestCheckResourceAttrSet(rName, "failed_times"),
					resource.TestCheckResourceAttrSet(rName, "is_shared"),
					resource.TestCheckResourceAttrSet(rName, "shared_count"),
					resource.TestCheckResourceAttrSet(rName, "is_free"),
				),
			},
			{
				Config: testModelartsService_basic_update(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a acc test updated"),
					resource.TestCheckResourceAttr(rName, "status", "stopped"),
					resource.TestCheckResourceAttr(rName, "config.0.envs.b", "3"),
					resource.TestCheckResourceAttr(rName, "additional_properties.0.smn_notification.0.events.0", "3"),
					resource.TestCheckResourceAttr(rName, "additional_properties.0.smn_notification.0.events.1", "7"),
					resource.TestCheckResourceAttr(rName, "schedule.0.type", "stop"),
					resource.TestCheckResourceAttr(rName, "schedule.0.duration", "2"),
					resource.TestCheckResourceAttr(rName, "schedule.0.time_unit", "HOURS"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"change_status_to"},
			},
		},
	})
}

func testModelartsService_basic(name, bucketName string) string {
	modelConfig := testModelartsModel_basic(name, bucketName)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "This is a acc test"
}

resource "huaweicloud_modelarts_service" "test" {
  name        = "%[2]s"
  description = "This is a acc test"
  infer_type  = "real-time"

  config {
    specification  = "modelarts.vm.gpu.p4u8.container"
    instance_count = 1
    weight         = 100
    model_id       = huaweicloud_modelarts_model.test.id
    envs = {
      "a" : "1",
      "b" : "2"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = [3]
    }
    log_report_channels {
      type = "LTS"
    }
  }

  schedule {
    type      = "stop"
    duration  = 1
    time_unit = "HOURS"
  }
}
`, modelConfig, name)
}

func testModelartsService_basic_update(name, bucketName string) string {
	modelConfig := testModelartsModel_basic(name, bucketName)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name         = "%[2]s"
  display_name = "This is a acc test"
}

resource "huaweicloud_modelarts_service" "test" {
  name        = "%[2]s"
  description = "This is a acc test updated"
  infer_type  = "real-time"

  config {
    specification  = "modelarts.vm.gpu.p4u8.container"
    instance_count = 1
    weight         = 100
    model_id       = huaweicloud_modelarts_model.test.id
    envs = {
      "a" : "1",
      "b" : "3"
    }
  }

  additional_properties {
    smn_notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = [3, 7]
    }
    log_report_channels {
      type = "LTS"
    }
  }

  schedule {
    type      = "stop"
    duration  = 2
    time_unit = "HOURS"
  }

  change_status_to = "stopped"
}
`, modelConfig, name)
}
