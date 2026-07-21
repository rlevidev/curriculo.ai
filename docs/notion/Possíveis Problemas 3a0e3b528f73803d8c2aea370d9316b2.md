# Possíveis Problemas

#### Cold start do Render.

- Primeiro acesso demora → Para resolver: Keep-alive ajuda, não elimina.

#### Bundle do TeXLive pesado.

- Deploy inicial no Render pode ser lento (imagem grande).

#### Sem escape de LaTeX quebra tudo.

- Nome com `&`, `%` ou `_` derruba compilação. Já mitigado no código, mas teste bastante.

#### pdflatex trava sem erro claro.

- Template malformado gera log gigante. Timeout de 20s ajuda, mas debug é chato.

#### CORS mal configurado.

- Esqueceu de atualizar `ALLOWED_ORIGIN` no Render? Frontend quebra silenciosamente.

#### URL do backend hardcoded no JS.

- Trocar de plano ou serviço exige editar `app.js` manualmente.

#### Sem rate limit.

- Qualquer pessoa pode spammar `/gerar-pdf`. Gasta as 750h grátis rápido.
Colocar um token bucket por IP (5 req/min).

#### Sem validação de tamanho de input.

- Descrição gigante de experiência pode gerar PDF de várias páginas ou travar o pdflatex.

#### Educação/Experiência fixas no frontend.

- Ainda não editável de verdade. Formulário incompleto.

#### Sem HTTPS mismatch, mas…

- Se esquecer `https://` no `API_URL`, erro de mixed content trava tudo.

#### Maior risco de todos: abandono do projeto.

- Portfolio só vale se terminar. Prioriza isso.

### Para evitar Bundle do TeXLive pesado esses são os pacotes do Dockerfile necessário para criar o currículo seguindo o template:

| Pacote | Tamanho | Necessário? |
| --- | --- | --- |
| `texlive-latex-base` | ~200 MB | ✅ **SIM** — tem geometry, fontenc, inputenc, hyperref, array |
| `texlive-latex-recommended` | ~200 MB | ✅ **SIM** — tem xcolor, parskip (e é puxado como dep do extra mesmo assim) |
| `texlive-latex-extra` | ~350 MB | ✅ **SIM** — tem enumitem, titlesec |
| `texlive-fonts-recommended` | ~250 MB | ✅ **SIM** — tem as fontes Libertine (`LinLibertine*.tfm`, `*.pfb`) — sem isso `\usepackage{libertine}` quebra no pdflatex |

**Total:** ~1 GB descomprimido na imagem Docker. Não dá pra cortar nada disso sem perder funcionalidade do template.

```docker
# TeX Live — mínimo pra compilar o template de currículo
RUN apt-get update && apt-get install -y --no-install-recommends \
    texlive-latex-base \
    texlive-latex-recommended \
    texlive-latex-extra \
    texlive-fonts-recommended \
    && rm -rf /var/lib/apt/lists/*
```