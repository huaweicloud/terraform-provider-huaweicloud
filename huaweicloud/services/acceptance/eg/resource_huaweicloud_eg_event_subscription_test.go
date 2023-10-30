package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eg/v1/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getEventSubscriptionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.EgV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG v1 client: %s", err)
	}

	return subscriptions.Get(client, state.Primary.ID)
}

func TestAccEventSubscription_basic(t *testing.T) {
	var (
		obj subscriptions.Subscription

		rName              = "huaweicloud_eg_event_subscription.test"
		name               = acceptance.RandomAccResourceName()
		uuidOfficialEG, _  = uuid.GenerateUUID()
		uuidOfficialSMN, _ = uuid.GenerateUUID()
		uuidCustomHTTPS, _ = uuid.GenerateUUID()

		rc = acceptance.InitResourceCheck(rName, &obj, getEventSubscriptionFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEventSubscription_basic_step1(name, uuidOfficialEG, uuidOfficialSMN),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					// Sources check
					resource.TestCheckResourceAttr(rName, "sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "sources.0.provider_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(rName, "sources.0.filter_rule"),
					// Targets check
					resource.TestCheckResourceAttr(rName, "targets.#", "2"),
				),
			},
			{
				Config: testAccEventSubscription_basic_step2(name, uuidOfficialEG, uuidCustomHTTPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					// Sources check
					resource.TestCheckResourceAttr(rName, "sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "sources.0.provider_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(rName, "sources.0.filter_rule"),
					// Targets length check
					resource.TestCheckResourceAttr(rName, "targets.#", "2"),
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

func testAccEventSubscription_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 2

  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s-target-${count.index}"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, count.index)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, count.index), 1)
}

resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_custom_event_channel" "target" {
  count = 2

  name = "%[1]s-target-${count.index}"
}
  
resource "huaweicloud_eg_custom_event_source" "test" {
  count = 2

  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[1]s-${count.index}"
}

resource "huaweicloud_eg_connection" "test" {
  count = 2

  name      = "%[1]s-${count.index}"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test[count.index].id
}

resource "huaweicloud_smn_topic" "test" {
  count = 2
  
  name = "%[1]s-${count.index}"
}
`, name)
}

func testAccEventSubscription_basic_step1(name, uuidOfficialEG, uuidOfficialSMN string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_subscription" "test" {
  // The behavior of parameter 'name' of the event source is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_eg_custom_event_source.test]

  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[2]s"
  description = "Created by acceptance test"

  sources {
    provider_type = "CUSTOM"
    name          = huaweicloud_eg_custom_event_source.test[0].name
    filter_rule   = jsonencode({
      "source": [{
        "op":"StringIn",
        "values":["${huaweicloud_eg_custom_event_source.test[0].name}"]
      }]
    })
  }

  targets {
    id            = "%[3]s"
    provider_type = "OFFICIAL"
    name          = "HC.EG"
    detail_name   = "eg_detail"
    detail        = jsonencode({
      "agency_name": "EG_TARGET_AGENCY",
      "cross_account": false,
      "cross_region": false,
      "target_channel_id": "${huaweicloud_eg_custom_event_channel.target[0].id}",
      "target_project_id": "%[5]s",
      "target_region": "%[6]s"
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }

  targets {
    id            = "%[4]s"
    provider_type = "OFFICIAL"
    name          = "HC.SMN"
    detail_name   = "smn_detail"
    detail        = jsonencode({
      "subject_transform": {
        "type": "CONSTANT",
        "value": "TEST_CONDTANT"
      },
      "urn": "${huaweicloud_smn_topic.test[0].topic_urn}",
      "agency_name": "EG_TARGET_AGENCY",
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }
}
`, testAccEventSubscription_base(name), name, uuidOfficialEG, uuidOfficialSMN,
		acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}

func testAccEventSubscription_basic_step2(name, uuidOfficialEG, uuidCustomHTTPS string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_subscription" "test" {
  // The behavior of parameter 'name' of the event source is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_eg_custom_event_source.test]

  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  name        = "%[2]s"
  description = "Created by acceptance test"

  sources {
    provider_type = "CUSTOM"
    name          = huaweicloud_eg_custom_event_source.test[1].name
    filter_rule   = jsonencode({
      "source": [{
        "op":"StringIn",
        "values":["${huaweicloud_eg_custom_event_source.test[1].name}"]
      }]
    })
  }

  targets {
    id            = "%[3]s"
    provider_type = "OFFICIAL"
    name          = "HC.EG"
    detail_name   = "eg_detail"
    detail        = jsonencode({
      "agency_name": "EG_TARGET_AGENCY",
      "cross_account": false,
      "cross_region": false,
      "target_channel_id": "${huaweicloud_eg_custom_event_channel.target[1].id}",
      "target_project_id": "%[5]s",
      "target_region": "%[6]s"
    })
    transform = jsonencode({
      type  = "ORIGINAL"
      value = ""
    })
  }

  targets {
    id            = "%[4]s"
    provider_type = "CUSTOM"
    name          = "HTTPS"
    connection_id = huaweicloud_eg_connection.test[0].id
    detail_name   = "detail"
    detail        = jsonencode({
      "url": "https://test.com/example",
    })
    transform = jsonencode({
      type  = "VARIABLE"
      value = "{\"name\":\"$.data.name\"}",
      template = "My name is $${name}"
    })
  }
}
`, testAccEventSubscription_base(name), name, uuidOfficialEG, uuidCustomHTTPS,
		acceptance.HW_PROJECT_ID, acceptance.HW_REGION_NAME)
}
