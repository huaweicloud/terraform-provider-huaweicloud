---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_models"
description: ""
---

# huaweicloud_modelarts_models

Use this data source to get models of ModelArts.

## Example Usage

```hcl
variable "model_name" {}

data "huaweicloud_modelarts_models" "test" {
  name        = var.model_name
  exact_match = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Model name. Fuzzy match is supported. Set **exact_match** to **true** to use exact match.

* `exact_match` - (Optional, String) Whether to use exact match.  
  Set this parameter to **true** to use exact match.

* `version` - (Optional, String) Model version.  

* `status` - (Optional, String) Model status.  
  Value options are as follows:
    + **publishing**: The model is being published.
    + **published**: The model has been published.
    + **failed**: Publishing the model failed.
    + **building**: The image is being created.
    + **building_failed**: Creating an image failed.

* `description` - (Optional, String) The description of the model. Fuzzy match is supported.  

* `workspace_id` - (Optional, String) Workspace ID, which defaults to 0.  

* `model_type` - (Optional, String) Model type, which is used for obtaining models of this type.  
  It can be **TensorFlow**, **MXNet**, **Caffe**, **Spark_MLlib**, **Scikit_Learn**,
  **XGBoost**, **Image**, **PyTorch**, or **Template**.
  Either **model_type** or **not_model_type** can be configured.

* `not_model_type` - (Optional, String) Model type, which is used for obtaining models of types except for this type.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `models` - The list of models.
  The [Models](#Models_Models) structure is documented below.

<a name="Models_Models"></a>
The `Models` block supports:

* `id` - Model ID.

* `name` - Model name.

* `version` - Model version.  

* `model_type` - Model type.  
  It can be **TensorFlow**, **MXNet**, **Caffe**, **Spark_MLlib**, **Scikit_Learn**,
  **XGBoost**, **Image**, **PyTorch**, or **Template**.

* `description` - Model description that consists of 1 to 100 characters.  

* `owner` - User ID of the tenant to which the model belongs.  

* `source_type` - Model source type.  
  This parameter is valid and its value is **auto** only if the model is deployed using ExeML.

* `model_source` - Model source.  
  Value options are as follows:
    + **auto**: ExeML.
    + **algos**: built-in algorithm.
    + **custom**: custom model.

* `install_type` - Deployment types supported by the model.  
  The value can be **real-time**, **edge**, or **batch**.

* `size` - Model size, in bytes.  

* `workspace_id` - Workspace ID.  
  Value 0 indicates the default workspace.

* `status` - Model status.

* `market_flag` - Whether the model is subscribed from AI Gallery.  

* `tunable` - Whether the model can be tuned.
  **true** indicates that the model can be tuned, and **false** indicates not.  

* `publishable_flag` - Whether the model can be published to AI Gallery.  

* `subscription_id` - Model subscription ID.  

* `extra` - Extended parameter.  

* `specification` - Minimum specifications for model deployment.  
  The [specification](#Models_ModelsSpecification) structure is documented below.

<a name="Models_ModelsSpecification"></a>
The `specification` block supports:

* `min_cpu` - Minimal CPU.

* `min_gpu` - Minimal GPU.

* `min_memory` - Minimum memory.

* `min_ascend` - Minimal Ascend.
