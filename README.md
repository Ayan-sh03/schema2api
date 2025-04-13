# Schema2API

A dynamic mock API generator that creates REST endpoints from JSON schemas. Perfect for testing and prototyping.

## Features

- Generate REST API endpoints from JSON schemas
- Supports both integer and string IDs
- Full CRUD operations (Create, Read, Update, Delete)
- Dynamic response generation based on schema types
- Containerized with Docker for easy deployment

## Quick Start

### Running Locally

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Ayan-sh03/schema2api.git
   cd schema2api
   ```

2. **Run the server:**
   ```bash
   go run main.go
   ```

3. **Upload a Schema:**
   Use a sample schema (e.g., `user_schema.json`):
   ```bash
   curl -X POST -H "Content-Type: application/json" --data @user_schema.json http://localhost:8081/upload
   ```

4. **Interact with the API:**
   - **GET List:**
     `curl http://localhost:8081/users`
   - **GET Single:**
     `curl http://localhost:8081/users/123`
   - **POST:**
     `curl -X POST -H "Content-Type: application/json" -d '{"name":"John", "email":"john@example.com"}' http://localhost:8081/users`
   - **PUT:**
     `curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Name"}' http://localhost:8081/users/123`
   - **DELETE:**
     `curl -X DELETE http://localhost:8081/users/123`

### Running with Docker

1. **Build the Docker image:**
   ```bash
   docker build -t schema2api .
   ```

2. **Run the container:**
   ```bash
   docker run -d -p 8081:8081 --name schema2api schema2api
   ```

## Testing

- **Go Unit Tests:**
  ```bash
  go test
  ```


## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add some amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.