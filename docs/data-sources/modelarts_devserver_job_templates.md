---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver_job_templates"
description: |-
  Use this data source to get the list of ModelArts DevServer job templates.
---

# huaweicloud_modelarts_devserver_job_templates

Use this data source to get the list of ModelArts DevServer job templates.

## Example Usage

### Query all DevServer job templates without any filter

```hcl
data "huaweicloud_modelarts_devserver_job_templates" "test" {}
```

### Query the DevServer job templates using type filter

```hcl
data "huaweicloud_modelarts_devserver_job_templates" "test" {
  type = "ASCEND_SYSTEM_CONFIG"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DevServer job templates are located.  
  If omitted, the provider-level region will be used.

* `template_id` - (Optional, String) Specifies the ID of the DevServer job template to be queried.

* `name` - (Optional, String) Specifies the name of the DevServer job template to be queried.

* `type` - (Optional, String) Specifies the type of the DevServer job template to be queried.

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of DevServer job templates that match the filter parameters.  
  The [templates](#modelarts_devserver_job_templates_attr) structure is documented below.

<a name="modelarts_devserver_job_templates_attr"></a>
The `templates` block supports:

* `id` - The ID of the DevServer job template.

* `name` - The name of the DevServer job template.

* `description` - The description of the DevServer job template.

* `type` - The type of the DevServer job template.

* `flavor_type` - The flavor type of the DevServer job template.

* `params` - The parameters of the DevServer job template.  
  The [params](#modelarts_devserver_job_templates_params_attr) structure is documented below.

<a name="modelarts_devserver_job_templates_params_attr"></a>
The `params` block supports:

* `name` - The name of the parameter.

* `description` - The description of the parameter.

* `value` - The value of the parameter.

* `visible` - Whether the parameter is visible in the console.

* `regex` - The regular expression for validating the parameter value.
