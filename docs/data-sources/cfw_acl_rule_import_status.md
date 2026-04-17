---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_acl_rule_import_status"
description: |-
  Use this data source to get the ACL rule import status within HuaweiCloud.
---

# huaweicloud_cfw_acl_rule_import_status

Use this data source to get the ACL rule import status within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

data "huaweicloud_cfw_acl_rule_import_status" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `object_id` - (Required, String) Specifies the protected object ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The ACL rule import status data.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The import task ID.

* `status` - The import status.  
  The valid values are as follows:
  + `0`: No task.
  + `1`: Task waiting.
  + `2`: Task executing.
  + `3`: Task success.
  + `4`: Task failed.
  + `5`: Task partially successful.
  + `6`: Task completely failed.
