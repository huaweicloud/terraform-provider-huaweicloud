package rms

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

func getRecorderResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	var (
		region             = acceptance.HW_REGION_NAME
		getRecorderHttpUrl = "v1/resource-manager/domains/{domain_id}/tracker-config"
		getRecorderProduct = "rms"
	)
	getRecorderClient, err := cfg.NewServiceClient(getRecorderProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS client: %s", err)
	}

	getRecorderPath := getRecorderClient.Endpoint + getRecorderHttpUrl
	getRecorderPath = strings.ReplaceAll(getRecorderPath, "{domain_id}", cfg.DomainID)

	getRecorderOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRecorderResp, err := getRecorderClient.Request("GET", getRecorderPath, &getRecorderOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS recorder: %s", err)
	}
	return utils.FlattenResponse(getRecorderResp)
}

func TestAccRecorder_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_rms_resource_recorder.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRecorderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Some test cases in RMS require the resource recorder to be enabled.
			// Skip this test case during batch execution to ensure other test cases execute smoothly.
			acceptance.TestAccPreCheckRMSResourceRecorder(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testRecorder_with_obs_partial(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "agency_name", "rms_tracker_agency"),
					resource.TestCheckResourceAttr(rName, "selector.0.all_supported", "false"),
					resource.TestCheckResourceAttr(rName, "selector.0.resource_types.#", "8"),
					resource.TestCheckResourceAttrSet(rName, "obs_channel.0.region"),
					resource.TestCheckResourceAttrPair(rName, "obs_channel.0.bucket",
						"huaweicloud_obs_bucket.test", "bucket"),
				),
			},
			{
				Config: testRecorder_with_obs_all(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "agency_name", "rms_tracker_agency"),
					resource.TestCheckResourceAttr(rName, "selector.0.all_supported", "true"),
					resource.TestCheckResourceAttr(rName, "selector.0.resource_types.#", "0"),
				),
			},
			{
				Config: testRecorder_with_smn(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "agency_name", "rms_tracker_agency"),
					resource.TestCheckResourceAttr(rName, "selector.0.all_supported", "true"),
					resource.TestCheckResourceAttrSet(rName, "smn_channel.0.region"),
					resource.TestCheckResourceAttrSet(rName, "smn_channel.0.project_id"),
					resource.TestCheckResourceAttrPair(rName, "smn_channel.0.topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
				),
			},
			{
				Config: testRecorder_with_all(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "smn_channel.0.topic_urn",
						"huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "obs_channel.0.bucket",
						"huaweicloud_obs_bucket.test", "bucket"),
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

func testRecorder_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  storage_class = "STANDARD"
  acl           = "private"
  force_destroy = true

  tags = {
    env = "rms_recorder_channel"
    key = "value"
  }
}

resource "huaweicloud_smn_topic" "test" {
  name         = "%[1]s"

  tags = {
    env = "rms_recorder_channel"
    key = "value"
  }
}
`, name)
}

func testRecorder_with_obs_partial(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported  = false

    resource_types = [
      "vpc.vpcs", "rds.instances", "dms.kafkas", "dms.rabbitmqs", "dms.queues",
      "config.trackers", "config.policyAssignments", "config.conformancePacks",
    ]
  }

  obs_channel {
    bucket = huaweicloud_obs_bucket.test.id
    region = "%s"
  }
}
`, testRecorder_base(name), acceptance.HW_REGION_NAME)
}

func testRecorder_with_obs_all(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported = true
  }

  obs_channel {
    bucket = huaweicloud_obs_bucket.test.id
    region = "%s"
  }
}
`, testRecorder_base(name), acceptance.HW_REGION_NAME)
}

func testRecorder_with_smn(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported = true
  }

  smn_channel {
    topic_urn = huaweicloud_smn_topic.test.topic_urn
  }
}
`, testRecorder_base(name))
}

func testRecorder_with_all(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported = true
  }

  obs_channel {
    bucket = huaweicloud_obs_bucket.test.id
    region = "%s"
  }
  smn_channel {
    topic_urn = huaweicloud_smn_topic.test.topic_urn
  }
}
`, testRecorder_base(name), acceptance.HW_REGION_NAME)
}
