# Configure the HuaweiCloud Provider
provider "huaweicloud" {
  region     = var.region[8]
  access_key = var.access_key
  secret_key = var.secret_key
  insecure   = true
  project_id = ""
}

######################################################################
# Manages an SMN topic resource
######################################################################

resource "huaweicloud_smn_topic" "topic_1" {
  name                     = "topic_1"
  display_name             = "The display name of topic_1"
  users_publish_allowed    = "urn:csp:iam:::root"
  services_publish_allowed = "obs,vod,cce"
  introduction             = "created by terraform"
}

######################################################################
# Manages a Function resource
######################################################################

resource "huaweicloud_fgs_function" "f_1" {
  name        = "func_1"
  app         = "default"
  #  agency      = "test"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}

######################################################################
# Manages a trigger resource
######################################################################

resource "huaweicloud_fgs_trigger" "test" {
  function_urn = huaweicloud_fgs_function.f_1.urn
  type         = "TIMER"

  timer {
    name          = "test"
    schedule_type = "Rate"
    schedule      = "1d"
  }
}

######################################################################
# Manages a RMS aggregator resource
######################################################################

resource "huaweicloud_rms_resource_aggregator" "account" {
  name        = var.name
  type        = "ACCOUNT"
  account_ids = var.source_account_list
}


######################################################################
# Manages a RMS aggregation authorization resource
######################################################################

resource "huaweicloud_rms_resource_aggregation_authorization" "test" {
  account_id = ""
}


######################################################################
# Manages a RMS aggregation authorization resource
######################################################################

resource "huaweicloud_rms_resource_recorder" "test" {
  agency_name = "rms_tracker_agency"

  selector {
    all_supported = true
  }

  smn_channel {
    topic_urn = var.topic_urn
  }
}


######################################################################
# Manages agency resource
######################################################################

resource "huaweicloud_identity_agency" "agency" {
  name                  = "test_agency"
  delegated_domain_name = "hwstaff_zhenguo"

  project_role {
    project = "cn-north-4"
    roles   = ["FullAccess"]
  }
}
