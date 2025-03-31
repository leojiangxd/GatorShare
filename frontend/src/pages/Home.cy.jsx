import React from 'react'
import Home from './Home'
import { MemoryRouter } from 'react-router-dom'

describe('<Home /> Front End', () => {
  it('renders without crashing', () => {
    cy.mount(
      <MemoryRouter>
        <Home />
      </MemoryRouter>
    )
    // Optionally, add frontend-specific assertions here.
    // For example, check if a NavBar or welcome message exists:
    cy.get('.navbar').should('exist')
  })
})
