import React from 'react'
import PostCard from './PostCard'
import { MemoryRouter } from 'react-router-dom'
import axios from 'axios'

describe('<PostCard />', () => {
  const post = {
    post_id: '123',
    title: 'Test Post Title',
    author: 'testuser',
    CreatedAt: '2023-03-02T10:00:00Z',
    comments: [{ content: 'Nice post' }],
    views: 100,
    likes: 20,
    dislikes: 5,
    content: 'This is the full content of the test post.',
    images: ['http://example.com/image1.jpg']
  }

  beforeEach(() => {
    // Stub window.alert for like/dislike tests.
    cy.stub(window, 'alert').as('alertStub')
  })

  it('renders in preview mode and increments view count on click', () => {
    // Intercept the PUT request for view count increment.
    cy.intercept('PUT', new RegExp(`/api/v1/post/${post.post_id}/increment-views`), {
      statusCode: 200
    }).as('incrementView')

    cy.mount(
      <MemoryRouter>
        <PostCard post={post} preview={true} />
      </MemoryRouter>
    )

    // In preview mode, the card is wrapped in a Link.
    cy.get('a').first().click()

    // Wait for the intercepted request and assert that the URL contains the correct endpoint.
    cy.wait('@incrementView').its('request.url').should('include', `/api/v1/post/${post.post_id}/increment-views`)
  })

  it('triggers like and dislike alerts when their buttons are clicked', () => {
    cy.mount(
      <MemoryRouter>
        <PostCard post={post} preview={false} />
      </MemoryRouter>
    )

    // Click the like button and verify the alert message.
    cy.get('button').contains(`${post.likes}`).click()
    cy.get('@alertStub').should('have.been.calledWith', 'liked post')

    // Click the dislike button and verify the alert message.
    cy.get('button').contains(`${post.dislikes}`).click()
    cy.get('@alertStub').should('have.been.calledWith', 'disliked post')
  })

  it('opens and closes image modal when an image is clicked', () => {
    cy.mount(
      <MemoryRouter>
        <PostCard post={post} preview={false} />
      </MemoryRouter>
    )

    // Click on the first image to open the modal.
    cy.get('img').first().click()

    // Use a specific selector for the modal overlay (assuming it uses "fixed inset-0").
    cy.get('div.fixed.inset-0').should('be.visible').within(() => {
      // Verify that the modal displays the correct image.
      cy.get('img').should('have.attr', 'src', post.images[0])
    })

    // Click on the modal overlay to close it.
    cy.get('div.fixed.inset-0').click({ force: true })

    // Assert that the modal overlay is removed.
    cy.get('div.fixed.inset-0').should('not.exist')
  })
})
