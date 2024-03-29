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

const ArtifactRegistryRepositoryAssetType string = "artifactregistry.googleapis.com/Repository"

func resourceConverterArtifactRegistryRepository() ResourceConverter {
	return ResourceConverter{
		AssetType: ArtifactRegistryRepositoryAssetType,
		Convert:   GetArtifactRegistryRepositoryCaiObject,
	}
}

func GetArtifactRegistryRepositoryCaiObject(d TerraformResourceData, config *Config) ([]Asset, error) {
	name, err := assetName(d, config, "//artifactregistry.googleapis.com/projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}")
	if err != nil {
		return []Asset{}, err
	}
	if obj, err := GetArtifactRegistryRepositoryApiObject(d, config); err == nil {
		return []Asset{{
			Name: name,
			Type: ArtifactRegistryRepositoryAssetType,
			Resource: &AssetResource{
				Version:              "v1",
				DiscoveryDocumentURI: "https://www.googleapis.com/discovery/v1/apis/artifactregistry/v1/rest",
				DiscoveryName:        "Repository",
				Data:                 obj,
			},
		}}, nil
	} else {
		return []Asset{}, err
	}
}

func GetArtifactRegistryRepositoryApiObject(d TerraformResourceData, config *Config) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	formatProp, err := expandArtifactRegistryRepositoryFormat(d.Get("format"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("format"); !isEmptyValue(reflect.ValueOf(formatProp)) && (ok || !reflect.DeepEqual(v, formatProp)) {
		obj["format"] = formatProp
	}
	descriptionProp, err := expandArtifactRegistryRepositoryDescription(d.Get("description"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	labelsProp, err := expandArtifactRegistryRepositoryLabels(d.Get("labels"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}
	kmsKeyNameProp, err := expandArtifactRegistryRepositoryKmsKeyName(d.Get("kms_key_name"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("kms_key_name"); !isEmptyValue(reflect.ValueOf(kmsKeyNameProp)) && (ok || !reflect.DeepEqual(v, kmsKeyNameProp)) {
		obj["kmsKeyName"] = kmsKeyNameProp
	}
	dockerConfigProp, err := expandArtifactRegistryRepositoryDockerConfig(d.Get("docker_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("docker_config"); !isEmptyValue(reflect.ValueOf(dockerConfigProp)) && (ok || !reflect.DeepEqual(v, dockerConfigProp)) {
		obj["dockerConfig"] = dockerConfigProp
	}
	mavenConfigProp, err := expandArtifactRegistryRepositoryMavenConfig(d.Get("maven_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("maven_config"); !isEmptyValue(reflect.ValueOf(mavenConfigProp)) && (ok || !reflect.DeepEqual(v, mavenConfigProp)) {
		obj["mavenConfig"] = mavenConfigProp
	}
	modeProp, err := expandArtifactRegistryRepositoryMode(d.Get("mode"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("mode"); !isEmptyValue(reflect.ValueOf(modeProp)) && (ok || !reflect.DeepEqual(v, modeProp)) {
		obj["mode"] = modeProp
	}
	virtualRepositoryConfigProp, err := expandArtifactRegistryRepositoryVirtualRepositoryConfig(d.Get("virtual_repository_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("virtual_repository_config"); !isEmptyValue(reflect.ValueOf(virtualRepositoryConfigProp)) && (ok || !reflect.DeepEqual(v, virtualRepositoryConfigProp)) {
		obj["virtualRepositoryConfig"] = virtualRepositoryConfigProp
	}
	remoteRepositoryConfigProp, err := expandArtifactRegistryRepositoryRemoteRepositoryConfig(d.Get("remote_repository_config"), d, config)
	if err != nil {
		return nil, err
	} else if v, ok := d.GetOkExists("remote_repository_config"); !isEmptyValue(reflect.ValueOf(remoteRepositoryConfigProp)) && (ok || !reflect.DeepEqual(v, remoteRepositoryConfigProp)) {
		obj["remoteRepositoryConfig"] = remoteRepositoryConfigProp
	}

	return obj, nil
}

func expandArtifactRegistryRepositoryFormat(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandArtifactRegistryRepositoryKmsKeyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryDockerConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedImmutableTags, err := expandArtifactRegistryRepositoryDockerConfigImmutableTags(original["immutable_tags"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedImmutableTags); val.IsValid() && !isEmptyValue(val) {
		transformed["immutableTags"] = transformedImmutableTags
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryDockerConfigImmutableTags(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryMavenConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowSnapshotOverwrites, err := expandArtifactRegistryRepositoryMavenConfigAllowSnapshotOverwrites(original["allow_snapshot_overwrites"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowSnapshotOverwrites); val.IsValid() && !isEmptyValue(val) {
		transformed["allowSnapshotOverwrites"] = transformedAllowSnapshotOverwrites
	}

	transformedVersionPolicy, err := expandArtifactRegistryRepositoryMavenConfigVersionPolicy(original["version_policy"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedVersionPolicy); val.IsValid() && !isEmptyValue(val) {
		transformed["versionPolicy"] = transformedVersionPolicy
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryMavenConfigAllowSnapshotOverwrites(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryMavenConfigVersionPolicy(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryMode(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryVirtualRepositoryConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedUpstreamPolicies, err := expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPolicies(original["upstream_policies"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedUpstreamPolicies); val.IsValid() && !isEmptyValue(val) {
		transformed["upstreamPolicies"] = transformedUpstreamPolicies
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPolicies(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedId, err := expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesId(original["id"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedId); val.IsValid() && !isEmptyValue(val) {
			transformed["id"] = transformedId
		}

		transformedRepository, err := expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesRepository(original["repository"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedRepository); val.IsValid() && !isEmptyValue(val) {
			transformed["repository"] = transformedRepository
		}

		transformedPriority, err := expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesPriority(original["priority"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedPriority); val.IsValid() && !isEmptyValue(val) {
			transformed["priority"] = transformedPriority
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesId(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryVirtualRepositoryConfigUpstreamPoliciesPriority(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedDescription, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigDescription(original["description"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDescription); val.IsValid() && !isEmptyValue(val) {
		transformed["description"] = transformedDescription
	}

	transformedDockerRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigDockerRepository(original["docker_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedDockerRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["dockerRepository"] = transformedDockerRepository
	}

	transformedMavenRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigMavenRepository(original["maven_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMavenRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["mavenRepository"] = transformedMavenRepository
	}

	transformedNpmRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigNpmRepository(original["npm_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedNpmRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["npmRepository"] = transformedNpmRepository
	}

	transformedPythonRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigPythonRepository(original["python_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPythonRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["pythonRepository"] = transformedPythonRepository
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigDockerRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPublicRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigDockerRepositoryPublicRepository(original["public_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublicRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["publicRepository"] = transformedPublicRepository
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigDockerRepositoryPublicRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigMavenRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPublicRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigMavenRepositoryPublicRepository(original["public_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublicRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["publicRepository"] = transformedPublicRepository
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigMavenRepositoryPublicRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigNpmRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPublicRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigNpmRepositoryPublicRepository(original["public_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublicRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["publicRepository"] = transformedPublicRepository
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigNpmRepositoryPublicRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigPythonRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedPublicRepository, err := expandArtifactRegistryRepositoryRemoteRepositoryConfigPythonRepositoryPublicRepository(original["public_repository"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPublicRepository); val.IsValid() && !isEmptyValue(val) {
		transformed["publicRepository"] = transformedPublicRepository
	}

	return transformed, nil
}

func expandArtifactRegistryRepositoryRemoteRepositoryConfigPythonRepositoryPublicRepository(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
