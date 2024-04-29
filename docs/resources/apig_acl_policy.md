---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_acl_policy"
description: ""
---

# huaweicloud_apig_acl_policy

Manages an ACL policy resource within HuaweiCloud.

## Example Usage

### Create an ACL policy with IP control

```hcl
variable "instance_id" {}
variable "policy_name" {}
variable "ip_addresses" {
  type = list(string)
}

resource "huaweicloud_apig_acl_policy" "ip_rule" {
  instance_id = var.instance_id
  name        = var.policy_name
  type        = "PERMIT"
  entity_type = "IP"
  value       = join(",", var.ip_addresses)
}
```

### Create an ACL policy with account control (via domain names)

```hcl
variable "instance_id" {}
variable "policy_name" {}
variable "domain_names" {
  type = list(string)
}

resource "huaweicloud_apig_acl_policy" "domain_rule" {
  instance_id = var.instance_id
  name        = var.policy_name
  type        = "PERMIT"
  entity_type = "DOMAIN"
  value       = join(",", var.domain_names)
}
```

### Create an ACL policy with account control (via domain IDs)

```hcl
variable "instance_id" {}
variable "policy_name" {}
variable "domain_ids" {
  type = list(string)
}

resource "huaweicloud_apig_acl_policy" "domain_id_rule" {
  instance_id = var.instance_id
  name        = var.policy_name
  type        = "PERMIT"
  entity_type = "DOMAIN_ID"
  value       = join(",", var.domain_ids)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the ACL policy is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the ACL
  policy belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the ACL policy.  
  The valid length is limited from `3` to `64`, only English letters, Chinese characters, digits and underscores (_) are
  allowed. The name must start with an English letter or Chinese character.

* `type` - (Required, String) Specifies the type of the ACL policy.
  The valid values are as follows:
  + **PERMIT**: Allow specific IPs or accounts to access API.
  + **DENY**: Forbid specific IPs or accounts to access API.

* `entity_type` - (Required, String, ForceNew) Specifies the entity type of the ACL policy.
  The valid values are as follows:
  + **IP**: This rule is specified to control access to the API for specific IPs.
  + **DOMAIN**: This rule is specified to control access to the API for specific accounts (specified by domain name).
  + **DOMAIN_ID**: This rule is specified to control access to the API for specific accounts (specified by domain ID).

  Changing this will create a new resource.

* `value` - (Required, String) Specifies one or more objects from which the access will be controlled.  
  Separate multiple objects with commas (,).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the ACL policy.

* `updated_at` - The latest update time of the ACL policy.

## Import

ACL Policies can be imported using their `id` and related dedicated instance ID, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_acl_policy.test <instance_id>/<id>
```
