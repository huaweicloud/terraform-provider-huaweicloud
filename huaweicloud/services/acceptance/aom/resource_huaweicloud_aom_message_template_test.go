package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getMessageTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	return aom.GetMessageTemplate(client, state.Primary.Attributes["name"])
}

func TestAccMessageTemplate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_message_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getMessageTemplateResourceFunc,
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
				Config: testMessageTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "locale", "en-us"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "templates.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testMessageTemplate_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "locale", "zh-cn"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "templates.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testMessageTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_message_template" "test" {
  name                  = "%[1]s"
  locale                = "en-us"
  enterprise_project_id = "%[2]s"
  description           = "test"

  templates {
    sub_type = "email"
    topic    = "$${region_name}[$${event_severity}_$${event_type}_$${clear_type}] have a new alert at $${starts_at}."
    content  = <<EOF
Alarm Name:$${event_name};
Alarm ID:$${id};
Occurred:$${starts_at};
Event Severity:$${event_severity};
Alarm Info:$${alarm_info};
Resource Identifier:$${resources_new};
Suggestion:$${alarm_fix_suggestion_zh};
EOF
  }

  templates {
    sub_type = "sms"
    content  = <<EOF
Alarm Name:$${event_name};
Alarm ID:$${id};
Occurred:$${starts_at};
Event Severity:$${event_severity};
Alarm Info:$${alarm_info};
Resource Identifier:$${resources_new};
Suggestion:$${alarm_fix_suggestion_zh};
EOF
  }
}`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testMessageTemplate_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_message_template" "test" {
  name                  = "%[1]s"
  locale                = "zh-cn"
  enterprise_project_id = "%[2]s"

  templates {
    sub_type = "email"
    content  = <<EOF
用户名:  $${domain_name};
上报类型:  $${clear_type};
事件名称:  $${event_name};
区域名称:  $${region_name};
事件级别:  $${event_severity};
事件源:  $${resource_provider};
消息:  $${message};
EOF
  }
}`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccMessageTemplate_lts(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_message_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getMessageTemplateResourceFunc,
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
				Config: testMessageTemplate_lts(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "locale", "en-us"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "source", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "templates.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testMessageTemplate_ltsUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "locale", "zh-cn"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "source", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "templates.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testMessageTemplate_lts(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_message_template" "test" {
  name                  = "%[1]s"
  locale                = "en-us"
  enterprise_project_id = "%[2]s"
  description           = "test"
  source                = "LTS"

  templates {
    sub_type = "email"
    topic    = "$${region_name}[$${event_severity}_$${event_type}_$${clear_type}] have a new alert at $${starts_at}."
    content  = <<EOF
Alarm Name:$${event_name};
Alarm ID:$${id};
Occurred:$${starts_at};
Event Severity:$${event_severity};
Alarm Info:$${alarm_info};
Resource Identifier:$${resources_new};
Suggestion:$${alarm_fix_suggestion_zh};
EOF
  }
}`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testMessageTemplate_ltsUpdate(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_message_template" "test" {
  name                  = "%[1]s"
  locale                = "zh-cn"
  enterprise_project_id = "%[2]s"
  source                = "LTS"

  templates {
    sub_type = "sms"
    content  = <<EOF
用户名:  $${domain_name};
上报类型:  $${clear_type};
事件名称:  $${event_name};
区域名称:  $${region_name};
事件级别:  $${event_severity};
事件源:  $${resource_provider};
消息:  $${message};
EOF
  }
}`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
