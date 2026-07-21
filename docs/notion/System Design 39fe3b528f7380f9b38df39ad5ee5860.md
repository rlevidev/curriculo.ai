# System Design

### **Fluxo de dados**:

- Usuário digita → App monta uma estrutura de dados (JSON) do currículo → Alimenta o preview e o gerador LaTeX.

### Componentes:

- Formulário de input
- Motor de preview (renderiza em tempo real)
- Motor de geração de LaTeX para PDF

### Modelo de dados (schema do currículo):

- Nome
- Contatos
- Educação
- Skills
- Experiências

### Posso utilizar o GitHub Pages?

Sim e não, pois para compilar LaTeX pra PDF precisaremos de um motor rodando em algum servidor.

Adotariamos um roteamonto de páginas estáticas:

- Github Pages (rlevidev.github.io/repo) → Frontend
- Render (Instalando o TeXLive via docker) → Backend
- A partir dai teremos um CORS no backend:

```go
w.Header().Set("Access-Control-Allow-Origin", "https://seuusuario.github.io")
w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

```json
fetch("https://criador-curriculo.onrender.com/gerar-pdf", { method: "POST", ... })
```

#### Problemática para pontuar:

- O Render "Dorme" após inatividade, cold start pode levar ~1 min
- Para resolver isto, podemos utilizar um workaround comum. Ele básicamente mantém o serviço “aquecido” com pings.
- Como fazer? Github Actions. Cria um workflow agendado que faz um `curl` no backend periodicamente:

```json
# .github/workflows/keep-alive.yml
name: Keep Render Alive
on:
  schedule:
    - cron: '*/14 * * * *'  # a cada 14 minutos
jobs:
  ping:
    runs-on: ubuntu-latest
    steps:
      - name: Ping backend
        run: curl -f https://criador-curriculo.onrender.com/health
