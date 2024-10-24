FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:latest

ARG USE_GKE_GCLOUD_AUTH_PLUGIN=true
    
COPY ./bin/buildey /usr/local/bin
