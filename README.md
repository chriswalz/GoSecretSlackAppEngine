#Running on App Engine 

Open cloud console 

go get -u github.com/chriswalz/GoSecretSlackAppEngine (don't use git clone even though that's suggested in the tutorial..)

go run . (to test)

gcloud app deploy (to deploy)




setting env variable 
env_variables:
  SLACK_TOKEN: ''
  SLACK_VERIF_TOKEN: ''
  
## Debugging 

gcloud app logs tail -s default