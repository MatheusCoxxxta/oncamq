# Changelog - Português (Brasil)

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

---

### [v0.1.7] - 2026-01-02
**Nome Amigável:** Ajustes de Precisão
**Sinopse:** Pequena correção técnica na configuração de distribuição (GoReleaser) para garantir que os builds sejam gerados corretamente.

- **Fix:** Correção de indentação no arquivo `.goreleaser.yaml`.

### [v0.1.6] / [v0.1.5] - 2026-01-02
**Nome Amigável:** Refinamento de Entrega
**Sinopse:** Simplificação do processo de release, removendo passos redundantes para tornar o ciclo de publicação mais ágil.

- **Fix:** Simplificação do passo do GoReleaser removendo argumentos desnecessários.

### [v0.1.4] - 2026-01-02
**Nome Amigável:** Automação de Binários
**Sinopse:** Introdução do GoReleaser para automatizar a criação de binários e releases multiplataforma.

- **Feat:** Adição da configuração inicial do `.goreleaser.yaml`.

### [v0.1.3] - 2026-01-02
**Nome Amigável:** Consolidação da Identidade
**Sinopse:** Padronização final do nome do pacote para `oncamq`, garantindo consistência entre o código e o repositório.

- **Fix:** Atualização do nome do pacote de `main` para `oncamq` no `worker.go`.

### [v0.1.2] - 2026-01-02
**Nome Amigável:** Ajuste de Escopo
**Sinopse:** Correção temporária no escopo do pacote para testes de importação interna.

- **Fix:** Mudança de pacote de `oncamq` para `main`.

### [v0.1.1] - 2026-01-02
**Nome Amigável:** Correção de Dependências
**Sinopse:** Resolução de conflitos de programas e pacotes no caminho de importação (Go import path).

- **CI:** Correção de conflito de pacote/programa no caminho de importação.
- **Fix:** Atualização da mensagem de versão de release no workflow do GitHub Actions.

### [v0.1.0] - 2026-01-02
**Nome Amigável:** O Nascimento do OncaMQ
**Sinopse:** Lançamento consolidado do core do consumer.

- **Feat:** Implementação de padrões idiomáticos de worker Go (remoção de estado global, contexto explícito).
- **Feat:** Ações de sucesso: gerenciamento de filas completas, valores de retorno e tentativas.
- **Feat:** Integração com GitHub Actions para publicação automática.
- **Docs:** Documentação expandida com guia de contribuição e exemplos realistas.
- **Chore:** Renomeação do módulo de `go-bullmq-consumer` para `oncamq`.
