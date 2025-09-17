package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getIssueResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetIssueTicketDetail(client, state.Primary.ID)
}

func TestAccIssue_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_issue.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getIssueResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: tesIssue_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "is_start_process_async"),
					resource.TestCheckResourceAttrSet(resourceName, "is_update_null"),
					resource.TestCheckResourceAttrSet(resourceName, "is_return_full_info"),
					resource.TestCheckResourceAttrSet(resourceName, "is_start_process"),
					resource.TestCheckResourceAttrSet(resourceName, "ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "regions_search"),
					resource.TestCheckResourceAttrSet(resourceName, "is_common_issue"),
					resource.TestCheckResourceAttrSet(resourceName, "is_need_change"),
					resource.TestCheckResourceAttrSet(resourceName, "is_enable_suspension"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "operator"),
					resource.TestCheckResourceAttrSet(resourceName, "real_ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "work_flow_status"),
					resource.TestCheckResourceAttrSet(resourceName, "baseline_status"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.#"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.is_deleted"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.main_ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.parent_ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.ticket_path"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.target_value"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.target_type"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.creator"),
					resource.TestCheckResourceAttrSet(resourceName, "sub_tickets.0.operator"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.#"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.is_deleted"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.match_type"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.name_en"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.biz_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.prop_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.model_id"),
					resource.TestCheckResourceAttrSet(resourceName, "meta_data_version"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "is_deleted"),
					resource.TestCheckResourceAttrSet(resourceName, "ticket_type_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"handle_time"},
			},
		},
	})
}

func tesIssue_basic(name string) string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_issue" "test" {
  title                    = "%[2]s"
  description              = "this is description"
  enterprise_project_id    = "0"
  ticket_type              = "issues_type_1000"
  virtual_schedule_type    = "issues_mgmt_virtual_schedule_type_2000"
  regions                  = ["%[3]s"]
  level                    = "issues_level_4000"
  root_cause_cloud_service = huaweicloud_coc_application.test.id
  source                   = "issues_mgmt_associated_type_1000"
  source_id                = huaweicloud_coc_incident.test.id
  found_time               = %[4]v
  issue_contact_person     = "%[5]s"
}
`, testIncident_basic(name), name, acceptance.HW_REGION_NAME, tenMinutesAgo, acceptance.HW_USER_ID)
}