```

- Como fallback podemos utilizar um loading state no frontend.

### Quais linguagens de programação adotar?

- Frontend: React ou HTML com JS puro com preview em tempo real
- Backend: Como iremos utilizar Docker, a linguagem que mais faz sentido é o Go. Ele perfoma bem e ajuda a criar container pequeno.

### Inputs necessários

#### Cabeçalho:

| Campo | Tipo | Exemplo | Opcional |
| --- | --- | --- | --- |
| `nome` | string | `"Allen Tran"` |  |
| `cargo` | string | `"Software Engineer"` |  |
| `contato.email` | string | `"allen.tran@email.com"` | Sim |
| `contato.telefone` | string | `"(555) 123-4567"` | Sim |
| `contato.linkedin_url` | string (URL completa) | `"<https://linkedin.com/in/allen-tran>"` | Sim |
| `contato.github_url` | string (URL completa) | `"<https://github.com/allentran>"` | Sim |
| `contato.localizacao` | string | `"California"` | Sim |

#### Educação (Opcional: Sim):

| Campo | Tipo | Opcional |
| --- | --- | --- |
| `instituicao` | string |  |
| `curso` | string |  |
| `periodo` | string |  |
| `notas` | lista de strings (bullets opcionais, ex: "Focus on AI") | Sim |

#### Skills:

| Campo | Tipo |
| --- | --- |
| `linguagens` | lista de strings — vira a linha "Languages:" |
| `tecnologias` | lista de strings — vira a linha "Technologies:" |

#### Experiência (Opcional: Sim):

| Campo | Tipo | Opcional |
| --- | --- | --- |
| `empresa` | string |  |
| `local` | string | Sim |
| `periodo` | string |  |
| `cargo` | string |  |
| `bullets` | lista de strings (cada uma vira um `\item`) | Sim |

#### Projetos:

| Campo | Tipo | Opcional |
| --- | --- | --- |
| `nome` | string |  |
| `link_url` | string (URL completa, opcional) |  |
| `link_label` | string — texto do link (ex: "GitHub Repository") | Sim, mas por padrão vai ser o texto “Link” |
| `bullets` | lista de strings | Sim |

#### Idiomas (Opcional: Sim):

| Campo | Tipo |
| --- | --- |
| `idioma` | string (ex: "Portuguese") |
| `nivel` | string (ex: "Native") |

#### Certificações (Opcional: Sim):

| Campo | Tipo |
| --- | --- |
| `certificacoes` | lista de strings |

#### Exemplo de Request JSON:

```json
{
  "nome": "Allen Tran",
  "cargo": "Software Engineer",
  "contato": {
    "email": "allen.tran@email.com",
    "telefone": "(555) 123-4567",
    "linkedin_url": "<https://linkedin.com/in/allen-tran>",
    "github_url": "<https://github.com/allentran>",
    "localizacao": "California"
  },
  "educacao": [
    { "instituicao": "Arizona State University", "curso": "Bachelor of Computer Science & Master of Computer Science", "periodo": "2017 -- Sep 2022", "notas": ["Focus on Artificial Intelligence"] }
  ],
  "linguagens": ["Python", "Go (Golang)", "Java"],
  "tecnologias": ["Node.js", "React.js", "AWS", "Docker"],
  "experiencias": [
    { "empresa": "Capital One", "local": "Chicago, IL", "periodo": "June 2022 -- August 2022", "cargo": "Software Engineer Intern",
      "bullets": ["Designed & developed full-stack application..."], "tecnologias": ["Go (Golang)", "AWS Lambda"] }
  ],
  "projetos": [
    { "nome": "Drop It | Google Drive Replica", "link_url": "<https://github.com/rlevidev/repository>", "link_label": "GitHub Repository", "bullets": ["Individually developed dropbox application..."] }
  ],
  "idiomas": [{ "idioma": "Portuguese", "nivel": "Native" }],
  "certificacoes": ["AWS Certified Solutions Architect Associate (2023)"]
}
```

### Observações de segurança:

- Nos campos de inputs de URLs que vai para `\href{…}`: Adicionar validação por allowlist de regex. Qualquer URL fora do padrão `https?://...` com caracteres seguros é rejeitada antes de chegar no template. Isso fecha a brecha de alguém injetar LaTeX arbitrário via esse campo.
- Rate Limiter tem vazamento de memória: Os buckets de IPs antigos nunca são removidos do mapa. Com tráfego suficiente (ou um scanner batendo com IPs variados), isso cresce para sempre até estourar a memória do container. Adicione uma goroutine de limpeza periódica.
- `X-Forwarded-For` confiado sem sanitização: O código pega o header inteiro como está, mas esse header é enviado pelo próprio cliente, então qualquer atacante pode mandar um valor diferente a cada requisição e furar o rate limit completamente, já que cada valor novo vira uma “chave” nova no mapa de buckets. O rate limit vira decorativo. É preciso achar uma forma de corrigir isto.
- `pdflatex` sem trava explícita de `shell-escape`: Por padrão o `pdflatex` já vem com `\write18` desabilitado, então comandos LaTeX maliciosos (tipo `\immediate\write18{rm -rf /}`) não deveriam executar mesmo que alguém conseguisse injetar LaTeX bruto em algum campo. Mas isso depende da configuração padrão da distribuição, e não custa nada travar explicitamente na chamada, em vez de confiar no padrão. É bom adicionar `-no-shell-escape` explicitamente na chamada, defesa em profundidade, não custa nada.
- Log de erro pode vazar dados pessoal: O `stderr` do `pdflatex` em caso de erro de compilação pode ecoar trechos do `.tex` (nome, email, telefone e etc…) na mensagem de erro. Isso vai direto pro log do Render, que persiste e nem sempre só você acessa. Só deve loga detalhe completo se `DEBUG_LOGS=1` estiver setado. Em produção fica só um alerta genérico.

### Botão de Download

- Fluxo completo

```json
[Usuário preenche form]
   ↓
[Clica "Baixar PDF"]
   ↓
[Frontend: JSON.stringify(dados) → POST /gerar-pdf]
   ↓
[Backend: valida → renderiza template LaTeX → pdflatex → PDF em memória/disco]
   ↓
[Backend: responde com bytes do PDF + headers corretos]
   ↓
[Frontend: recebe blob → cria link temporário → dispara download]
```

- Mais ou menos como se espera o código do Front-End

