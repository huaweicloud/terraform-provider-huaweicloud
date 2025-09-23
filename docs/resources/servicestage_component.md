---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestage_component"
description: ""
---

# huaweicloud_servicestage_component

This resource is used to manage a component under specified application within HuaweiCloud ServiceStage service.

## Example Usage

### Create a Web component using GitHub repository

```hcl
variable "application_id"
variable "component_name"
variable "token_auth_name"
variable "repo_url"
variable "repo_ref"
variable "repo_namespace"
variable "organization_name"
variable "cluster_id"

resource "huaweicloud_servicestage_component" "test" {
  application_id = var.application_id
  name           = var.component_name
  type           = "Webapp"
  runtime        = "Nodejs14"
  framework      = "Web"

  source {
    type           = "GitHub"
    authorization  = var.token_auth_name
    url            = var.repo_url
    repo_ref       = var.repo_ref
    repo_namespace = var.repo_namespace
  }

  builder {
    organization = var.organization_name
    cluster_id   = var.cluster_id

    node_label = {
      owner = "terraform"
    }
  }
}
```

### Create a MicroService Docker component

```hcl
variable "application_id"
variable "component_name"

resource "huaweicloud_servicestage_component" "test" {
  application_id = var.application_id
  name           = var.component_name
  type           = "MicroService"
  runtime        = "Docker"
  framework      = "Mesher"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application and component are located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new component.

* `application_id` - (Required, String, ForceNew) Specifies the application ID to which the component belongs.
  Changing this parameter will create a new component.

* `name` - (Required, String) Specifies the authorization name.
  The name can contain of `2` to `64` characters, only letters, digits, underscores (_) and hyphens (-) are allowed,
  and the name must start with a letter and end with a letter or digit.

* `type` - (Required, String, ForceNew) Specifies the component type. The valid values are as follows:
  + **Webapp**
  + **MicroService**
  + **Common**

  Changing this parameter will create a new component.

* `runtime` - (Required, String, ForceNew) Specifies the component runtime, such as **Docker**, **Java8**, etc.
  Changing this parameter will create a new component.

* `framework` - (Optional, String, ForceNew) Specifies the component framework.
  + The framework of type **Webapp** is **Web**.
  + The framework of type **MicroService** supports: **Java Classis**, **Go Classis**, **Mesher**, **Spring Cloud**,
  **Dubbo**.
  + The framework of type **Common** can be empty.

  Changing this parameter will create a new component.

-> For the runtime and framework corresponding to each type of component, please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-servicestage/servicestage_user_0411.html).

* `source` - (Optional, List) Specifies the repository source.
  The [object](#servicestage_component_source) structure is documented below.

* `builder` - (Optional, List) Specifies the component builder.
  The [object](#servicestage_component_builder) structure is documented below.

<a name="servicestage_component_source"></a>
The `source` block supports:

* `type` - (Required, String) Specifies the type of repository source or storage.
  The valid values are **GitHub**, **GitLab**, **Gitee**, **Bitbucket** and **package**.

* `url` - (Required, String) Specifies the URL of the repository or package storage.

* `authorization` - (Optional, String) Specifies the authorization name.
  This parameter and `storage_type` are alternative.

* `repo_ref` - (Optional, String) Specifies the name of the branch of the code repository.
  The default value is `master`.

* `repo_namespace` - (Optional, String) Specifies the namespace name.

* `storage_type` - (Optional, String) Specifies the storage type, such as **obs**, **swr**.
  This parameter is conflict with `repo_ref` and `repo_namespace`.

* `properties` - (Optional, List) Specifies the component builder's properties.
  The [object](#servicestage_component_properties) structure is documented below.

<a name="servicestage_component_builder"></a>
The `builder` block supports:

* `organization` - (Required, String) Specifies the organization name.
  The organization is usually **domain name**. You can find out in the organization management of SWR.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `cluster_name` - (Optional, String) Specifies the cluster Name.

* `cluster_type` - (Optional, String) Specifies the cluster type.

* `cmd` - (Optional, String) Specifies the build command. If omitted, the default command will be used.
  + About the  default command or script: build.sh in the root directory will be preferentially executed.
    If build.sh does not exist, the code will be compiled using the common method of the selected language,
    for example, mvn clean package for Java.
  + About the custom command: Commands will be customized using the selected language.
    Alternatively, the default command or script will be used after build.sh is modified.

* `dockerfile_path` - (Optional, String) Specifies the file path for dockerfile.

* `use_public_cluster` - (Optional, Bool) Specifies whether to use the public cluster.

* `node_label` - (Optional, Map) Specifies the filter labels for CCE nodes.

-> Before using the label, please make sure that the node is bound to the EIP and can access the public network.

<a name="servicestage_component_properties"></a>
The `properties` block supports:

* `endpoint` - (Optional, String) Specifies the endpoint of obs.

* `bucket` - (Optional, String) Specifies the bucket name of obs.

* `key` - (Optional, String) Specifies the key of obs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Components can be imported using their `application_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_servicestage_component.test dd7a1ce2-c48c-4f41-85bb-d0d09969eec9/9ab8ef79-d318-4de5-acf9-e1e1e25a0395
```
