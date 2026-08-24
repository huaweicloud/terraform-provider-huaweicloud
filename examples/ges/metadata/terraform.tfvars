bucket_name          = "tf-test-ges-metadata-bucket"
metadata_name        = "tf_test_ges_metadata"
metadata_description = "This is a demo"
metadata_schema_file = "schema_demo.xml"
metadata_properties  = [{
    dataType    = "char"
    name        = "sex"
    cardinality = "single"
}]
