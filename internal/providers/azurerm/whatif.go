package azurerm

import (
	"encoding/json"
	"log"

	"github.com/tidwall/gjson"
)

type ChangeType string
type PropertyChangeType string

const (
	Create      ChangeType = "Create"
	Delete                 = "Delete"
	Deploy                 = "Deploy"
	Ignore                 = "Ignore"
	Modify                 = "Modify"
	NoChange               = "NoChange"
	Unsupported            = "Unsupported"
)

const (
	PropCreate   PropertyChangeType = "Create"
	PropDelete                      = "Delete"
	PropArray                       = "Array"
	PropModify                      = "Modify"
	PropNoEffect                    = "NoEffect"
)

// Struct for serializing the JSON response of a whatif call
// Modeled after the schema of deployments/whatIf in the AzureRM REST API
// see: https://learn.microsoft.com/en-us/rest/api/resources/deployments/what-if-at-subscription-scope
type WhatIf struct {
	Status     string           `json:"status"`
	Properties WhatifProperties `json:"properties"`
	Error      ErrorResponse    `json:"error,omitempty"`
}

type WhatifProperties struct {
	CorrelationId string         `json:"correlationId"`
	Changes       []WhatifChange `json:"changes,omitempty"`
}

type WhatifChange struct {
	ResourceId string     `json:"resourceId"`
	ChangeType ChangeType `json:"changeType"`
	// Before/After include several fields that are always present (resourceId, type etc.)
	// A resource's 'properties' field differs greatly, so serialize as raw JSON
	Before json.RawMessage `json:"before,omitempty"`
	After  json.RawMessage `json:"after,omitempty"`
	// TODO: Should be of type WhatIfChange
	Delta             json.RawMessage `json:"delta,omitempty"`
	UnsupportedReadon string          `json:"unsupportedReason,omitempty"`
}

func (this *WhatifChange) MarshalAfter() gjson.Result {
	marshal, err := this.After.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed marshalling After")
	}

	return gjson.ParseBytes(marshal)

	gjson.ParseBytes(marshal)
}

func (this *WhatifChange) MarshalBefore() {
	marshal, err := this.Before.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed marshalling Before")
	}

	return gjson.ParseBytes(marshal)
}

type WhatIfPropertyChange struct {
	After  json.RawMessage `json:"after,omitempty"`
	Before json.RawMessage `json:"before,omitempty"`
	// TODO: this should be of type []WhatIfPropertyChange
	// go lang structs can't do self-referring, references work?
	Children []*WhatIfPropertyChange `json:"children,omitempty"`
	Path     string                  `json:"path,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
}

type ErrorAdditionalInfo struct {
	Info string `json:"info"`
	Type string `json:"type"`
}
