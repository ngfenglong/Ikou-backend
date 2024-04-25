build:
	@echo Building...
	@go build -o dist/api ./cmd
	@echo Backend built!

start: build
	@echo Starting the Backend ...
	./dist/api
	@echo Backend running!

stop:
	@echo Stopping the Backend...
	@pkill -f api
	@echo Stopped Backend