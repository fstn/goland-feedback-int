package googleextapiservice

import (
	"encoding/json"
	"moul.io/http2curl"
	"net/http"
	googleextapiservice "pixelme/_features/adsplateformintegration/externalapi/google/service/interface"
	apitools "pixelme/_servers/common/apitools"
	exectx "pixelme/_servers/common/context"
	commonapitools "pixelme/_servers/common/httpclienttools"
	"strings"
)

func (self GoogleExternalAPIService) GetRedirectURL(ctx exectx.Ctx) string {
	body := strings.NewReader("")
	req, err := http.NewRequest("POST", BASE_URL+"/redirect-url", body)
	if apitools.ReplyErrorStrCodeErr("Error during writing google request", http.StatusInternalServerError, err, ctx) {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	curl, _ := http2curl.GetCurlCommand(req)
	googleResp, err := http.DefaultClient.Do(req)
	if apitools.ReplyErrorStrCodeErr("Error during google call", http.StatusInternalServerError, err, ctx) {
		return ""
	}
	if apitools.ReplyErrorStrCodeErr("Error during google call", http.StatusInternalServerError, err, ctx) {
		return ""
	}
	var data struct {
		Url string `json:"url"`
	}
	dec := json.NewDecoder(googleResp.Body)
	err = dec.Decode(&data)
	if apitools.ReplyErrorStrCodeErr("Error during reading google response", http.StatusInternalServerError, err, ctx) {
		return ""
	}
	if apitools.SetWrappedStatusAndNotify("Error during oauth validation", *curl, *req, googleResp, ctx) {
		return ""
	}
	return data.Url
}

////////////////////////////////////////////////// RefreshToken //////////////////////////////////////////////////

func (self GoogleExternalAPIService) GetRefreshToken(ctx exectx.Ctx, code string, responseData *googleextapiservice.RefreshTokenResponse) bool {
	redirectURL := self.GetRedirectURL(ctx)
	type Data = struct {
		Code    string `json:"code"`
		BaseUri string `json:"baseUri"`
	}
	requestData := Data{
		code,
		redirectURL,
	}

	done := commonapitools.Do(ctx,
		"Google Refresh Token Request",
		"POST",
		BASE_URL+"/refresh-token",
		nil,
		&requestData,
		&responseData)
	if done {
		return true
	}

	return false
}
