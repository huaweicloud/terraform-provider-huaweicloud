---
subcategory: "Cloud Service Engine (CSE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cse_nacos_namespaces"
description: |-
  Use this data source to query available Nacos namespaces within HuaweiCloud.
---

# huaweicloud_cse_nacos_namespaces

Use this data source to query available Nacos namespaces within HuaweiCloud.

## Example Usage

```hcl
variable "nacos_engine_id" {}
variable "enterprise_project_id" {}

data "huaweicloud_cse_nacos_namespaces" "test" {
  engine_id             = var.nacos_engine_id
  # If the Nacos engine is located under a specific enterprise project, you can only retrieve the list of namespaces
  # under the target engine by specifying the corresponding enterprise project.
  enterprise_project_id = var.enterprise_project_id 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Nacos namespaces are located.  
  If omitted, the provider-level region will be used.

* `engine_id` - (Required, String) Specifies the ID of the Nacos microservice engine to which the namespaces belong.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the Nacos namespaces
  belong.  
  If the Nacos engine belongs to the non-default enterprise project, this parameter is required and is only valid
  for enterprise users.  
  If omitted, the provider-level enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `namespaces` - All queried Nacos namespaces that match the filter parameters.
  The [namespaces](#cse_nacos_namespaces_attr) structure is documented below.

<a name="cse_nacos_namespaces_attr"></a>
The `namespaces` block supports:

* `id` - The ID of the Nacos namespace.

  -> The reserved namespace (**public**) does not have the ID.

* `name` - The name of the Nacos namespace.
