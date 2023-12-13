package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/dayu/v1/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkspaces_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	var dayuInstance instances.Instance
	resourceName := "huaweicloud_dataarts_studio_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&dayuInstance,
		getInstanceResourceFunc,
	)

	dataSourceName := "data.huaweicloud_dataarts_studio_workspaces.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)        // enterprise project ID is required for creating instance
			acceptance.TestAccPreCheckChargingMode(t) // the resource instance only supports pre-paid charging mode
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataArtsStudioWorkspaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "workspaces.#"),
				),
			},
		},
	})
}

func testAccDataArtsStudioWorkspaces_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dataarts_studio_workspaces" "all" {
  instance_id = huaweicloud_dataarts_studio_instance.test.id
}
`, testAccInstance_basic(rName))
}

func TestAccDataSourceWorkspaces_filer(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	var dayuInstance instances.Instance
	resourceName := "huaweicloud_dataarts_studio_instance.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&dayuInstance,
		getInstanceResourceFunc,
	)

	byEpsId := "data.huaweicloud_dataarts_studio_workspaces.filter_by_eps_id"
	byEpsIdNotFound := "data.huaweicloud_dataarts_studio_workspaces.filter_by_eps_id_not_found"
	dcbyEpsId := acceptance.InitDataSourceCheck(byEpsId)
	dcbyEpsIdNotFound := acceptance.InitDataSourceCheck(byEpsIdNotFound)

	byName := "data.huaweicloud_dataarts_studio_workspaces.filter_by_name"
	byNameNotFound := "data.huaweicloud_dataarts_studio_workspaces.filter_by_name_not_found"
	dcbyName := acceptance.InitDataSourceCheck(byName)
	dcbyNameNotFound := acceptance.InitDataSourceCheck(byNameNotFound)

	byWorkspaceId := "data.huaweicloud_dataarts_studio_workspaces.filter_by_workspace_id"
	byWorkspaceIdNotFound := "data.huaweicloud_dataarts_studio_workspaces.filter_by_workspace_id_not_found"
	dcbyWorkspaceId := acceptance.InitDataSourceCheck(byWorkspaceId)
	dcbyWorkspaceIdNotFound := acceptance.InitDataSourceCheck(byWorkspaceIdNotFound)

	byCreatedBy := "data.huaweicloud_dataarts_studio_workspaces.filter_by_created_by"
	byCreatedByNotFound := "data.huaweicloud_dataarts_studio_workspaces.filter_by_created_by_not_found"
	dcbyCreatedBy := acceptance.InitDataSourceCheck(byCreatedBy)
	dcbyCreatedByNotFound := acceptance.InitDataSourceCheck(byCreatedByNotFound)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)        // enterprise project ID is required for creating instance
			acceptance.TestAccPreCheckChargingMode(t) // the resource instance only supports pre-paid charging mode
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataArtsStudioWorkspaces_filter(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),

					dcbyEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcbyEpsIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful_not_found", "true"),

					dcbyName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcbyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),

					dcbyWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					dcbyWorkspaceIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful_not_found", "true"),

					dcbyCreatedBy.CheckResourceExists(),
					resource.TestCheckOutput("is_created_by_filter_useful", "true"),
					dcbyCreatedByNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_created_by_filter_useful_not_found", "true"),
				),
			},
		},
	})
}

func testAccDataArtsStudioWorkspaces_filter(rName string) string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

// filter_by_eps_id
data "huaweicloud_dataarts_studio_workspaces" "filter_by_eps_id" {
  instance_id           = huaweicloud_dataarts_studio_instance.test.id
  enterprise_project_id = "%[2]s"
}

data "huaweicloud_dataarts_studio_workspaces" "filter_by_eps_id_not_found" {
  instance_id           = huaweicloud_dataarts_studio_instance.test.id
  enterprise_project_id = "%[3]s"
}

locals {
  filter_result = [for v in data.huaweicloud_dataarts_studio_workspaces.filter_by_eps_id.workspaces[*].name :
                   v == "default"]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result) == 1
}

output "is_eps_id_filter_useful_not_found" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_eps_id_not_found.workspaces) == 0
}

// filter_by_name
data "huaweicloud_dataarts_studio_workspaces" "filter_by_name" {
  instance_id = huaweicloud_dataarts_studio_instance.test.id
  name        = "default"
}

data "huaweicloud_dataarts_studio_workspaces" "filter_by_name_not_found" {
  instance_id = huaweicloud_dataarts_studio_instance.test.id
  name        = "%[4]s"
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_name.workspaces) == 1
}

output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_name_not_found.workspaces) == 0
}

// filter_by_workspace_id
data "huaweicloud_dataarts_studio_workspaces" "filter_by_workspace_id" {
  instance_id  = huaweicloud_dataarts_studio_instance.test.id
  workspace_id = data.huaweicloud_dataarts_studio_workspaces.filter_by_name.workspaces[0].id
}

data "huaweicloud_dataarts_studio_workspaces" "filter_by_workspace_id_not_found" {
  instance_id  = huaweicloud_dataarts_studio_instance.test.id
  workspace_id = "%[3]s"
}

output "is_workspace_id_filter_useful" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_workspace_id.workspaces) == 1
}

output "is_workspace_id_filter_useful_not_found" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_workspace_id_not_found.workspaces) == 0
}

// filter_by_created_by
data "huaweicloud_dataarts_studio_workspaces" "filter_by_created_by" {
  instance_id = huaweicloud_dataarts_studio_instance.test.id
  created_by  = data.huaweicloud_dataarts_studio_workspaces.filter_by_name.workspaces[0].created_by
}

data "huaweicloud_dataarts_studio_workspaces" "filter_by_created_by_not_found" {
  instance_id = huaweicloud_dataarts_studio_instance.test.id
  created_by  = "%[4]s"
}

output "is_created_by_filter_useful" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_created_by.workspaces) == 1
}

output "is_created_by_filter_useful_not_found" {
  value = length(data.huaweicloud_dataarts_studio_workspaces.filter_by_created_by_not_found.workspaces) == 0
}
`, testAccInstance_basic(rName), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randUUID, rName)
}
