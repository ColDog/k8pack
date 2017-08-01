# Auth Plugin

This plugin handles creating Kubeconfig files based on authentication with an oauth2 (openid connect?) provider. Kubeconfig files are downloaded locally to the developers machine configured with the correct groups and common name. All authentication is handled through oauth2.

## Flow

1. Developer visits auth login endpoint. The login endpoint includes the oauth2 provider and the api server host as parameters.
2. Browser is redirected to oauth2 provider, normal oauth2 flow takes place.
3. Browser redirects to auth callback endpoint. This endpoint retrieves user name and user groups from the oauth2 provider.
4. Browser is redirected to `localhost:6129` with a nonce to contact the auth api. This is a local docker container spun up for the duration of the login flow.
5. Local api contacts the auth api with nonce and retrieves a kubeconfig file.
6. Local api writes kubeconfig file to `~/.kube/config`.
7. Success page is shown.

## Components

1. `auth-api`: This is a server running remotely that creates kubeconfig files.
2. `local-api`: This is a local server that handles browser redirects and writes local kubeconfig files.
