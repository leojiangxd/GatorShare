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
- Follower system


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
