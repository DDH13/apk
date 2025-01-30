/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package datastore

import (
	"log"
	"sync"

	api "github.com/wso2/apk/adapter/pkg/discovery/api/wso2/discovery/api"
	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/util"
)

// APIStore is a thread-safe store for APIs.
type APIStore struct {
	apis        map[string]*requestconfig.API
	mu          sync.RWMutex
	configStore *ConfigStore
}

// NewAPIStore creates a new instance of APIStore.
func NewAPIStore(configStore *ConfigStore) *APIStore {
	return &APIStore{
		configStore: configStore,
		// apis: make(map[string]*api.Api, 0),
	}
}

// AddAPIs adds a list of APIs to the store.
// This method is thread-safe.
func (s *APIStore) AddAPIs(apis []*api.Api) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.apis = make(map[string]*requestconfig.API, len(apis))
	for _, api := range apis {
		customAPI := requestconfig.API{
			Name:                   api.Title,
			Version:                api.Version,
			Vhost:                  api.Vhost,
			BasePath:               api.BasePath,
			APIType:                api.ApiType,
			EnvType:                api.EnvType,
			APILifeCycleState:      api.ApiLifeCycleState,
			AuthorizationHeader:    "", // You might want to set this field if applicable
			OrganizationID:         api.OrganizationId,
			UUID:                   api.Id,
			Tier:                   api.Tier,
			DisableAuthentication:  api.DisableAuthentications,
			DisableScopes:          api.DisableScopes,
			Resources:              make([]requestconfig.Resource, 0),
			IsMockedAPI:            false, // You can add logic to determine if the API is mocked
			MutualSSL:              api.MutualSSL,
			TransportSecurity:      api.TransportSecurity,
			ApplicationSecurity:    api.ApplicationSecurity,
			JwtConfigurationDto:    convertBackendJWTTokenInfoToJWTConfig(api.BackendJWTTokenInfo),
			SystemAPI:              api.SystemAPI,
			APIDefinition:          api.ApiDefinitionFile,
			Environment:            api.Environment,
			SubscriptionValidation: api.SubscriptionValidation,
			// Endpoints:              api.Endpoints,
			// EndpointSecurity:       convertSecurityInfoToEndpointSecurity(api.EndpointSecurity),
			AiProvider:                        convertAIProviderToDTO(api.Aiprovider),
			AIModelBasedRoundRobin:            convertAIModelBasedRoundRobinToDTO(api.AiModelBasedRoundRobin),
			DoSubscriptionAIRLInHeaderReponse: api.Aiprovider != nil && api.Aiprovider.PromptTokens != nil && api.Aiprovider.PromptTokens.In == dto.InHeader,
			DoSubscriptionAIRLInBodyReponse:   api.Aiprovider != nil && api.Aiprovider.PromptTokens != nil && api.Aiprovider.PromptTokens.In == dto.InBody,
		}
		for _, resource := range api.Resources {
			for _, operation := range resource.Methods {
				resource := buildResource(operation, resource.Path, convertAIModelBasedRoundRobinToDTO(resource.AiModelBasedRoundRobin), func() []*requestconfig.EndpointSecurity {
					endpointSecurity := make([]*requestconfig.EndpointSecurity, len(resource.EndpointSecurity))
					for i, es := range resource.EndpointSecurity {
						endpointSecurity[i] = &requestconfig.EndpointSecurity{
							Password:         es.Password,
							Enabled:          es.Enabled,
							Username:         es.Username,
							SecurityType:     es.SecurityType,
							CustomParameters: es.CustomParameters,
						}
					}
					return endpointSecurity
				}())
				customAPI.Resources = append(customAPI.Resources, resource)
			}
		}
		log.Printf("Adding API: %+v", customAPI.JwtConfigurationDto)
		s.apis[util.PrepareAPIKey(api.Vhost, api.BasePath, api.Version)] = &customAPI
	}
}

// GetAPIs retrieves the list of APIs from the store.
// This method is thread-safe.
func (s *APIStore) GetAPIs() map[string]*requestconfig.API {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis
}

