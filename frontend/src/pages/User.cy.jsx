import React from 'react'
import User from './User'
import { MemoryRouter, Routes, Route } from 'react-router-dom'

describe('User Page Front End', () => {
  context('When posts are returned', () => {
    // Updated dummy posts include "content" and "CreatedAt"
    const dummyPosts = [
      {
        post_id: "1",
        title: "First Post",
        likes: 10,
        dislikes: 2,
        content: "Content of the first post.",
        CreatedAt: "2023-03-02T10:00:00Z",
      },
      {
        post_id: "2",
        title: "Second Post",
        likes: 20,
        dislikes: 1,
        content: "Content of the second post.",
        CreatedAt: "2023-03-02T11:00:00Z",
      },
    ];

    beforeEach(() => {
      // Intercept the GET request for posts and return dummy data.
      cy.intercept('GET', '**/api/v1/member/*/posts', {
        statusCode: 200,
        body: { data: dummyPosts },
      }).as('getPosts');
    });

    it('renders posts when API returns data (no search term)', () => {
      // Use an initial entry with a valid id so that safeId is not empty.
      cy.mount(
        <MemoryRouter initialEntries={['/user/Alice']}>
          <Routes>
            <Route path="/user/:id" element={<User />} />
          </Routes>
        </MemoryRouter>
      );
      // Check that the profile card shows the correct id (e.g., "Alice")
      cy.contains('Alice').should('exist');
      // Verify that both post titles are rendered.
      cy.contains('First Post').should('exist');
      cy.contains('Second Post').should('exist');
    });
  });

  context('When no posts are returned', () => {
    beforeEach(() => {
      // Intercept the GET request for posts and return an empty array.
      cy.intercept('GET', '**/api/v1/member/*/posts', {
        statusCode: 200,
        body: { data: [] },
      }).as('getPosts');
    });

    it('renders "No posts found" when there are no posts', () => {
      cy.mount(
        <MemoryRouter initialEntries={['/user/Alice']}>
          <Routes>
            <Route path="/user/:id" element={<User />} />
          </Routes>
        </MemoryRouter>
      );
      // Verify that the "No posts found" message is displayed.
      cy.contains('No posts found').should('exist');
    });
  });

  context('Basic rendering of the posts list', () => {
    const dummyPosts = [
      {
        post_id: "1",
        title: "First Post",
        likes: 10,
        dislikes: 2,
        content: "Content of the first post.",
        CreatedAt: "2023-03-02T10:00:00Z",
      },
      {
        post_id: "2",
        title: "Second Post",
        likes: 20,
        dislikes: 1,
        content: "Content of the second post.",
        CreatedAt: "2023-03-02T11:00:00Z",
      },
    ];

    beforeEach(() => {
      cy.intercept('GET', '**/api/v1/member/*/posts', {
        statusCode: 200,
        body: { data: dummyPosts },
      }).as('getPosts');
    });

    it('renders the posts list correctly', () => {
      cy.mount(
        <MemoryRouter initialEntries={['/user/Alice']}>
          <Routes>
            <Route path="/user/:id" element={<User />} />
          </Routes>
        </MemoryRouter>
      );
      // Ensure that both posts are rendered.
      cy.contains('First Post').should('exist');
      cy.contains('Second Post').should('exist');
    });
  });
});
