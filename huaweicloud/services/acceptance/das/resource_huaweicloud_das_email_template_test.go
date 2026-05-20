package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/das"
)

func getEmailTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}

	datastoreType := state.Primary.Attributes["datastore_type"]
	return das.GetEmailTemplateById(client, datastoreType, state.Primary.ID)
}

func TestAccEmailTemplate_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_das_email_template.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getEmailTemplateResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSmnSubscribedTopicId(t)
			acceptance.TestAccPrecheckSmnSubscribedTopicUrn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEmailTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "datastore_type", "mysql"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "groups.#", "1"),
					resource.TestCheckResourceAttr(rName, "health_rank.#", "2"),
					resource.TestCheckResourceAttr(rName, "health_rank.0", "dangerous"),
					resource.TestCheckResourceAttr(rName, "health_rank.1", "sub_healthy"),
					resource.TestCheckResourceAttr(rName, "inspection_time", "00:00-00:00"),
					resource.TestCheckResourceAttr(rName, "send_time", "08:00-10:00"),
					resource.TestCheckResourceAttr(rName, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(rName, "email", "test@example.com"),
					resource.TestCheckResourceAttrSet(rName, "user_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAccEmailTemplate_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "email", ""),
					resource.TestCheckResourceAttr(rName, "groups.#", "2"),
					resource.TestCheckResourceAttr(rName, "health_rank.#", "4"),
					resource.TestCheckResourceAttr(rName, "health_rank.0", "dangerous"),
					resource.TestCheckResourceAttr(rName, "health_rank.1", "sub_healthy"),
					resource.TestCheckResourceAttr(rName, "health_rank.2", "healthy"),
					resource.TestCheckResourceAttr(rName, "health_rank.3", "high_risk"),
					resource.TestCheckResourceAttr(rName, "inspection_time", "12:00-12:00"),
					resource.TestCheckResourceAttr(rName, "send_time", "10:00-12:00"),
					resource.TestCheckResourceAttr(rName, "topic", acceptance.HW_SMN_SUBSCRIBED_TOPIC_ID),
					resource.TestCheckResourceAttr(rName, "topic_urn", acceptance.HW_SMN_SUBSCRIBED_TOPIC_URN),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEmailTemplateImportIdFunc(rName),
				ImportStateVerifyIgnore: []string{
					"email",
					"topic",
					"topic_urn",
					"obs_bucket_name",
				},
			},
		},
	})
}

func testAccEmailTemplateImportIdFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["datastore_type"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["datastore_type"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["datastore_type"], rs.Primary.ID), nil
	}
}

func testAccEmailTemplate_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_das_instance_group" "test" {
  count = 2  

  datastore_type = "mysql"
  group_name     = "%[1]s-${count.index}"
  description    = "Created by terraform script"
}

locals {
  group_ids = huaweicloud_das_instance_group.test[*].id
}
`, name)
}

func testAccEmailTemplate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_email_template" "test" {
  datastore_type  = "mysql"
  name            = "%[2]s"
  groups          = slice(local.group_ids, 0, 1)
  health_rank     = ["dangerous", "sub_healthy"]
  inspection_time = "00:00-00:00"
  send_time       = "08:00-10:00"
  time_zone       = "Asia/Shanghai"
  email           = "test@example.com"
}
`, testAccEmailTemplate_base(name), name)
}

func testAccEmailTemplate_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_email_template" "test" {
  datastore_type  = "mysql"
  name            = "%[2]s_update"
  groups          = local.group_ids
  health_rank     = ["dangerous", "sub_healthy", "healthy", "high_risk"]
  inspection_time = "12:00-12:00"
  send_time       = "10:00-12:00"
  time_zone       = "Asia/Shanghai"
  topic           = "%[3]s"
  topic_urn       = "%[4]s"
}
`, testAccEmailTemplate_base(name), name, acceptance.HW_SMN_SUBSCRIBED_TOPIC_ID, acceptance.HW_SMN_SUBSCRIBED_TOPIC_URN)
}