```jsx
const btnBaixar = document.getElementById('btn-baixar-pdf');

async function baixarCurriculo(dadosCurriculo) {
  const btn = document.getElementById('btn-baixar-pdf');
  const textoOriginal = btn.textContent;

  btn.disabled = true;
  btn.textContent = 'Gerando PDF...';

  try {
    const response = await fetch('https://criador-curriculo.onrender.com/gerar-pdf', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(dadosCurriculo),
    });

    if (!response.ok) {
      // Backend deve devolver JSON de erro nesses casos, não PDF
      const erro = await response.json().catch(() => null);
      throw new Error(erro?.mensagem || `Erro ${response.status} ao gerar PDF`);
    }

    const blob = await response.blob();

    // Nome do arquivo: tenta pegar do header, senão usa fallback com nome do usuário
    const nomeArquivo = extrairNomeArquivo(response.headers.get('Content-Disposition'))
      || `curriculo-${slugify(dadosCurriculo.nome)}.pdf`;

    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = nomeArquivo;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);

  } catch (err) {
    console.error(err);
    mostrarErroUsuario(err.message.includes('Failed to fetch')
      ? 'Servidor demorando para responder (cold start). Tente novamente em ~1 minuto.'
      : err.message);
  } finally {
    btn.disabled = false;
    btn.textContent = textoOriginal;
  }
}

function extrairNomeArquivo(contentDisposition) {
  if (!contentDisposition) return null;
  const match = contentDisposition.match(/filename="?([^"]+)"?/);
  return match ? match[1] : null;
}

function slugify(texto) {
  return texto.toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '').replace(/\s+/g, '-');
}
```

Ponto importante: `fetch` não trata `4xx/5xx` como erro automaticamente, precisa checar `response.ok` manualmente, senão o `.blob()` vira um PDF corrompido contendo uma mensagem de erro JSON.

- Mais ou menos como se espera o Back-End

```go
func gerarPDFHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOWED_ORIGIN"))

    var dados CurriculoInput
    if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
        respondJSON(w, http.StatusBadRequest, ErrorResponse{Mensagem: "JSON inválido"})
        return
    }

    if err := validar(dados); err != nil {
        respondJSON(w, http.StatusBadRequest, ErrorResponse{Mensagem: err.Error()})
        return
    }

    pdfBytes, err := compilarLatex(dados) // com timeout de 20s já previsto
    if err != nil {
        log.Printf("erro ao compilar: %v", sanitizarLog(err)) // sem vazar dados pessoais
        respondJSON(w, http.StatusInternalServerError, ErrorResponse{Mensagem: "Falha ao gerar PDF"})
        return
    }

    nomeArquivo := fmt.Sprintf("curriculo-%s.pdf", slugify(dados.Nome))
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, nomeArquivo))
    w.Header().Set("Content-Length", strconv.Itoa(len(pdfBytes)))
    w.WriteHeader(http.StatusOK)
    w.Write(pdfBytes)
}
```

- Detalhes para lembrar
1. Estado de loading visível - cold start pode levar ~1min mesmo com keep-alive (o ping de 14 em 14 min não cobre 100% dos casos, ex: primeiro deploy ou serviço reiniciado). O botão precisa comunicar isso, não só “Gerando…”.
2. Content-Disposition é essencial - sem ele, o navegador tenta abrir o PDF inline ou usa um nome genérico tipo download.pdf 
3. Erro precisa vir como JSON, PDF só no caminho feliz - senão o frontend não consegue distinguir “PDF gerado” de “erro” só olhando o blob 
4. Timeout no fetch do frontend - sem isso, se o backend travar (pdflatex sem erro claro, que vocês já mapearam), o usuário fica com “Gerando PDF…” infinito. Vale um AbortController com uns 30s.
5. Double-click guard - desabilitar o botão durante a requisição evita múltiplas chamadas simultâneas (relevante já que não tem rate limit robusto ainda).

### Preview pode Sobrecarregar o Back-End

Para resolver isto, é recomendado um preview 100% client-side.

A ideia central: o preview não precisa ser o PDF real. Ele precisa só parecer com o resultado final o suficiente pra o usuário confiar no que está escrevendo.

```json
[Usuário digita] → [JS monta o preview em HTML/CSS que imita o layout do template LaTeX]
                     (client-side, zero requisição ao backend)

[Usuário clica "Baixar PDF"] → [Só aí chama o backend, compila LaTeX de verdade, gera o PDF real]
```

Prático de fazer porque:

- O template de vocês já é bem estruturado (header centralizado, seções com `\titlerule`, bullets, `\hfill` pra alinhar datas à direita). Isso reproduz bem em HTML/CSS com flexbox — é literalmente o mesmo padrão visual que vocês já usam no `index.html` do portfólio.
- Reaproveita o schema JSON que já está definido no System Design — o mesmo objeto que alimenta o preview alimenta o request pro backend.
- Debounce de ~150-300ms no preview client-side já é suficiente pra não travar a digitação, sem custo nenhum de rede.

