---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_acl_policy_associate"
description: ""
---

# huaweicloud_apig_acl_policy_associate

Use this resource to bind the APIs to the ACL policy within HuaweiCloud.

-> An ACL policy can only create one `huaweicloud_apig_acl_policy_associate` resource.

## Example Usage

```hcl
variable "instance_id" {}
variable "policy_id" {}
variable "api_publish_ids" {
  type = list(string)
}

resource "huaweicloud_apig_acl_policy_associate" "test" {
  instance_id = var.instance_id
  policy_id   = var.policy_id
  publish_ids = var.api_publish_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ACL policy and the APIs are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the APIs and the
  ACL policy belong.  
  Changing this will create a new resource.

* `policy_id` - (Required, String, ForceNew) Specifies the ACL Policy ID for APIs binding.  
  Changing this will create a new resource.

* `publish_ids` - (Required, List) Specifies the publish IDs corresponding to the APIs bound by the ACL policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Resource ID. The format is `<instance_id>/<policy_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
* `update` - Default is 3 minutes.
* `delete` - Default is 3 minutes.

## Import

Associate resources can be imported using their `policy_id` and the APIG dedicated instance ID to which the policy
belongs, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_acl_policy_associate.test <instance_id>/<policy_id>
```
