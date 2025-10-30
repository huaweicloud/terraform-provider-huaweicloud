---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_dashboard"
description:  |-
  Manages an AOM dashboard resource within HuaweiCloud.
---

# huaweicloud_aom_dashboard

Manages an AOM dashboard resource within HuaweiCloud.

## Example Usage

<!--markdownlint-disable MD033-->
```hcl
variable "dashboard_title " {}
variable "folder_title " {}
variable "dashboard_type " {}

resource "huaweicloud_aom_dashboard" "test" {
  dashboard_title = var.dashboard_title
  folder_title    = var.folder_title
  dashboard_type  = var.dashboard_type

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
}
```
<!-- markdownlint-enable MD033 -->

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `dashboard_title` - (Required, String) Specifies the dashboard title.

* `folder_title` - (Required, String) Specifies the folder title.

* `dashboard_type` - (Required, String) Specifies the dashboard type. It's customized by user.

* `is_favorite` - (Optional, Bool) Specifies whether to favorite the dashboard. Defaults to **false**.

* `dashboard_tags` - (Optional, List) Specifies the dashboard tags. It's an array of map.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the dashboards
  belongs. Defaults to **0**. Changing this parameter will create a new resource.

* `charts` - (Optional, String) Specifies the dashboard charts. It's in json format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

## Import

The AOM dashboard resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_dashboard.test <id>
```
