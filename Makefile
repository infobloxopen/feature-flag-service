include Makefile.vars
include Makefile.common

# to override a target from Makefile.common just redefine the target.
# you can also chain the original atlas target by adding
# -atlas to the dependency of the redefined target

include Makefile.controller

CHART := feature-flag
CHART_VERSION := $(IMAGE_VERSION)
CHART_FILE := $(CHART)-$(CHART_VERSION).tgz
HELM_IMAGE ?= infoblox/helm:2.14.3-1

.helm-lint:
	cd deploy && helm lint $(CHART) -f $(CHART)/minikube-values.yaml

.PHONY: push-chart
push-chart: AWS_ACCESS_KEY_ID?=`aws configure get aws_access_key_id`
push-chart: AWS_SECRET_ACCESS_KEY?=`aws configure get aws_secret_access_key`
push-chart: AWS_REGION?=`aws configure get region`
push-chart: .helm-lint deploy/$(CHART_FILE) build.properties
	docker run -e AWS_REGION=${AWS_REGION} \
		-e AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
		-e AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
		-v $(PWD)/deploy:/pkg \
		${HELM_IMAGE} s3 push /pkg/$(CHART_FILE) infobloxcto

deploy/$(CHART_FILE):
	cd deploy && helm package $(CHART) --version $(CHART_VERSION)

build.properties: build.properties.in
	@sed 's/{CHART_FILE}/$(CHART_FILE)/g' build.properties.in > build.properties
