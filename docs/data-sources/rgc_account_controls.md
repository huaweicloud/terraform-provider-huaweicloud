---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account_controls"
description: |-
  Use this data source to list controls for an managed account in Resource Governance Center.
---

# huaweicloud_rgc_account_controls

Use this data source to list controls for an managed account in Resource Governance Center.

## Example Usage

```hcl
variable "managed_account_id" {}
data "huaweicloud_rgc_account_controls" "test" {
  managed_account_id = var.managed_account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_account_id` - (Required, String) The ID of the managed account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control_summaries` - A list of control summaries. Each summary contains details about a specific control policy.
  The [control_summaries](#control_summaries) structure is documented below.

<a name="control_summaries"></a>
The `control_summaries` block supports:

* `behavior` - The type of control policy. It includes Proactive, Detective, and Preventive control policies.

* `control_identifier` - The identifier of the control policy.

* `control_objective` - The objective of the control policy.

* `description` - The description information of the control policy.

* `guidance` - The necessity of the control policy.

* `implementation` - The service control policy (SCP) or configuration rule.

* `manage_account_id` - The ID of the managed account.

* `name` - The name of the control policy.

* `owner` - The source of the managed account's creation. It includes CUSTOM and RGC.

* `regional_preference` - The regional preference, which can be either regional or global.

* `service` - The service to which the control policy belongs.

* `state` - The enablement status of the control policy.

* `version` - The current version number of the control policy.
