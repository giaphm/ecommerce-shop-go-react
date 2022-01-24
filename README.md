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

web_1                        | audited 658 packages in 63.776s
web_1                        |
web_1                        | > web-next@ dev /web
web_1                        | > next dev
web_1                        |
web_1                        | ready - started server on 0.0.0.0:3000, url: http://localhost:3000
web_1                        | event - compiled client and server successfully in 7.9s (469 modules)
```

### Testing
```go
> make test
```
