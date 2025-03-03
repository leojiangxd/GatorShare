import React from 'react'
import Login from './Login'
import { MemoryRouter } from 'react-router-dom'

describe('Login Component Front End', () => {
  it('renders all elements correctly', () => {
    cy.mount(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    )
    
    // Check that NavBar is rendered (assuming NavBar renders an element with class "navbar")
    cy.get('.navbar').should('exist')
    
    // Verify that the username and password inputs are present.
    cy.get('input[placeholder="Username"]').should('exist')
    cy.get('input[placeholder="Password"]').should('exist')
    
    // Check that the Login button (specifically the <button> element) is present and initially disabled.
    cy.get('button').contains('Login').should('be.disabled')
    
    // Check that the register link is rendered.
    cy.contains('Need an account? Register here!').should('exist')
  })

  it('enables the Login button when username and password are provided', () => {
    cy.mount(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    )
    
    // Enter valid username and password.
    cy.get('input[placeholder="Username"]').type('TestUser')
    cy.get('input[placeholder="Password"]').type('TestPassword')
    
    // The Login button should now be enabled.
    cy.get('button').contains('Login').should('not.be.disabled')
  })

  it('toggles password visibility when the Eye icon is clicked', () => {
    cy.mount(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    )
    
    // Initially, the password input should be of type "password".
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'password')
    
    // Click the Eye icon (rendered with a "cursor-pointer" class) to toggle visibility.
    cy.get('svg.cursor-pointer').first().click()
    
    // After clicking, the password input should change to type "text".
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'text')
    
    // Click again to toggle back.
    cy.get('svg.cursor-pointer').first().click()
    cy.get('input[placeholder="Password"]').should('have.attr', 'type', 'password')
  })
})
