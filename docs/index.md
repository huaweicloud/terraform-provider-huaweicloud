# HuaweiCloud Provider

The HuaweiCloud provider is used to interact with the many resources supported by HuaweiCloud. The provider needs to be
configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

-> **NOTE:** You can also use this provider to interact with resources supported by HCSO if the API is compatible
with HuaweiCloud.

## Example Usage

Terraform 0.13 and later:

```hcl
terraform {
  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.36.0"
    }
  }
}

# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  region     = "cn-north-4"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "huaweicloud_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

Terraform 0.12 and earlier:

```hcl
# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  version    = ">= 1.36.0"
  region     = "cn-north-4"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}

# Create a VPC
resource "huaweicloud_vpc" "example" {
  name = "my_vpc"
  cidr = "192.168.0.0/16"
}
```

## Authentication

The Huawei Cloud provider offers a flexible means of providing credentials for authentication. The following methods are
supported, in this order, and explained below:

* Static credentials
* Environment variables
* Shared configuration file
* ECS Instance Metadata Service

The Huawei Cloud Provider supports assuming role with IAM agency, either in the provider configuration
block parameter assume_role or shared configuration file.

### Static credentials

!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage
should this file ever be committed to a public version control system.

Static credentials can be provided by adding an `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "huaweicloud" {
  region     = "cn-north-4"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

Temporary security credentials can be provided by adding an `security_token` with `access_key` and `secret_key`
in-line in the provider block:

Usage:

```hcl
provider "huaweicloud" {
  region         = "cn-north-4"
  access_key     = "my-access-key"
  secret_key     = "my-secret-key"
  security_token = "my-security-token"
}
```

### Environment variables

You can provide your credentials via the `HW_ACCESS_KEY` and
`HW_SECRET_KEY` environment variables, representing your Huawei Cloud Access Key and Secret Key, respectively.
For temporary security credentials, need to add one more environment variable `HW_SECURITY_TOKEN`.

```hcl
provider "huaweicloud" {}
```

Usage:

```sh
$ export HW_ACCESS_KEY="anaccesskey"
$ export HW_SECRET_KEY="asecretkey"
$ export HW_REGION_NAME="cn-north-4"
$ terraform plan
```

### Shared Configuration File

You can use a
[HuaweiCloud CLI configuration file](https://support.huaweicloud.com/intl/en-us/usermanual-hcli/hcli_03_002.html)
to specify your credentials. You need to specify a location in the Terraform configuration by providing the
`shared_config_file` argument or using the `HW_SHARED_CONFIG_FILE` environment variable.
This method also supports a `profile` configuration and matching `HW_PROFILE` environment variable:

-> **NOTE:** The CLI configuration file can not be used directly by terraform, you need to skip encrypting
authentication information in the configuration file by running the following command:

```sh
hcloud configure set --cli-auth-encrypt=false
```

Usage:

```terraform
provider "huaweicloud" {
  shared_config_file = "/home/tf_user/.hcloud/config.json"
  profile            = "customprofile"
}
```

### ECS Instance Metadata Service

If you're running Terraform from an ECS instance with Agency configured, Terraform will just ask
[the metadata API](https://support.huaweicloud.com/intl/en-us/usermanual-ecs/ecs_03_0166.html)
for credentials.

This is a preferred approach over any other when running in ECS as you can avoid
hard coding credentials. Instead these are leased on-the-fly by Terraform
which reduces the chance of leakage.

Usage:

```hcl
provider "huaweicloud" {
  region = "cn-north-4"
}
```

### Assume role

If provided with an IAM agency, Terraform will attempt to assume this role using the supplied credentials.

Usage:

```hcl
provider "huaweicloud" {
  region     = "cn-north-4"
  access_key = "my-access-key"
  secret_key = "my-secret-key"

  assume_role {
    agency_name = "agency"
    domain_name = "agency_domain"
  }
}
```

## Configuration Reference

The following arguments are supported:

* `region` - (Optional) This is the Huawei Cloud region. It must be provided when using `static credentials`
  authentication, but it can also be sourced from the `HW_REGION_NAME` environment variables.

* `access_key` - (Optional) The access key of the HuaweiCloud to use. If omitted, the `HW_ACCESS_KEY` environment
  variable is used.

* `secret_key` - (Optional) The secret key of the HuaweiCloud to use. If omitted, the `HW_SECRET_KEY` environment
  variable is used.

* `shared_config_file` - (Optional) The path to the shared config file. If omitted, the `HW_SHARED_CONFIG_FILE` environment
  variable is used.

* `profile` - (Optional) The profile name as set in the shared config file. If omitted, the `HW_PROFILE` environment
  variable is used. Defaults to the `current` profile in the shared config file.

* `assume_role` - (Optional) Configuration block for an assumed role. See below. Only one assume_role
  block may be in the configuration.

* `project_name` - (Optional) The Name of the project to login with. If omitted, the `HW_PROJECT_NAME` environment
  variable or `region` is used.

* `domain_name` - (Optional) The [Account name](https://support.huaweicloud.com/en-us/usermanual-iam/iam_01_0552.html)
  of IAM to scope to. If omitted, the `HW_DOMAIN_NAME` environment variable is used.

* `security_token` - (Optional) The security token to authenticate with a
  [temporary security credential](https://support.huaweicloud.com/intl/en-us/iam_faq/iam_01_0620.html). If omitted,
  the `HW_SECURITY_TOKEN` environment variable is used.

* `cloud` - (Optional) The endpoint of the cloud provider. If omitted, the
  `HW_CLOUD` environment variable is used. Defaults to `myhuaweicloud.com`.

* `auth_url` - (Optional, Required before 1.14.0) The Identity authentication URL. If omitted, the
  `HW_AUTH_URL` environment variable is used. Defaults to `https://iam.{{region}}.{{cloud}}/v3`.

* `insecure` - (Optional) Trust self-signed SSL certificates. If omitted, the
  `HW_INSECURE` environment variable is used.

* `max_retries` - (Optional) This is the maximum number of times an API call is retried, in the case where requests are
  being throttled or experiencing transient failures. The delay between the subsequent API calls increases
  exponentially. The default value is `5`. If omitted, the `HW_MAX_RETRIES` environment variable is used.

* `enterprise_project_id` - (Optional) Default Enterprise Project ID for supported resources. Please see the
  documentation
  at [EPS](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs/data-sources/enterprise_project).
  If omitted, the `HW_ENTERPRISE_PROJECT_ID` environment variable is used.

* `signing_algorithm` - (Optional) The signing algorithm for authentication. Valid values are **HmacSHA256**,
  **HmacSM3**, **EcdsaP256SHA256**, **SM2SM3**.
  If omitted, the `HW_SIGNING_ALGORITHM` environment variable is used.

* `regional` - (Optional) Whether the service endpoints are regional. The default value is `false`.

* `skip_check_website_type` - (Optional) Whether to skip website type check. The default value is `false`.

* `skip_check_upgrade` - (Optional) Whether to skip upgrade check. The default value is `false`.

* `endpoints` - (Optional) Configuration block in key/value pairs for customizing service endpoints.
  The [endpoints](#block--endpoints) block to support custom endpoints is documented below.
  An example provider configuration:

```hcl
provider "huaweicloud" {
  ...
  endpoints = {
    ecs = "https://ecs-customizing-endpoint.com"
  }
}
```

* `default_tags` - (Optional) The default tags of resources managed by this provider.
  The `default_tags` only works in resources that support the `tags` argument.
  The `default_tags` can be overridden if `tags` in the resource has new values for matching keys.

The `assume_role` block supports:

* `agency_name` - (Required) The name of the agency for assume role.
  If omitted, the `HW_ASSUME_ROLE_AGENCY_NAME` environment variable is used.

* `domain_name` - (Required) The name of the agency domain for assume role.
  If omitted, the `HW_ASSUME_ROLE_DOMAIN_NAME` environment variable is used.

<a name="block--endpoints"></a>
The `endpoints` block supports:

* `aad` - (Optional) Use this to override the default endpoint URL. It's used to customize **AAD** endpoints.

* `accessanalyzer` - (Optional) Use this to override the default endpoint URL. It's used to customize **AccessAnalyzer**
  endpoints.

* `anti-ddos` - (Optional) Use this to override the default endpoint URL. It's used to customize **Anti-DDoS** endpoints.

* `aom` - (Optional) Use this to override the default endpoint URL. It's used to customize **AOM** endpoints.

* `apig` - (Optional) Use this to override the default endpoint URL. It's used to customize **APIG** endpoints.

* `apm` - (Optional) Use this to override the default endpoint URL. It's used to customize **APM** endpoints.

* `appstream` - (Optional) Use this to override the default endpoint URL. It's used to customize **WorkspaceApp** endpoints.

* `asm` - (Optional) Use this to override the default endpoint URL. It's used to customize **ASM** endpoints.

* `autoscaling` - (Optional) Use this to override the default endpoint URL. It's used to customize **AS** endpoints.

* `bcs` - (Optional) Use this to override the default endpoint URL. It's used to customize **BCS** endpoints.

* `bms` - (Optional) Use this to override the default endpoint URL. It's used to customize **BMS** endpoints.

* `bss` - (Optional) Use this to override the default endpoint URL. It's used to customize **CBC** endpoints.

* `cae` - (Optional) Use this to override the default endpoint URL. It's used to customize **CAE** endpoints.

* `cbh` - (Optional) Use this to override the default endpoint URL. It's used to customize **CBH** endpoints.

* `cbr` - (Optional) Use this to override the default endpoint URL. It's used to customize **CBR** endpoints.

* `cc` - (Optional) Use this to override the default endpoint URL. It's used to customize **CC** endpoints.

* `cce` - (Optional) Use this to override the default endpoint URL. It's used to customize **CCE** endpoints.

* `cci` - (Optional) Use this to override the default endpoint URL. It's used to customize **CCI** endpoints.

* `ccm` - (Optional) Use this to override the default endpoint URL. It's used to customize **CCM** endpoints.

* `cdm` - (Optional) Use this to override the default endpoint URL. It's used to customize **CDM** endpoints.

* `cdn` - (Optional) Use this to override the default endpoint URL. It's used to customize **CDN** endpoints.

* `ces` - (Optional) Use this to override the default endpoint URL. It's used to customize **CES** endpoints.

* `cfw` - (Optional) Use this to override the default endpoint URL. It's used to customize **CFW** endpoints.

* `cloudtable` - (Optional) Use this to override the default endpoint URL. It's used to customize **CloudTable**
  endpoints.

* `cms` - (Optional) Use this to override the default endpoint URL. It's used to customize **CMS** endpoints.

* `coc` - (Optional) Use this to override the default endpoint URL. It's used to customize **COC** endpoints.

* `codearts_build` - (Optional) Use this to override the default endpoint URL. It's used to customize **CodeArtsBuild**
  endpoints.

* `codearts_deploy` - (Optional) Use this to override the default endpoint URL. It's used to customize **CodeArtsDeploy**
  endpoints.

* `codearts_pipeline` - (Optional) Use this to override the default endpoint URL. It's used to customize
  **CodeArtsPipeline** endpoints.

* `codehub` - (Optional) Use this to override the default endpoint URL. It's used to customize **CodeHub** endpoints.

* `cph` - (Optional) Use this to override the default endpoint URL. It's used to customize **CPH** endpoints.

* `cpts` - (Optional) Use this to override the default endpoint URL. It's used to customize **CPTS** endpoints.

* `cse` - (Optional) Use this to override the default endpoint URL. It's used to customize **CSE** endpoints.

* `css` - (Optional) Use this to override the default endpoint URL. It's used to customize **CSS** endpoints.

* `cts` - (Optional) Use this to override the default endpoint URL. It's used to customize **CTS** endpoints.

* `dataarts` - (Optional) Use this to override the default endpoint URL. It's used to customize **DataArtsStudio**
  endpoints.

* `dbss` - (Optional) Use this to override the default endpoint URL. It's used to customize **DBSS** endpoints.

* `dc` - (Optional) Use this to override the default endpoint URL. It's used to customize **DC** endpoints.

* `dcs` - (Optional) Use this to override the default endpoint URL. It's used to customize **DCS** endpoints.

* `ddm` - (Optional) Use this to override the default endpoint URL. It's used to customize **DDM** endpoints.

* `dds` - (Optional) Use this to override the default endpoint URL. It's used to customize **DDS** endpoints.

* `deh` - (Optional) Use this to override the default endpoint URL. It's used to customize **DEH** endpoints.

* `dis` - (Optional) Use this to override the default endpoint URL. It's used to customize **DIS** endpoints.

* `dli` - (Optional) Use this to override the default endpoint URL. It's used to customize **DLI** endpoints.

* `dms` - (Optional) Use this to override the default endpoint URL. It's used to customize **DMS** endpoints.

* `dns` - (Optional) Use this to override the default endpoint URL. It's used to customize **DNS** endpoints.

* `drs` - (Optional) Use this to override the default endpoint URL. It's used to customize **DRS** endpoints.

* `dsc` - (Optional) Use this to override the default endpoint URL. It's used to customize **DSC** endpoints.

* `dws` - (Optional) Use this to override the default endpoint URL. It's used to customize **DWS** endpoints.

* `ecs` - (Optional) Use this to override the default endpoint URL. It's used to customize **ECS** endpoints.

* `eg` - (Optional) Use this to override the default endpoint URL. It's used to customize **EG** endpoints.

* `elb` - (Optional) Use this to override the default endpoint URL. It's used to customize **ELB** endpoints.

* `eps` - (Optional) Use this to override the default endpoint URL. It's used to customize **EPS** endpoints.

* `er` - (Optional) Use this to override the default endpoint URL. It's used to customize **ER** endpoints.

* `evs` - (Optional) Use this to override the default endpoint URL. It's used to customize **EVS** endpoints.

* `fgs` - (Optional) Use this to override the default endpoint URL. It's used to customize **FunctionGraph**
  endpoints.

* `ga` - (Optional) Use this to override the default endpoint URL. It's used to customize **GA** endpoints.

* `gaussdb` - (Optional) Use this to override the default endpoint URL. It's used to customize **GaussDBforMySQL**
  endpoints.

* `geip` - (Optional) Use this to override the default endpoint URL. It's used to customize **GEIP** endpoints.

* `geminidb` - (Optional) Use this to override the default endpoint URL. It's used to customize **GaussDBforNoSQL**
  endpoints.

* `ges` - (Optional) Use this to override the default endpoint URL. It's used to customize **GES** endpoints.

* `hss` - (Optional) Use this to override the default endpoint URL. It's used to customize **HSS** endpoints.

* `iam` - (Optional) Use this to override the default endpoint URL. It's used to customize **IAM** endpoints.

* `identitycenter` - (Optional) Use this to override the default endpoint URL. It's used to customize **IdentityCenter**
  endpoints.

* `identitystore` - (Optional) Use this to override the default endpoint URL. It's used to customize **IdentityCenter**
  endpoints.

* `iec` - (Optional) Use this to override the default endpoint URL. It's used to customize **IEC** endpoints.

* `ims` - (Optional) Use this to override the default endpoint URL. It's used to customize **IMS** endpoints.

* `iotda` - (Optional) Use this to override the default endpoint URL. It's used to customize **IoTDA** endpoints.

* `kms` - (Optional) Use this to override the default endpoint URL. It's used to customize **DEW** endpoints.

* `live` - (Optional) Use this to override the default endpoint URL. It's used to customize **Live** endpoints.

* `lts` - (Optional) Use this to override the default endpoint URL. It's used to customize **LTS** endpoints.

* `meeting` - (Optional) Use this to override the default endpoint URL. It's used to customize **Meeting** endpoints.

* `mkt` - (Optional) Use this to override the default endpoint URL. It's used to customize **KooGallery** endpoints.

* `modelarts` - (Optional) Use this to override the default endpoint URL. It's used to customize **ModelArts** endpoints.

* `mpc` - (Optional) Use this to override the default endpoint URL. It's used to customize **MPC** endpoints.

* `mrs` - (Optional) Use this to override the default endpoint URL. It's used to customize **MRS** endpoints.

* `nat` - (Optional) Use this to override the default endpoint URL. It's used to customize **NAT** endpoints.

* `oms` - (Optional) Use this to override the default endpoint URL. It's used to customize **OMS** endpoints.

* `opengauss` - (Optional) Use this to override the default endpoint URL. It's used to customize **GaussDB** endpoints.

* `organizations` - (Optional) Use this to override the default endpoint URL. It's used to customize **Organizations**
  endpoints.

* `projectman` - (Optional) Use this to override the default endpoint URL. It's used to customize **ProjectMan** endpoints.

* `ram` - (Optional) Use this to override the default endpoint URL. It's used to customize **RAM** endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL. It's used to customize **RDS** endpoints.

* `rfs` - (Optional) Use this to override the default endpoint URL. It's used to customize **RFS** endpoints.

* `rgc` - (Optional) Use this to override the default endpoint URL. It's used to customize **RGC** endpoints.

* `rts` - (Optional) Use this to override the default endpoint URL. It's used to customize **RTS** endpoints.

* `scm` - (Optional) Use this to override the default endpoint URL. It's used to customize **CCM** endpoints.

* `sdrs` - (Optional) Use this to override the default endpoint URL. It's used to customize **SDRS** endpoints.

* `secmaster` - (Optional) Use this to override the default endpoint URL. It's used to customize **SecMaster** endpoints.

* `servicestage` - (Optional) Use this to override the default endpoint URL. It's used to customize **ServiceStage**
  endpoints.

* `sfs` - (Optional) Use this to override the default endpoint URL. It's used to customize **SFS** endpoints.

* `sfs-turbo` - (Optional) Use this to override the default endpoint URL. It's used to customize **SFSTurbo** endpoints.

* `smn` - (Optional) Use this to override the default endpoint URL. It's used to customize **SMN** endpoints.

* `sts` - (Optional) Use this to override the default endpoint URL. It's used to customize **IAM** endpoints.

* `swr` - (Optional) Use this to override the default endpoint URL. It's used to customize **SWR** endpoints.

* `tms` - (Optional) Use this to override the default endpoint URL. It's used to customize **TMS** endpoints.

* `ucs` - (Optional) Use this to override the default endpoint URL. It's used to customize **UCS** endpoints.

* `vbs` - (Optional) Use this to override the default endpoint URL. It's used to customize **VBS** endpoints.

* `vod` - (Optional) Use this to override the default endpoint URL. It's used to customize **VOD** endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL. It's used to customize **VPC** endpoints.

* `vpcep` - (Optional) Use this to override the default endpoint URL. It's used to customize **VPCEP** endpoints.

* `vpn` - (Optional) Use this to override the default endpoint URL. It's used to customize **VPN** endpoints.

* `vss` - (Optional) Use this to override the default endpoint URL. It's used to customize **CodeArtsInspector**
  endpoints.

* `waf` - (Optional) Use this to override the default endpoint URL. It's used to customize **WAF** endpoints.

* `workspace` - (Optional) Use this to override the default endpoint URL. It's used to customize **WorkspaceDesktop** endpoints.

## Testing and Development

In order to run the Acceptance Tests for development, the following environment variables must also be set:

* `HW_REGION_NAME` - The region in which to create the resources.

* `HW_ACCESS_KEY` - The access key of the HuaweiCloud to use.

* `HW_SECRET_KEY` - The secret key of the HuaweiCloud to use.

You should be able to use any HuaweiCloud environment to develop on as long as the above environment variables are set.
