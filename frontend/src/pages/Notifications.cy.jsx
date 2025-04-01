import React from 'react'
import Notifications from './Notifications'
import { MemoryRouter } from 'react-router-dom'

describe('<Notifications />', () => {
  it('renders page title and "Mark All as Read" button', () => {
    cy.mount(
      <MemoryRouter>
        <Notifications />
      </MemoryRouter>
    )
    cy.contains("Notifications").should("exist")
    cy.contains("Mark All as Read").should("exist")
  })

  it('renders notifications grouped by date', () => {
    cy.mount(
      <MemoryRouter>
        <Notifications />
      </MemoryRouter>
    )
    // Based on your sample data
    cy.contains("March 31, 2025").should("exist")
    cy.contains("March 30, 2025").should("exist")
    cy.contains("March 29, 2025").should("exist")
  })

  it('toggles read status of a notification', () => {
    cy.mount(
      <MemoryRouter>
        <Notifications />
      </MemoryRouter>
    )
    // Find an unread notification's toggle button and click it.
    // Initially, it should show "Mark as Read"
    cy.get("button.btn.btn-outline.btn-xs")
      .contains("Mark as Read")
      .first()
      .click()
    // After clicking, it should now display "Mark as Unread"
    cy.get("button.btn.btn-outline.btn-xs")
      .contains("Mark as Unread")
      .should("exist")
  })

  it('marks all notifications as read', () => {
    cy.mount(
      <MemoryRouter>
        <Notifications />
      </MemoryRouter>
    )
    // Ensure there is at least one "Mark as Read" button present.
    cy.get("button.btn.btn-outline.btn-xs")
      .contains("Mark as Read")
      .should("exist")
    // Click the "Mark All as Read" button.
    cy.contains("Mark All as Read").click()
    // Verify that no "Mark as Read" buttons remain and that each shows "Mark as Unread".
    cy.get("button.btn.btn-outline.btn-xs")
      .contains("Mark as Read")
      .should("not.exist")
    cy.get("button.btn.btn-outline.btn-xs").each(($btn) => {
      cy.wrap($btn).should("contain", "Mark as Unread")
    })
  })
})
