# Development Plan

This document outlines potential future work and improvements for the Hexya NG project.

## Short-Term Goals

- **Database Integration:** Replace the in-memory data stores with a persistent database like PostgreSQL. This will involve:
  - Choosing a database driver and ORM.
  - Updating the services to use the database.
  - Managing database migrations.

- **Configuration Management:** Move hardcoded values like service URLs and ports to configuration files or environment variables.

- **Automated Testing:**
  - Add unit tests for each service.
  - Add integration tests to verify the communication between services.
  - Set up a CI/CD pipeline to run tests automatically.

## Medium-Term Goals

- **Authentication and Authorization:**
  - Implement a user management service.
  - Add authentication (e.g., JWT) to the API gateway.
  - Implement role-based access control (RBAC) to restrict access to certain endpoints.

- **Service Discovery:** Implement a service discovery mechanism (e.g., Consul, etcd) to allow services to find each other dynamically.

- **Observability:**
  - Add structured logging to all services.
  - Implement distributed tracing to track requests as they flow through the system.
  - Set up monitoring and alerting to track the health of the services.

## Long-Term Goals

- **Web Frontend:** Create a web-based user interface to interact with the system. This could be a single-page application (SPA) built with a modern framework like React or Vue.

- **Asynchronous Communication:** Introduce a message broker (e.g., RabbitMQ, Kafka) for asynchronous communication between services. This would be useful for tasks like sending email notifications or processing long-running jobs.

- **Framework Extraction:** Refactor the code to extract a generic framework that can be used to build other applications. This would involve creating a set of libraries and tools for common tasks like database access, service communication, and configuration.
