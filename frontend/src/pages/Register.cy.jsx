import React from 'react'
import Register from './Register'
import { MemoryRouter } from 'react-router-dom'

describe('Register Component Frontend Tests', () => {
  it('renders all elements correctly', () => {
    cy.mount(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    )

    // Check for the NavBar (assuming NavBar renders an element with a class "navbar")
    cy.get('.navbar').should('exist')

    // Check that the username, email, and password inputs are present.
    cy.get('input[placeholder="Username"]').should('exist')
    cy.get('input[placeholder="example@ufl.edu"]').should('exist')
    cy.get('input[placeholder="Password"]').should('exist')

    // Check that the Register button is present and initially disabled.
    cy.contains('Register').should('be.disabled')

    // Check for the login link.
    cy.contains('Have an account? Login here!').should('exist')
  })

  it('enables the Register button with valid inputs', () => {
    cy.mount(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    )

    // Enter a valid username.
    cy.get('input[placeholder="Username"]').type('TestUser')
    // Enter a valid UF email.
    cy.get('input[placeholder="example@ufl.edu"]').type('test@ufl.edu')
    // Enter a valid password: must be at least 8 characters, with one digit, one lowercase, one uppercase.
    cy.get('input[placeholder="Password"]').type('Test1234')

    // After valid input, the Register button should be enabled.
    cy.contains('Register').should('not.be.disabled')
  })

  it('toggles password visibility when the Eye icon is clicked', () => {
    cy.mount(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    )

    // Initially, the password input type should be "password".
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'password')

    // Click the Eye icon (rendered with a "cursor-pointer" class).
    cy.get('svg.cursor-pointer').first().click()

    // After clicking, the password input type should change to "text".
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'text')

    // Click again to toggle back.
    cy.get('svg.cursor-pointer').first().click()
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'password')
  })
})
