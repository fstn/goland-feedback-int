package googleextapiservice

import (
	adscampaignsdto "pixelme/_features/adsplateformintegration/crud/campaigns/dto"
	adsprovidersdto "pixelme/_features/adsplateformintegration/crud/providers/dto"
	adsprovidersmodel "pixelme/_features/adsplateformintegration/crud/providers/model"
	exectx "pixelme/_servers/common/context"
	"time"
)

type GoogleExtAPIMetricsDTO struct {
	ExternalCampaignID string    `json:"external_campaign_id"`
	From               time.Time `json:"from"`
	To                 time.Time `json:"to"`
	RefreshToken       string    `json:"refresh_token"`
	GoogleCustomerID   string    `json:"google_customer_id"`
}

type GoogleExtAPICampaignDTO struct {
	PixelMeID        string            `json:"pixelMeId"`
	GoogleCustomerID string            `json:"googleCustomerId"`
	Campaign         GoogleAPICampaign `json:"campaign"`
	AdGroup          GoogleAPIAdGroup  `json:"adGroup"`
	Ad               GoogleAPIAd       `json:"ad"`
	RefreshToken     string            `json:"refreshToken"`
}

type GoogleExtAPICustomerDTO struct {
	ID              int
	DescriptiveName string
	ResourceName    string
}

type GoogleAPICampaign struct {
	ExternalID          string                                    `json:"externalId"`
	Name                string                                    `json:"name"`
	Status              int                                       `json:"status"`
	GeoCriteria         []adscampaignsdto.GoogleGeoCriterion      `json:"geoCriteria"`
	LanguagesCriteria   []adscampaignsdto.GoogleLanguageCriterion `json:"languagesCriteria"`
	BudgetMicros        int                                       `json:"budgetMicros"`
	Start               time.Time                                 `json:"start"`
	End                 time.Time                                 `json:"end"`
	Urls                []string                                  `json:"urls"`
	BiddingStrategyType int                                       `json:"biddingStrategyType"`
}

type GoogleLanguageCriterion struct {
	Code string `json:"code" validate:"required"`
}

type GoogleGeoCriterion struct {
	ID              int         `json:"id" validate:"required"`
	Name            string      `json:"name" validate:"required"`
	Active          interface{} `json:"active" validate:"required"`
	Status          interface{} `json:"status" validate:"required"`
	Negative        bool        `json:"negative" validate:"required"`
	Targettype      string      `json:"targetType" validate:"required"`
	Countrycode     string      `json:"countryCode" validate:"required"`
	Resourcename    interface{} `json:"resourceName" validate:"required"`
	Canonicalname   string      `json:"canonicalName" validate:"required"`
	Parentgeotarget string      `json:"parentGeoTarget" validate:"-"`
}

type GoogleAPIAdGroup struct {
	Name             string                                    `json:"name"`
	KeywordsCriteria []adscampaignsdto.GoogleKeywordsCriterion `json:"keywordsCriteria"` //Eg. keyword = Broad match , "keyword" = Phrase match , [keyword] = Exact match, -keyword = Negative match
	CpcBidMicros     int                                       `json:"cpcBidMicros"`
	ExternalID       string                                    `json:"externalId"`
}

type GoogleAPIAd struct {
	DisplayUrl      string                 `json:"displayUrl"`
	Headlines       []string               `json:"headlines"`
	Description     string                 `json:"description"`
	Description2    string                 `json:"description2"`
	DevicesCriteria map[string]interface{} `json:"devicesCriteria"`
	Path1           string                 `json:"path1"`
	Path2           string                 `json:"path2"`
}

type Customer struct {
	ID              int    `json:"id"`
	DescriptiveName string `json:"descriptive_name"`
	ResourceName    string `json:"resource_name"`
}

type RefreshTokenResponse struct {
	ClientID     string      `json:"clientId"`
	ClientSecret string      `json:"clientSecret"`
	RefreshToken string      `json:"refreshToken"`
	AccessToken  AccessToken `json:"accessToken"`
	UserID       string      `json:"userId"`
	Email        string      `json:"email"`
}

type AccessToken struct {
	TokenValue           string `json:"tokenValue"`
	ExpirationTimeMillis int    `json:"expirationTimeMillis"`
}

type IGoogleExternalAPIService interface {
	CreateCampaign(ctx exectx.Ctx, adsProvider adsprovidersdto.Provider, createCampaignRequest adscampaignsdto.Campaign, result *struct {
		Campaign struct {
			ResourceName string `json:"resourceName"`
		} `json:"campaign"`
		AdGroup struct {
			ResourceName string `json:"resourceName"`
		} `json:"adGroup"`
	}) bool
	GetMetricsCampaign(ctx exectx.Ctx, getMetricsCampaignRequest GoogleExtAPIMetricsDTO, response *struct {
		Impressions int
		Clicks      int
		CostMicros  int
	}) bool
	GetProfiles(ctx exectx.Ctx, requestData struct {
		RefreshToken string `json:"refreshToken"`
	}, Profiles *[]Customer) bool
	GetProfile(ctx exectx.Ctx, requestData struct {
		RefreshToken string `json:"refreshToken"`
	},
		customerID string, Profiles *Customer) bool
	DeleteCampaign(ctx exectx.Ctx, provider adsprovidersdto.Provider, campaign adscampaignsdto.Campaign,
		result *struct {
			RemovedAd              []string `json:"removedAd"`
			RemovedCampaigns       []string `json:"removedCampaigns"`
			RemovedAdGroups        []string `json:"removedAdGroups"`
			RemovedCampaignBudgets []string `json:"removedCampaignBudgets"`
		}) bool
	GetRedirectURL(ctx exectx.Ctx) string
	GetRefreshToken(ctx exectx.Ctx, code string, responseData *RefreshTokenResponse) bool
	SearchGeoLocations(ctx exectx.Ctx, name string, responseData *[]adsprovidersmodel.GeoCriteria) bool
}
