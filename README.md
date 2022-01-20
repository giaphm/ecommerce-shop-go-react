# Ecommerce shop

Ecommerce shop is a project using golang with Clean Architecture, DDD, CQRS, and React, Nextjs, Typescript

### Directories

- [api](api/) OpenAPI and gRPC definitions
- [docker](docker/) Dockerfiles
- [internal](internal/) application code
- [scripts](scripts/) deployment and development scripts
- [terraform](terraform/) - infrastructure definition
- [web](web/) - frontend JavaScript code(React, Nextjs, Typescript)

### Running locally

```go
> docker-compose up

# ...

web_1             |  INFO  Starting development server...
web_1             |  DONE  Compiled successfully in 6315ms11:18:26 AM
web_1             |
web_1             |
web_1             |   App running at:
web_1             |   - Local:   http://localhost:3000/
web_1             |
web_1             |   It seems you are running Vue CLI inside a container.
web_1             |   Access the dev server via http://localhost:<your container's external mapped port>/
web_1             |
web_1             |   Note that the development build is not optimized.
web_1             |   To create a production build, run yarn build.
```
