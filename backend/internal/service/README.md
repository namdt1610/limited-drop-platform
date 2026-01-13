# Service (Symbicode)

This package contains business logic and is the single source of truth for Symbicode behaviour.

Guidelines:
- Keep one implementation for Symbicode generation, verification, and QR helpers here (DRY).
- External integrations (queues, email, 3rd-party APIs) should remain in `internal/integrations` and call into this service when necessary (SOLID).
- Prefer small, well-tested service methods over duplicated free functions spread across packages (KIS).
