package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getIncidentResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetIncident(client, state.Primary.ID)
}

func TestAccResourceIncident_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_incident.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getIncidentResourceFunc,
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
				Config: testIncident_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "incident_title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project.0"),
					resource.TestCheckResourceAttrSet(resourceName, "current_cloud_service.0"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_level"),
					resource.TestCheckResourceAttrSet(resourceName, "is_service_interrupt"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_type"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_title"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_description"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_source"),
					resource.TestCheckResourceAttrSet(resourceName, "incident_assignee.0"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.#"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.filed_key"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.enum_key"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(resourceName, "enum_data_list.0.name_en"),
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

func testIncident_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_incident" "test" {
  incident_level        = "level_50"
  is_service_interrupt  = false
  incident_type         = "inc_type_p_security_issues"
  incident_title        = "%[2]s"
  incident_source       = "incident_source_forwarding"
  creator               = "%[3]s"
  incident_assignee     = ["%[3]s"]
  enterprise_project    = ["0"]
  current_cloud_service = [huaweicloud_coc_application.test.id]
  incident_description  = "%[2]s"
}`, testAccApplication_basic(name), name, acceptance.HW_USER_ID)
}
