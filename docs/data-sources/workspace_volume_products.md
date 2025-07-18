---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_volume_products"
description: |-
  Use this data source to get the list of Workspace volume products within HuaweiCloud.
---

# huaweicloud_workspace_volume_products

Use this data source to get the list of Workspace volume products within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_volume_products" "test" {}
```

### Filter by availability zone

```hcl
variable "availability_zone" {}

data "huaweicloud_workspace_volume_products" "test" {
  availability_zone = var.availability_zone
}
```

### Filter by volume type

```hcl
variable "volume_type" {}

data "huaweicloud_workspace_volume_products" "test" {
  volume_type = var.volume_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the volume products are located.  
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the availability zone where the volume products are located.

* `volume_type` - (Optional, String) Specifies the type of volume products.  
  The valid values are as follows:
  + **SATA**: Common I/O disk.
  + **SAS**: High I/O disk.
  + **SSD**: Ultra-high I/O disk.
  + **GPSSD**: General Purpose SSD Disk.
  + **ESSD**: Extreme SSD Disk.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `volume_products` - The list of volume products that matched filter parameters.  
  The [volume_products](#workspace_volume_product_attr) structure is documented below.

<a name="workspace_volume_product_attr"></a>
The `volume_products` block supports:

* `resource_spec_code` - The ID of volume product.

* `volume_type` - The volume type of volume product.

* `volume_product_type` - The product type of volume product.

* `resource_type` - The resource type of volume product.

* `cloud_service_type` - The cloud service type of volume product.

* `domain_ids` - The list of domain IDs that support this volume.

* `names` - The list of volume product name information.  
  The [names](#workspace_volume_product_name) structure is documented below.

* `status` - The status of the volume product.

<a name="workspace_volume_product_name"></a>
The `names` block supports:

* `language` - The language of volume product name.

* `value` - The volume product name in this language.
