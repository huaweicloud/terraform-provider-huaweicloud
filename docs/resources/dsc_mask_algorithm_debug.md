---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_mask_algorithm_debug"
description: |-
  Use this resource to test a DSC mask algorithm within HuaweiCloud.
---

# huaweicloud_dsc_mask_algorithm_debug

Use this resource to test a DSC mask algorithm within HuaweiCloud.

-> This resource is a one-time action resource used to test a DSC mask algorithm. Deleting this resource will not
  clear the corresponding test record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "data" {}
variable "parameter" {}

resource "huaweicloud_dsc_mask_algorithm_debug" "test" {
  data           = var.data
  algorithm      = "PRESNM"
  algorithm_type = "MASK_BY_OVERWRITE"
  parameter      = var.parameter
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the mask algorithm to be tested is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `algorithm` - (Required, String, NonUpdatable) Specifies the mask algorithm.  
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

* `algorithm_type` - (Required, String, NonUpdatable) Specifies the type of the mask algorithm.  
  The valid values are as follows:
  + **MASK_BY_HASH**
  + **MASK_BY_ENCRYPT**
  + **MASK_BY_OVERWRITE**
  + **MASK_BY_KEYWORDS_EXCHANGE**
  + **MASK_BY_NULL**
  + **MASK_BY_BRIEF**
  + **MASK_BY_DECRYPT**
  + **MASK_BY_FAKE**

* `data` - (Required, String, NonUpdatable) Specifies the data content to be processed by the mask algorithm.

* `parameter` - (Optional, String, NonUpdatable) Specifies the configuration parameters of the mask algorithm, in JSON
  format.  
  This parameter is **required** only when the `algorithm_type` is set to **MASK_BY_OVERWRITE**, **MASK_BY_KEYWORDS_EXCHANGE**,
  **MASK_BY_BRIEF**, or **MASK_BY_FAKE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `processed_data` - The data content after the mask algorithm processing.
