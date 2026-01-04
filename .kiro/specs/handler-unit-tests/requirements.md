# Requirements Document

## Introduction

This feature adds comprehensive automated unit testing for all REST API handlers using Go's table-driven test pattern with httptest. The goal is to ensure all handler endpoints are thoroughly tested with consistent patterns, covering success cases, error conditions, and edge cases.

## Glossary

- **Handler**: HTTP request handler that receives requests, validates input, and calls services
- **Table_Driven_Test**: A testing pattern where test cases are defined as a slice of structs, each containing inputs and expected outputs
- **Mock_Service**: A test double that implements a service interface for isolated handler testing
- **HTTP_Test**: Go's httptest package for testing HTTP handlers without a running server
- **Test_Case**: A single entry in a table-driven test containing name, setup, input, and expected output

## Requirements

### Requirement 1: Auth Handler Login Tests

**User Story:** As a developer, I want comprehensive tests for the login endpoint, so that I can verify authentication works correctly.

#### Acceptance Criteria

1.1. WHEN a valid login request is submitted, THE Auth_Handler SHALL return a 200 status with a JWT token and user data
1.2. WHEN an invalid JSON body is submitted to login, THE Auth_Handler SHALL return a 400 status with an error message
1.3. WHEN login credentials fail validation, THE Auth_Handler SHALL return a 422 status with validation errors
1.4. WHEN invalid credentials are provided, THE Auth_Handler SHALL return a 401 status with an unauthorized message
1.5. WHEN the auth service returns an unexpected error, THE Auth_Handler SHALL return a 500 status

### Requirement 2: Auth Handler Me Endpoint Tests

**User Story:** As a developer, I want tests for the /auth/me endpoint, so that I can verify user context retrieval works correctly.

#### Acceptance Criteria

2.1. WHEN an authenticated request is made to /auth/me, THE Auth_Handler SHALL return a 200 status with user_id, email, and role from context
2.2. WHEN context values are set, THE Auth_Handler SHALL include all context values in the response

### Requirement 3: User Handler Create Tests

**User Story:** As a developer, I want comprehensive tests for user creation, so that I can verify registration works correctly.

#### Acceptance Criteria

3.1. WHEN a valid user creation request is submitted, THE User_Handler SHALL return a 201 status with the created user data
3.2. WHEN an invalid JSON body is submitted, THE User_Handler SHALL return a 400 status with an error message
3.3. WHEN user input fails validation, THE User_Handler SHALL return a 422 status with validation errors
3.4. WHEN the email already exists, THE User_Handler SHALL return a 400 status with a duplicate email error
3.5. WHEN the user service returns an unexpected error, THE User_Handler SHALL return a 500 status

### Requirement 4: User Handler FindByID Tests

**User Story:** As a developer, I want comprehensive tests for fetching users by ID, so that I can verify user retrieval works correctly.

#### Acceptance Criteria

4.1. WHEN a valid user ID is requested, THE User_Handler SHALL return a 200 status with the user data
4.2. WHEN a non-existent user ID is requested, THE User_Handler SHALL return a 404 status with a not found error
4.3. WHEN the user service returns an unexpected error, THE User_Handler SHALL return a 500 status

### Requirement 5: User Handler FindAll Tests

**User Story:** As a developer, I want comprehensive tests for listing users with pagination, so that I can verify pagination works correctly.

#### Acceptance Criteria

5.1. WHEN users are requested without pagination params, THE User_Handler SHALL return a 200 status with default pagination (page 1, 10 per page)
5.2. WHEN users are requested with valid pagination params, THE User_Handler SHALL return a 200 status with the specified pagination
5.3. WHEN invalid page number is provided (less than 1), THE User_Handler SHALL normalize it to page 1
5.4. WHEN invalid per_page is provided (less than 1 or greater than 100), THE User_Handler SHALL normalize it to 10
5.5. WHEN the user service returns an error, THE User_Handler SHALL return a 500 status

### Requirement 6: User Handler Update Tests

**User Story:** As a developer, I want comprehensive tests for updating users, so that I can verify user updates work correctly.

#### Acceptance Criteria

6.1. WHEN a valid update request is submitted, THE User_Handler SHALL return a 200 status with the updated user data
6.2. WHEN an invalid JSON body is submitted, THE User_Handler SHALL return a 400 status with an error message
6.3. WHEN update input fails validation, THE User_Handler SHALL return a 422 status with validation errors
6.4. WHEN updating a non-existent user, THE User_Handler SHALL return a 404 status with a not found error
6.5. WHEN the user service returns an unexpected error, THE User_Handler SHALL return a 500 status

### Requirement 7: User Handler Delete Tests

**User Story:** As a developer, I want comprehensive tests for deleting users, so that I can verify user deletion works correctly.

#### Acceptance Criteria

7.1. WHEN a valid delete request is submitted, THE User_Handler SHALL return a 204 status with no content
7.2. WHEN deleting a non-existent user, THE User_Handler SHALL return a 404 status with a not found error
7.3. WHEN the user service returns an unexpected error, THE User_Handler SHALL return a 500 status

### Requirement 8: Table-Driven Test Pattern

**User Story:** As a developer, I want all tests to follow the table-driven pattern, so that tests are consistent and maintainable.

#### Acceptance Criteria

8.1. THE Test_Suite SHALL define test cases as a slice of structs with name, setup function, request builder, and expected assertions
8.2. THE Test_Suite SHALL iterate over test cases using t.Run with the test case name as subtest name
8.3. THE Test_Suite SHALL use testify/assert for assertions and testify/mock for service mocking
8.4. THE Test_Suite SHALL isolate each test case with fresh mock instances
