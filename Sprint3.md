# Sprint 3

## Progress
### Frontend
Completed:
- Editing user profile (change bio, username, password, email and delete)
- Liking/disliking posts
- Pagination of posts on home page
- Editing and deleting posts/comments
- Sorting posts by creation, likes, dislikes, views, and comments
- Searching for posts by author, title, or content
- Messaging/notification system

Incomplete:
- Liking/disliking comments


### Backend
Completed:
- Sorting members by creation, username, bio, or email
- Searching for members by username, email, or bio
- Sorting posts by creation, likes, dislikes, views, and comments
- Searching for posts by author, title, or content
- Seperating member or post search results into slices based the limit and offset parameters
- Added total count of members or posts in the search query (not including limits or offsets)
- Updated unit tests and documentation for sorting and searching members and posts
- Liking or dislinking posts
- Getting list of posts liked by the current user

Incomplete:
- Messaging API
- Liking/disliking for comments

## Unit Tests

**Frontend Unit Tests:**
- **NavBar Component:**
    - Renders correctly and displays expected navigation elements.
    - Integrates properly with React Router.
    - Added testing for any API calls related to user session and navigation state.
- **Login Component:**
    - Renders username and password fields, Login button, and register link.
    - Login button is initially disabled and becomes enabled when both fields are filled.
    - Password visibility toggles correctly via the Eye/EyeOff icon.
    - Added testing for the authentication API call.
- **Register Component:**
    - Renders form fields for username, email, and password.
    - Validates input against required patterns (e.g., alphanumeric username, UF email format, password criteria).
    - Enables/disables the Register button based on form validity.
    - Added testing for the registration API endpoint.
- **PostCard Component:**
    - Renders post details (title, content, likes/dislikes, images).
    - Handles like/dislike interactions (e.g., triggering alerts).
    - Displays image modal when an image is clicked.
    - Added testing for API calls associated with post interactions.
- **CommentCard Component:**
    - Renders comment details including author, content, likes, and dislikes.
    - Handles like/dislike interactions appropriately.
    - Added testing for API calls associated with comment interactions.
- **Post Page:**
    - Renders post details and its associated comments.
    - Handles URL parameters via React Router correctly.
    - Added testing for API calls related to retrieving post details and associated comments.
- **Create Page:**
    - Renders a post creation form with title, text, and image upload.
    - Previews selected images and allows removal of images.
    - Added testing for API calls involved in post creation.
- **User Profile Page:**
    - Renders user profile information (avatar, bio, edit options).
    - Displays the list of posts from the specific user.
    - Added testing for API calls related to fetching and updating user data.
- **Home (User) Component:**
    - Renders without errors, showing NavBar and a list of posts.
    - Added testing for API calls that fetch posts data.
- **Notifications Page:**
    - Renders all notifications, sorted properly by date, and user is able to mark as read/unread.
    - Added testing for API calls related to notification retrieval and updates.
- **Edit Post Page:**
    - Renders all components, allowing the user to update the title and add attachments to the existing post.
    - Added testing for API calls related to post editing and attachment management.

**Backend Unit Tests:**
- **TestRegister:**
  - Creates a new user
  - This test user is used consistently in future unit tests
- **TestGetMembers:**
  - Searches for the first 10 members containing 'test', sorted by created_at descending
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
  - Searches for the first 10 posts containing 'test', sorted by created_at descending
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
- **TestLikePost:**
  - Creates a new post
  - Adds a like to the post under the test user
- **TestGetUserLikedPosts:**
  - Retrieves a list of all posts liked by the test user
- **TestDislikePost:**
  - Creates a new post
  - Adds a dislike to the post under the test user

## Documentation
- In the backend, install Swagger using `go install github.com/swaggo/swag/cmd/swag@latest`
- If docs.go is throwing an error, run `go get -u github.com/swaggo/swag`
- Then build the backend application using `go build` and run `platform.exe`
- Once the backend application is running, the documentation can be found at: http://localhost:8080/swagger/index.html#/member

- PDF version of the swagger documentation in the github repo as api-documentation.pdf

## Sprint Video
