build_arm:
	docker build -t vdub-core -f Dockerfile . --build-arg ARCH=arm64

build_amd:
	docker build -t vdub-core -f Dockerfile . --build-arg ARCH=amd64

docker_run:
	docker start vdub-core || docker run -it -d -p 29000:29000 -v ./shared:/root/shared -v .:/root/vdub --name vdub-core vdub-core go run .

docker_run_win:
	docker start vdub-core || docker run -it -d -p 29000:29000 -v /home/umarkotak/umar/personal_projects/go/vdub/shared:/root/shared -v /home/umarkotak/umar/personal_projects/go/vdub:/root/vdub --name vdub-core vdub-core go run .

docker_stop:
	docker stop vdub-core

docker_rerun:
	docker stop vdub-core
	docker start vdub-core

docker_run_raw:
	docker run -dit -p 29000:29000 -v ./shared:/root/shared -v .:/root/vdub --name vdub-core vdub-core

docker_ssh:
	docker exec -it vdub-core bash

docker_rm:
	docker rm -f vdub-core
