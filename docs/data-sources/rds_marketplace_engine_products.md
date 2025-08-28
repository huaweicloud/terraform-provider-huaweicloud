---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_marketplace_engine_products"
description: |-
  Use this data source to get the list of RDS marketplace engine products.
---

# huaweicloud_rds_marketplace_engine_products

Use this data source to get the list of RDS marketplace engine products.

## Example Usage

```hcl
variable "bp_domain_id" {}

data "huaweicloud_rds_marketplace_engine_products" "test" {
  bp_domain_id = var.bp_domain_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `bp_domain_id` - (Required, String) Specifies the service provider ID.

* `engine_id` - (Optional, String) Specifies the engine ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `marketplace_engine_products` - Indicates the list of cloud marketplace engine products.

  The [marketplace_engine_products](#marketplace_engine_products_struct) structure is documented below.

<a name="marketplace_engine_products_struct"></a>
The `marketplace_engine_products` block supports:

* `engine_id` - Indicates the engine ID.

* `engine_version` - Indicates the engine version.

* `spec_code` - Indicates the flavor ID.

* `product_id` - Indicates the product ID.

* `bp_name` - Indicates the service provider name.

* `bp_domain_id` - Indicates the service provider ID.

* `instance_mode` - Indicates the instance mode. The value can be:
  + **Single**: Single instance
  + **Ha**: Ha instance
  + **Replica**: Replica instance
  + **All**: All of the above are supported

* `image_id` - Indicates the image ID.

* `user_license_agreement` - Indicates the user license agreement.

* `agreements` - Indicates the agreements.
  The [agreements](#marketplace_engine_products_agreements_struct) structure is documented below.

<a name="marketplace_engine_products_agreements_struct"></a>
 The `agreements` block supports:

* `id` - Indicates the ID of the agreement.

* `name` - Indicates the name of the agreement.

* `language` - Indicates the language of the agreement.

* `version` - Indicates the version of the agreement.

* `provision_url` - Indicates the provision url of the agreement.
