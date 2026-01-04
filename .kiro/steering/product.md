# Product Overview

A production-ready REST API built with Go Fiber framework. The API provides user management and authentication functionality with JWT-based security.

## Core Features

- User registration and management (CRUD operations)
- JWT authentication with role-based access control
- Swagger/OpenAPI documentation
- Health check endpoint with database status
- Pagination support for list endpoints

## API Structure

- Base path: `/api/v1`
- Auth endpoints: `/auth/login`, `/auth/me`
- User endpoints: `/users` (CRUD)
- Documentation: `/swagger/*`
- Health: `/health`
