resource "huaweicloud_s3_bucket" "bucket" {
  bucket = "tf-test-bucket"
  acl = "public-read"
}
resource "huaweicloud_smn_topic_v2" "topic_1" {
  name		  = "topic_check"
  display_name    = "The display name of topic_1"
}

resource "huaweicloud_cts_tracker_v1" "tracker_v1" {
  bucket_name      = "${huaweicloud_s3_bucket.bucket.bucket}"
  file_prefix_name      = "yO8Q"
  is_support_smn = true
  topic_id = "${huaweicloud_smn_topic_v2.topic_1.id}"
  is_send_all_key_operation = false
  operations = ["delete","create","login"]
  need_notify_user_list = ["user1"]
}