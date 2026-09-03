import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from '../App';

describe('Export Error Handling', () => {
  const originalFetch = window.fetch;

  beforeEach(() => {
    // Reset fetch mock before each test
    window.fetch = vi.fn();
    // Add fake localstorage data so it passes validation (name and title required)
    localStorage.setItem('resume-draft-v2', JSON.stringify({
      name: 'Test Name',
      title: 'Test Title'
    }));
  });

  afterEach(() => {
    window.fetch = originalFetch;
    localStorage.clear();
  });

  it('displays a friendly error message when server returns HTML (502 Bad Gateway)', async () => {
    const user = userEvent.setup();
    
    // Mock health check and export response
    (window.fetch as any).mockImplementation((url: string) => {
      if (url.includes('/health')) {
        return Promise.resolve(new Response(null, { status: 200 }));
      }
      
      // Simulate NGINX 502 returning HTML instead of JSON
      if (url.includes('/generate-pdf')) {
        return Promise.resolve(new Response(
          '<html><body>502 Bad Gateway</body></html>',
          { 
            status: 502,
            headers: new Headers({ 'Content-Type': 'text/html' })
          }
        ));
      }
      return Promise.resolve(new Response(null, { status: 404 }));
    });

    render(<App />);

    // Wait for the UI to be ready
    await screen.findByText('Test Name');

    // Find and click the export button in the TopBar
    const exportButtons = screen.getAllByText(/Exportar PDF/i);
    await user.click(exportButtons[0]);

    // Should display a friendly error instead of JSON parse error (Unexpected token < in JSON at position 0)
    await waitFor(() => {
      expect(screen.getByText(/Erro no servidor/i)).toBeTruthy();
    });
  });
});