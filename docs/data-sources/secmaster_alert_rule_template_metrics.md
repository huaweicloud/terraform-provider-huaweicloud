---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule_template_metrics"
description: |-
  Use this data source to get the metrics of SecMaster alert rule templates.
---

# huaweicloud_secmaster_alert_rule_template_metrics

Use this data source to get the metrics of SecMaster alert rule templates.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_alert_rule_template_metrics" "example" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics_attribute` - The metrics information of alert rule templates in JSON format.
