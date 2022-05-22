# Usage

Setup:

```
% mkdir -p tmp
% cp config.default.yaml tmp/config.yaml
% cat tmp/config.yaml
OAuthSetup:
    Facebook:
        ClientID: YourFBAPPOAuthClientID
        ClientSecret: YourFBAPPOAuthClientSecret
        RedirectURL: /oauth/facebook/callback
        DoneURL: /oauth/facebook/profile
        CookieNameAccessToken: "oaFBAccessToken"

% go run .
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] GET    /privacy                  --> main.main.func2 (3 handlers)
[GIN-debug] GET    /oauth/facebook/          --> main.oauthFacebookProfileHandler (3 handlers)
[GIN-debug] GET    /oauth/facebook/profile   --> main.oauthFacebookProfileHandler (4 handlers)
[GIN-debug] GET    /oauth/facebook/callback  --> main.oauthFacebookCallbackHandler (3 handlers)
[GIN-debug] Listening and serving HTTPS on 
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
...
```

Browse the webpage: "https://localhost/oauth/facebook/profile"

