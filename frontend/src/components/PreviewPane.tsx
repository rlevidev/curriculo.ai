import { useEffect } from 'react';

interface PreviewPaneProps {
  isMobileOpen?: boolean;
  onCloseMobile?: () => void;
  onExportPdf?: () => void;
  isExporting?: boolean;
  resumeData: any;
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
  // Function to escape HTML entities (for security, though we're not rendering HTML in preview)
  const escapeHtml = (text: string): string => {
    return text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  };

  // Function to check for overflow
  useEffect(() => {
    if (!previewRef.current) return;

    const checkOverflow = () => {
      const previewContainer = previewRef.current;
      if (!previewContainer) return;

      // We'll let the parent handle the warning display
    };

    const observer = new ResizeObserver(checkOverflow);
    observer.observe(previewRef.current);

    return () => observer.disconnect();
  }, [resumeData]); // Re-run when resume data changes

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
            <p className="cv-name">{escapeHtml(resumeData.name || '')}</p>
            <p className="cv-role">{escapeHtml(resumeData.title || '')}</p>
            <p className="cv-contacts">
              {resumeData.email && <span id="pvEmail">{escapeHtml(resumeData.email)}</span>}
              {resumeData.location && (
                <span id="pvCity">
                  {resumeData.email ? ' · ' : ''}{escapeHtml(resumeData.location)}
                </span>
              )}
              {resumeData.linkedin && (
                <span id="pvLinkedin">
                  {(resumeData.email || resumeData.location) ? ' · ' : ''}{escapeHtml(resumeData.linkedin)}
                </span>
              )}
              {resumeData.github && (
                <span id="pvGithub">
                  {(resumeData.email || resumeData.location || resumeData.linkedin) ? ' · ' : ''}{escapeHtml(resumeData.github)}
                </span>
              )}
            </p>
          </div>

          {/* Education Section */}
          {resumeData.education.length > 0 && (
            <div className="cv-section">
              <h2>Educação</h2>
              {resumeData.education.map((edu: any, index: number) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{escapeHtml(edu.institution || '')}</span>
                    <span className="cv-entry-meta">{escapeHtml(edu.period || '')}</span>
                  </div>
                  <p className="cv-entry-sub">{escapeHtml(edu.degree || '')}</p>
                  {edu.notes && edu.notes.length > 0 && (
                    <ul>
                      {edu.notes.map((note: string, noteIndex: number) => (
                        <li key={noteIndex}>{escapeHtml(note)}</li>
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
                {[...resumeData.skills.languages, ...resumeData.skills.technologies]
                  .map((skill: string) => escapeHtml(skill))
                  .join(', ')}
              </p>
            </div>
          )}

          {/* Experience Section */}
          {resumeData.experiences.length > 0 && (
            <div className="cv-section">
              <h2>Experiência profissional</h2>
              {resumeData.experiences.map((exp: any, index: number) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{escapeHtml(exp.company || '')}</span>
                    <span className="cv-entry-meta">{escapeHtml(exp.period || '')}</span>
                  </div>
                  <p className="cv-entry-sub">{escapeHtml(exp.role || '')}</p>
                  {exp.location && (
                    <p className="cv-entry-sub" style={{ marginTop: '4px' }}>
                      {escapeHtml(exp.location)}
                    </p>
                  )}
                  {exp.bullets && exp.bullets.length > 0 && (
                    <ul>
                      {exp.bullets.map((bullet: string, bulletIndex: number) => (
                        <li key={bulletIndex}>
                          {escapeHtml(bullet)}
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
              {resumeData.projects.map((proj: any, index: number) => (
                <div className="cv-entry" key={index}>
                  <div className="cv-entry-top">
                    <span className="cv-entry-title">{escapeHtml(proj.name || '')}</span>
                    <span className="cv-entry-meta">{escapeHtml(proj.link_url || '')}</span>
                  </div>
                  {proj.link_url && (
                    <p style={{ marginTop: '4px' }}>
                      <a href={proj.link_url} target="_blank" rel="noopener noreferrer">
                        {escapeHtml(proj.link_label || 'Link')}
                      </a>
                    </p>
                  )}
                  {proj.bullets && proj.bullets.length > 0 && (
                    <ul>
                      {proj.bullets.map((bullet: string, bulletIndex: number) => (
                        <li key={bulletIndex}>
                          {escapeHtml(bullet)}
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
                {resumeData.spokenLanguages.map((lang: any, index: number) => (
                  <span key={index} style={{ marginRight: index < resumeData.spokenLanguages.length - 1 ? '8px' : 0 }}>
                    <b>{escapeHtml(lang.language || '')}:</b> {escapeHtml(lang.level || '')}
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
                {resumeData.certifications.map((cert: string, index: number) => (
                  <span key={index} style={{ marginRight: index < resumeData.certifications.length - 1 ? '8px' : 0 }}>
                    {escapeHtml(cert)}
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