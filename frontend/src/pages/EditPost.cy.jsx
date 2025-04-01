import React from "react"
import EditPost from "./EditPost.jsx" // Adjust the import path as necessary
import { MemoryRouter, Route, Routes } from "react-router-dom"

describe("<EditPost />", () => {
  beforeEach(() => {
    // Stub the GET request so the component loads with dummy data.
    cy.intercept("GET", "**/api/v1/post/*", {
      statusCode: 200,
      body: {
        data: {
          title: "Original Title",
          content: "Original content",
          images: [],
          author: "currentUser"
        }
      }
    }).as("getPost")
  })

  const mountWithRouter = (initialEntry = "/edit-post/1") => {
    cy.mount(
      <MemoryRouter initialEntries={[initialEntry]}>
        <Routes>
          <Route path="/edit-post/:id" element={<EditPost />} />
        </Routes>
      </MemoryRouter>
    )
  }

  it("renders input fields and buttons", () => {
    mountWithRouter()
    cy.wait("@getPost")

    // Check that the title input, textarea, and both buttons exist.
    cy.get('input[placeholder="Title"]').should("exist")
    cy.get('textarea[placeholder="Type here"]').should("exist")
    cy.get("button.btn.btn-primary").should("exist") // Update Post button
    cy.get("button.btn.btn-secondary").should("exist") // Paperclip button
  })

  it("allows updating title and text", () => {
    mountWithRouter()
    cy.wait("@getPost")

    // Type new values into the inputs and verify they update.
    cy.get('input[placeholder="Title"]')
      .clear()
      .type("New Title")
      .should("have.value", "New Title")

    cy.get('textarea[placeholder="Type here"]')
      .clear()
      .type("New Content")
      .should("have.value", "New Content")
  })

  it("opens file selector when clicking the paperclip button", () => {
    mountWithRouter()
    cy.wait("@getPost")

    // Spy on the file input's click event.
    cy.get('input[type="file"]').then(($input) => {
      cy.stub($input[0], "click").as("fileInputClick")
    })

    // Click the paperclip button.
    cy.get("button.btn.btn-secondary").click()

    // Verify that the hidden file input was triggered.
    cy.get("@fileInputClick").should("have.been.called")
  })
})
