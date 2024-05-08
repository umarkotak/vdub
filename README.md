# VDUB Project
This app receive youtube video url as an input then it will generate a new video using Bahasa Indonesia as an output

## API

POST /vdub/api/dubb/start
body:
```
{
  "task_name": "task-1" // task must unique
  "youtube_url": "",    // this is youtube video url
}
```

GET /vdub/api/dubb/:task_name/status
response:
```
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
