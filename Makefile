build_base:
	docker build -t vdub-base -f DockerfileBase .

build_arm:
	docker build -t vdub-core-arm -f DockerfileArm .

build:
	docker build -t vdub-core -f Dockerfile .

docker_run_arm:
	docker start vdub-core-arm || docker run -it -d -p 29000:29000 -v ./shared:/root/shared -v .:/root/vdub --name vdub-core-arm vdub-core-arm go run .

docker_run:
	docker start vdub-core || docker run -it -d -p 29000:29000 -v ./shared:/root/shared -v .:/root/vdub --name vdub-core vdub-core go run .

docker_run_win:
	docker start vdub-core || docker run -it -d -p 29000:29000 -v /home/umarkotak/umar/personal_projects/go/vdub/shared:/root/shared -v /home/umarkotak/umar/personal_projects/go/vdub:/root/vdub --name vdub-core vdub-core go run .

docker_run_win_gpu:
	docker start vdub-core || docker run -it -d --gpus all -p 29000:29000 -v /home/umarkotak/umar/personal_projects/go/vdub/shared:/root/shared -v /home/umarkotak/umar/personal_projects/go/vdub:/root/vdub --name vdub-core vdub-core go run .

docker_stop_arm:
	docker stop vdub-core-arm

docker_stop:
	docker stop vdub-core

docker_rerun_arm:
	docker stop vdub-core-arm
	docker start vdub-core-arm

docker_rerun:
	docker stop vdub-core
	docker start vdub-core

docker_run_raw:
	docker run -dit -p 29000:29000 -v ./shared:/root/shared -v .:/root/vdub --name vdub-core vdub-core

docker_ssh_arm:
	docker exec -it vdub-core-arm bash

docker_ssh:
	docker exec -it vdub-core bash

docker_rm_arm:
	docker rm -f vdub-core-arm

docker_rm:
	docker rm -f vdub-core
