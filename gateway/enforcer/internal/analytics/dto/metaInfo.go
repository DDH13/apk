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

package dto

// MetaInfo represents meta information attributes in an analytics event.
type MetaInfo struct {
	CorrelationID string `json:"correlationId"`
	RegionID      string `json:"regionId"`
	GatewayType   string `json:"gatewayType"`
}

// GetCorrelationID returns the correlation ID.
func (m *MetaInfo) GetCorrelationID() string {
	return m.CorrelationID
}

// SetCorrelationID sets the correlation ID.
func (m *MetaInfo) SetCorrelationID(correlationID string) {
	m.CorrelationID = correlationID
}

// GetRegionID returns the region ID.
func (m *MetaInfo) GetRegionID() string {
	return m.RegionID
}

// SetRegionID sets the region ID.
func (m *MetaInfo) SetRegionID(regionID string) {
	m.RegionID = regionID
}

// GetGatewayType returns the gateway type.
func (m *MetaInfo) GetGatewayType() string {
	return m.GatewayType
}

// SetGatewayType sets the gateway type.
func (m *MetaInfo) SetGatewayType(gatewayType string) {
	m.GatewayType = gatewayType
}
