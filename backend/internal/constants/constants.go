package constants

const AppName = "Linked"

const AccessTokenKey = "access_token"

const RedirectUrlKey = "redirect_url"

const GoogleOAuthTokenInfoUrl = "https://oauth2.googleapis.com/tokeninfo?access_token="

const GoogleOAuthUserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type authKey_ struct{}

var AuthKey = authKey_{}
