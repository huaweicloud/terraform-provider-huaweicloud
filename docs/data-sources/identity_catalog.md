---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_catalog"
description: |-
  Use this data source to query the service catalog within HuaweiCloud.
---

# huaweicloud_identity_catalog

Use this data source to query the service catalog within HuaweiCloud.

## Example Usage

```hcl
variable "account_name" {}
variable "user_name" {}
variable "password" {}
variable "project_name" {}

resource "huaweicloud_identity_user_token" "test" {
  account_name = var.account_name
  user_name    = var.user_name
  password     = var.password
  project_name = var.project_name
}

data "huaweicloud_identity_catalog" "test" {
  project_token = huaweicloud_identity_user_token.test.token
}
```

## Argument Reference

The following arguments are supported:

* `project_token` - (Required, String) Specifies project scope token for IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `catalog` - Indicates service catalog information list.
  The [catalog](#IdentityCatalog_Catalog) structure is documented below.

<a name="IdentityCatalog_Catalog"></a>
The `catalog` block supports:

* `endpoints` - Indicates service endpoint information.
  The [endpoints](#IdentityCatalog_Enpoints) structure is documented below.

* `id` - Indicates the service id.

* `name` - Indicates the service name.

* `type` - Indicates the service type.

* <a name="IdentityCatalog_Enpoints"></a>
The `endpoints` block supports:

* `id` - Indicates the endpoint id.

* `interface` - Indicates the endpoint interface.

* `region` - Indicates the endpoint region.

* `region_id` - Indicates the endpoint region id.

* `url` - Indicates the endpoint url.
