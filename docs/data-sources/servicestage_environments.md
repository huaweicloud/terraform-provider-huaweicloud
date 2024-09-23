---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_environments"
description: |-
  Use this data source to query available environments within HuaweiCloud.
---

# huaweicloud_servicestage_environments

Use this data source to query available environments within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestage_environments" "test" {}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region where the environments are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `environments` - All queried environments.
  The [environments](#servicestage_environments) structure is documented below.

<a name="servicestage_environments"></a>
The `environments` block contains:

* `id` - The ID of the environment.

* `name` - The name of the environment.

* `description` - The description of the environment.

* `deploy_mode` - The deploy mode of the environment.

* `vpc_id` - The VPC ID to which the environment belongs.

* `basic_resources` - The basic resources.
  The [basic_resources](#resources_associated_environment) structure is documented below.

* `optional_resources` - The optional resources.
  The [optional_resources](#resources_associated_environment) structure is documented below.

* `creator` - The creator name.

* `created_at` - The creation time of the environment, in RFC3339 format.

* `updated_at` - The latest update time of the environment, in RFC3339 format.

<a name="resources_associated_environment"></a>
The `basic_resources` and `optional_resources` block contains:

* `id` - The resource ID.

* `type` - The resource type.
