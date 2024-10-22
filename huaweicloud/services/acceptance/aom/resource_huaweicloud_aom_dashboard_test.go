package aom

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

func getDashboardResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v2/{project_id}/aom/dashboards"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(state),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dashboards: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening dashboards: %s", err)
	}

	jsonPath := fmt.Sprintf("dashboards[?dashboard_id=='%s']|[0]", state.Primary.ID)
	dashboard := utils.PathSearch(jsonPath, listRespBody, nil)
	if dashboard == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dashboard, nil
}

func TestAccDashboard_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	newName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_dashboard.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDashboardResourceFunc,
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
				Config: testDashboard_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dashboard_title", rName),
					resource.TestCheckResourceAttrPair(resourceName, "folder_title", "huaweicloud_aom_dashboards_folder.test", "folder_title"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_type", "dashboard"),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_aom_dashboards_folder.test", "enterprise_project_id"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_tags.0.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "charts"),
				),
			},
			{
				Config: testDashboard_update(newName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "dashboard_title", newName),
					resource.TestCheckResourceAttrPair(resourceName, "folder_title", "huaweicloud_aom_dashboards_folder.test", "folder_title"),
					resource.TestCheckResourceAttr(resourceName, "dashboard_type", "custom"),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "enterprise_project_id",
						"huaweicloud_aom_dashboards_folder.test", "enterprise_project_id"),
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

//nolint:revive
func testDashboard_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_dashboard" "test" {
  depends_on = [huaweicloud_aom_dashboards_folder.test]

  dashboard_title       = "%[2]s"
  folder_title          = huaweicloud_aom_dashboards_folder.test.folder_title
  dashboard_type        = "dashboard"
  is_favorite           = true
  enterprise_project_id = huaweicloud_aom_dashboards_folder.test.enterprise_project_id
  dashboard_tags        = [
    {
      key = "value"
    }
  ]

  charts = jsonencode(
    [
      {
        "definition": {
          "requests": {
            "promql": [
              "label_replace(avg_over_time(actual_workload{version=\"latest\"}[59999ms]),\"__name__\",\"actual_workload\",\"\",\"\")"
            ],
            "copyPromql": [
              "label_replace(avg_over_time(actual_workload{version=\"latest\"}[59999ms]),\"__name__\",\"actual_workload\",\"\",\"\")"
            ],
            "sql": "label_replace(avg_over_time(actual_workload{version=\"latest\"}[59999ms]),\"__name__\",\"actual_workload\",\"\",\"\")"
          },
          "requests_datasource": "prometheus",
          "requests_type": "metric",
          "type": "line",
          "chart_title": "test",
          "config": "{\"chartConfig\":{},\"data\":[{\"namespace\":\"\",\"metricName\":\"actual_workload\",\"alias\":\"\",\"isShowCharts\":true}],\"metricSelectConfig\":{\"metricData\":[{\"code\":\"a\",\"metricName\":\"actual_workload\",\"period\":60000,\"statisticRule\":{\"aggregation_type\":\"average\",\"operator\":\">\",\"thresholdNum\":1},\"aggregate_type\":{\"aggregate_type\":\"by\",\"groupByDimension\":[]},\"triggerRule\":3,\"alarmLevel\":\"Critical\",\"conditionOption\":[{\"id\":\"first\",\"dimension\":\"version\",\"conditionValue\":[{\"name\":\"latest\"}],\"conditionList\":[{\"name\":\"latest\"}],\"addMode\":\"first\",\"conditionCompare\":\"=\",\"regExpress\":null}],\"isShowCharts\":true,\"alias\":\"\",\"query\":\"label_replace({statisticMethod}_over_time(actual_workload{version=\\\"latest\\\"}[59999ms]),\\\"__name__\\\",\\\"actual_workload\\\",\\\"\\\",\\\"\\\")\",\"metircV3OriginData\":{\"metricName\":\"actual_workload\",\"label\":\"actual_workload\",\"namespace\":\"\",\"unit\":\"count\",\"help\":\"\"},\"promql\":\"label_replace(avg_over_time(actual_workload{version=\\\"latest\\\"}[59999ms]),\\\"__name__\\\",\\\"actual_workload\\\",\\\"\\\",\\\"\\\")\",\"transformPromql\":\"label_replace(avg_over_time(actual_workload{version=\\\"latest\\\"}[59999ms]),\\\"__name__\\\",\\\"actual_workload\\\",\\\"\\\",\\\"\\\")\"}],\"mixValue\":{\"mixValue\":\"\",\"statisticRule\":{\"aggregation_type\":\"average\",\"operator\":\">\",\"thresholdNum\":1},\"triggerRule\":3,\"alarmLevel\":\"Critical\",\"isShowCharts\":true,\"alias\":\"\"},\"type\":\"single\"}}",
          "period": 60000,
          "currentTime": "-1.-1.30",
          "promMethod": "avg",
          "statsMethod": "average",
          "operationType": "edit",
          "chart_id": "5k7n66zwew1xoxupkxray7dp"
        },
        "chart_layout": {
          "width": 6,
          "x": 0,
          "y": 0,
          "height": 4
        },
        "chart_id": "5k7n66zwew1xoxupkxray7dp",
        "chart_title": "test"
      }
    ]
  )
}`, testDashboardsFolder_basic(name, false), name)
}

func testDashboard_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_dashboard" "test" {
  depends_on = [huaweicloud_aom_dashboards_folder.test]

  dashboard_title       = "%[2]s"
  folder_title          = huaweicloud_aom_dashboards_folder.test.folder_title
  dashboard_type        = "custom"
  is_favorite           = false
  enterprise_project_id = huaweicloud_aom_dashboards_folder.test.enterprise_project_id
}`, testDashboardsFolder_basic(name, false), name)
}
