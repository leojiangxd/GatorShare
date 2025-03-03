import React from 'react'
import User from './Home'
import { MemoryRouter } from 'react-router-dom'

describe('<User /> Front End', () => {
  it('renders without crashing', () => {
    cy.mount(
      <MemoryRouter>
        <User />
      </MemoryRouter>
    )
    // Optionally, add frontend-specific assertions here.
    // For example, check if a NavBar or welcome message exists:
    cy.get('.navbar').should('exist')
  })
})
