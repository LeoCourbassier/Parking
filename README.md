# Install and Run instructions on INSTALL.MD

To use this webserver in k8s, you need to modify a few archives.

## Change ./kubernetes/mlabs/deployment
* Replace ACCOUNT_ID with your account id
* Replace REPO_NAME with your repo name

## Change secrets
* Add AWS_ACCESS_KEY_ID
* Add AWS_SECRET_ACCESS_KEY
* Add KUBE_CONFIG_DATA