Vamos estruturar isso em camadas: schema → template de render → função de update. A ideia é que o mesmo JSON do System Design vire tanto o preview quanto o payload pro backend.

- Estrutura HTML base (esqueleto fixo, mimetiza o `\documentclass` + geometry)

```html
<div id="preview-container" class="preview-page">
  <div id="preview-header" class="preview-header"></div>
  <div id="preview-education" class="preview-section"></div>
  <div id="preview-skills" class="preview-section"></div>
  <div id="preview-experience" class="preview-section"></div>
  <div id="preview-projects" class="preview-section"></div>
  <div id="preview-languages" class="preview-section"></div>
  <div id="preview-certifications" class="preview-section"></div>
</div>
```

```css
.preview-page {
  /* A4 em pixels a 96dpi: 210mm x 297mm ≈ 794px x 1123px */
  width: 794px;
  min-height: 1123px;
  padding: 57.6px; /* margin=0.6in do template */
  background: white;
  color: #000;
  font-family: Arial, Helvetica, sans-serif; /* equivalente ao \usepackage{helvet} */
  font-size: 11pt;
  box-shadow: 0 4px 24px rgba(0,0,0,0.3); /* só estética, não é parte do PDF */
}

.preview-section-title {
  font-size: 1.1em;
  font-weight: bold;
  text-transform: uppercase;
  border-bottom: 1px solid #000; /* equivalente ao \titlerule */
  margin: 12px 0 6px;
}

.preview-row-split {
  display: flex;
  justify-content: space-between; /* equivalente ao \hfill entre dois elementos */
}
```

O ponto chave aqui: **794px de largura fixa simula a proporção do A4**, então o que quebra de linha no preview quebra de forma parecida no PDF real. Não é perfeito (fontes do navegador ≠ métricas do LaTeX), mas é suficiente pra dar confiança visual.

- Função de render por seção (pega o JSON, devolve HTML)

```jsx
function renderHeader(dados) {
  const contatos = [
    dados.contato?.email,
    dados.contato?.telefone,
    dados.contato?.linkedin_url ? `<a href="${dados.contato.linkedin_url}">LinkedIn</a>` : null,
    dados.contato?.github_url ? `<a href="${dados.contato.github_url}">GitHub</a>` : null,
    dados.contato?.localizacao,
  ].filter(Boolean).join(' &nbsp;|&nbsp; ');

  return `
    <div style="text-align:center">
      <div style="font-size:2em; font-weight:bold">${escapeHtml(dados.nome || '')}</div>
      <div style="margin-top:4px">${escapeHtml(dados.cargo || '')} &nbsp;|&nbsp; ${contatos}</div>
    </div>
  `;
}

function renderExperiencia(experiencias = []) {
  if (!experiencias.length) return '';
  const items = experiencias.map(exp => `
    <div style="margin-bottom:8px">
      <div class="preview-row-split">
        <strong>${escapeHtml(exp.empresa)}</strong>
        <em>${escapeHtml(exp.periodo)}</em>
      </div>
      <div class="preview-row-split">
        <em>${escapeHtml(exp.cargo)}</em>
        <span>${escapeHtml(exp.local || '')}</span>
      </div>
      <ul>
        ${(exp.bullets || []).map(b => `<li>${escapeHtml(b)}</li>`).join('')}
      </ul>
    </div>
  `).join('');

  return `<div class="preview-section-title">Professional Experience</div>${items}`;
}
```

**`escapeHtml` é obrigatório aqui** — não pelo motivo do LaTeX (aquele é outro escape, do lado do backend), mas porque sem isso um usuário digitando `<script>` no campo de bullet quebra o preview ou abre XSS. Repara que os dois problemas de "caractere perigoso" são independentes: o frontend escapa pra HTML, o backend escapa pra LaTeX. São camadas diferentes, cada uma cuida do seu formato.

```jsx
function escapeHtml(str = '') {
  const div = document.createElement('div');
  div.textContent = str;
  return div.innerHTML;
}
```

- Orquestrador com debounce

