package main

import (
    "fmt"
    "strings"
    "time"
    "net/http"
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
)

var oauthFacebookConfig *oauth2.Config
var oauthFacebookConnectDoneEndpoint string
var oauthFacebookCookieNameAccessToken string

func oauthFacebookSetup(c *oauth2.Config, oauthConnectDoneEndpoint string, oauthCookieNameAccessToken string) {
    oauthFacebookConfig = c
    oauthFacebookConnectDoneEndpoint = oauthConnectDoneEndpoint
    oauthFacebookCookieNameAccessToken = oauthCookieNameAccessToken
}

func oauthFacebookInitUrl(host string) string {
    if oauthFacebookConfig == nil {
        return ""
    }
    // Add hostname when callback url is not full url
    if !strings.HasPrefix(oauthFacebookConfig.RedirectURL, "https://") && !strings.HasPrefix(oauthFacebookConfig.RedirectURL, "http://")  {
        oauthFacebookConfig.RedirectURL = "https://" + host + oauthFacebookConfig.RedirectURL
    }
    return oauthFacebookConfig.AuthCodeURL("state")
}

func oauthFacebookRequiredMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("oauthFacebookRequiredMiddleware:", oauthFacebookInitUrl(c.Request.Host))
        code := c.DefaultQuery("code", "")
        if code == "" {
            c.Redirect(http.StatusFound, oauthFacebookInitUrl(c.Request.Host))
        }
        c.Next()
    }
}

func oauthFacebookRequiredAccessTokenMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("oauthFacebookRequiredAccessTokenMiddleware:", c.Request.Host)
        if cookie, err := c.Cookie(oauthFacebookCookieNameAccessToken); err == nil {
            fmt.Println("access_token:", cookie)
        } else {
            fmt.Println("oauthFacebookProfileHandler:",err)
            c.Redirect(http.StatusFound, oauthFacebookInitUrl(c.Request.Host))
            return
        }
        c.Next()
    }
}

func oauthFacebookCallbackHandler(c *gin.Context) {
    code := c.DefaultQuery("code", "")
    if code == "" {
        c.Redirect(http.StatusFound, oauthFacebookInitUrl(c.Request.Host))
        return
    }
    tok, err := oauthFacebookConfig.Exchange(oauth2.NoContext, code)
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
    c.SetCookie(oauthFacebookCookieNameAccessToken, tok.AccessToken, maxAge, "/", c.Request.Host, true, true)

    if oauthFacebookConnectDoneEndpoint != "" {
        c.Redirect(http.StatusFound, oauthFacebookConnectDoneEndpoint)
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": true, "token": tok,})
}

func oauthFacebookProfileHandler(c *gin.Context) {
    var accessToken string
    if cookie, err := c.Cookie(oauthFacebookCookieNameAccessToken); err == nil {
        accessToken = cookie
    } else {
        fmt.Println("oauthFacebookProfileHandler:",err)
    }
    c.JSON(http.StatusOK, gin.H{
        "status": true,
        "host": c.Request.Host,
        "access_token": accessToken,
    })
}
