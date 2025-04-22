# Sprint 4

## Progress
### Frontend
Completed:
- Creating notifications
- Marking notifications as read/unread
- Marking all notifications as read/unread
- Deleting notifications
- Liking/disliking comments

Incomplete:
- Messaging system

### Backend
Completed:
- Creating notification API
- Updating notifications API
- Getting notifications API
- Deleting notifications API
- Liking/disliking comments API
- Follow/Unfollow API

Incomplete:
- Messaging API

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
- **TestGetComments:**
  - Creates a new post and comment
  - Gets all comments from the created post
- **TestGetCommentById:**
  - Creates a new post and comment
  - Gets the created comment based on its commentId
- **TestCreateComment:**
  - Creates a new post and comment
- **TestUpdateComment:**
  - Creates a new post and comment
  - Updates the comment content
- **TestDeleteComment:**
  - Creates a new post and comment
  - Deletes the created comment
- **TestLikeOrDislikeComment:**
  - Creates a new post and comment
  - Adds a like and a dislike to the comment
- **TestSendNotification:**
  - Creates a new notification
- **TestGetNotifications:**
  - Gets all notifications of the test user
- **TestGetNotificationById:**
  - Creates a new notification
  - Fetches the notificationId and then gets the notification by the id
- **TestUpdateNotification:**
  - Creates a new notification
  - Fetches the notificationId and then updates the read status to true
- **TestUpdateNotifications:**
  - Creates a new notification
  - Updates the read status to true for all of the test users notifications
- **TestDeleteNotification:**
  - Creates a new notification
  - Fetches the notificationId and then deletes the notification
- **TestFollowUnfollowSystem:**
  - Creates a new test user 'gus'
  - The original test user 'saul' follows 'gus'
  - Gets users followed by 'saul'
  - Gets users following 'gus'
  - Unfollows 'gus' and validates that the followers and following sets are empty


## Documentation
- In the backend, install Swagger using `go install github.com/swaggo/swag/cmd/swag@latest`
- If docs.go is throwing an error, run `go get -u github.com/swaggo/swag`
- Then build the backend application using `go build` and run `platform.exe`
- Once the backend application is running, the documentation can be found at: http://localhost:8080/swagger/index.html#/member

- PDF version of the swagger documentation in the github repo as api-documentation.pdf

## Sprint Video
https://youtu.be/4minjXcBX1s