---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_authorized_cloud_services"
description: |-
  Use this data source to get the list of authorized cloud services.
---

# huaweicloud_smn_authorized_cloud_services

Use this data source to get the list of authorized cloud services.

## Example Usage

```hcl
data "huaweicloud_smn_authorized_cloud_services" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cloud_services` - Indicates the list of authorized cloud services.

  The [cloud_services](#cloud_services_struct) structure is documented below.

<a name="cloud_services_struct"></a>
The `cloud_services` block supports:

* `name` - Indicates the cloud service name.

* `show_name` - Indicates the display name of a cloud service.
