# Avaliação do projeto curriculo.ai

> Avaliação técnica realizada em 23/08/2026.
>
> Verificado: `go vet` limpo, `go test` passando, `tsc && vite build` OK, oxlint com 1 warning. Binário e `.env.local` não estão no git ✓.

---

## 🔴 Bugs críticos

### 1. Injeção de LaTeX — `texEscape` existe mas nunca é usado

O backend define e testa `texEscape()` (`backend/main.go:99`), mas o template nunca o chama — todos os campos vão crus para o `.tex`:

```go
\section*{<[ .Name ]>}   // deveria ser <[ texEscape .Name ]>
```

**Impacto:** usuário digita `{`, `\input{/etc/passwd}` ou aspas desbalanceadas → compilação quebra (500), conteúdo de arquivos renderizado no PDF, DoS. O `-no-shell-escape` mitiga execução, mas não leitura de arquivos.

**Solução:** escapar na entrada, uma vez só:

```go
data.Name = texEscape(data.Name)
data.Title = texEscape(data.Title)
// ... ou um normalize() que percorre todos os campos/listas antes do Execute
```

### 2. Rate limiter quebrado nos dois sentidos

`backend/main.go:33` usa `r.RemoteAddr` como chave — que é `IP:porta`. Cada nova conexão TCP tem porta efêmera diferente → **visitante novo com 5 tokens frescos → limite trivialmente burlável**. E quando a porta coincide (keep-alive), **os tokens nunca repõem** → usuário legítimo bloqueado pra sempre até reiniciar. De bônus, o mapa `visitors` cresce sem limite (memory leak) e `lastSeen` nunca é usado.

**Solução:** chavear só pelo IP + repor tokens por janela de tempo:

```go
func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		mu.Lock()
		v, ok := visitors[ip]
		if !ok || time.Since(v.lastSeen) > time.Hour {
			v = &Visitor{tokens: 5}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		if v.tokens <= 0 {
			mu.Unlock()
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		v.tokens--
		mu.Unlock()
		next(w, r)
	}
}
```

(ou `golang.org/x/time/rate` por IP com cleanup periódico)

### 3. O PDF gerado ignora metade dos dados

O template (`main.go:138-167`) só renderiza nome, título, email, telefone, experiências, idiomas e certificações. **Faltam:** educação, projetos, skills, LinkedIn, GitHub e localização — tudo que o editor coleta e o preview mostra. Além disso, os ranges colam itens sem separador (`Inglês (Fluente)Espanhol (Básico)`). **O preview mente sobre o PDF final.**

**Solução:** completar o template com as seções faltantes e separadores (` · ` / `\item`).

---

## 🟠 Bugs médios

### 4. Campo de skills "come" vírgulas enquanto você digita

`EditorPane.tsx:229-247`: o onChange faz `split(',') → filter(Boolean)` e o estado re-renderiza com `join(', ')`. Digitar `"Go,"` vira `"Go"` — a vírgula desaparece antes de você digitar a próxima linguagem, e `"Go"` + `"Rust"` viram `"GoRust"`. Impossível usar sem colar texto pronto.

**Solução:** estado local de string no componente, commitando a lista no `onBlur`.

### 5. Draft antigo do localStorage pode quebrar o app

`App.tsx:48` faz merge raso (`{...default, ...parsed}`). Se um draft salvo por versão anterior tiver `skills: {languages: [...]}` sem `technologies`, ou experiência sem `bullets`, o app crasha em pontos sem guard: `PreviewPane.tsx:95` (`skills.languages.length`), `EditorPane.tsx:305` (`exp.bullets.map`) e `types.ts:114`.

**Solução:** função `normalize(parsed)` que valida campo a campo contra defaults (e versiona a key: `resume-draft-v2`).

### 6. Checagem de tamanho de payload é burlável

`main.go:117` confia em `r.ContentLength`, que é `-1` em requests chunked → bypass; o decode lê corpo ilimitado.

**Solução:** `r.Body = http.MaxBytesReader(w, r.Body, 50<<10)` antes do decode (a checagem de ContentLength pode sair).

### 7. Testes acoplados a estado global (flake latente)

`TestGeneratePdfHandler_Success` consome 1 token do IP default `192.0.2.1`; se rodar antes de `TestRateLimiter` (que espera 5 tokens do mesmo IP), falha. Hoje passa porque o `pdflatex` não está instalado e o teste dá skip — instale o texlive e o suite quebra.

**Solução:** helper `resetVisitors()` chamado no início de cada teste, ou injetar o limiter em vez de usar global.

---

## 🟡 Menores

| Problema | Local | Solução |
|---|---|---|
| PATH pessoal vazou no Dockerfile (`/home/rlevidev/.mimocode/bin`...) | `backend/Dockerfile:3` | `ENV PATH="/usr/local/go/bin:${PATH}"` |
| Build do Go dentro da imagem texlive (imagem gigante, cold start pior no Render free tier) | Dockerfile inteiro | Multi-stage: binário no `golang:1.24-alpine`, copiar só o binário pro estágio texlive |
| Erro de export chama `response.json()` sem guard — 502 HTML do Render vira "Unexpected token..." pro usuário | `App.tsx:254` | Checar content-type ou usar `await response.text()` com fallback |
| Timeout de 10s morto no export | `App.tsx:240-242` | Deletar |
| Warning `exhaustive-deps` (único do lint) | `App.tsx:56` | Mover interval pra effect próprio com `[serverStatus]` já cobre, ou incluir callback |
| `ioutil.*` deprecated desde Go 1.16 | `main.go` | `os.TempDir`/`os.ReadFile`/`os.WriteFile` |
| Sem timeouts no servidor (Slowloris) | `main.go:237` | `&http.Server{ReadHeaderTimeout: 5 * time.Second, ...}` |
| Tabs "Templates"/"Histórico" e `previewRef` não fazem nada | TopBar/App | Remover até existirem |

---

## Trade-offs aceitáveis (documentar, não mudar)

- **Render free tier + texlive**: cold start de ~1min — a UI até avisa ("servidor acordando..."). O multi-stage acima reduz o tamanho da imagem e melhora isso.
- **ATS score heurístico** (`/\d/` = métrica): ingênuo, mas honesto pra v1. Os pontos somam exatamente 100.
- **Só localStorage**: privacidade como feature, ao custo de não sincronizar entre dispositivos. OK.
- **`key={index}`** nas listas: aceitável aqui; pode causar glitch de foco ao remover item do meio.
- **Endpoint sem auth**: ok pra ferramenta pessoal, mas o rate limiter é a única defesa — e está quebrado (#2), por isso ele é prioridade junto com #1.

## Ordem sugerida de correção

1. `texEscape` no template (#1) — segurança
2. Rate limiter por IP com refill (#2) — segurança + usabilidade
3. Completar template do PDF (#3) — o produto mentir pro usuário é o pior bug de UX
4. Skills textarea (#4) e normalização do draft (#5) — estabilidade percebida
5. Resto conforme sobrar tempo
