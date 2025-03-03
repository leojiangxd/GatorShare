import React from 'react'
import CommentCard from './CommentCard'
import { MemoryRouter } from 'react-router-dom'

describe('<CommentCard />', () => {
  const comment = {
    author: 'testuser',
    CreatedAt: '2023-03-02T10:00:00Z',
    likes: 10,
    dislikes: 3,
    content: 'This is a test comment',
  }

  beforeEach(() => {
    cy.stub(window, 'alert').as('alertStub')
  })

  it('renders comment details correctly', () => {
    cy.mount(
      <MemoryRouter>
        <CommentCard comment={comment} />
      </MemoryRouter>
    )
    
    cy.contains(comment.author).should('be.visible')
    cy.get('a').should('have.attr', 'href', `/user/${comment.author}`)
    
    cy.contains(comment.content).should('be.visible')
    
    cy.contains(`${comment.likes}`).should('be.visible')
    cy.contains(`${comment.dislikes}`).should('be.visible')
    

  })

  it('triggers the "liked comment" alert when the like button is clicked', () => {
    cy.mount(
      <MemoryRouter>
        <CommentCard comment={comment} />
      </MemoryRouter>
    )

    cy.get('button').first().click()
    cy.get('@alertStub').should('have.been.calledWith', 'liked comment')
  })

  it('triggers the "disliked comment" alert when the dislike button is clicked', () => {
    cy.mount(
      <MemoryRouter>
        <CommentCard comment={comment} />
      </MemoryRouter>
    )

    cy.get('button').last().click()
    cy.get('@alertStub').should('have.been.calledWith', 'disliked comment')
  })
})
