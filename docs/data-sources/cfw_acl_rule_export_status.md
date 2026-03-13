---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_acl_rule_export_status"
description: |-
  Use this data source to get the ACL rule export status within HuaweiCloud.
---

# huaweicloud_cfw_acl_rule_export_status

Use this data source to get the ACL rule export status within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_acl_rule_export_status" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The ACL rule export status data.

  The [data](#cfw_acl_rule_export_status_data) structure is documented below.

<a name="cfw_acl_rule_export_status_data"></a>
The `data` block supports:

* `id` - The protected object ID.

* `status` - The export status.  
  The valid values are as follows:
  + **0**: No task.
  + **1**: Task waiting.
  + **2**: Task executing.
  + **3**: Task success.
  + **4**: Task failed.
