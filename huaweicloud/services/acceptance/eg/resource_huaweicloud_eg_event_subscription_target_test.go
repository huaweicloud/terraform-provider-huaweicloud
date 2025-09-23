package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/eg"
)

func getEventSubscriptionTargetFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("eg", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	return eg.GetEventSubscriptionTargetById(client, state.Primary.Attributes["subscription_id"], state.Primary.ID)
}

func TestAccEventSubscriptionTarget_basic(t *testing.T) {
	var (
		object interface{}

		resourceName = "huaweicloud_eg_event_subscription_target.test"
		rc           = acceptance.InitResourceCheck(resourceName, &object, getEventSubscriptionTargetFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckFgsFunctionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testEventSubscriptionTarget_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "HC.FunctionGraph"),
					resource.TestCheckResourceAttr(resourceName, "provider_type", "OFFICIAL"),
					resource.TestCheckResourceAttr(resourceName, "retry_times", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testEventSubscriptionTarget_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "retry_times", "13"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEventSubscriptionTargetImportStateFunc(resourceName),
			},
		},
	})
}

func testAccEventSubscriptionTargetImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rsName, rs)
		}

		var (
			subscription_id = rs.Primary.Attributes["subscription_id"]
			target_id       = rs.Primary.ID
		)
		if subscription_id == "" || target_id == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<subscription_id>/<id>',  got '%s/%s'",
				subscription_id, target_id)
		}

		return fmt.Sprintf("%s/%s", subscription_id, target_id), nil
	}
}

func testEventSubscriptionTarget_base(name string) string {
	targetId, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_custom_event_channel" "target" {
  name = "%[1]s-target"
}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[1]s"
}

resource "huaweicloud_eg_event_subscription" "test" {
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
    huaweicloud_eg_custom_event_channel.test,
    huaweicloud_eg_custom_event_channel.target,
  ]

  name        = "%[1]s"
  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  description = "Created by acceptance test"

  sources {
    provider_type = "CUSTOM"
    name          = huaweicloud_eg_custom_event_source.test.name
    filter_rule = jsonencode({
      "source" : [{
        "op" : "StringIn",
        "values" : ["${huaweicloud_eg_custom_event_source.test.name}"]
      }]
    })
  }

  targets {
    id            = "%[2]s"
    provider_type = "OFFICIAL"
    name          = "HC.EG"
    detail_name   = "eg_detail"
    detail = jsonencode({
      "agency_name" : "EG_TARGET_AGENCY",
      "cross_account" : false,
      "cross_region" : false,
      "target_channel_id" : "${huaweicloud_eg_custom_event_channel.target.id}",
      "target_project_id": "%[3]s",
      "target_region": "%[4]s"
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }

  lifecycle {
    ignore_changes = [sources, targets]
  }
}`, name, targetId, acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testEventSubscriptionTarget_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_subscription_target" "test" {
  subscription_id  = huaweicloud_eg_event_subscription.test.id
  name             = "HC.FunctionGraph"
  provider_type    = "OFFICIAL"
  retry_times      = 3
  
  detail = jsonencode({
    "agency_name": "EG_TARGET_AGENCY",
    "invoke_type": "SYNC",
    "urn": "urn:fss:%[3]s:%[2]s:function:default:%[4]s:latest"
  })

  key_transform {
    type = "ORIGINAL"
  }
}
`, testEventSubscriptionTarget_base(name), acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME,
		acceptance.HW_FGS_FUNCTION_NAME)
}

func testEventSubscriptionTarget_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_subscription_target" "test" {
  subscription_id  = huaweicloud_eg_event_subscription.test.id
  name             = "HC.FunctionGraph"
  provider_type    = "OFFICIAL"
  retry_times      = 13
  
  detail = jsonencode({
    "agency_name": "EG_TARGET_AGENCY",
    "invoke_type": "SYNC",
    "urn": "urn:fss:%[3]s:%[2]s:function:default:receiver:latest",
	"terraform_test": "Terraform_Test"
  })

  key_transform {
    type = "ORIGINAL"
  }
}
`, testEventSubscriptionTarget_base(name), acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}
