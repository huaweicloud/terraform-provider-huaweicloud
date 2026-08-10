function_name         = "tf-test-function"
function_code         = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'statusCode': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
private_provider_name = "tf-test-provider"
