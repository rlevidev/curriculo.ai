import React from 'react';
import type { ResumeData } from '../types';

interface EditorPaneProps {
  resumeData: ResumeData;
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => void;
  onAddItem: <T>(fieldName: keyof ResumeData, newItem: T) => void;
  onRemoveItem: (fieldName: keyof ResumeData, index: number) => void;
  onUpdateItemField: <T extends object>(fieldName: keyof ResumeData, index: number, fieldKey: keyof T, value: any) => void;
  onUpdateNestedList: (outerFieldName: keyof ResumeData, outerIndex: number, innerFieldName: keyof any, innerIndex: number, value: string) => void;
  onRemoveNestedItem: (outerFieldName: keyof ResumeData, outerIndex: number, innerFieldName: keyof any, innerIndex: number) => void;
  onUpdateArrayField: (fieldName: keyof ResumeData, value: any[]) => void;
}

const EditorPane: React.FC<EditorPaneProps> = ({
  resumeData,
  onChange,
  onAddItem,
  onRemoveItem,
  onUpdateItemField,
  onUpdateNestedList,
  onRemoveNestedItem,
  onUpdateArrayField
}) => {
  return (
    <div className="editor">
      <p className="editor-hint">
        // cada seção mapeia direto pra uma parte do <span className="accent">template .tex</span> — o preview ao lado é o PDF real, atualizado a cada edição.
      </p>

      {/* Header Section */}
      <details className="section-block" open>
        <summary className="section-head">
          <span className="section-num">01</span>
          <span className="section-name">Cabeçalho</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          <div className="field">
            <label>Nome</label>
            <input
              type="text"
              name="name"
              value={resumeData.name || ''}
              onChange={onChange}
              maxLength={100}
            />
          </div>
          <div className="field">
            <label>Cargo / subtítulo</label>
            <input
              type="text"
              name="title"
              value={resumeData.title || ''}
              onChange={onChange}
              maxLength={100}
            />
          </div>
        </div>
      </details>

      {/* Contact Section */}
      <details className="section-block" open>
        <summary className="section-head">
          <span className="section-num">02</span>
          <span className="section-name">Contato</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          <div className="field-row">
            <div className="field">
              <label>Email</label>
              <input
                type="email"
                name="email"
                value={resumeData.email || ''}
                onChange={onChange}
                maxLength={100}
              />
            </div>
            <div className="field">
              <label>Telefone</label>
              <input
                type="tel"
                name="phone"
                value={resumeData.phone || ''}
                onChange={onChange}
                maxLength={50}
              />
            </div>
          </div>
          <div className="field-row">
            <div className="field">
              <label>LinkedIn</label>
              <input
                type="url"
                name="linkedin"
                value={resumeData.linkedin || ''}
                onChange={onChange}
                maxLength={100}
              />
            </div>
            <div className="field">
              <label>GitHub</label>
              <input
                type="url"
                name="github"
                value={resumeData.github || ''}
                onChange={onChange}
                maxLength={100}
              />
            </div>
          </div>
          <div className="field">
            <label>Localização</label>
            <input
              type="text"
              name="location"
              value={resumeData.location || ''}
              onChange={onChange}
              maxLength={100}
            />
          </div>
        </div>
      </details>

      {/* Education Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">03</span>
          <span className="section-name">Educação</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          {resumeData.education.map((edu: any, index: number) => (
            <div className="entry-card" key={index}>
              <div className="field">
                <label>Instituição</label>
                <input
                  type="text"
                  value={edu.institution || ''}
                  onChange={(e) => onUpdateItemField('education', index, 'institution', e.target.value)}
                  maxLength={100}
                />
              </div>
              <div className="field-row">
                <div className="field">
                  <label>Curso</label>
                  <input
                    type="text"
                    value={edu.degree || ''}
                    onChange={(e) => onUpdateItemField('education', index, 'degree', e.target.value)}
                    maxLength={100}
                  />
                </div>
                <div className="field">
                  <label>Período</label>
                  <input
                    type="text"
                    value={edu.period || ''}
                    onChange={(e) => onUpdateItemField('education', index, 'period', e.target.value)}
                    maxLength={50}
                  />
                </div>
              </div>
              <div className="field">
                <label>Observações (máx. 3)</label>
                {(edu.notes || []).map((note: string, noteIndex: number) => (
                  <div className="field" key={noteIndex}>
                    <textarea
                      value={note}
                      onChange={(e) => onUpdateNestedList('education', index, 'notes', noteIndex, e.target.value)}
                      rows={2}
                      maxLength={500}
                    />
                    <button
                      type="button"
                      onClick={() => onRemoveNestedItem('education', index, 'notes', noteIndex)}
                      className="btn-remove"
                    >
                      remover
                    </button>
                  </div>
                ))}
                <button
                  type="button"
                  onClick={() => {
                    const currentNotes = edu.notes || [];
                    if (currentNotes.length < 3) {
                      onUpdateItemField('education', index, 'notes', [...currentNotes, '']);
                    }
                  }}
                  className="btn-add-nested"
                >
                  + adicionar observação
                </button>
              </div>
              <button 
                type="button" 
                className="btn-remove-entry"
                onClick={() => onRemoveItem('education', index)}
              >
                remover formação
              </button>
            </div>
          ))}
          {resumeData.education.length < 8 && (
            <button className="add-entry" onClick={() => onAddItem('education', {
              institution: '',
              degree: '',
              period: '',
              notes: []
            })}>+ adicionar educação</button>
          )}
        </div>
      </details>

      {/* Skills Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">04</span>
          <span className="section-name">Skills</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          <div className="field">
            <label>Linguagens de programação</label>
            <textarea
              value={resumeData.skills.languages?.join(', ') || ''}
              onChange={(e) => {
                const languages = e.target.value.split(',').map(lang => lang.trim()).filter(Boolean);
                onUpdateItemField('skills', 0, 'languages', languages);
              }}
              rows={2}
              maxLength={200}
            />
          </div>
          <div className="field">
            <label>Tecnologias</label>
            <textarea
              value={resumeData.skills.technologies?.join(', ') || ''}
              onChange={(e) => {
                const technologies = e.target.value.split(',').map(tech => tech.trim()).filter(Boolean);
                onUpdateItemField('skills', 0, 'technologies', technologies);
              }}
              rows={2}
              maxLength={200}
            />
          </div>
        </div>
      </details>

      {/* Experience Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">05</span>
          <span className="section-name">Experiência profissional</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          {resumeData.experiences.map((exp: any, index: number) => (
            <div className="entry-card" key={index}>
              <div className="field-row">
                <div className="field">
                  <label>Empresa</label>
                  <input
                    type="text"
                    value={exp.company || ''}
                    onChange={(e) => onUpdateItemField('experiences', index, 'company', e.target.value)}
                    maxLength={100}
                  />
                </div>
                <div className="field">
                  <label>Local (opcional)</label>
                  <input
                    type="text"
                    value={exp.location || ''}
                    onChange={(e) => onUpdateItemField('experiences', index, 'location', e.target.value)}
                    maxLength={100}
                  />
                </div>
              </div>
              <div className="field-row">
                <div className="field">
                  <label>Período</label>
                  <input
                    type="text"
                    value={exp.period || ''}
                    onChange={(e) => onUpdateItemField('experiences', index, 'period', e.target.value)}
                    maxLength={50}
                  />
                </div>
                <div className="field">
                  <label>Cargo</label>
                  <input
                    type="text"
                    value={exp.role || ''}
                    onChange={(e) => onUpdateItemField('experiences', index, 'role', e.target.value)}
                    maxLength={100}
                  />
                </div>
              </div>
              <div className="field">
                <label>Conquistas (máx. 8)</label>
                {exp.bullets.map((bullet: string, bulletIndex: number) => (
                  <div className="field" key={bulletIndex}>
                    <textarea
                      value={bullet}
                      onChange={(e) => onUpdateNestedList('experiences', index, 'bullets', bulletIndex, e.target.value)}
                      rows={2}
                      maxLength={500}
                    />
                    <button
                      type="button"
                      onClick={() => onRemoveNestedItem('experiences', index, 'bullets', bulletIndex)}
                      className="btn-remove"
                    >
                      remover
                    </button>
                  </div>
                ))}
                <button
                  type="button"
                  onClick={() => {
                    if (exp.bullets.length < 8) {
                      const newBullets = [...exp.bullets, ''];
                      onUpdateItemField('experiences', index, 'bullets', newBullets);
                    }
                  }}
                  className="btn-add-nested"
                >
                  + adicionar conquista
                </button>
              </div>
              <button 
                type="button" 
                className="btn-remove-entry"
                onClick={() => onRemoveItem('experiences', index)}
              >
                remover experiência
              </button>
            </div>
          ))}
          {resumeData.experiences.length < 8 && (
            <button className="add-entry" onClick={() => onAddItem('experiences', {
              company: '',
              location: '',
              period: '',
              role: '',
              bullets: []
            })}>+ adicionar experiência</button>
          )}
        </div>
      </details>

      {/* Projects Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">06</span>
          <span className="section-name">Projetos</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          {resumeData.projects.map((proj: any, index: number) => (
            <div className="entry-card" key={index}>
              <div className="field">
                <label>Nome</label>
                <input
                  type="text"
                  value={proj.name || ''}
                  onChange={(e) => onUpdateItemField('projects', index, 'name', e.target.value)}
                  maxLength={100}
                />
              </div>
              <div className="field-row">
                <div className="field">
                  <label>URL (opcional)</label>
                  <input
                    type="text"
                    value={proj.link_url || ''}
                    onChange={(e) => onUpdateItemField('projects', index, 'link_url', e.target.value)}
                    maxLength={200}
                  />
                </div>
                <div className="field">
                  <label>Label do link (padrão: "Link")</label>
                  <input
                    type="text"
                    value={proj.link_label || 'Link'}
                    onChange={(e) => onUpdateItemField('projects', index, 'link_label', e.target.value)}
                    maxLength={50}
                  />
                </div>
              </div>
              <div className="field">
                <label>Conquistas (máx. 8)</label>
                {proj.bullets.map((bullet: string, bulletIndex: number) => (
                  <div className="field" key={bulletIndex}>
                    <textarea
                      value={bullet}
                      onChange={(e) => onUpdateNestedList('projects', index, 'bullets', bulletIndex, e.target.value)}
                      rows={2}
                      maxLength={500}
                    />
                    <button
                      type="button"
                      onClick={() => onRemoveNestedItem('projects', index, 'bullets', bulletIndex)}
                      className="btn-remove"
                    >
                      remover
                    </button>
                  </div>
                ))}
                <button
                  type="button"
                  onClick={() => {
                    if (proj.bullets.length < 8) {
                      const newBullets = [...proj.bullets, ''];
                      onUpdateItemField('projects', index, 'bullets', newBullets);
                    }
                  }}
                  className="btn-add-nested"
                >
                  + adicionar conquista
                </button>
              </div>
              <button 
                type="button" 
                className="btn-remove-entry"
                onClick={() => onRemoveItem('projects', index)}
              >
                remover projeto
              </button>
            </div>
          ))}
          {resumeData.projects.length < 8 && (
            <button className="add-entry" onClick={() => onAddItem('projects', {
              name: '',
              link_url: '',
              link_label: 'Link',
              bullets: []
            })}>+ adicionar projeto</button>
          )}
        </div>
      </details>

      {/* Languages Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">07</span>
          <span className="section-name">Idiomas</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          {resumeData.spokenLanguages.map((lang: any, index: number) => (
            <div className="entry-card" key={index}>
              <div className="field-row">
                <div className="field">
                  <label>Idioma</label>
                  <input
                    type="text"
                    value={lang.language || ''}
                    onChange={(e) => onUpdateItemField('spokenLanguages', index, 'language', e.target.value)}
                    maxLength={50}
                  />
                  <span className="field-hint">
                    (ex: Inglês, Espanhol)
                  </span>
                </div>
                <div className="field">
                  <label>Nível</label>
                  <input
                    type="text"
                    value={lang.level || ''}
                    onChange={(e) => onUpdateItemField('spokenLanguages', index, 'level', e.target.value)}
                    maxLength={50}
                  />
                  <span className="field-hint">
                    (ex: Básico, Intermediário, Avançado, Fluente)
                  </span>
                </div>
              </div>
              <button 
                type="button" 
                className="btn-remove-entry"
                onClick={() => onRemoveItem('spokenLanguages', index)}
              >
                remover idioma
              </button>
            </div>
          ))}
          {resumeData.spokenLanguages.length < 8 && (
            <button className="add-entry" onClick={() => onAddItem('spokenLanguages', {
              language: '',
              level: ''
            })}>+ adicionar idioma</button>
          )}
        </div>
      </details>

      {/* Certifications Section */}
      <details className="section-block">
        <summary className="section-head">
          <span className="section-num">08</span>
          <span className="section-name">Certificados</span>
          <span className="section-chevron">▸</span>
        </summary>
        <div className="section-body">
          {resumeData.certifications.map((cert: string, index: number) => (
            <div className="field" key={index}>
              <input
                type="text"
                value={cert}
                onChange={(e) => {
                  const newCerts = [...resumeData.certifications];
                  newCerts[index] = e.target.value;
                  onUpdateArrayField('certifications', newCerts);
                }}
                maxLength={200}
              />
              <button
                type="button"
                onClick={() => {
                  const newCerts = [...resumeData.certifications];
                  newCerts.splice(index, 1);
                  onUpdateArrayField('certifications', newCerts);
                }}
                className="btn-remove"
              >
                remover
              </button>
            </div>
          ))}
          {resumeData.certifications.length < 8 && (
            <button className="add-entry" onClick={() => onAddItem('certifications', '')}>
              + adicionar certificado
            </button>
          )}
        </div>
      </details>
    </div>
  );
};

export default EditorPane;