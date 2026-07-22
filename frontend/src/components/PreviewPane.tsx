import React from 'react';
import type { ResumeData } from '../types';

interface PreviewPaneProps {
  isMobileOpen?: boolean;
  onCloseMobile?: () => void;
  onExportPdf?: () => void;
  isExporting?: boolean;
  resumeData: ResumeData;
  previewRef: React.RefObject<HTMLDivElement | null>;
  overflowWarning: string | null;
}

const PreviewPane: React.FC<PreviewPaneProps> = ({
  resumeData,
  previewRef,
  overflowWarning,
  isMobileOpen,
  onCloseMobile,
  onExportPdf,
  isExporting
}) => {
  return (
    <div className={`preview-pane ${isMobileOpen ? 'mobile-open' : ''}`} ref={previewRef}>
      <div className="preview-modal-topbar">
        <button className="btn-back" onClick={onCloseMobile}>← Voltar</button>
        <button
          className="btn-export"
          onClick={onExportPdf}
          disabled={isExporting}
        >
          {isExporting ? 'Exportando...' : 'Exportar PDF ↓'}
        </button>
      </div>
      <div className="preview-scroll-area">
      <div className="paper-wrap">
        <div className="paper-toprow">
          <span className="paper-filename">curriculo.pdf</span>
          <span className="ats-badge">
            {/* ATS badge will be handled in TopBar */}
          </span>
        </div>

        <div className="paper">
          <span className="sprocket left">
            <span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span>
          </span>
          <span className="sprocket right">
            <span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span>
          </span>

          {/* Header */}
          <div className="cv-header">
            <p className="cv-name">{resumeData.name || ''}</p>
            <p className="cv-role">{resumeData.title || ''}</p>
            <p className="cv-contacts">
              {resumeData.email && <span id="pvEmail">{resumeData.email}</span>}
              {resumeData.location && (
                <span id="pvCity">
                  {resumeData.email ? ' · ' : ''}{resumeData.location}
                </span>
              )}
              {resumeData.linkedin && (
                <span id="pvLinkedin">
                  {(resumeData.email || resumeData.location) ? ' · ' : ''}{resumeData.linkedin}
                </span>
              )}
              {resumeData.github && (
                <span id="pvGithub">
                  {(resumeData.email || resumeData.location || resumeData.linkedin) ? ' · ' : ''}{resumeData.github}
                </span>
              )}
            </p>
          </div>

          {/* Education Section */}
          {resumeData.education.length > 0 && (
            <div className="cv-section">
              <h2>Educação</h2>
              {resumeData.education.map((edu, index) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{edu.institution || ''}</span>
                    <span className="cv-entry-meta">{edu.period || ''}</span>
                  </div>
                  <p className="cv-entry-sub">{edu.degree || ''}</p>
                  {edu.notes && edu.notes.length > 0 && (
                    <ul>
                      {edu.notes.map((note, noteIndex) => (
                        <li key={noteIndex}>{note}</li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Skills Section */}
          {(resumeData.skills.languages.length > 0 || resumeData.skills.technologies.length > 0) && (
            <div className="cv-section">
              <h2>Skills técnicas</h2>
              <p className="cv-skills">
                <b>Linguagens & tecnologias:</b>{' '}
                {[...resumeData.skills.languages, ...resumeData.skills.technologies].join(', ')}
              </p>
            </div>
          )}

          {/* Experience Section */}
          {resumeData.experiences.length > 0 && (
            <div className="cv-section">
              <h2>Experiência profissional</h2>
              {resumeData.experiences.map((exp, index) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{exp.company || ''}</span>
                    <span className="cv-entry-meta">{exp.period || ''}</span>
                  </div>
                  <p className="cv-entry-sub">{exp.role || ''}</p>
                  {exp.location && (
                    <p className="cv-entry-sub" style={{ marginTop: '4px' }}>
                      {exp.location}
                    </p>
                  )}
                  {exp.bullets && exp.bullets.length > 0 && (
                    <ul>
                      {exp.bullets.map((bullet, bulletIndex) => (
                        <li key={bulletIndex}>
                          {bullet}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Projects Section */}
          {resumeData.projects.length > 0 && (
            <div className="cv-section">
              <h2>Projetos</h2>
              {resumeData.projects.map((proj, index) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{proj.name || ''}</span>
                    <span className="cv-entry-meta">{proj.link_url || ''}</span>
                  </div>
                  {proj.link_url && (
                    <p style={{ marginTop: '4px' }}>
                      <a href={!/^(javascript|data):/i.test(proj.link_url) ? proj.link_url : '#'} target="_blank" rel="noopener noreferrer">
                        {proj.link_label || 'Link'}
                      </a>
                    </p>
                  )}
                  {proj.bullets && proj.bullets.length > 0 && (
                    <ul>
                      {proj.bullets.map((bullet, bulletIndex) => (
                        <li key={bulletIndex}>
                          {bullet}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* Languages Section */}
          {resumeData.spokenLanguages.length > 0 && (
            <div className="cv-section">
              <h2>Idiomas</h2>
              <p className="cv-skills">
                {resumeData.spokenLanguages.map((lang, index) => (
                  <span key={index} style={{ marginRight: index < resumeData.spokenLanguages.length - 1 ? '8px' : 0 }}>
                    <b>{lang.language || ''}:</b> {lang.level || ''}
                  </span>
                ))}
              </p>
            </div>
          )}

          {/* Certifications Section */}
          {resumeData.certifications.length > 0 && (
            <div className="cv-section">
              <h2>Certificados</h2>
              <p className="cv-skills">
                {resumeData.certifications.map((cert, index) => (
                  <span key={index} style={{ marginRight: index < resumeData.certifications.length - 1 ? '8px' : 0 }}>
                    {cert}
                  </span>
                ))}
              </p>
            </div>
          )}

          <p className="paper-footnote">
            renderizado com pdflatex · texlive-latex-extra
          </p>
        </div>
      </div>

      {/* Overflow warning */}
      {overflowWarning && (
        <div className="overflow-warning" style={{
          position: 'absolute',
          bottom: '20px',
          left: '50%',
          transform: 'translateX(-50%)',
          background: 'var(--amber)',
          color: '#04211c',
          padding: '8px 16px',
          borderRadius: '4px',
          fontFamily: 'var(--font-mono)',
          fontSize: '12px',
          whiteSpace: 'nowrap',
          zIndex: 1000
        }}>
          {overflowWarning}
        </div>
      )}
    </div>
    </div>
  );
};

export default PreviewPane;