---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_web_app_service_statistics"
description: |-
  Use this data source to get the statistics of the web service, web application and database.
---

# huaweicloud_hss_asset_web_app_service_statistics

Use this data source to get the statistics of the web service, web application and database.

## Example Usage

```hcl
variable "category" {}
variable "catalogue" {}

data "huaweicloud_hss_asset_web_app_service_statistics" "test" {
  category  = var.category
  catalogue = var.catalogue
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the asset category.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `catalogue` - (Required, String) Specifies the asset type.
  The valid values are as follows:
  + **web-app**: Web application.
  + **web-service**: Web service.
  + **database**: Database.

* `name` - (Optional, String) Specifies the web application, web service or database name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The list of servers.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `name` - The web application, web service or database name.

* `num` - The host number.
