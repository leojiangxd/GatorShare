import React from 'react'
import Create from './Create'
import { MemoryRouter } from 'react-router-dom'

describe('<Create />', () => {
  it('renders the Create component correctly', () => {
    cy.mount(
      <MemoryRouter>
        <Create />
      </MemoryRouter>
    )

    // Verify that the NavBar is rendered (assuming it renders an element with a "navbar" class)
    cy.get('.navbar').should('exist')

    // Verify that the title input, textarea, and buttons are present.
    cy.get('input[placeholder="Title"]').should('exist')
    cy.get('textarea').should('exist')
    cy.contains('Create Post').should('exist')
  })
})
