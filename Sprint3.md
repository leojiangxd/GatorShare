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

Incomplete:
- Liking/disliking comments
- Messaging/notification system

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
