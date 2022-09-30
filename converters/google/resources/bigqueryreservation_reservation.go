// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import "reflect"

const BigqueryReservationReservationAssetType string = "bigqueryreservation.googleapis.com/Reservation"

func resourceConverterBigqueryReservationReservation() ResourceConverter {
	return ResourceConverter{
		AssetType: BigqueryReservationReservationAssetType,
		Convert:   GetBigqueryReservationReservationCaiObject,
	}
}

func GetBigqueryReservationReservationCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//bigqueryreservation.googleapis.com/projects/{{project}}/locations/{{location}}/reservations/{{name}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetBigqueryReservationReservationApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: BigqueryReservationReservationAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/bigqueryreservation/v1/rest",
				DiscoveryName:        "Reservation",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetBigqueryReservationReservationApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	slotCapacityProp, err := expandBigqueryReservationReservationSlotCapacity(d.Get("slot_capacity"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("slot_capacity"); !isEmptyValue(reflect.ValueOf(slotCapacityProp)) && (ok || !reflect.DeepEqual(v, slotCapacityProp)) {
		obj["slotCapacity"] = slotCapacityProp
	}
	ignoreIdleSlotsProp, err := expandBigqueryReservationReservationIgnoreIdleSlots(d.Get("ignore_idle_slots"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("ignore_idle_slots"); !isEmptyValue(reflect.ValueOf(ignoreIdleSlotsProp)) && (ok || !reflect.DeepEqual(v, ignoreIdleSlotsProp)) {
		obj["ignoreIdleSlots"] = ignoreIdleSlotsProp
	}
	concurrencyProp, err := expandBigqueryReservationReservationConcurrency(d.Get("concurrency"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("concurrency"); !isEmptyValue(reflect.ValueOf(concurrencyProp)) && (ok || !reflect.DeepEqual(v, concurrencyProp)) {
		obj["concurrency"] = concurrencyProp
	}
	multiRegionAuxiliaryProp, err := expandBigqueryReservationReservationMultiRegionAuxiliary(d.Get("multi_region_auxiliary"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("multi_region_auxiliary"); !isEmptyValue(reflect.ValueOf(multiRegionAuxiliaryProp)) && (ok || !reflect.DeepEqual(v, multiRegionAuxiliaryProp)) {
		obj["multiRegionAuxiliary"] = multiRegionAuxiliaryProp
	}

	return obj, nil
}

func expandBigqueryReservationReservationSlotCapacity(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryReservationReservationIgnoreIdleSlots(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryReservationReservationConcurrency(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandBigqueryReservationReservationMultiRegionAuxiliary(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
