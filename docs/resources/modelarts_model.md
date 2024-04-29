---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_model"
description: ""
---

# huaweicloud_modelarts_model

Manages a Modelarts model resource within HuaweiCloud.  

## Example Usage

### Import a model from OBS

```hcl
variable "source_obs_path" {}

resource "huaweicloud_modelarts_model" "test" {
  name            = "demo"
  version         = "0.0.2"
  description     = "This is a demo"
  source_location = var.source_obs_path
  model_type      = "TensorFlow"
  runtime         = "python3.6"
}
```

### Import a model from OBS and override the configuration file

```hcl
variable "source_obs_path" {}

resource "huaweicloud_modelarts_model" "test" {
  name            = "demo"
  version         = "0.0.2"
  description     = "This is a demo"
  source_location = var.source_obs_path
  model_type      = "TensorFlow"
  runtime         = "python3.6"

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  initial_config = jsonencode(
    {
      "model_algorithm" : "object_detection",
      "metrics" : {},
      "apis" : [
        {
          "url" : "/",
          "method" : "post",
          "request" : {
            "Content-type" : "multipart/form-data",
            "data" : {
              "type" : "object",
              "properties" : {
                "images" : {
                  "type" : "file"
                }
              }
            }
          },
          "response" : {
            "Content-type" : "application/json",
            "data" : {
              "type" : "object",
              "properties" : {
                "mnist_result" : {
                  "type" : "array",
                  "item" : [
                    {
                      "type" : "string"
                    }
                  ]
                }
              }
            }
          }
        }
      ]
    }
  )
}
```

### Import a model from template

```hcl
variable "template_id" {}
variable "infer_format" {}
variable "template_obs_path" {}

resource "huaweicloud_modelarts_model" "test" {
  name        = "demo"
  version     = "0.0.1"
  description = "This is a demo"
  model_type  = "Template"

  template {
    template_id  = var.template_id
    infer_format = var.infer_format
    template_inputs {
      input_id = "model_folder"
      input    = var.template_obs_path
    }
  }

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }
}
```

### Import a model trained by a ModelArts training job

```hcl
variable "source_obs_path" {}
variable "execution_code_path" {}
variable "source_job_id" {}


resource "huaweicloud_modelarts_model" "test" {
  name            = "demo_train_model"
  version         = "0.0.1"
  description     = "This is a demo import from train"
  source_location = var.source_obs_path
  source_job_id   = var.source_job_id
  model_type      = "TensorFlow"
  runtime         = "tf1.13-python3.7-cpu"
  execution_code  = var.execution_code_path

  install_type = ["real-time", "batch"]

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  initial_config = jsonencode(
    {
      "model_algorithm" : "image_classification",
      "model_type" : "TensorFlow",
      "metrics" : {
        "f1" : 0.2,
        "recall" : 0,
        "precision" : 0,
        "accuracy" : 0
      },
      "apis" : [
        {
          "url" : "/",
          "method" : "post",
          "request" : {
            "data" : {
              "type" : "object",
              "properties" : {
                "images" : {
                  "type" : "file"
                }
              }
            },
            "Content-type" : "multipart/form-data"
          },
          "response" : {
            "data" : {
              "type" : "object",
              "required" : [
                "predicted_label",
                "scores"
              ],
              "properties" : {
                "predicted_label" : {
                  "type" : "string"
                },
                "scores" : {
                  "type" : "array",
                  "items" : {
                    "type" : "array",
                    "minItems" : 2,
                    "maxItems" : 2,
                    "items" : [
                      {
                        "type" : "string"
                      },
                      {
                        "type" : "number"
                      }
                    ]
                  }
                }
              }
            },
            "Content-type" : "multipart/form-data"
          }
        }
      ],
      "dependencies" : [
        {
          "installer" : "pip",
          "packages" : [
            {
              "package_name" : "numpy",
              "package_version" : "1.17.0",
              "restraint" : "EXACT"
            },
            {
              "package_name" : "h5py",
              "package_version" : "2.8.0",
              "restraint" : "EXACT"
            },
            {
              "package_name" : "Pillow",
              "package_version" : "5.2.0",
              "restraint" : "EXACT"
            },
            {
              "package_name" : "scipy",
              "package_version" : "1.2.1",
              "restraint" : "EXACT"
            },
            {
              "package_name" : "resampy",
              "package_version" : "0.2.1",
              "restraint" : "EXACT"
            },
            {
              "package_name" : "scikit-learn",
              "package_version" : "0.22.2",
              "restraint" : "EXACT"
            }
          ]
        }
      ],
      "runtime" : "tf1.13-python3.7-cpu",
      "model_source" : "algos",
      "tunable" : false
    }
  )
}
```

