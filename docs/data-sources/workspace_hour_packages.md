---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_hour_packages"
description: |-
  Use this data source to get the list of Workspace hour packages within Huaweicloud.
---

# huaweicloud_workspace_hour_packages

Use this data source to get the list of Workspace hour packages within Huaweicloud.

## Example Usage

```hcl
data "huaweicloud_workspace_hour_packages" "test" {}
```

### Filter by desktop resource spec code

```hcl
variable "desktop_resource_spec_code" {}

data "huaweicloud_workspace_hour_packages" "test" {
  desktop_resource_spec_code = var.desktop_resource_spec_code
}
```

### Filter by resource spec code

```hcl
variable "resource_spec_code" {}

data "huaweicloud_workspace_hour_packages" "test" {
  resource_spec_code = var.resource_spec_code
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the hour packages.  
  If omitted, the provider-level region will be used.

* `desktop_resource_spec_code` - (Optional, String) Specifies the specification code of desktop resource to be queried.

* `resource_spec_code` - (Optional, String) Specifies the specification code of hour package to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hour_packages` - The list of hour package information that matched filter parameters.  
  The [hour_packages](#workspace_desktop_hour_packages) structure is documented below.

<a name="workspace_desktop_hour_packages"></a>
The `hour_packages` block supports:

* `cloud_service_type` - The type of cloud service.

* `resource_type` - The type of resource.

* `resource_spec_code` - The ID of hour package.

* `desktop_resource_spec_code` - The ID of desktop resource.

* `descriptions` - The descriptions of hour package.  
  The [descriptions](#workspace_desktop_hour_package_descriptions) structure is documented below.

* `package_duration` - The duration of hour package.

* `domain_ids` - The list of domain IDs supported by the hour package.

* `status` - The status of hour package.

<a name="workspace_desktop_hour_package_descriptions"></a>
The `descriptions` block supports:

* `zh_cn` - The Chinese description of hour package.

* `en_us` - The English description of hour package.
