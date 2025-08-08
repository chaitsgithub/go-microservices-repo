# 🧰 Go Monorepo Utilities

This repository hosts shared utilities and tooling for all services in this Go monorepo. The goal is to ensure consistency, reuse, and ease of development across the platform.

---

## 🧱 Utility Development Roadmap

This roadmap lists the utilities to be developed, in the recommended order of priority:

---

### 🔰 Phase 1: Core Foundations

1. **Logger**  
   - Structured logging using `zap` or `logrus`
   - Consistent format across services

2. **Configuration Loader**  
   - Load from `.env`, `.yaml`, `.json`, etc. using libraries like `viper`
   - Support for override hierarchy

3. **Error Handling**  
   - Custom error types
   - Error wrapping with context and stack trace support

4. **HTTP Client Wrapper**  
   - Unified HTTP client with retry, timeout, logging, and tracing

5. **Database Utility**  
   - DB connection setup
   - Transaction helpers and migrations (e.g., with `golang-migrate`)

---

### 📊 Phase 2: Observability & Monitoring

6. **Metrics Utility**  
   - Prometheus integration for counters, gauges, histograms

7. **Tracing Utility**  
   - Distributed tracing using OpenTelemetry
   - Exporters like Jaeger or Zipkin

8. **Middleware**  
   - Request logging, panic recovery, CORS, rate limiting

---

### 🔐 Phase 3: Security & Communication

9. **JWT Utility**  
   - Token generation, validation, and claims extraction

10. **gRPC Wrappers & Interceptors**  
   - Middleware for logging, tracing, recovery, and metrics

11. **Secrets Manager Client**  
   - Integration with Vault, AWS Secrets Manager, etc.

---

### 📦 Phase 4: Advanced Utilities & Tooling

12. **Redis/Cache Wrapper**  
   - TTL support, key namespacing, retry logic

13. **Message Queue Wrapper**  
   - Kafka/NATS consumers/producers with logging and retry

14. **Mock Generators & Test Helpers**  
   - Interface mocks using `mockery`, `testify`, or `counterfeiter`

15. **Integration Test Setup**  
   - Docker test containers for DBs, MQs, etc.

16. **Code Generator Scripts**  
   - Automate protobuf, Swagger, and mock generation

17. **CLI & Task Runner**  
   - Tools for local dev workflows using `mage` or `task`

---

### 🧩 Bonus: Dev Experience Utilities

18. **Service Bootstrap Generator**  
   - CLI to scaffold new service templates with standard code structure

19. **Versioning Utility**  
   - Track and manage service/module versions in the monorepo

---

## 📁 Suggested Folder Structure
- /internal
- /logger
- /config
- /errors
- /httpclient
- /db
- /metrics
- /tracing
- /middleware
- /jwt
- /grpc
- /secrets
- /cache
- /mq
- /testutils
- /codegen
- /cli
- /bootstrap
---

### Design Patterns
- Functional Options Pattern
