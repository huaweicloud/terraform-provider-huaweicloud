---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_ip_blacklist_retry"
description: |-
  Manages a resource to retry the failed IP blacklist import within HuaweiCloud.
---

# huaweicloud_cfw_ip_blacklist_retry

Manages a resource to retry the failed IP blacklist import within HuaweiCloud.

-> 1. This resource is a one-time action resource used to retry the failed IP blacklist import. Deleting this resource
  will not clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. When importing the IP blacklist fails, you can use this resource to retry the failed IP blacklist import.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "name" {}

resource "huaweicloud_cfw_ip_blacklist_retry" "test" {
  fw_instance_id = var.fw_instance_id
  name           = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `name` - (Required, String, NonUpdatable) Specifies the name of the failed IP blacklist to retry.  
  The value can be obtained using the `data_stource_huaweicloud_cfw_ic_blacklist` data source.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.