### Import a model from a container image

```hcl
variable "swr_imag_path" {}

resource "huaweicloud_modelarts_model" "test" {
  name            = "demo_swr"
  version         = "0.0.1"
  description     = "This is a demo import from swr"
  source_location = var.swr_imag_path
  model_type      = "Image"
  source_copy     = true

  install_type = ["real-time", "batch"]


  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  initial_config = jsonencode(
    {
      "protocol" : "https",
      "port" : 90,
      "model_type" : "Image",
      "algorithm" : "unknown_algorithm",
      "health" : {
        "check_method" : "EXEC",
        "command" : "echo 1",
        "period_seconds" : "1",
        "failure_threshold" : "2",
        "initial_delay_seconds" : "1"
      }
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Model name, which consists of 1 to 64 characters.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed.

  Changing this parameter will create a new resource.

* `model_type` - (Required, String, ForceNew) Model type.  
  It can be **TensorFlow**, **MXNet**, **Caffe**, **Spark_MLlib**, **Scikit_Learn**,
  **XGBoost**, **Image**, **PyTorch**, or **Template** read from the configuration file.

  Changing this parameter will create a new resource.
  
* `version` - (Required, String, ForceNew) Model version in the format of Digit.Digit.Digit.  
  Each digit is a one-digit or two-digit positive integer, but cannot start with 0.
  For example, 01.01.01 is not allowed.

  Changing this parameter will create a new resource.

* `source_location` - (Required, String, ForceNew) OBS path where the model is located or the SWR image location.  

  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Model description that consists of 1 to 100 characters.  
  The following special characters cannot be contained: **&!'"<>=**.

  Changing this parameter will create a new resource.

* `model_docs` - (Optional, List, ForceNew) List of model description documents. A maximum of three documents are supported.

  Changing this parameter will create a new resource.
  The [ModelDocs](#ModelartsModel_ModelDocs) structure is documented below.

* `template` - (Optional, List, ForceNew) Configuration items in a template.  
  This parameter is mandatory when `model_type` is set to **Template**.

  Changing this parameter will create a new resource.
  The [Template](#ModelartsModel_Template) structure is documented below.

* `source_copy` - (Optional, String, ForceNew) Whether to enable image replication.  
  This parameter is valid only when `model_type` is set to **Image**.
  Value options are as follows:
    + **true**: Default value, indicating that image replication is enabled.
              After this function is enabled, AI applications cannot be rapidly created, and modifying
              or deleting an image in the SWR source directory will not affect service deployment.
    + **false**: Image replication is not enabled.
                After this function is disabled, AI applications can be rapidly created, but modifying
                or deleting an image in the SWR source directory will affect service deployment.

  Changing this parameter will create a new resource.

* `execution_code` - (Optional, String, ForceNew) OBS path for storing the execution code.  
  The name of the execution code file is consistently to be **customize_service.py**.
  The inference code file must be stored in the model directory.
  This parameter can be left blank. Then, the system will automatically identify the inference
  code in the model directory.

  Changing this parameter will create a new resource.

* `source_job_id` - (Optional, String, ForceNew) ID of the source training job.  
  If the model is generated from a training job, input this parameter for source tracing.
  If the model is imported from a third-party meta model, leave this parameter blank.

  Changing this parameter will create a new resource.

* `source_job_version` - (Optional, String, ForceNew) Version of the source training job.  
  If the model is generated from a training job, input this parameter for source tracing.
  If the model is imported from a third-party meta model, leave this parameter blank.

  Changing this parameter will create a new resource.
  
* `source_type` - (Optional, String, ForceNew) Model source type, which can only be **auto**,
  indicating an ExeML model (model download is not allowed).
  If the model is obtained from a training job, leave this parameter blank.

  Changing this parameter will create a new resource.

* `workspace_id` - (Optional, String, ForceNew) Workspace ID, which defaults to 0.  

  Changing this parameter will create a new resource.

* `install_type` - (Optional, List, ForceNew) Deployment type. Only lowercase letters are supported.
  The value can be **real-time**, **edge**, or **batch**. Default value: [real-time, edge, batch].

  Changing this parameter will create a new resource.
  
* `initial_config` - (Optional, String, ForceNew) The model configuration file describes the model usage,
  computing framework, precision, inference code dependency package, and model API.
  The fields such as `model_algorithm`, `model_type`, `runtime`, `swr_location`, `metrics`, `apis`,
  `dependencies`, and `health` in the configuration file config.json.
  For details, see [Specifications for Writing the Model Configuration File](https://support.huaweicloud.com/intl/en-us/inference-modelarts/inference-modelarts-0056.html)

  Changing this parameter will create a new resource.

* `model_algorithm` - (Optional, String, ForceNew) Model algorithm.  
  If the algorithm is read from the configuration file, this parameter can be left blank.
  The value can be **predict_analysis**, **object_detection**, **image_classification**, or **unknown_algorithm**.

  Changing this parameter will create a new resource.

* `runtime` - (Optional, String, ForceNew) Model runtime environment.  
  Its possible values are determined based on model_type.
  For details, see [Supported AI engines and their runtime](https://support.huaweicloud.com/intl/en-us/inference-modelarts/inference-modelarts-0003.html#section3)

  Changing this parameter will create a new resource.

* `metrics` - (Optional, String, ForceNew) Model precision.  
  If the value is read from the configuration file, this parameter can be left blank.

  Changing this parameter will create a new resource.

* `dependencies` - (Optional, List, ForceNew) Package required for inference code and model.  
  If the package is read from the configuration file, this parameter can be left blank.
  The [Dependency](#ModelartsModel_Dependency) structure is documented below.

  Changing this parameter will create a new resource.
  
<a name="ModelartsModel_ModelDocs"></a>
The `ModelDocs` block supports:

* `doc_url` - (Optional, String, ForceNew) HTTP(S) link of the document.
  Changing this parameter will create a new resource.

* `doc_name` - (Optional, String, ForceNew) Document name, which must start with a letter. Enter 1 to 48 characters.  
  Only letters, digits, hyphens (-), and underscores (_) are allowed.
  Changing this parameter will create a new resource.

<a name="ModelartsModel_Template"></a>
The `Template` block supports:

* `template_id` - (Required, String, ForceNew) ID of the used template.  
  The template has a built-in input and output mode.
  Changing this parameter will create a new resource.

* `template_inputs` - (Required, List, ForceNew) Template input configuration,
  specifying the source path for configuring a model.
  The [TemplateInput](#ModelartsModel_TemplateInput) structure is documented below.
  Changing this parameter will create a new resource.

* `infer_format` - (Optional, String, ForceNew) ID of the input and output mode.  
  When this parameter is used, the input and output mode built in the template does not take effect.
  Changing this parameter will create a new resource.
  
<a name="ModelartsModel_TemplateInput"></a>
The `TemplateInput` block supports:

* `input` - (Required, String, ForceNew) Template input path, which can be a path to an OBS file or directory.  
  When you use a template with multiple input items to create a model,
  if the target paths input_properties specified in the template are the same,
  the OBS directory or OBS file name entered here must be unique to prevent files from being overwritten.
  Changing this parameter will create a new resource.

* `input_id` - (Required, String, ForceNew) Input item ID, which is obtained from template details.
  Changing this parameter will create a new resource.

<a name="ModelartsModel_Dependency"></a>
The `Dependency` block supports:

* `installer` - (Required, String, ForceNew) Installation mode. Only **pip** is supported.
  Changing this parameter will create a new resource.

* `packages` - (Required, List, ForceNew) Collection of dependency packages.
  The [package](#ModelartsModel_package) structure is documented below.
  Changing this parameter will create a new resource.

<a name="ModelartsModel_package"></a>
The `package` block supports:

* `package_version` - (Optional, String, ForceNew) Version of a dependency package.
  If this parameter is left blank, the latest version is installed by default.
  Chinese characters and special characters (&!'"<>=) are not allowed.
  Changing this parameter will create a new resource.

* `package_name` - (Required, String, ForceNew) Name of a dependency package.
  Ensure that the package name is correct and available.
  Chinese characters and special characters (&!'"<>=) are not allowed.
  Changing this parameter will create a new resource.

* `restraint` - (Optional, String, ForceNew) Version restriction, which can be **EXACT**, **ATLEAST**, or **ATMOST**.
  This parameter is mandatory only when package_version is available.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `schema_doc` - Download address of the model schema file.

* `image_address` - Image path generated after model packaging.

* `model_size` - Model size, in bytes.

* `status` - Model status.

* `model_source` - Model source.  
  Value options are as follows:
    + **auto**: ExeML.
    + **algos**: built-in algorithm.
    + **custom**: custom model.

* `tunable` - Whether a model can be tuned.  

* `market_flag` - Whether a model is subscribed from AI Gallery.  

* `publishable_flag` - Whether a model is subscribed from AI Gallery.  

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.

## Import

The modelarts model can be imported using id, e.g.

```bash
$ terraform import huaweicloud_modelarts_model.test 635a2d50-0546-469d-b45d-0204b9ad4f14
```
