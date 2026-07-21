import { useState, useMemo, useEffect } from 'react'
import './App.css'

// ... (interfaces remain the same) ...
interface Education {
  institution: string
  degree: string
  period: string
  notes: string[]
}

interface Experience {
  company: string
  role: string
  period: string
  bullets: string[]
}

interface Project {
  name: string
  link: string
  bullets: string[]
}

interface Language {
  name: string
  proficiency: string
}

interface Skills {
  languages: string[]
  technologies: string[]
}

interface ResumeData {
  name: string
  title: string
  email: string
  phone: string
  linkedin: string
  github: string
  location: string
  education: Education[]
  experiences: Experience[]
  projects: Project[]
  spokenLanguages: Language[]
  certifications: string[]
  skills: Skills
}

function calculateATS(data: ResumeData): number {
  let score = 0
  if (data.name) score += 5
  if (data.email) score += 10
  if (data.experiences.length > 0) score += 15
  if (data.skills.languages.length > 0 && data.skills.technologies.length > 0) score += 10
  if (data.education.length > 0) score += 10
  if (data.linkedin) score += 10
  if (data.github) score += 5
  if (data.projects.length > 0) score += 10
  if (data.name && data.email && data.experiences.length > 0 && data.education.length > 0) score += 10
  return score
}

function App() {
  const [data, setData] = useState<ResumeData>(() => {
    const saved = localStorage.getItem('resume-draft')
    return saved ? JSON.parse(saved) : {
        name: '', title: '', email: '', phone: '', linkedin: '', github: '', location: '',
        education: [], experiences: [], projects: [], spokenLanguages: [], certifications: [],
        skills: { languages: [], technologies: [] }
    }
  })
  const [status, setStatus] = useState<'pending' | 'online' | 'offline'>('pending')

  useEffect(() => {
    localStorage.setItem('resume-draft', JSON.stringify(data))
  }, [data])

  useEffect(() => {
    fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/health`)
      .then(() => setStatus('online'))
      .catch(() => setStatus('offline'))
  }, [])

  const score = useMemo(() => calculateATS(data), [data])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setData({ ...data, [e.target.name]: e.target.value })
  }

  const exportPdf = async () => {
    const response = await fetch(`${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/generate-pdf`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    if (response.ok) {
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'curriculo.pdf'
      a.click()
    } else {
      alert('Failed to generate PDF')
    }
  }

  return (
    <div>
        <div style={{ padding: '10px', background: '#333', color: 'white' }}>
            Status: {status} | ATS Score: {score}/100
        </div>
        <div style={{ display: 'flex', gap: '20px', padding: '20px' }}>
            <div style={{ flex: 1 }}>
                <h2>Form</h2>
                {['name', 'title', 'email', 'phone'].map((key) => (
                <div key={key} style={{ marginBottom: '10px' }}>
                    <label>{key}</label>
                    <input name={key} value={(data as any)[key]} onChange={handleChange} />
                </div>
                ))}
                <button onClick={exportPdf}>Exportar PDF</button>
            </div>
            {/* ... preview ... */}
        </div>
    </div>
  )
}

export default App
