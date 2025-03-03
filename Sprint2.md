# Sprint 2

## Progress
- Implemented extensive front‑end testing using Cypress component tests.
- Achieved high coverage for UI rendering and user interactions across multiple components.

## Unit Tests

- **List unit tests for frontend**
  - **NavBar Component Tests**
    - Renders correctly and displays navigation links.
    - Provides expected UI elements for user interactions.
  - **PostCard Component Tests**
    - Renders post details (title, content, likes, dislikes, images).
    - Handles like/dislike button interactions (alert messages).
    - Supports image modal display and closing.
  - **CommentCard Component Tests**
    - Renders comment details including author, content, likes, and dislikes.
    - Triggers expected actions (alert messages) on like/dislike buttons.
  - **Post Page Tests**
    - Renders the post title and content.
    - Displays associated comments correctly.
    - Properly handles dynamic content when no comments are available.
  - **Login Page Tests**
    - Renders username and password inputs, Login button, and register link.
    - Validates that the Login button is disabled when fields are empty and enabled when filled.
    - Toggles password visibility when the Eye/EyeOff icon is clicked.
  - **Register Page Tests**
    - Renders form fields for username, email, and password.
    - Validates input fields against required patterns (e.g., email format, password criteria).
    - Ensures the Register button is disabled/enabled based on form validity.
  - **Create Page Tests**
    - Renders input fields for title, text, and image selection.
    - Validates that images are previewed and can be removed on click.
    - Confirms that UI elements render properly without triggering actual API calls.

- **List unit tests for backend**
  - 

## Documentation
- [Swagger API Documentation for Member Endpoints](http://localhost:8080/swagger/index.html#/member)

## Frontend Unit Test Results
- [Video for Frontend Tests](https://drive.google.com/file/d/1hIu5GaK_aPSH3oM1sUGyMFDxQYmiuMam/view?usp=sharing)
