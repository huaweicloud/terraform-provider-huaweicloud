---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_model_templates"
description: ""
---

# huaweicloud_modelarts_model_templates

Use this data source to get model templates of ModelArts.

## Example Usage

```hcl
data "huaweicloud_modelarts_model_templates" "test" {
  type = "Classification"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) The type of model. The valid values are **Classification** and **Common**.  

* `engine` - (Optional, String) The AI engine.  
  The valid values are **Caffe**, **Caffe1.0 CPU**, **Caffe1.0 GPU**, **MXNet**, **MXNet1.2.1**,
   **MindSpore**, **PyTorch**, **PyTorch1.0**, **TensorFlow**, and **TensorFlow1.8**.

* `environment` - (Optional, String) Model runtime environment.  
  The valid values are **ascend-arm-py2.7**, **python2.7**, and **python3.6**.

* `keyword` - (Optional, String) Keywords to search in name or description. Fuzzy match is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of model templates.
  The [templates](#ModelTemplate_templates) structure is documented below.

<a name="ModelTemplate_templates"></a>
The `templates` block supports:

* `id` - Template ID.

* `name` - Template name.

* `description` - Template description.  

* `arch` - Architecture type. The valid values are **X86_64** and **AARCH64**.

* `type` - The type of model. The valid values are **Classification** and **Common**.  

* `engine` - The AI engine.  
  The valid values are **Caffe**, **MXNet**, **MindSpore**, **PyTorch**, and **TensorFlow**.

  -> It may not be equal to the filter argument `engine`, it does not have a version information suffix.

* `environment` - Model runtime environment.  
  The valid values are **aarch64**, **python2.7**, **python2.7-cpu**, **python2.7-gpu**, **python3.6**,
   **python3.6-gpu**, and **python3.6-gpu**.

  -> It may not be equal to the filter argument `environment`.

* `template_docs` - List of template description documents.  
  The [template_docs](#ModelTemplate_TemplatestemplateDocs) structure is documented below.

* `template_inputs` - List of input parameters for the model.
  The [template_inputs](#ModelTemplate_TemplatestemplateInputs) structure is documented below.

<a name="ModelTemplate_TemplatestemplateDocs"></a>
The `template_docs` block supports:

* `doc_url` - HTTP(S) link of the document.

* `doc_name` - Document name.  

<a name="ModelTemplate_TemplatestemplateInputs"></a>
The `template_inputs` block supports:

* `id` - The ID of the input parameter.  

* `type` - The type of the input parameter.  

* `name` - The name of the input parameter.  

* `description` - The description of the input parameter.
