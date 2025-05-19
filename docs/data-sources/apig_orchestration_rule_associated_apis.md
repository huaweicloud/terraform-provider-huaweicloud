---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_orchestration_rule_associated_apis"
description: |-
  Use this data source to get the list of associated APIs under specified orchestration rule within HuaweiCloud.
---

# huaweicloud_apig_orchestration_rule_associated_apis

Use this data source to get the list of associated APIs under specified orchestration rule within HuaweiCloud.

## Example Usage

### Query all associated APIs under a specified orchestration rule

```hcl
variable "instance_id" {}
variable "orchestration_rule_id" {}

data "huaweicloud_apig_orchestration_rule_associated_apis" "test" {
  instance_id = var.instance_id
  rule_id     = var.orchestration_rule_id
}
```

### Query a specified associated API information under a specified orchestration rule

```hcl
variable "instance_id" {}
variable "orchestration_rule_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_orchestration_rule_associated_apis" "test" {
  instance_id = var.instance_id
  rule_id     = var.orchestration_rule_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the orchestration rule is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the orchestration rule belongs.

* `rule_id` - (Required, String) Specifies the ID of the orchestration rule to be queried.

* `api_id` - (Optional, String) Specifies the associated API ID under the orchestration rule.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - All API list that match the filter parameters under a specified orchestration rule.  
  The [apis](#orchestration_rule_associated_apis) structure is documented below.

<a name="orchestration_rule_associated_apis"></a>
The `apis` block supports:

* `api_id` - The ID of the associated API.

* `api_name` - The name of the associated API.

* `req_method` - The request method of the associated API.

* `req_uri` - The request URI of the associated API.

* `auth_type` - The auth type of the associated API.

* `match_mode` - The match mode of the associated API.

* `group_id` - The group ID to which the associated API belongs.

* `group_name` - The group name to which the associated API belongs.

* `attached_time` - The attached time of the associated API.
