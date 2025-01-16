---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_alert_rule_metrics"
description: |-
  Use this data source to get the list of SecMaster alert rule metrics.
---

# huaweicloud_secmaster_alert_rule_metrics

Use this data source to get the list of SecMaster alert rule metrics.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_alert_rule_metrics" "test" {
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

* `metrics` - The metrics value.
