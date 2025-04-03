package model

import (
	"one-api/common"
	"one-api/setting/operation_setting"
	"sync"
	"time"
)

type Pricing struct {
	ModelName       string   `json:"model_name"`
	QuotaType       int      `json:"quota_type"`
	ModelRatio      float64  `json:"model_ratio"`
	ModelPrice      float64  `json:"model_price"`
	OwnerBy         string   `json:"owner_by"`
	CompletionRatio float64  `json:"completion_ratio"`
	EnableGroup     []string `json:"enable_groups,omitempty"`
	Tags            []string `json:"tags,omitempty"`
}

var (
	pricingMap         []Pricing
	lastGetPricingTime time.Time
	updatePricingLock  sync.Mutex
)

func GetPricing() []Pricing {
	updatePricingLock.Lock()
	defer updatePricingLock.Unlock()

	if time.Since(lastGetPricingTime) > time.Second*10 || len(pricingMap) == 0 {
		updatePricing()
	}
	//if group != "" {
	//	userPricingMap := make([]Pricing, 0)
	//	models := GetGroupModels(group)
	//	for _, pricing := range pricingMap {
	//		if !common.StringsContains(models, pricing.ModelName) {
	//			pricing.Available = false
	//		}
	//		userPricingMap = append(userPricingMap, pricing)
	//	}
	//	return userPricingMap
	//}
	return pricingMap
}

func updatePricing() {
	//modelRatios := common.GetModelRatios()
	enableAbilities := GetAllEnableAbilities()
	modelGroupsMap := make(map[string][]string)
	modelTagsMap := make(map[string][]string)
	for _, ability := range enableAbilities {
		groups := modelGroupsMap[ability.Model]
		tags := modelTagsMap[ability.Model]
		if groups == nil {
			groups = make([]string, 0)
			tags = make([]string, 0)
		}
		if !common.StringsContains(groups, ability.Group) {
			groups = append(groups, ability.Group)
			if ability.Tag != nil {
				tags = append(tags, *ability.Tag)
			}
		}
		modelGroupsMap[ability.Model] = groups
		modelTagsMap[ability.Model] = tags
	}

	pricingMap = make([]Pricing, 0)
	for model, groups := range modelGroupsMap {
		pricing := Pricing{
			ModelName:   model,
			EnableGroup: groups,
			Tags:        modelTagsMap[model],
		}
		modelPrice, findPrice := operation_setting.GetModelPrice(model, false)
		if findPrice {
			pricing.ModelPrice = modelPrice
			pricing.QuotaType = 1
		} else {
			modelRatio, _ := operation_setting.GetModelRatio(model)
			pricing.ModelRatio = modelRatio
			pricing.CompletionRatio = operation_setting.GetCompletionRatio(model)
			pricing.QuotaType = 0
		}
		pricingMap = append(pricingMap, pricing)
	}
	lastGetPricingTime = time.Now()
}
