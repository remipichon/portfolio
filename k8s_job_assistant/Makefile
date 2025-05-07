############ Contribute

backend:  ## Run the Go app
	cd goapp && go run main.go

frontend:  ## Run the React app in dev mode
	cd reactapp && npm install && npm run dev

build : ## Docker Build for local use
	docker build -t kja:latest .

deploy: build ## Deploy into cluster (from kubectl and your defined kubeconfig)
	kubectl apply -k kustomize/overlays/demo

portfwd: deploy ## Create the port forward to access KJA from localhost:8080
	kubectl port-forward -n kja-demo service/kube-job-assistant 8080:8080

############ Commons

INFO_COLOR = \033[0;36m
NO_COLOR   = \033[m

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: help
