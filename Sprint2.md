# Sprint 2

## Progress
### Frontend
Completed:
- Integrate login, register, logout api with frontend
- Integrate post creation and comment api with frontend
- Add and view images in posts
- Search for posts by titles
- View posts from a specific user on their profile
- Increment view count for posts
- Unit testing with cypress

Incomplete:
- Integrating backend and frontend
- Update user profile implementation
- Like/dislike posts and comments implementation
- Messaging system
- Notification system
### Backend
Completed:
- Created database schema for Posts and Comments
- Get current logged in user api
- Create posts, get all posts, get specific post, get all posts from a specific user, updating a post, deleting a specifc post
- Create comments, get all comments in a post, get a specific comment, deleting a specific comment

Incomplete:
- Unit tests
- Swagger Documentation

## Unit Tests

**Frontend Unit Tests:**
- **NavBar Component:**
  - Renders correctly and displays expected navigation elements.
  - Integrates properly with React Router.
- **Login Component:**
  - Renders username and password fields, Login button, and register link.
  - Login button is initially disabled and becomes enabled when both fields are filled.
  - Password visibility toggles correctly via the Eye/EyeOff icon.
- **Register Component:**
  - Renders form fields for username, email, and password.
  - Validates input against required patterns (e.g., alphanumeric username, UF email format, password criteria).
  - Enables/disables the Register button based on form validity.
- **PostCard Component:**
  - Renders post details (title, content, likes/dislikes, images).
  - Handles like/dislike interactions (e.g., triggering alerts).
  - Displays image modal when an image is clicked.
- **CommentCard Component:**
  - Renders comment details including author, content, likes, and dislikes.
  - Handles like/dislike interactions appropriately.
- **Post Page:**
  - Renders post details and its associated comments.
  - Handles URL parameters via React Router correctly.
- **Create Page:**
  - Renders a post creation form with title, text, and image upload.
  - Previews selected images and allows removal of images.
- **User Profile Page:**
  - Renders user profile information (avatar, bio, edit options).
  - Displays the list of posts from the specific user.
- **Home (User) Component:**
  - Renders without errors, showing NavBar and a list of posts.

- List unit tests for backend

## Documentation
- In the backend, install Swagger using `go install github.com/swaggo/swag/cmd/swag@latest`
- Then build the backend application using `go build` and run `platform.exe`
- Once the backend application is running, the documentation can be found at: http://localhost:8080/swagger/index.html#/member

## Demos
- [Integrated Application Demo](https://youtu.be/WVI0q96o_bc?si=bK06mV_s4whQu8LY)
- [Frontend Unit Testing Demo](https://drive.google.com/file/d/1hIu5GaK_aPSH3oM1sUGyMFDxQYmiuMam/view?usp=sharing)
