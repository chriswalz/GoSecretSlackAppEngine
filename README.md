### Warning
Using the Slack Client is inherently less secure than using gosecret.io directly. For very sensitive information we recommend using gosecret.io directly

# GoSecret Slack Client 

##### Setting up local

`git clone https://github.com/chriswalz/GoSecretSlackAppEngine.git`

These two constants need to be created in sensitive/config.go:
`SLACK_TOKEN, SLACK_VERIF_TOKEN`

##### Running on App Engine 

`cd GoSecretSlackAppEngine` (where app.yaml file is)

`gcloud app deploy --quiet`

##### Debugging App Engine deployment

`gcloud app logs tail`





  
