import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import EditorPane from './EditorPane';

const defaultResumeData = {
  name: '',
  title: '',
  email: '',
  phone: '',
  linkedin: '',
  github: '',
  location: '',
  education: [],
  skills: { languages: ['Go', 'TypeScript'], technologies: [] },
  experiences: [],
  projects: [],
  spokenLanguages: [],
  certifications: []
};

describe('EditorPane - Skills', () => {
  it('allows typing free text with commas and commits on blur', async () => {
    const onUpdateItemField = vi.fn();
    const user = userEvent.setup();
    
    render(
      <EditorPane
        resumeData={defaultResumeData}
        onChange={vi.fn()}
        onAddItem={vi.fn()}
        onRemoveItem={vi.fn()}
        onUpdateItemField={onUpdateItemField}
        onUpdateNestedList={vi.fn()}
        onRemoveNestedItem={vi.fn()}
        onUpdateArrayField={vi.fn()}
      />
    );

    const textarea = screen.getByLabelText(/Linguagens de programação/i);
    
    // Type a comma and another word
    await user.type(textarea, ', Rust');
    
    // It should not trigger onUpdateItemField immediately (we want it on blur)
    expect(onUpdateItemField).not.toHaveBeenCalled();
    
    // It should have the text we typed
    expect(textarea).toHaveValue('Go, TypeScript, Rust');
    
    // Now trigger blur
    await user.tab();
    
    // NOW it should trigger update with the array
    expect(onUpdateItemField).toHaveBeenCalledWith(
      'skills', -1, 'languages', ['Go', 'TypeScript', 'Rust']
    );
  });
});