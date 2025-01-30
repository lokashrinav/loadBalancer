# Golang Load Balancer

This is a simple load balancer implemented in Go that distributes incoming HTTP requests across multiple backend servers using a round-robin algorithm. It also performs periodic health checks to ensure backend availability.

## Features
- Round-robin load balancing
- Health checks for backend servers
- Concurrent request handling

## Installation
Ensure you have [Go installed](https://go.dev/dl/). Clone this repository and navigate to the project directory:

```sh
git clone https://github.com/yourusername/golang-load-balancer.git
cd golang-load-balancer
```

## Usage
1. Start multiple backend servers on different ports (e.g., `5001`, `5002`, `5003`).
2. Update the `servers` list in `main.go` with your backend URLs.
3. Run the load balancer:

```sh
go run main.go
```

4. Send requests to `http://localhost:8080/` to see the load balancer in action.

## Backend Health Checks
The load balancer periodically checks the health of backend servers by sending requests to `/health`. Ensure your backend servers implement a `/health` endpoint that returns a `200 OK` status when they are operational.

## Example Backend Server
Run a simple backend server using Python:

```sh
python -m http.server 5001
```

Repeat for additional ports (e.g., `5002`, `5003`).

## Contributing
Feel free to open issues or submit pull requests to improve this project.

## License
This project is licensed under the MIT License.

