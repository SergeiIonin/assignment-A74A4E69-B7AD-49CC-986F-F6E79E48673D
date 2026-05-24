# Todo Dashboard Test Assignment

This is a test assignment called assignment-A74A4E69-B7AD-49CC-986F-F6E79E48673D which is quite a compelling name because secrecy is our top priority!

This project aggregates and zips user data with their todo list into the dashboard similar to
```json
{
  "id": 3,
  "full_name": "Sophia Brown",
  "status": "Rookie",
  "pending_task_count": 1,
  "next_urgent_task": "Attend a local cultural festival",
  "error_warning": null
}
```

## How to Run This Project
Whatever approach is chosen, checkout the project first.

### Plain approach
1) Run `go run cmd/main.go`
2) Hit the endpoint in browser `http://localhost:8080/dashboard/:id` or run `curl http://localhost:8080/dashboard/:id`

### Run Locally
This way you will assert that dashboard endpoint serve requests in <= 2s
1) Set `chmod +x scripts/test.sh`
2) Run `./scripts/test.sh`
3) Review logs that flooded console or `.log` file created by each test run
4) Report github issue if you noticed that any of 3 test runs took more than 2 seconds

### Docker
Portability and ease of use are our second top priority! This project is distributed as Docker container (everything is in containers nowadays, right?) and for testing you only need to follow these simple steps:
1) Be Linux or MacOS user to build Docker image without pain
2) Set `chmod +x scripts/testDocker.sh`
3) Check your intuition for always inspecting executable files
4) Run `./scripts/testDocker.sh`
5) Same as 3 in `Run Locally`
