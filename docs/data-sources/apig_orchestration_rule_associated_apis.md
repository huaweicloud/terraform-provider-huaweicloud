---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_orchestration_rule_associated_apis"
description: |-
  Use this data source to get API list associated with the specified orchestration rule within HuaweiCloud.
---

# huaweicloud_apig_orchestration_rule_associated_apis

Use this data source to get API list associated with the specified orchestration rule within HuaweiCloud.

## Example Usage

### Querying all APIs associated with the specified orchestration rule

```hcl
variable "instance_id" {}
variable "orchestration_rule_id" {}

data "huaweicloud_apig_orchestration_rule_associated_apis" "test" {
  instance_id = var.instance_id
  rule_id     = var.orchestration_rule_id
}
```

### Querying API associated with the orchestration rule using specified API ID

```hcl
variable "instance_id" {}
variable "orchestration_rule_id" {}
variable "api_id" {}

data "huaweicloud_apig_orchestration_rule_associated_apis" "test" {
  instance_id = var.instance_id
  rule_id     = var.orchestration_rule_id
  api_id      = var.api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the orchestration rule belongs.

* `rule_id` - (Required, String) Specifies the ID of the orchestration rule.

* `api_id` - (Optional, String) Specifies the ID of the API associated with the orchestration rule.

* `api_name` - (Optional, String) Specifies the name of the API associated with the orchestration rule,
  fuzzy matching is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `apis` - All associated APIs that match the filter parameters.  
  The [apis](#orchestration_rule_associated_apis) structure is documented below.

<a name="orchestration_rule_associated_apis"></a>
The `apis` block supports:

* `api_id` - The ID of the API.

* `api_name` - The name of the API.

* `req_uri` - The request address of the API.

* `req_method` - The request method of the API.

* `auth_type` - The security authentication mode of the API request.

* `match_mode` - The matching mode of the API.

* `group_id` - The ID of the API group to which the API belongs.

* `group_name` - The name of the API group to which the API belongs.

* `attached_time` - The time when the orchestration rule is associated with the API, in RFC3339 format.
