# Create an IoTDA Resource Space and Product

This example provides best practice code for using Terraform to create an IoTDA resource
space and a product with service capabilities in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An IoTDA **standard** or **enterprise** edition instance has been created in the console
  (the free unit `S0` is sufficient for this example, all resources in this example are
  metadata and do not incur extra fees)
* The HTTPS application access address of the instance has been obtained: login to the IoTDA
  console, choose the instance **Overview** and click **Access Details**

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IoTDA service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `space_name` - The name of the resource space
* `product_name` - The name of the product
* `iotda_access_address` - The HTTPS application access address of the IoTDA instance

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  space_name           = "tf_test_iotda_space"
  product_name         = "tf_test_iotda_product"
  iotda_access_address = "https://9bc34xxxxx.st1.iotda-app.cn-north-4.myhuaweicloud.com"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Notes

* When accessing an IoTDA standard or enterprise edition instance, the instance access
  address must be specified through the `endpoints.iotda` argument of the provider block,
  otherwise the API requests will be sent to the public service endpoint and fail
* The name of the resource space must not exceed `64` characters, only letters, digits,
  hyphens (`-`), underscores (`_`) and the following special characters are allowed:
  `?'#().,&%@!`
* The name of the product must not exceed `64` characters, the `device_type` must not
  exceed `32` characters, and the character set is the same as the space name
* The valid values of the `protocol` argument are **MQTT**, **CoAP**, **HTTP**, **HTTPS**,
  **Modbus**, **ONVIF**, **OPC-UA**, **OPC-DA** and **Other**, the valid values of the
  `data_type` argument are **json** and **binary**, and the `method` argument of the
  service property only supports **RW**, **W** and **R** in the current provider version
* The `space_id` argument of the product cannot be updated, changing it will recreate the
  product, and the product must be deleted before the resource space it belongs to can
  be deleted
* The creation and update of the product both submit the full set of the `services` block,
  do not modify the same product in the console at the same time, otherwise the service
  capabilities will be overwritten by Terraform
* The minimum valid value and maximum valid value of the int or decimal property must be
  strings, e.g. `min = "-40"`
* Both resources support import using the `id`

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.38.0 |
