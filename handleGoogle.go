package main

import (
    "fmt"
    "strings"
    "time"
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
)

var oauthGoogleConfig *oauth2.Config
var oauthGoogleConnectDoneEndpoint string
var oauthGoogleCookieNameAccessToken string

func oauthGoogleSetup(c *oauth2.Config, oauthConnectDoneEndpoint string, oauthCookieNameAccessToken string) {
    oauthGoogleConfig = c
    oauthGoogleConnectDoneEndpoint = oauthConnectDoneEndpoint
    oauthGoogleCookieNameAccessToken = oauthCookieNameAccessToken
}

func oauthGoogleInitUrl(host string) string {
    if oauthGoogleConfig == nil {
        return ""
    }
    // Add hostname when callback url is not full url
    if !strings.HasPrefix(oauthGoogleConfig.RedirectURL, "https://") && !strings.HasPrefix(oauthGoogleConfig.RedirectURL, "http://")  {
        oauthGoogleConfig.RedirectURL = "https://" + host + oauthGoogleConfig.RedirectURL
    }
    return oauthGoogleConfig.AuthCodeURL("state")
}

func oauthGoogleRequiredAccessTokenMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("oauthGoogleRequiredAccessTokenMiddleware:", c.Request.Host)
        if cookie, err := c.Cookie(oauthGoogleCookieNameAccessToken); err == nil {
            fmt.Println("access_token:", cookie)
        } else {
            fmt.Println("oauthGoogleProfileHandler:",err)
            c.Redirect(http.StatusFound, oauthGoogleInitUrl(c.Request.Host))
            return
        }
        c.Next()
    }
}

func oauthGoogleCallbackHandler(c *gin.Context) {
    code := c.DefaultQuery("code", "")
    if code == "" {
        c.Redirect(http.StatusFound, oauthGoogleInitUrl(c.Request.Host))
        return
    }
    tok, err := oauthGoogleConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"status": false, "error": err})
        return 
    }
    // https://pkg.go.dev/golang.org/x/oauth2#Token
    // https://pkg.go.dev/net/http#Cookie => skip
    // https://pkg.go.dev/github.com/gin-gonic/gin#Context.SetCookie
    // https://pkg.go.dev/time#Time.Sub
    // https://pkg.go.dev/time#Duration
    maxAge := int(tok.Expiry.Sub(time.Now()).Seconds())
    if maxAge > 7200 {
        maxAge = maxAge - 3600
    }
    fmt.Println("cookie host:", c.Request.Host)
    c.SetCookie(oauthGoogleCookieNameAccessToken, tok.AccessToken, maxAge, "/", c.Request.Host, true, true)

    if oauthGoogleConnectDoneEndpoint != "" {
        c.Redirect(http.StatusFound, oauthGoogleConnectDoneEndpoint)
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": true, "token": tok,})
}

func oauthGoogleProfileHandler(c *gin.Context) {
    var accessToken string
    if cookie, err := c.Cookie(oauthGoogleCookieNameAccessToken); err == nil {
        accessToken = cookie
    } else {
        fmt.Println("oauthGoogleProfileHandler:",err)
    }
    c.JSON(http.StatusOK, gin.H{
        "status": true,
        "host": c.Request.Host,
        "access_token": accessToken,
    })
}
