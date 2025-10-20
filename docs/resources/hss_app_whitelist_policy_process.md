---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_whitelist_policy_process"
description: |-
  Manages an app whitelist policy process status marking task within HuaweiCloud.
---

# huaweicloud_hss_app_whitelist_policy_process

Manages an app whitelist policy process status marking task within HuaweiCloud.

-> This resource is only a one-time action resource for HSS app whitelist policy process status marking. Deleting
   this resource will not affect the marking status, but will only remove the resource information from the tf
   state file.

## Example Usage

```hcl
variable "policy_id" {}
variable "process_hash1" {}
variable "process_hash2" {}

resource "huaweicloud_hss_app_whitelist_policy_process" "test" {
  policy_id             = var.policy_id
  process_status        = "trust"
  enterprise_project_id = "0"

  process_hash_list = [
    var.process_hash1,
    var.process_hash2,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the whitelist policy.

* `process_status` - (Required, String, NonUpdatable) Specifies the trust status of the process.
  The valid values are as follows:
  + **trust**: Trusted
  + **suspicious**: Suspicious
  + **malicious**: Malicious
  + **unknown**: Unknown

* `process_hash_list` - (Required, List, NonUpdatable) Specifies the list of process hash values to mark.
  The process hash values must be valid hash strings.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the ID of the enterprise project that the server
  belongs to. The value **0** indicates the default enterprise project. To query servers in all enterprise projects,
  set this parameter to **all_granted_eps**. If you have only the permission on an enterprise project, you need to
  transfer the enterprise project ID to query the server in the enterprise project.
  Otherwise, an error is reported due to insufficient permission.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
