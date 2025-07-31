---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_service"
description: |-
  Use this data source to get the service information of Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_service

Use this data source to get the service information of Workspace APP within HuaweiCloud.

-> Before using this data source, please ensure that the Workspace service is enabled.

## Example Usage

```hcl
data "huaweicloud_workspace_app_service" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `service_status` - The status of the APP service.
  + **active**
  + **inactive**

* `open_with_ad` - Whether the APP service is connected to AD.

* `tenant_domain_id` - The domain ID to which the APP service belongs.

* `tenant_domain_name` - The domain name to which the APP service belongs.

* `created_at` - The creation time of the Workspace service, in RFC3339 format.
