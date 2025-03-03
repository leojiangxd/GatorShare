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
- Unit tests for user and post API calls
- Swagger documentation generated based on formatted comments for each API call

Incomplete:
- Unit tests for comment API
- Messaging API
- Like/dislike API

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

**Backend Unit Tests:**
- **TestRegister:**
  - Creates a new user
  - This test user is used consistently in future unit tests
- **TestGetMembers:**
  - Gets the first 10 members
- **TestGetMemberByUsername:**
  - Gets the user with username 'saul'
  - Validates response against the test user (email, username, password, bio)
- **TestLogout:**
  - Logs out the current user
- **TestLogin:**
  - Logs in a user based on the passed in username and password
- **TestUpdateMember:**
  - Updates the test user's information
  - Validates response against updated test user (email, username, password, bio)
- **TestGetCurrentUser:**
  - Gets the username of the currently logged in user based on the session token
- **TestDeleteMember:**
  - Deletes the user whose username is passed in
- **TestCreatePost:**
  - Creates a new test user and logs them in for authorization purposes
  - Creates a new post with the test user as the author
- **TestGetPosts:**
  - Gets posts in descending order based on CreatedAt (latest first)
- **TestGetPostById:**
  - Creates a new post
  - Gets the post using the postID generated
- **TestUpdatePost:**
  - Creates a new post
  - Updates the post using the postID generated
- **TestDeletePost:**
  - Creates a new post
  - Deletes the post using the postID generated
- **TestGetUserPosts:**
  - Gets all the posts created by the test user
- **TestIncrementPostViews:**
  - Creates a new post
  - Increments the view count of the post

## Documentation
- In the backend, install Swagger using `go install github.com/swaggo/swag/cmd/swag@latest`
- If docs.go is throwing an error, run `go get -u github.com/swaggo/swag`
- Then build the backend application using `go build` and run `platform.exe`
- Once the backend application is running, the documentation can be found at: http://localhost:8080/swagger/index.html#/member

## Demos
- [Integrated Application Demo](https://youtu.be/WVI0q96o_bc?si=bK06mV_s4whQu8LY)
- [Frontend Unit Testing Demo](https://drive.google.com/file/d/1hIu5GaK_aPSH3oM1sUGyMFDxQYmiuMam/view?usp=sharing)
- [FullVideo] (COMBINED LINK HERE)
