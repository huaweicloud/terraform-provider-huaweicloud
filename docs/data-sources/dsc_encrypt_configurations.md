---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_encrypt_configurations"
description: |-
  Use this data source to get the list of DSC encryption configurations within HuaweiCloud.
---

# huaweicloud_dsc_encrypt_configurations

Use this data source to get the list of DSC encryption configurations within HuaweiCloud.

## Example Usage

### Query all encryption configurations under the specified algorithm type

```hcl
variable "algorithm_type" {}

data "huaweicloud_dsc_encrypt_configurations" "test" {
  algorithm_type = var.algorithm_type
}
```

### Query encryption configurations by name

```hcl
variable "configuration_name" {}
variable "algorithm_type" {}

data "huaweicloud_dsc_encrypt_configurations" "test" {
  algorithm_type     = var.algorithm_type
  configuration_name = var.configuration_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the encryption configurations are located.  
  If omitted, the provider-level region will be used.

* `algorithm_type` - (Required, String) Specifies the type of the encryption algorithm.  
  The valid values are as follows:
  + **AES**
  + **SM4**

* `configuration_name` - (Optional, String) Specifies the name of the encryption configuration.  
  Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configurations` - The list of the encryption configurations.  
  The [configurations](#dsc_encrypt_configurations) structure is documented below.

* `access_permission` - Whether the user has the access permission.

<a name="dsc_encrypt_configurations"></a>
The `configurations` block supports:

* `id` - The ID of the encryption configuration.

* `configuration_name` - The name of the encryption configuration.

* `algorithm_name` - The name of the encryption algorithm.

* `algorithm_type` - The type of the encryption algorithm.

* `enable_rotate` - Whether the key rotation is enabled.

* `encrypt_mode` - The encryption mode.

* `filling_method` - The filling method used for encryption masking.

* `kms_context` - The KMS context information.  
  The [kms_context](#dsc_encrypt_configurations_kms_context) structure is documented below.

* `mask_task_num` - The number of the masking tasks.

* `rotate_period` - The key rotation period, in days.

<a name="dsc_encrypt_configurations_kms_context"></a>
The `kms_context` block supports:

* `kms_key_alias` - The alias of the KMS key.

* `kms_key_id` - The ID of the KMS key.

* `kms_region` - The region where the KMS key is located.
