run:
	go build | go install | task_manager

db:
	docker build -t my-sqlite-image . && docker run -it -v sqlite-data:/app my-sqlite-image
