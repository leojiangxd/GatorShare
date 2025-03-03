import React from 'react'
import NavBar from './NavBar'
import { MemoryRouter } from 'react-router-dom'

describe('<NavBar />', () => {
  it('renders', () => {
    cy.mount(
      <MemoryRouter>
        <NavBar />
      </MemoryRouter>
    )
  })
})
