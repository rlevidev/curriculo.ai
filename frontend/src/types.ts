export interface Education {
  institution: string;
  degree: string;
  period: string;
  notes: string[];
}

export interface Experience {
  company: string;
  location: string;
  period: string;
  role: string;
  bullets: string[];
}

export interface Project {
  name: string;
  link_url: string;
  link_label: string;
  bullets: string[];
}

export interface Language {
  language: string;
  level: string;
}

export interface Skills {
  languages: string[];
  technologies: string[];
}

export interface ResumeData {
  name: string;
  title: string;
  email: string;
  phone: string;
  linkedin: string;
  github: string;
  location: string;
  education: Education[];
  experiences: Experience[];
  projects: Project[];
  spokenLanguages: Language[];
  certifications: string[];
  skills: Skills;
}

export interface ATSCriteria {
  label: string;
  passed: boolean;
  points: number;
}

export interface ATSResult {
  score: number;
  criteria: ATSCriteria[];
}

// Calculate ATS score based on resume completeness
export function calculateATSScore(data: ResumeData): ATSResult {
  const criteria: ATSCriteria[] = [];
  let score = 0;

  // Nome preenchido (+5)
  const namePassed = !!data.name.trim();
  if (namePassed) {
    score += 5;
  }
  criteria.push({
    label: 'Nome preenchido',
    passed: namePassed,
    points: 5
  });

  // Email (+10)
  const emailPassed = !!data.email.trim();
  if (emailPassed) {
    score += 10;
  }
  criteria.push({
    label: 'Email preenchido',
    passed: emailPassed,
    points: 10
  });

  // ≥1 experiência (+15)
  const experiencePassed = data.experiences.length > 0;
  if (experiencePassed) {
    score += 15;
  }
  criteria.push({
    label: 'Pelo menos uma experiência',
    passed: experiencePassed,
    points: 15
  });

  // bullets com métricas/números (+15)
  const bulletsWithMetrics = data.experiences.some(exp =>
    exp.bullets.some(bullet => /\d/.test(bullet))
  ) || data.projects.some(proj =>
    proj.bullets.some(bullet => /\d/.test(bullet))
  );
  if (bulletsWithMetrics) {
    score += 15;
  }
  criteria.push({
    label: 'Bullets com métricas/números',
    passed: bulletsWithMetrics,
    points: 15
  });

  // ambas listas de skills preenchidas (+10)
  const skillsPassed = data.skills.languages.length > 0 && data.skills.technologies.length > 0;
  if (skillsPassed) {
    score += 10;
  }
  criteria.push({
    label: 'Ambas listas de skills preenchidas',
    passed: skillsPassed,
    points: 10
  });

  // educação preenchida (+10)
  const educationPassed = data.education.length > 0;
  if (educationPassed) {
    score += 10;
  }
  criteria.push({
    label: 'Educação preenchida',
    passed: educationPassed,
    points: 10
  });

  // LinkedIn (+10)
  const linkedinPassed = !!data.linkedin.trim();
  if (linkedinPassed) {
    score += 10;
  }
  criteria.push({
    label: 'LinkedIn preenchido',
    passed: linkedinPassed,
    points: 10
  });

  // GitHub (+5)
  const githubPassed = !!data.github.trim();
  if (githubPassed) {
    score += 5;
  }
  criteria.push({
    label: 'GitHub preenchido',
    passed: githubPassed,
    points: 5
  });

  // ≥1 projeto (+10)
  const projectsPassed = data.projects.length > 0;
  if (projectsPassed) {
    score += 10;
  }
  criteria.push({
    label: 'Pelo menos um projeto',
    passed: projectsPassed,
    points: 10
  });

  // todas as seções visíveis preenchidas (+10)
  // A section is considered visible if it has data
  const visibleSections = [
    data.name.trim() !== '',
    data.title.trim() !== '',
    data.email.trim() !== '',
    data.education.length > 0,
    data.experiences.length > 0,
    data.projects.length > 0,
    data.skills.languages.length > 0,
    data.skills.technologies.length > 0,
    data.spokenLanguages.length > 0,
    data.certifications.length > 0
  ];
  const allSectionsPassed = visibleSections.every(Boolean);
  if (allSectionsPassed) {
    score += 10;
  }
  criteria.push({
    label: 'Todas as seções visíveis preenchidas',
    passed: allSectionsPassed,
    points: 10
  });

  // Cap score at 100
  score = Math.min(score, 100);

  return {
    score,
    criteria
  };
}