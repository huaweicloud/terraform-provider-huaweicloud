---
subcategory: "Cloud Application Engine (CAE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cae_component"
description: ""
---

# huaweicloud_cae_component

Manages a component resource within HuaweiCloud.

-> A maximum of `50` components can be created on the same environment. If you need more quotas,
   please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ticket/topic_0065264094.html)
   to submit a service ticket.

## Example Usage

### Create a component

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "component_name" {}
variable "image_url" {}

resource "huaweicloud_cae_component" "test" {
  environment_id = var.environment_id
  application_id = var.application_id

  metadata {
    name = var.component_name

    annotations = {
      version = "1.0.0"
    }
  }

  spec {
    runtime = "Docker"
    replica = 1

    source {
      type = "image"
      url  = var.image_url
    }

    resource_limit {
      cpu    = "500m"
      memory = "1Gi"
    }
  }
}
```

### Create and deploy a component

```hcl
variable "environment_id" {}
variable "application_id" {}
variable "component_name" {}
variable "image_url" {}

resource "huaweicloud_cae_component" "test" {
  environment_id = var.environment_id
  application_id = var.application_id

  metadata {
    name = var.component_name

    annotations = {
      version = "1.0.0"
    }
  }

  spec {
    runtime = "Docker"
    replica = 1

    source {
      type = "image"
      url  = var.image_url
    }

    resource_limit {
      cpu    = "500m"
      memory = "1Gi"
    }
  }

  action = "deploy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `environment_id` - (Required, String, ForceNew) Specifies the ID of the environment to which the application and the
  component belongs.
  Changing this creates a new resource.

* `application_id` - (Required, String, ForceNew) Specifies the ID of the application to which the component belongs.
  Changing this creates a new resource.

* `metadata` - (Required, List) Specifies the metadata of the component.
  The [metadata](#component_metadata) structure is documented below.

* `spec` - (Required, List) Specifies the configuration information of the component.
  The [spec](#component_spec) structure is documented below.

* `action` - (Optional, String) Specifies operation type of the component.  
  The valid values are as follows:
  + **deploy**: Deploy component. Only valid for undeployed component.
  + **configure**: Configurations of effesctive component. Only valid for deployed component.
  + **upgrade**: Upgrade component. Only valid for deployed component.

* `configurations` - (Optional, List) Specifies the list of configurations of the component.  
  The [configurations](#component_configurations) structure is documented below.  

  -> This parameter must be used together with `action` parameter.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the
  component belongs.  
  If the `application_id` belongs to the non-default enterprise project, this parameter is required and is only valid
  for enterprise users.

<a name="component_metadata"></a>
The `metadata` block supports:

* `name` - (Required, String) Specifies the name of the component.
  The name can contain `4` to `32` characters, only lowercase letters, digits, and hyphens (-) allowed.
  The name must start with a lowercase letter and end with lowercase letters and digits.

* `annotations` - (Required, Map) Specifies the key/value pairs parameters related to the component.
  Currently, only `version` is supported and required.
  The format is `A.B.C` or `A.B.C.D`, A, B, C and D must be integer. e.g.`1.0.0` or `1.0.0.0`

<a name="component_spec"></a>
The `spec` block supports:

* `replica` - (Required, Int) Specifies the instance number of the component. The valid value ranges from `1` to `99`.

* `runtime` - (Required, String) Specifies the component runtime to match. The valid values are **Docker**, **Java8**,
  **Java11**, **Java17**, **Tomcat8**, **Tomcat9**, **Python3**, **Nodejs8**, **Nodejs14**, **Nodejs16**, and **Php7**.

* `source` - (Required, List) Specifies the code source configuration information corresponding to the component.
  The [source](#component_spec_source) structure is documented below.

* `resource_limit` - (Required, List) Specifies instance specification corresponding to the component.
  The [resource_limit](#component_spec_resource_limit) structure is documented below.

* `build` - (Optional, List) Specifies the build information of the code source corresponding to the component.
  The [build](#component_spec_build) structure is documented below.

<a name="component_spec_resource_limit"></a>
The `resource_limit` block supports:

* `cpu` - (Required, String) Specifies CPU core. The valid values are **500m**, **1000m** and **2000m**.

* `memory` - (Required, String) Specifies memory size. The valid values are **1Gi**, **2Gi** and **4Gi**.

  -> If `cpu` parameter is set to `500m`, this parameter cannot be set to `4Gi`.

<a name="component_spec_source"></a>
The `source` block supports:

* `type` - (Required, String) Specifies code source type corresponding to the component.
  The valid values are **image**, **code** and **softwarePackage**.

* `url` - (Required, String) Specifies code source URL corresponding to the component.
  + When `type` is **image**, the URL represents image URL.
  + When `type` is **code**, the URL represents Git URL.
  + When `type` is **softwarePackage**, the URL represents software package URL.

* `sub_type` - (Optional, String) Specifies the subtype corresponding to the code source.
  If the `source.type` is set to `code`, the `sub_type` parameter means different code repositories.
  The valid values are `DevCloud`, `GitHub`, `GitLab`, `Gitee` and `Bitbucket`.
  If the `source.type` is set to `softwarePackage`, the `sub_type` parameter means different software package repositories.
  The valid values are `BinObs` and `BinDevCloud`.

  -> The parameter is required when `source.type` is set to `code` or `softwarePackage`.

* `code` - (Optional, List) Specifies code source repository.
  The [code](#component_spec_source_code) structure is documented below.

<a name="component_spec_source_code"></a>
The `code` block supports:

* `auth_name` - (Required, String) Specifies the name of authorization corresponding to the code source.

* `branch` - (Required, String) Specifies the branch name of code source repository.

* `namespace` - (Required, String) Specifies the username or organization corresponding to the code source repository.

<a name="component_spec_build"></a>
The `build` block supports:

* `archive` - (Required, List) Specifies product configuration after building the code source corresponding to component.
  The [archive](#component_spec_build_archive) structure is documented below.

* `parameters` - (Required, Map) Specifies the key/value pairs configuration information required to build the code source
  corresponding to the component.
  It is required when `source.type` is **code** or **softwarePackage**.
  + **base_image**: Base image address.
  + **build_cmd**: Custom build command.
  + **dockerfile_content**: Custom dockerfile content.
  + **dockerfile_path**: Custom dockerfile file path.
  + **artifact_name**: Select and run the specified JAR package from multiple JAR packages generated during Maven build.
  The JAR package end with **.jar**. Fuzzy match is supported. e.g. `demo-1.0.jar`, `demo*.jar`.

  -> `build_cmd`, `dockerfile_path` and `artifact_name` parameters are valid only when `source.type` is set to `code`.
     `dockerfile_path` and `artifact_name` parameters can't be set at the same time.
     `dockerfile_content` is valid only when `source.type` is set to `softwarePackage`.

<a name="component_spec_build_archive"></a>
The `archive` block supports:

* `artifact_namespace` - (Required, String) Specifies the name of the SWR organization after the code source
  corresponding to component is built.

<a name="component_configurations"></a>
The `configurations` block supports:

* `type` - (Required, String) Specifies the type of the component configuration.  
  Please following [reference documentation](https://support.huaweicloud.com/api-cae/CreateComponentConfiguration.html#CreateComponentConfiguration__request_ConfigurationItem).

* `data` - (Required, String) Specifies the configuration detail, in JSON format.  
  Please following [reference documentation](https://support.huaweicloud.com/api-cae/CreateComponentConfiguration.html#CreateComponentConfiguration__request_ConfigurationData).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `available_replica` - The number of available instances under the component.

* `status` - The current status of the component.
  + **running**
  + **paused**
  + **notReady**: The component deployed but not ready.
  + **created**: The component was not deployed.

* `created_at` - The creation time of the component.

* `updated_at` - The latest update time of the component.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The component can be imported using `environment_id`, `application_id` and `id`, separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_cae_component.test <environment_id>/<application_id>/<id>
```

For the component with the `enterprise_project_id`, its enterprise project ID need to be specified additionanlly when
importing. All fields are separated by slashes (/), e.g.

```bash
$ terraform import huaweicloud_cae_component.test <environment_id>/<application_id>/<id>/<enterprise_project_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `metadata.0.annotations`, `spec.0.build.0.parameters`, `action`, `configurations`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cae_component" "test" {
  ...

  lifecycle {
    ignore_changes = [
      metadata.0.annotations, spec.0.build.0.parameters, action, configurations,
    ]
  }
}
```
