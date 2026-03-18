---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_export_acl_rule"
description: |-
  Manages a resource to export CFW ACL rule within HuaweiCloud.
---

# huaweicloud_cfw_export_acl_rule

Manages a resource to export CFW ACL rule within HuaweiCloud.

-> This resource is a one-time action resource used to export ACL rule. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "object_id" {}

resource "huaweicloud_cfw_export_acl_rule" "test" {
  object_id = var.object_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.
  You can get this value from data source `huaweicloud_cfw_firewalls`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
