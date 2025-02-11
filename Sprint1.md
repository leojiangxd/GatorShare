# Sprint 1

## User Stories
- We identified 15 key issues to address in the project, primarily focusing on design and user requirements.
- Additional user stories may be added as the project scope evolves.

## Planned Issues
- For Sprint 1, we aimed to complete 8 critical issues:
  - User Authentication & Management: Login, registration, logout, profile updates, and user deletion.
  - Core Features: Feed, comments, and posts.

## Progress & Completed Issues
- The frontend team initialized the React application and installed the necessary libraries and dependencies, including Tailwind CSS.
- The UI for all 8 planned issues has been successfully designed.
- The backend team initialized the sqlite database and set up a basic structure for the api in go
- 4 api calls were fully implemented
  - GET user
  - GET user/:id 
  - POST user
  - DELETE user/:id

## Incomplete Issues
- The frontend team prioritized defining the website’s overall look before committing to full implementation.
  - While the UI is mostly complete, more work still needs to be done to integrate backend functionality and the current UI is subject to change.
- While some api structure was made for the user table, we didn't set aside enough time to fully develop all the planned functionality
  - PUT user/:id was exhibiting odd behavior just before the demo and we were unable to debug in time
  - A table for the posts as well as one for comments was discussed but not implemented
  - API was not integrated into the frontend due to the lack of many needed calls

## Demos
- [Frontend](https://drive.google.com/file/d/1vd5z2KagNu0IWx2A81BlvunXyC8bc0lM/view?usp=drive_link)
- [Backend](https://youtu.be/AS4bCdQDr0s)
