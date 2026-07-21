---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_mask_algorithm"
description: |-
  Manages a DSC mask algorithm resource within HuaweiCloud.
---

# huaweicloud_dsc_mask_algorithm

Manages a DSC mask algorithm resource within HuaweiCloud.

## Example Usage

### Create a character mask algorithm

```hcl
variable "algorithm_name" {}

resource "huaweicloud_dsc_mask_algorithm" "test" {
  algorithm_name = var.algorithm_name
  algorithm      = "PRESNM"
  algorithm_type = "MASK_BY_OVERWRITE"
  category       = "BUILT_SELF"
  parameter      = jsonencode({
    type   = "CHAR"
    first  = 6
    second = 4
    method = "*"
  })
}
```

### Create a keyword replace mask algorithm

```hcl
variable "algorithm_name" {}

resource "huaweicloud_dsc_mask_algorithm" "test" {
  algorithm_name = var.algorithm_name
  algorithm      = "KEYWORD"
  algorithm_type = "MASK_BY_KEYWORDS_EXCHANGE"
  category       = "BUILT_SELF"
  parameter      = jsonencode({
    key    = "keyword"
    target = "replace content"
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the mask algorithm is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `algorithm_name` - (Required, String) Specifies the name of the mask algorithm.  
  The maximum length is `255` characters.  
  Only Chinese Characters, letters, digits, underscores (_), or hyphens (-) are allowed.  
  The name must be unique.

* `algorithm` - (Required, String) Specifies the encryption mask algorithm type.  
  The valid values are as follows:
  + **SHA256**
  + **SHA512**
  + **PRESNM**
  + **MASKNM**
  + **PRESXY**
  + **MASKXY**
  + **SYMBOL**
  + **KEYWORD**
  + **NULL**
  + **EMPTY**
  + **DATE**
  + **NUMERIC**
  + **AES**
  + **EMBED**
  + **SM4**
  + **DECRYPT**
  + **FAKE**

* `algorithm_type` - (Required, String) Specifies the type of the mask algorithm.  
  The valid values are as follows:
  + **MASK_BY_HASH**
  + **MASK_BY_ENCRYPT**
  + **MASK_BY_OVERWRITE**
  + **MASK_BY_KEYWORDS_EXCHANGE**
  + **MASK_BY_NULL**
  + **MASK_BY_BRIEF**
  + **MASK_BY_DECRYPT**
  + **MASK_BY_FAKE**

* `category` - (Required, String, NonUpdatable) Specifies the category of the mask algorithm.  
  The valid values are as follows:
  + **BUILT_SELF**

* `parameter` - (Required, String) Specifies the configuration parameters of the mask algorithm, in JSON format.

* `data` - (Optional, String) Specifies the data content processed by the mask algorithm.  
  This parameter is available when the algorithm type is `MASK_BY_OVERWRITE`.

* `processed_data` - (Optional, String) Specifies the data content after masking.  
  This parameter is available when the algorithm type is `MASK_BY_OVERWRITE`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_dsc_mask_algorithm.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `algorithm_type`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the mask algorithm, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dsc_mask_algorithm" "test" {
  ...

  lifecycle {
    ignore_changes = [
      algorithm_type,
    ]
  }
}
```
