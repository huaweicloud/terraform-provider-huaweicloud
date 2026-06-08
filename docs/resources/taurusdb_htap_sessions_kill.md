---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_sessions_kill"
description: |-
  Use this resource to delete sessions of a TaurusDB HTAP instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_sessions_kill

Use this resource to delete sessions of a TaurusDB HTAP instance within HuaweiCloud.

-> This resource is a one-time action resource for deleting sessions of an HTAP instance. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "htap_instance_id" {}
variable "session_ids" {
  type = list(string)
}

resource "huaweicloud_taurusdb_htap_sessions_kill" "test" {
  instance_id  = var.htap_instance_id
  process_list = var.session_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the HTAP instance ID.

* `process_list` - (Required, List, NonUpdatable) Specifies the session IDs of the HTAP instance to be deleted.
  A maximum of 20 session IDs can be specified at a time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The response result of the session deletion.
  The valid values are as follows:
  + **success**
  + **failed**

* `msg` - The response error information of the session deletion.
