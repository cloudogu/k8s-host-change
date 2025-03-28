#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# this function will be sourced from release.sh and be called from release_functions.sh
update_versions_modify_files() {
  newReleaseVersion="${1}"
  valuesYAML=k8s/helm/values.yaml
  componentPatchTplYAML=k8s/helm/component-patch-tpl.yaml

  ./.bin/yq -i ".job.image.tag = \"${newReleaseVersion}\"" "${valuesYAML}"
  ./.bin/yq -i ".values.images.job = \"cloudogu/k8s-host-change:${newReleaseVersion}\" " "${componentPatchTplYAML}"
}

update_versions_stage_modified_files() {
  valuesYAML=k8s/helm/values.yaml
  componentPatchTplYAML=k8s/helm/component-patch-tpl.yaml

  git add "${valuesYAML}" "${componentPatchTplYAML}"
}
