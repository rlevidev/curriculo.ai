# Exemplo de front-end

```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>curriculo.ai — mockup de interface</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=IBM+Plex+Sans:wght@400;500;600&family=Source+Serif+4:opsz,wght@8..60,400;8..60,600;8..60,700&display=swap" rel="stylesheet">
<style>
  :root{
    --ink:#0a0d12;
    --ink-surface:#12151c;
    --ink-surface-2:#171b24;
    --ink-border:#20242e;
    --ink-border-strong:#2c313d;
    --text:#e8eaf0;
    --text-muted:#838b9c;
    --text-faint:#545b6b;
    --accent:#5eead4;
    --accent-dim:rgba(94,234,212,0.12);
    --amber:#fbbf24;

    --paper:#fbf9f4;
    --paper-edge:#efece3;
    --paper-line:#e1dccd;
    --paper-ink:#211f1a;
    --paper-ink-muted:#7a715f;

    --font-mono:'IBM Plex Mono', ui-monospace, Menlo, monospace;
    --font-sans:'IBM Plex Sans', -apple-system, 'Segoe UI', sans-serif;
    --font-serif:'Source Serif 4', Georgia, serif;
  }

  *{box-sizing:border-box;}
  html,body{margin:0;height:100%;}
  body{
    background:var(--ink);
    color:var(--text);
    font-family:var(--font-sans);
    overflow:hidden;
  }

  a{color:inherit;}
  button{font-family:inherit;}
  button:focus-visible, input:focus-visible, textarea:focus-visible, summary:focus-visible{
    outline:2px solid var(--accent);
    outline-offset:2px;
  }

  /* ---------- top bar ---------- */
  .topbar{
    height:52px;
    display:flex;
    align-items:center;
    justify-content:space-between;
    padding:0 18px;
    border-bottom:1px solid var(--ink-border);
    background:var(--ink-surface);
    gap:16px;
  }
  .brand{
    font-family:var(--font-mono);
    font-weight:700;
    font-size:15px;
    letter-spacing:-0.01em;
    display:flex;
    align-items:center;
    gap:8px;
    white-space:nowrap;
  }
  .brand .dot{color:var(--accent);}

  .tabs{
    display:flex;
    gap:2px;
    font-family:var(--font-mono);
    font-size:12.5px;
  }
  .tabs button{
    background:none;
    border:1px solid transparent;
    color:var(--text-muted);
    padding:7px 12px;
    border-radius:6px;
    cursor:pointer;
  }
  .tabs button.active{
    color:var(--text);
    background:var(--ink-surface-2);
    border-color:var(--ink-border-strong);
  }

  .topbar-right{
    display:flex;
    align-items:center;
    gap:12px;
  }

  .compile-status{
    font-family:var(--font-mono);
    font-size:12px;
    color:var(--text-faint);
    display:flex;
    align-items:center;
    gap:7px;
    white-space:nowrap;
  }
  .compile-status .dot{
    width:6px;height:6px;border-radius:50%;
    background:var(--text-faint);
    flex-shrink:0;
  }
  .compile-status.ok{color:#34d399;}
  .compile-status.ok .dot{background:#34d399;}
  .compile-status.busy{color:var(--amber);}
  .compile-status.busy .dot{background:var(--amber); animation:blink 0.9s ease-in-out infinite;}
  @keyframes blink{50%{opacity:0.25;}}

  .btn-export{
    font-family:var(--font-mono);
    font-size:12.5px;
    font-weight:600;
    background:var(--accent);
    color:#04211c;
    border:1px solid var(--accent);
    padding:8px 14px;
    border-radius:6px;
    cursor:pointer;
    white-space:nowrap;
  }
  .btn-export:hover{background:transparent; color:var(--accent);}

  /* ---------- layout ---------- */
  .layout{
    display:grid;
    grid-template-columns: minmax(340px, 460px) 1fr;
    height:calc(100% - 52px);
  }

  /* ---------- editor pane ---------- */
  .editor{
    background:var(--ink);
    border-right:1px solid var(--ink-border);
    overflow-y:auto;
    padding:18px 0 60px;
  }
  .editor::-webkit-scrollbar{width:9px;}
  .editor::-webkit-scrollbar-thumb{background:var(--ink-border-strong); border-radius:5px;}

  .editor-hint{
    font-family:var(--font-mono);
    font-size:11.5px;
    color:var(--text-faint);
    padding:0 20px 16px;
    line-height:1.6;
  }
  .editor-hint .accent{color:var(--accent);}

  .section-block{
    border-bottom:1px solid var(--ink-border);
  }

  .section-head{
    display:flex;
    align-items:center;
    gap:12px;
    padding:13px 20px;
    cursor:pointer;
    user-select:none;
    list-style:none;
  }
  .section-head::-webkit-details-marker{display:none;}

  .section-num{
    font-family:var(--font-mono);
    font-size:12px;
    color:var(--text-faint);
    width:20px;
    flex-shrink:0;
  }
  details[open] .section-num{color:var(--accent);}

  .section-name{
    font-family:var(--font-mono);
    font-size:13.5px;
    font-weight:600;
    color:var(--text);
    flex-grow:1;
  }

  .section-chevron{
    font-size:11px;
    color:var(--text-faint);
    transition:transform 0.15s ease;
  }
  details[open] .section-chevron{transform:rotate(90deg); color:var(--accent);}

  .section-body{
    padding:2px 20px 20px 52px;
    display:flex;
    flex-direction:column;
    gap:12px;
  }

  .field{
    display:flex;
    flex-direction:column;
    gap:5px;
  }
  .field label{
    font-family:var(--font-mono);
    font-size:10.5px;
    text-transform:uppercase;
    letter-spacing:0.06em;
    color:var(--text-faint);
  }
  .field input, .field textarea{
    background:var(--ink-surface);
    border:1px solid var(--ink-border-strong);
    border-radius:6px;
    color:var(--text);
    font-family:var(--font-sans);
    font-size:13.5px;
    padding:9px 10px;
    resize:vertical;
  }
  .field input::placeholder, .field textarea::placeholder{color:var(--text-faint);}
  .field input:focus, .field textarea:focus{border-color:var(--accent); outline:none;}

  .field-row{display:flex; gap:10px;}
  .field-row .field{flex:1;}

  .entry-card{
    border:1px solid var(--ink-border);
    border-radius:8px;
    padding:12px;
    display:flex;
    flex-direction:column;
    gap:10px;
    background:var(--ink-surface);
  }

  .add-entry{
    font-family:var(--font-mono);
    font-size:12px;
    color:var(--accent);
    background:none;
    border:1px dashed var(--ink-border-strong);
    border-radius:6px;
    padding:9px;
    cursor:pointer;
    text-align:left;
  }
  .add-entry:hover{border-color:var(--accent); background:var(--accent-dim);}

  /* ---------- preview pane ---------- */
  .preview-pane{
    background:
      radial-gradient(ellipse at top, #14171f 0%, var(--ink) 60%);
    overflow-y:auto;
    display:flex;
    justify-content:center;
    padding:36px 24px 80px;
  }
  .preview-pane::-webkit-scrollbar{width:9px;}
  .preview-pane::-webkit-scrollbar-thumb{background:var(--ink-border-strong); border-radius:5px;}

  .paper-wrap{
    width:100%;
    max-width:620px;
  }

  .paper-toprow{
    display:flex;
    justify-content:space-between;
    align-items:baseline;
    margin-bottom:10px;
    padding:0 6px;
  }
  .paper-filename{
    font-family:var(--font-mono);
    font-size:11.5px;
    color:var(--text-faint);
  }
  .ats-badge{
    font-family:var(--font-mono);
    font-size:11px;
    color:var(--text-muted);
    display:flex;
    align-items:center;
    gap:6px;
  }

  /* the "paper": torn perforated edge + sprocket holes, like tractor-feed printer paper */
  .paper{
    position:relative;
    background:var(--paper);
    color:var(--paper-ink);
    box-shadow: 0 20px 50px -12px rgba(0,0,0,0.55), 0 0 0 1px var(--paper-edge);
    clip-path: polygon(
      0% 14px, 2% 8px, 4% 14px, 6% 8px, 8% 14px, 10% 8px, 12% 14px, 14% 8px,
      16% 14px, 18% 8px, 20% 14px, 22% 8px, 24% 14px, 26% 8px, 28% 14px, 30% 8px,
      32% 14px, 34% 8px, 36% 14px, 38% 8px, 40% 14px, 42% 8px, 44% 14px, 46% 8px,
      48% 14px, 50% 8px, 52% 14px, 54% 8px, 56% 14px, 58% 8px, 60% 14px, 62% 8px,
      64% 14px, 66% 8px, 68% 14px, 70% 8px, 72% 14px, 74% 8px, 76% 14px, 78% 8px,
      80% 14px, 82% 8px, 84% 14px, 86% 8px, 88% 14px, 90% 8px, 92% 14px, 94% 8px,
      96% 14px, 98% 8px, 100% 14px,
      100% 100%, 0% 100%
    );
    padding:38px 44px 46px;
    font-family:var(--font-serif);
  }
  .sprocket{
    position:absolute;
    top:18px; bottom:18px;
    width:20px;
    display:flex;
    flex-direction:column;
    justify-content:space-between;
    pointer-events:none;
  }
  .sprocket.left{left:-10px;}
  .sprocket.right{right:-10px;}
  .sprocket span{
    width:8px; height:8px;
    border-radius:50%;
    background:var(--ink);
    box-shadow: inset 0 0 0 1px rgba(0,0,0,0.15);
  }

  .cv-header{
    text-align:center;
    border-bottom:2px solid var(--paper-ink);
    padding-bottom:14px;
    margin-bottom:18px;
  }
  .cv-name{
    font-size:28px;
    font-weight:700;
    margin:0 0 4px;
    letter-spacing:-0.01em;
  }
  .cv-role{
    font-family:var(--font-sans);
    font-size:13px;
    color:var(--paper-ink-muted);
    margin:0 0 8px;
  }
  .cv-contacts{
    font-family:var(--font-sans);
    font-size:11.5px;
    color:var(--paper-ink-muted);
  }
  .cv-contacts span:not(:last-child)::after{content:" · ";}

  .cv-section{margin-bottom:16px;}
  .cv-section:last-child{margin-bottom:0;}
  .cv-section h2{
    font-family:var(--font-sans);
    font-size:11px;
    letter-spacing:0.1em;
    text-transform:uppercase;
    font-weight:600;
    color:var(--paper-ink);
    border-bottom:1px solid var(--paper-line);
    padding-bottom:4px;
    margin:0 0 9px;
  }

  .cv-entry{margin-bottom:10px;}
  .cv-entry:last-child{margin-bottom:0;}
  .cv-entry-top{
    display:flex;
    justify-content:space-between;
    gap:12px;
    font-size:14px;
  }
  .cv-entry-title{font-weight:700;}
  .cv-entry-meta{
    font-family:var(--font-sans);
    font-size:11.5px;
    color:var(--paper-ink-muted);
    white-space:nowrap;
  }
  .cv-entry-sub{
    font-style:italic;
    font-size:12.5px;
    color:var(--paper-ink-muted);
    margin:1px 0 4px;
  }
  .cv-entry ul{
    margin:0;
    padding-left:16px;
    font-family:var(--font-sans);
    font-size:12px;
    line-height:1.55;
    color:#3a362e;
  }
  .cv-entry ul li{margin-bottom:2px;}

  .cv-skills{
    font-family:var(--font-sans);
    font-size:12px;
    line-height:1.7;
  }
  .cv-skills b{font-weight:600;}

  .caret{
    display:inline-block;
    width:7px; height:14px;
    background:var(--accent);
    vertical-align:-2px;
    margin-left:2px;
    animation:blink 1s step-end infinite;
  }

  /* ---------- footer note under paper ---------- */
  .paper-footnote{
    font-family:var(--font-mono);
    font-size:11px;
    color:var(--text-faint);
    text-align:center;
    margin-top:16px;
  }

  /* ---------- responsive ---------- */
  @media (max-width: 880px){
    .layout{grid-template-columns:1fr; height:auto; overflow:visible;}
    body{overflow:auto;}
    .editor, .preview-pane{overflow:visible; height:auto;}
    .editor{border-right:none; border-bottom:1px solid var(--ink-border);}
    .tabs{display:none;}
  }

  @media (prefers-reduced-motion: reduce){
    .compile-status.busy .dot{animation:none;}
    .caret{animation:none;}
  }
</style>
</head>
<body>

  <div class="topbar">
    <div class="brand"><span class="dot">◆</span> curriculo<span style="color:var(--text-faint)">.ai</span></div>
    <nav class="tabs">
      <button class="active">Editor</button>
      <button>Templates</button>
      <button>Histórico</button>
    </nav>
    <div class="topbar-right">
      <div class="compile-status ok" id="compileStatus">
        <span class="dot"></span> <span id="compileLabel">compilado · 0.7s</span>
      </div>
      <button class="btn-export">Exportar PDF ↓</button>
    </div>
  </div>

  <div class="layout">

    <!-- EDITOR -->
    <div class="editor">
      <p class="editor-hint">// cada seção mapeia direto pra uma parte do <span class="accent">template .tex</span> — o preview ao lado é o PDF real, atualizado a cada edição.</p>

      <details class="section-block" open>
        <summary class="section-head">
          <span class="section-num">01</span>
          <span class="section-name">Cabeçalho</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <div class="field">
            <label>Nome</label>
            <input type="text" id="inName" value="Raí Levi" oninput="sync()">
          </div>
          <div class="field">
            <label>Cargo / subtítulo</label>
            <input type="text" id="inRole" value="Backend Engineer · Go & Java" oninput="sync()">
          </div>
        </div>
      </details>

      <details class="section-block" open>
        <summary class="section-head">
          <span class="section-num">02</span>
          <span class="section-name">Contato</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <div class="field-row">
            <div class="field">
              <label>Email</label>
              <input type="text" id="inEmail" value="rlevi.dev@gmail.com" oninput="sync()">
            </div>
            <div class="field">
              <label>Cidade</label>
              <input type="text" id="inCity" value="Fortaleza, CE" oninput="sync()">
            </div>
          </div>
          <div class="field">
            <label>LinkedIn</label>
            <input type="text" id="inLinkedin" value="linkedin.com/in/rlevidev" oninput="sync()">
          </div>
          <div class="field">
            <label>GitHub</label>
            <input type="text" id="inGithub" value="github.com/rlevidev" oninput="sync()">
          </div>
          <button class="add-entry">+ adicionar contato</button>
        </div>
      </details>

      <details class="section-block">
        <summary class="section-head">
          <span class="section-num">03</span>
          <span class="section-name">Educação</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <div class="entry-card">
            <div class="field">
              <label>Instituição</label>
              <input type="text" id="inSchool" value="Universidade Federal do Ceará" oninput="sync()">
            </div>
            <div class="field-row">
              <div class="field">
                <label>Curso</label>
                <input type="text" id="inDegree" value="Ciência da Computação" oninput="sync()">
              </div>
              <div class="field">
                <label>Período</label>
                <input type="text" id="inEduDate" value="2021 — 2025" oninput="sync()">
              </div>
            </div>
          </div>
          <button class="add-entry">+ adicionar educação</button>
        </div>
      </details>

      <details class="section-block">
        <summary class="section-head">
          <span class="section-num">04</span>
          <span class="section-name">Skills</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <div class="field">
            <label>Linguagens & tecnologias</label>
            <textarea id="inSkills" rows="2" oninput="sync()">Go, Java, PostgreSQL, Docker, LaTeX, REST, Git</textarea>
          </div>
        </div>
      </details>

      <details class="section-block" open>
        <summary class="section-head">
          <span class="section-num">05</span>
          <span class="section-name">Experiência profissional</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <div class="entry-card">
            <div class="field-row">
              <div class="field">
                <label>Empresa</label>
                <input type="text" id="inCompany" value="Rinha de Backend" oninput="sync()">
              </div>
              <div class="field">
                <label>Período</label>
                <input type="text" id="inJobDate" value="2026" oninput="sync()">
              </div>
            </div>
            <div class="field">
              <label>Função</label>
              <input type="text" id="inTitle" value="Desenvolvedor Backend" oninput="sync()">
            </div>
            <div class="field">
              <label>Conquista mensurável</label>
              <textarea id="inBullet" rows="2" oninput="sync()">Implementei busca por similaridade vetorial contra 100k registros sob limites estritos de CPU e memória</textarea>
            </div>
          </div>
          <button class="add-entry">+ adicionar experiência</button>
        </div>
      </details>

      <details class="section-block">
        <summary class="section-head">
          <span class="section-num">06</span>
          <span class="section-name">Projetos</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <button class="add-entry">+ adicionar projeto</button>
        </div>
      </details>

      <details class="section-block">
        <summary class="section-head">
          <span class="section-num">07</span>
          <span class="section-name">Certificados</span>
          <span class="section-chevron">▸</span>
        </summary>
        <div class="section-body">
          <button class="add-entry">+ adicionar certificado</button>
        </div>
      </details>
    </div>

    <!-- PREVIEW -->
    <div class="preview-pane">
      <div class="paper-wrap">
        <div class="paper-toprow">
          <span class="paper-filename">curriculo.pdf</span>
          <span class="ats-badge">● ATS score 92</span>
        </div>

        <div class="paper">
          <span class="sprocket left"><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span></span>
          <span class="sprocket right"><span></span><span></span><span></span><span></span><span></span><span></span><span></span><span></span></span>

          <div class="cv-header">
            <p class="cv-name" id="pvName">Raí Levi</p>
            <p class="cv-role" id="pvRole">Backend Engineer · Go & Java</p>
            <p class="cv-contacts">
              <span id="pvEmail">rlevi.dev@gmail.com</span>
              <span id="pvCity">Fortaleza, CE</span>
              <span id="pvLinkedin">linkedin.com/in/rlevidev</span>
              <span id="pvGithub">github.com/rlevidev</span>
            </p>
          </div>

          <div class="cv-section">
            <h2>Educação</h2>
            <div class="cv-entry">
              <div class="cv-entry-top">
                <span class="cv-entry-title" id="pvSchool">Universidade Federal do Ceará</span>
                <span class="cv-entry-meta" id="pvEduDate">2021 — 2025</span>
              </div>
              <p class="cv-entry-sub" id="pvDegree">Ciência da Computação</p>
            </div>
          </div>

          <div class="cv-section">
            <h2>Skills técnicas</h2>
            <p class="cv-skills"><b>Linguagens & tecnologias:</b> <span id="pvSkills">Go, Java, PostgreSQL, Docker, LaTeX, REST, Git</span></p>
          </div>

          <div class="cv-section">
            <h2>Experiência profissional</h2>
            <div class="cv-entry">
              <div class="cv-entry-top">
                <span class="cv-entry-title" id="pvCompany">Rinha de Backend</span>
                <span class="cv-entry-meta" id="pvJobDate">2026</span>
              </div>
              <p class="cv-entry-sub" id="pvTitle">Desenvolvedor Backend</p>
              <ul><li id="pvBullet">Implementei busca por similaridade vetorial contra 100k registros sob limites estritos de CPU e memória<span class="caret"></span></li></ul>
            </div>
          </div>
        </div>

        <p class="paper-footnote">renderizado com pdflatex · texlive-latex-extra</p>
      </div>
    </div>

  </div>

<script>
  function sync(){
    document.getElementById('pvName').textContent = document.getElementById('inName').value;
    document.getElementById('pvRole').textContent = document.getElementById('inRole').value;
    document.getElementById('pvEmail').textContent = document.getElementById('inEmail').value;
    document.getElementById('pvCity').textContent = document.getElementById('inCity').value;
    document.getElementById('pvLinkedin').textContent = document.getElementById('inLinkedin').value;
    document.getElementById('pvGithub').textContent = document.getElementById('inGithub').value;
    document.getElementById('pvSchool').textContent = document.getElementById('inSchool').value;
    document.getElementById('pvEduDate').textContent = document.getElementById('inEduDate').value;
    document.getElementById('pvDegree').textContent = document.getElementById('inDegree').value;
    document.getElementById('pvSkills').textContent = document.getElementById('inSkills').value;
    document.getElementById('pvCompany').textContent = document.getElementById('inCompany').value;
    document.getElementById('pvJobDate').textContent = document.getElementById('inJobDate').value;
    document.getElementById('pvTitle').textContent = document.getElementById('inTitle').value;
    document.getElementById('pvBullet').firstChild.textContent = document.getElementById('inBullet').value;

    // fake compile cycle, like a real pdflatex run reacting to the edit
    const status = document.getElementById('compileStatus');
    const label = document.getElementById('compileLabel');
    clearTimeout(window._compileTimer);
    status.classList.remove('ok');
    status.classList.add('busy');
    label.textContent = 'compilando…';
    window._compileTimer = setTimeout(() => {
      status.classList.remove('busy');
      status.classList.add('ok');
      const t = (Math.random()*0.6+0.4).toFixed(1);
      label.textContent = 'compilado · ' + t + 's';
    }, 550);
  }
</script>

</body>
</html>
```