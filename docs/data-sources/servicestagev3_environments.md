---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_environments"
description: |-
  Use this data source to query the list of environments within HuaweiCloud.
---

# huaweicloud_servicestagev3_environments

Use this data source to query the list of environments within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestagev3_environments" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the environments are located.  
  If omitted, the provider-level region will be used.

* `environment_id` - (Optional, String) Specifies the ID of the environment to be queried.

* `name` - (Optional, String) Specifies the name of the environment to be queried.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the environments belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `environments` - All environments that match the filter parameters.  
  The [environments](#servicestage_v3_environments) structure is documented below.

<a name="servicestage_v3_environments"></a>
The `environments` block supports:

* `id` - The environment ID.

* `name` - The environment name.

* `description` - The description of the environment.

* `enterprise_project_id` - The ID of the enterprise project to which the environment belongs.

* `deploy_mode` - The deploy mode of the environment.
  + **container**
  + **virtualmachine**
  + **mixed**

* `vpc_id` - The ID of the VPC to which the environment belongs.

* `version` - The version number of the environment.

* `tags` - The key/value pairs to associate with the environment.

* `creator` - The creator name of the environment.

* `created_at` - The creation time of the environment, in RFC3339 format.

* `updated_at` - The latest update time of the environment, in RFC3339 format.
