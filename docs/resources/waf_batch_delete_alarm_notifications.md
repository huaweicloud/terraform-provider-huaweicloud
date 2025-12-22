---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_batch_delete_alarm_notifications"
description: |-
  Manages a resource to batch delete the alarm notifications within HuaweiCloud.
---

# huaweicloud_waf_batch_delete_alarm_notifications

Manages a resource to batch delete the alarm notifications within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

-> This resource is only a one-time action resource using to batch delete alarm notifications. Deleting this resource
  will not clear the corresponding request record, but will only remove the resource information from the tf state
  file.

## Example Usage

```hcl
variable "enterprise_project_id" {}
variable "config_id"  {}

resource "huaweicloud_waf_batch_delete_alarm_notifications" "test" {
  enterprise_project_id = var.enterprise_project_id

  alert_notice_configs {
    id = var.config_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `alert_notice_configs` - (Required, List, NonUpdatable) Specifies the alarm notification details.
  The [alert_notice_configs](#alert_notice_configs) structure is documented below.

<a name="alert_notice_configs"></a>
The `alert_notice_configs` block supports:

* `id` - (Required, String, NonUpdatable) Specifies the ID of the alarm notification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
