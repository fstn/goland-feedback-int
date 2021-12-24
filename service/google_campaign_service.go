package googleextapiservice

import (
	"github.com/globalsign/mgo/bson"
	"net/url"
	adscampaignsdto "pixelme/_features/adsplateformintegration/crud/campaigns/dto"
	adscampaignsfacade "pixelme/_features/adsplateformintegration/crud/campaigns/facade"
	adsproductsfacade "pixelme/_features/adsplateformintegration/crud/products/facade"
	adsprovidersdto "pixelme/_features/adsplateformintegration/crud/providers/dto"
	adsprovidersfacade "pixelme/_features/adsplateformintegration/crud/providers/facade"
	adsprovidersmodel "pixelme/_features/adsplateformintegration/crud/providers/model"
	googleextapiconv "pixelme/_features/adsplateformintegration/externalapi/google/service/conv"
	googleextapiservice "pixelme/_features/adsplateformintegration/externalapi/google/service/interface"
	apitools "pixelme/_servers/common/apitools"
	exectx "pixelme/_servers/common/context"
	commonapitools "pixelme/_servers/common/httpclienttools"
	redirectsservice "pixelme/crud/audiencebuilder/redirects/service"
	"pixelme/database/mongodb"
	"pixelme/notifications"
	"time"
)

const BASE_URL = "http://localhost:8080"

type GoogleExternalAPIService struct {
	db                 mongodb.DBManager
	AdsCampaignsFacade adscampaignsfacade.IFacade
	AdsProvidersFacade adsprovidersfacade.IFacade
	AdsProductsFacade  adsproductsfacade.IFacade
	Notifier           notifications.Service
	RedirectsService   redirectsservice.IService
}

func NewService(db mongodb.DBManager, AdsProvidersFacade adsprovidersfacade.IFacade, AdsCampaignsFacade adscampaignsfacade.IFacade, AdsProductsFacade adsproductsfacade.IFacade, RedirectsService redirectsservice.IService, Notifier notifications.Service) googleextapiservice.IGoogleExternalAPIService {
	return GoogleExternalAPIService{
		db,
		AdsCampaignsFacade,
		AdsProvidersFacade,
		AdsProductsFacade,
		Notifier,
		RedirectsService,
	}
}

//CreateCampaign  on Google Ads API
func (self GoogleExternalAPIService) CreateCampaign(
	ctx exectx.Ctx,
	adsProviderToCreateCampaign adsprovidersdto.Provider,
	newCampaignToCreate adscampaignsdto.Campaign,
	responseData *struct {
		Campaign struct {
			ResourceName string `json:"resourceName"`
		} `json:"campaign"`
		AdGroup struct {
			ResourceName string `json:"resourceName"`
		} `json:"adGroup"`
	}) bool {

	googleAPICreateCampaign, err := googleextapiconv.ToExtAPICampaignDTO(newCampaignToCreate, adsProviderToCreateCampaign)
	if apitools.TryConv(err,
		"Unable to convert campaign request to external campaign request",
		ctx) {
		return true
	}

	done := commonapitools.Do(ctx, "Google Create Campaign Request", "POST", BASE_URL+"/post/campaign", nil, &googleAPICreateCampaign, &responseData)
	if done {
		return done
	}

	return false
}

//DeleteCampaign on Google Ads API
func (self GoogleExternalAPIService) DeleteCampaign(
	ctx exectx.Ctx,
	adsProviderToCreateCampaign adsprovidersdto.Provider,
	newCampaignToCreate adscampaignsdto.Campaign,
	responseData *struct {
		RemovedAd              []string `json:"removedAd"`
		RemovedCampaigns       []string `json:"removedCampaigns"`
		RemovedAdGroups        []string `json:"removedAdGroups"`
		RemovedCampaignBudgets []string `json:"removedCampaignBudgets"`
	}) bool {
	googleAPICreateCampaign, err := googleextapiconv.ToExtAPICampaignDTO(newCampaignToCreate, adsProviderToCreateCampaign)
	if apitools.TryConv(err,
		"Unable to convert campaign request to external campaign request",
		ctx) {
		return true
	}

	done := commonapitools.Do(ctx, "Google Create Campaign Request", "POST", BASE_URL+"/delete/campaign", nil, &googleAPICreateCampaign, &responseData)
	if done {
		return done
	}

	return false
}

// Retrieve customers list from Google API
func (self GoogleExternalAPIService) GetProfiles(
	ctx exectx.Ctx,
	requestData struct {
		RefreshToken string `json:"refreshToken"`
	}, profiles *[]googleextapiservice.Customer) bool {
	responseData := []googleextapiservice.GoogleExtAPICustomerDTO{}
	done := commonapitools.Do(ctx, "Google Get profiles Request", "POST", BASE_URL+"/get/customers", nil, &requestData, &responseData)
	if done {
		return done
	}

	for _, customer := range responseData {
		*profiles = append(*profiles, googleextapiservice.Customer{
			customer.ID,
			customer.DescriptiveName,
			customer.ResourceName,
		})

	}
	return false
}

func (self GoogleExternalAPIService) GetProfile(
	ctx exectx.Ctx,
	requestData struct {
		RefreshToken string `json:"refreshToken"`
	},
	customerID string, profile *googleextapiservice.Customer) bool {
	customer := googleextapiservice.GoogleExtAPICustomerDTO{}
	done := commonapitools.Do(ctx,
		"Google Get Profile Request",
		"POST",
		BASE_URL+"/get/customers/"+customerID,
		nil,
		&requestData,
		&customer)
	if done {
		return done
	}
	*profile = googleextapiservice.Customer{
		customer.ID,
		customer.DescriptiveName,
		customer.ResourceName,
	}
	return false
}

// GetMetricsCampaign  on Google Ads API
func (self GoogleExternalAPIService) GetMetricsCampaign(
	ctx exectx.Ctx,
	getMetricsCampaignRequest googleextapiservice.GoogleExtAPIMetricsDTO,
	responseData *struct {
		Impressions int
		Clicks      int
		CostMicros  int
	}) bool {

	requestData := url.Values{
		"googleCustomerId": []string{getMetricsCampaignRequest.GoogleCustomerID},
		"refreshToken":     []string{getMetricsCampaignRequest.RefreshToken},
		"from":             []string{getMetricsCampaignRequest.From.Format(time.RFC3339)},
		"to":               []string{getMetricsCampaignRequest.To.Format(time.RFC3339)},
	}
	done := commonapitools.Do(ctx,
		"Google Get Metrics Request",
		"POST", BASE_URL+"/get/campaign-metrics/{campaignId}",
		nil,
		&requestData,
		&responseData)
	if done {
		return done
	}

	return false
}

func (self GoogleExternalAPIService) SearchGeoLocations(ctx exectx.Ctx, name string, responseData *[]adsprovidersmodel.GeoCriteria) bool {
	database := self.db.Database()
	if apitools.TryDB(database.C("google_geo_targets").Find(bson.M{"canonical_name": bson.M{"$regex": ".*" + name + ".*", "$options": "i"}}).Limit(100).All(responseData),
		"Unable to load Google Geo locations",
		ctx) {
		return true
	}

	return false
}
