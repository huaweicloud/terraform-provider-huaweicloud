---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_component_used_resources"
description: |-
  Use this data source to query the list of component used resources within HuaweiCloud.
---

# huaweicloud_servicestagev3_component_used_resources

Use this data source to query the component used resources within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_servicestagev3_component_used_resources" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the components are located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `applications` - The application list that component used.
  The [applications](#servicestage_v3_component_used_applications) structure is documented below.

* `enterprise_projects` - The ID list of the enterprise projects that component used.

* `environments` - The environment list that component used.
  The [environments](#servicestage_v3_component_used_environments) structure is documented below.

<a name="servicestage_v3_component_used_applications"></a>
The `applications` block supports:

* `id` - The ID of the application that component used.

* `label` - The name of the application that component used.

<a name="servicestage_v3_component_used_environments"></a>
The `environments` block supports:

* `id` - The ID of the environment that component used.

* `label` - The name of the environment that component used.
