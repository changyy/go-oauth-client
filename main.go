package main

import (
    "log"
    "io/ioutil"
    "net/http"
    "gopkg.in/yaml.v3"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/facebook"
    //"golang.org/x/oauth2/google"
)

//
// % cat config.default.yaml
//
type YAML_OAuthClientSetup struct {
    ClientID string                 `yaml:"ClientID"`
    ClientSecret string             `yaml:"ClientSecret"`
    RedirectURL string              `yaml:"RedirectURL"`
    DoneURL string                  `yaml:"DoneURL"`
    CookieNameAccessToken string    `yaml:"CookieNameAccessToken"`
}

type YAML_MainConfig struct {
    OAuthSetup map[string]YAML_OAuthClientSetup `yaml:"OAuthSetup"`
}

func main() {
    //mainConfig := make(map[interface{}]interface{})
    var mainConfig YAML_MainConfig 
    if yfile, err := ioutil.ReadFile("tmp/config.yaml"); err == nil {
        if err := yaml.Unmarshal(yfile, &mainConfig) ; err != nil {
            log.Fatalln(err)
        }
    }

    r := gin.Default()

    r.GET("/", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": true})
    })
    r.GET("/privacy", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": true})
    })

    oauthEndPoint := r.Group("oauth")
    {
        //
        // /oatuh/facebook
        // /oatuh/facebook/profile
        // /oatuh/facebook/callback
        //
        if oauthClientInfo, isSetup := mainConfig.OAuthSetup["Facebook"]; isSetup {
            oauthFacebookSetup(
                &oauth2.Config{
                    ClientID: oauthClientInfo.ClientID,
                    ClientSecret: oauthClientInfo.ClientSecret,
                    //RedirectURL: "https://localhost/oauth/facebook/callback",
                    // add host later
                    RedirectURL: oauthClientInfo.RedirectURL,
                    Endpoint: facebook.Endpoint,
                },
                // finishOAuthPage
                oauthClientInfo.DoneURL,
                // AccessTokenCookieName
                oauthClientInfo.CookieNameAccessToken,
            )
            oauthFacebook := oauthEndPoint.Group("facebook")
            //oauthFacebook.Use(oauthFacebookRequiredMiddleware())
            {
                oauthFacebook.GET("/", oauthFacebookProfileHandler)
                oauthFacebook.GET("/profile", oauthFacebookRequiredAccessTokenMiddleware(), oauthFacebookProfileHandler)
                oauthFacebook.GET("/callback", oauthFacebookCallbackHandler)
            }
        }

        if _, isSetup := mainConfig.OAuthSetup["Google"]; isSetup {
            oauthGoogle := oauthEndPoint.Group("google")
            {
                oauthGoogle.GET("/", nil)
                oauthGoogle.GET("/callback", nil)
            }
        }
    }

    //r.Run()
    r.RunTLS("", "./testdata/cert.pem", "./testdata/key.pem")
}