// convertAIModelBasedRoundRobinToDTO converts AIModelBasedRoundRobin to DTO.
func convertAIModelBasedRoundRobinToDTO(aiModelBasedRoundRobin *api.AIModelBasedRoundRobin) *dto.AIModelBasedRoundRobin {
	if aiModelBasedRoundRobin == nil {
		return nil
	}
	return &dto.AIModelBasedRoundRobin{
		Enabled:                      aiModelBasedRoundRobin.Enabled,
		OnQuotaExceedSuspendDuration: int(aiModelBasedRoundRobin.OnQuotaExceedSuspendDuration),
		Models:                       convertModelWeights(aiModelBasedRoundRobin.Models),
	}
}

// convertModelWeights converts []*api.ModelWeight to []dto.ModelWeight.
func convertModelWeights(apiModelWeights []*api.ModelWeight) []dto.ModelWeight {
	dtoModelWeights := make([]dto.ModelWeight, len(apiModelWeights))
	for i, modelWeight := range apiModelWeights {
		dtoModelWeights[i] = dto.ModelWeight{
			Model:  modelWeight.Model,
			Weight: int(modelWeight.Weight),
		}
	}
	return dtoModelWeights
}

// convertAIProviderToDTO converts AIProvider to DTO.
func convertAIProviderToDTO(aiProvider *api.AIProvider) *dto.AIProvider {
	if aiProvider == nil {
		return nil
	}
	return &dto.AIProvider{
		ProviderName:       aiProvider.ProviderName,
		ProviderAPIVersion: aiProvider.ProviderAPIVersion,
		Organization:       aiProvider.Organization,
		Enabled:            aiProvider.Enabled,
		SupportedModels:    aiProvider.SupportedModels,
		Model:              convertValueDetailsPtr(aiProvider.Model),
		PromptTokens:       convertValueDetailsPtr(aiProvider.PromptTokens),
		CompletionToken:    convertValueDetailsPtr(aiProvider.CompletionToken),
		TotalToken:         convertValueDetailsPtr(aiProvider.TotalToken),
	}
}

// convertValueDetailsPtr converts *api.ValueDetails to *dto.ValueDetails.
func convertValueDetailsPtr(valueDetails *api.ValueDetails) *dto.ValueDetails {
	if valueDetails == nil {
		return nil
	}
	return &dto.ValueDetails{
		In:    valueDetails.In,
		Value: valueDetails.Value,
	}
}

// GetMatchedAPI retrieves the API that matches the given API key.
// GetMatchedAPI retrieves the API that matches the given API key.
func (s *APIStore) GetMatchedAPI(apiKey string) *requestconfig.API {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.apis[apiKey]
}

// ConvertBackendJWTTokenInfoToJWTConfig converts BackendJWTTokenInfo to JWTConfiguration.
func convertBackendJWTTokenInfoToJWTConfig(info *api.BackendJWTTokenInfo) *dto.JWTConfiguration {
	if info == nil {
		return nil
	}

	// Convert CustomClaims from map[string]*Claim to map[string]ClaimValue
	customClaims := make(map[string]dto.ClaimValue)
	for key, claim := range info.CustomClaims {
		if claim != nil {
			customClaims[key] = dto.ClaimValue{
				Value: claim.Value,
				Type:  claim.Type,
			}
		}
	}

	return &dto.JWTConfiguration{
		Enabled:                 info.Enabled,
		JWTHeader:               info.Header,
		ConsumerDialectURI:      "", // Add a default value or fetch if needed
		SignatureAlgorithm:      info.SigningAlgorithm,
		Encoding:                info.Encoding,
		GatewayJWTGeneratorImpl: "",                               // Add a default value or fetch if needed
		TokenIssuerDtoMap:       make(map[string]dto.TokenIssuer), // Populate if required
		JwtExcludedClaims:       make(map[string]bool),            // Populate if required
		PublicCert:              nil,                              // Add conversion logic if needed
		PrivateKey:              nil,                              // Add conversion logic if needed
		TTL:                     int64(info.TokenTTL),             // Convert int32 to int64
		CustomClaims:            customClaims,
	}
}
