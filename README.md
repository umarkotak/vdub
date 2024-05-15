# VDUB Project
This app receive youtube video url as an input then it will generate a new video using Bahasa Indonesia as an output

## API

```
POST /vdub/api/dubb/start
body:
{
  "task_name": "task-1" // task must unique
  "youtube_url": "",    // this is youtube video url
}
```

```
GET /vdub/api/dubb/:task_name/status
response:
{
  "status": "initialized",
  "progress": [
    {
      "status": "initialized",   //
      "progress": "done"         // incompleted, processing, completed
    }
  ]
}
```

## Tasking Structure
for every task will be stored on their respective folder

```
_ /base/dir
|__ /task-1
  |__
|__ /task-2
|__ /task-n
```


## Docker Commands

```
# 1 To run the docker

-- On macbook
docker run -dit -p 29000:29000 -v /Users/umarramadhana/umar/personal_projects/vdub/shared:/root/shared -v /Users/umarramadhana/umar/personal_projects/vdub/go-be:/root/go-be --name vdub-core vdub-core

-- On windows
docker run -dit -p 29000:29000 -v /home/umarkotak/umar/personal_projects/vdub/shared:/root/shared -v /home/umarkotak/umar/personal_projects/vdub/go-be:/root/go-be --name vdub-core vdub-core

# 2 To ssh into the container

docker exec -it vdub-core bash

# 3
```

python -m bark --text "Hello, my name is Suno." --output_filename "example.wav"
