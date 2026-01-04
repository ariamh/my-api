# Implementation Plan: Handler Unit Tests

## Overview

Implement comprehensive table-driven unit tests for all REST API handlers using Go's httptest package and testify for assertions/mocking.

## Tasks

- [x] 1. Create auth handler test file with mock service
  - [x] 1.1 Create MockAuthService implementing service.AuthService interface
    - Implement Login method with mock.Mock
    - _Requirements: 1.1-1.5_
  - [x] 1.2 Create setupAuthTestApp helper function
    - Initialize validator
    - Create Fiber app with auth routes (/auth/login, /auth/me)
    - _Requirements: 1.1, 2.1_

- [x] 2. Implement auth handler login tests
  - [x] 2.1 Implement tests for Login endpoint
    - Test case: valid login returns 200 with token and user
    - Test case: invalid JSON body returns 400
    - Test case: validation failure returns 422
    - Test case: invalid credentials returns 401
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 2.2 Add missing Login service error test
    - Test case: service error returns 500
    - _Requirements: 1.5_

- [x] 3. Implement auth handler me endpoint tests
  - [x] 3.1 Implement table-driven tests for Me endpoint
    - Test case: returns user context values with 200 status
    - Test case: includes all context fields (user_id, email, role)
    - _Requirements: 2.1, 2.2_

- [x] 4. Checkpoint - Verify auth handler tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Refactor user handler Create tests to table-driven pattern
  - [x] 5.1 Refactor Create endpoint tests to table-driven
    - Convert existing tests (TestUserHandler_Create_Success, TestUserHandler_Create_ValidationError) to table-driven structure
    - Add missing test cases: invalid JSON returns 400, duplicate email returns 400, service error returns 500
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 6. Refactor user handler FindByID tests to table-driven pattern
  - [x] 6.1 Refactor FindByID endpoint tests to table-driven
    - Convert existing tests (TestUserHandler_FindByID_Success, TestUserHandler_FindByID_NotFound) to table-driven structure
    - Add missing test case: service error returns 500
    - _Requirements: 4.1, 4.2, 4.3_

- [x] 7. Implement user handler FindAll tests
  - [x] 7.1 Implement table-driven tests for FindAll endpoint
    - Test case: default pagination (no params) returns 200
    - Test case: custom pagination params returns 200
    - Test case: invalid page (< 1) normalized to 1
    - Test case: invalid per_page (< 1 or > 100) normalized to 10
    - Test case: service error returns 500
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_
  - [ ]* 7.2 Write property test for pagination normalization
    - **Property 2: Pagination Parameter Normalization**
    - **Validates: Requirements 5.3, 5.4**

- [x] 8. Implement user handler Update tests
  - [x] 8.1 Implement table-driven tests for Update endpoint
    - Test case: valid update returns 200
    - Test case: invalid JSON returns 400
    - Test case: validation failure returns 422
    - Test case: not found returns 404
    - Test case: service error returns 500
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 9. Refactor user handler Delete tests to table-driven pattern
  - [x] 9.1 Refactor Delete endpoint tests to table-driven
    - Convert existing test (TestUserHandler_Delete_Success) to table-driven structure
    - Add missing test cases: not found returns 404, service error returns 500
    - _Requirements: 7.1, 7.2, 7.3_

- [ ] 10. Final checkpoint - Verify all tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - Run `make test` to verify full test suite

## Notes

- Each task references specific requirements for traceability
- Existing MockUserService in user_handler_test.go will be reused
- Property tests use rapid library for Go property-based testing
- Tasks marked with `*` are optional and can be skipped for faster MVP