```jsx
let debounceTimer;

function onFormChange(dados) {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => atualizarPreview(dados), 200);
}

function atualizarPreview(dados) {
  document.getElementById('preview-header').innerHTML = renderHeader(dados);
  document.getElementById('preview-education').innerHTML = renderEducacao(dados.educacao);
  document.getElementById('preview-skills').innerHTML = renderSkills(dados.linguagens, dados.tecnologias);
  document.getElementById('preview-experience').innerHTML = renderExperiencia(dados.experiencias);
  document.getElementById('preview-projects').innerHTML = renderProjetos(dados.projetos);
  document.getElementById('preview-languages').innerHTML = renderIdiomas(dados.idiomas);
  document.getElementById('preview-certifications').innerHTML = renderCertificacoes(dados.certificacoes);
}
```

Cada campo do form dispara `onFormChange(dadosAtuais)` no `input` event. O debounce de 200ms evita re-render a cada tecla mas ainda parece instantâneo pro usuário.

- Ponto de atenção: overflow de página

O template LaTeX tem `\raggedbottom` e penalidades contra quebra de página, mas o preview em HTML normalmente não simula páginas múltiplas — ele só cresce verticalmente. Para resolver isto, podemos deixar o preview crescer livremente e adicionar um aviso visual tipo "conteúdo pode ocupar mais de 1 página" quando `preview-page.scrollHeight > 1123px` (isso resolve também o problema que vocês já mapearam de "sem validação de tamanho de input").

```jsx
function verificarOverflow() {
  const preview = document.getElementById('preview-container');
  const aviso = document.getElementById('aviso-overflow');
  
  const excedeu = preview.scrollHeight > 1123;
  aviso.style.display = excedeu ? 'flex' : 'none';
  
  if (excedeu) {
    const paginasEstimadas = Math.ceil(preview.scrollHeight / 1123);
    aviso.textContent = `⚠ Conteúdo pode ocupar ~${paginasEstimadas} páginas. Considere reduzir o texto.`;
  }
}

// chama no fim do atualizarPreview(), depois do DOM já ter sido atualizado
function atualizarPreview(dados) {
  // ...renders de cada seção...
  requestAnimationFrame(verificarOverflow); // garante que o layout já recalculou
}
```

O `requestAnimationFrame` importa aqui — ler `scrollHeight` imediatamente após trocar `innerHTML` às vezes pega o valor antigo, porque o navegador ainda não recalculou o layout.

Um detalhe: isso é um **aviso**, não um bloqueio. Não trava o botão de baixar PDF, só avisa. Quem decide se aceita 3 páginas é o usuário — o job do frontend é só não deixar ele ser surpreendido.

- **Atenção**
1. **O aviso client-side não substitui a validação do backend.** Já foi anotado isso no "Possíveis Problemas" ("sem validação de tamanho de input"), mas vale reforçar: são duas camadas com propósitos diferentes. O aviso no preview é UX (evita surpresa). O backend precisa rejeitar ou truncar de verdade — porque alguém pode simplesmente chamar `/gerar-pdf` direto via curl, ignorando o frontend inteiro. Um `len(bullet) > N` e `len(bullets) > M` no Go, antes de montar o `.tex`, fecha esse buraco.
2. **Rascunho não persiste.** Se o usuário recarrega a página ou a aba crasha, perde tudo que digitou. Um `localStorage.setItem('curriculo-draft', JSON.stringify(dados))` no mesmo debounce que já dispara o preview resolve isso de graça — mesmo evento, mesmo timing, custo perto de zero. Dado o risco que vocês mesmos marcaram como "maior de todos: abandono do projeto", eu preferiria fechar isso cedo: é o tipo de coisa pequena que, se faltar, faz o usuário desistir no meio do preenchimento.
3. **Divergência preview vs. PDF real vai acontecer, e tudo bem.** Arial no navegador não tem exatamente a métrica do Helvetica do LaTeX. Em algum momento o número de páginas do preview e do PDF final vai divergir um pouco. Vale um texto pequeno e honesto perto do preview tipo "prévia aproximada — o PDF final pode variar levemente", só pra ninguém abrir um issue achando que é bug.
4. **Prioridade dos itens do "Possíveis Problemas".** Vocês já têm um bom inventário de riscos ali, mas está tudo no mesmo nível. Eu separaria em "bloqueia lançamento" (rate limiter com vazamento de memória, `X-Forwarded-For` sem sanitização, escape de LaTeX) vs "pode esperar" (CORS, cold start já tem workaround, URL hardcoded). Os dois primeiros são vulnerabilidade de verdade, não só chateação operacional.