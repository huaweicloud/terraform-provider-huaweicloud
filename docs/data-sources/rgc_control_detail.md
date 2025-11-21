---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_control_detail"
description: |-
  Use this data source to get control detail in Resource Governance Center.
---

# huaweicloud_rgc_control_detail

Use this data source to get control detail in Resource Governance Center.

## Example Usage

```hcl
variable control_id {}

data "huaweicloud_rgc_control_detail" "test" {
  control_id = var.control_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Required, String) The ID of the control policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `aliases` - The aliases of the control policy.

* `behavior` - The type of control policy. Includes Proactive, Detective, and Preventive.

* `control_objective` - The objective of the control policy.

* `framework` - The framework from which the governance policy originates.

* `guidance` - The necessity of the control policy.

* `identifier` - The ID of the control policy.

* `implementation` - Service control policy (SCP), configuration rules.

* `owner` - The source of the managed account's creation, including CUSTOM and RGC.

* `release_date` - The release date of the control policy.

* `resource` - The resources governed.

* `service` - The service to which the control policy belongs.

* `severity` - The severity of the control policy.

* `version` - The version of the control policy.

* `artifacts` - Information about the artifacts of the control policy.

  The [artifacts](#artifacts) structure is documented below.

<a name="artifacts"></a>
The `artifacts` block supports:

* `type` - The type of policy.

* `content` - Information about the content of the control policy.

  The [content](#content) structure is documented below.

<a name="content"></a>
The `content` block supports:

* `ch` - The Chinese content of the control policy.

* `en` - The English content of the control policy.
