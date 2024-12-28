output "secrets_layer_arn" {
    value = aws_lambda_layer_version.retrieve_secret_lambda_layer.arn
}