import { describe, it, expect } from 'vitest';
import { normalizeResumeData } from './normalize';

describe('normalizeResumeData', () => {
  it('should fill missing arrays with defaults for an old draft', () => {
    // Simulate an old draft format where some arrays/objects are missing
    const oldDraft: any = {
      name: "John Doe",
      title: "Developer",
      skills: {
        technologies: ["React", "TypeScript"]
        // languages is missing
      },
      experiences: [
        {
          company: "Acme Corp",
          location: "Remote",
          period: "2020 - 2022",
          role: "Frontend Developer"
          // bullets missing
        }
      ],
      education: [
        {
          institution: "University",
          degree: "BSc",
          period: "2015 - 2019"
          // notes missing
        }
      ],
      projects: [
        {
          name: "Project A",
          link_url: "https://example.com",
          link_label: "Live"
          // bullets missing
        }
      ]
      // spokenLanguages, certifications, github, etc. might be missing
    };

    const normalized = normalizeResumeData(oldDraft);

    // Should preserve existing valid data
    expect(normalized.name).toBe("John Doe");
    expect(normalized.skills.technologies).toEqual(["React", "TypeScript"]);
    
    // Should fill missing skills.languages
    expect(normalized.skills.languages).toEqual([]);

    // Should fill missing arrays in lists
    expect(normalized.experiences[0].bullets).toEqual([]);
    expect(normalized.education[0].notes).toEqual([]);
    expect(normalized.projects[0].bullets).toEqual([]);

    // Should fill missing root arrays/strings
    expect(normalized.spokenLanguages).toEqual([]);
    expect(normalized.certifications).toEqual([]);
    expect(normalized.github).toBe("");
    expect(normalized.linkedin).toBe("");
    expect(normalized.email).toBe("");
    expect(normalized.phone).toBe("");
    expect(normalized.location).toBe("");
  });

  it('should handle completely missing skills object', () => {
    const oldDraft: any = {
      name: "John Doe"
    };
    const normalized = normalizeResumeData(oldDraft);
    expect(normalized.skills).toEqual({ languages: [], technologies: [] });
  });

  it('should return a full default object if draft is null or empty', () => {
    const normalized = normalizeResumeData(null);
    expect(normalized.name).toBe("");
    expect(normalized.skills).toEqual({ languages: [], technologies: [] });
    expect(normalized.experiences).toEqual([]);
    expect(normalized.education).toEqual([]);
    expect(normalized.projects).toEqual([]);
    expect(normalized.spokenLanguages).toEqual([]);
    expect(normalized.certifications).toEqual([]);
  });
});