import type { ResumeData } from '../types';

export function normalizeResumeData(data: any): ResumeData {
  if (!data || typeof data !== 'object') {
    data = {};
  }

  return {
    name: data.name || "",
    title: data.title || "",
    email: data.email || "",
    phone: data.phone || "",
    linkedin: data.linkedin || "",
    github: data.github || "",
    location: data.location || "",
    education: (data.education || []).map((edu: any) => ({
      institution: edu.institution || "",
      degree: edu.degree || "",
      period: edu.period || "",
      notes: Array.isArray(edu.notes) ? edu.notes : [],
    })),
    experiences: (data.experiences || []).map((exp: any) => ({
      company: exp.company || "",
      location: exp.location || "",
      period: exp.period || "",
      role: exp.role || "",
      bullets: Array.isArray(exp.bullets) ? exp.bullets : [],
    })),
    projects: (data.projects || []).map((proj: any) => ({
      name: proj.name || "",
      link_url: proj.link_url || "",
      link_label: proj.link_label || "",
      bullets: Array.isArray(proj.bullets) ? proj.bullets : [],
    })),
    spokenLanguages: (data.spokenLanguages || []).map((lang: any) => ({
      language: lang.language || "",
      level: lang.level || "",
    })),
    certifications: Array.isArray(data.certifications) ? data.certifications : [],
    skills: {
      languages: Array.isArray(data.skills?.languages) ? data.skills.languages : [],
      technologies: Array.isArray(data.skills?.technologies) ? data.skills.technologies : [],
    },
  };
}