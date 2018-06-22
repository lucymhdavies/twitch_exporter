.PHONY: build up push clean
.DEFAULT_GOAL: build


build:
	@echo ========================================
	@echo Building Docker Images
	@echo ========================================
	
	docker-compose build --pull

up:
	@echo ========================================
	@echo Running Locally
	@echo ========================================
	
	docker-compose up

push:
	@echo ========================================
	@echo Pushing Docker Images
	@echo ========================================
	
	docker-compose push

clean:
	@echo ========================================
	@echo Cleaning Up
	@echo ========================================
	
	docker-compose rm -s -f
