# Changelog - Português (Brasil)

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

---

### [v0.1.7] Correção na Configuração de CI/Release (2026-01-02)

- **Fix:** Correção de indentação no arquivo `.goreleaser.yaml`.

### [v0.1.6] / [v0.1.5] Otimização do Pipeline de Build (2026-01-02)

- **Fix:** Simplificação do passo do GoReleaser removendo argumentos desnecessários.

### [v0.1.4] Suporte para Distribuição Multi-plataforma (2026-01-02)

- **Feat:** Adição da configuração inicial do `.goreleaser.yaml`.

### [v0.1.3] Padronização do Namespace do Pacote (2026-01-02)

- **Fix:** Atualização do nome do pacote de `main` para `oncamq` no `worker.go`.

### [v0.1.2] Ajuste de Visibilidade de Escopo (2026-01-02)

- **Fix:** Mudança de pacote de `oncamq` para `main`.

### [v0.1.1] Resolução de Conflitos de Importação e Workflow (2026-01-02)

- **CI:** Correção de conflito de pacote/programa no caminho de importação.
- **Fix:** Atualização da mensagem de versão de release no workflow do GitHub Actions.

### [v0.1.0] Implementação Base e Padrões Idiomáticos (Gopher Way) (2026-01-02)

- **Feat:** Implementação de padrões idiomáticos de worker Go (remoção de estado global, contexto explícito).
- **Feat:** Ações de sucesso: gerenciamento de filas completas, valores de retorno e tentativas.
- **Feat:** Integração com GitHub Actions para publicação automática.
- **Docs:** Documentação expandida com guia de contribuição e exemplos realistas.
- **Chore:** Renomeação do módulo de `go-bullmq-consumer` para `oncamq`.
